// Package network provides a Terraform resource.
package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govdctypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkIsolatedResource{}
	_ resource.ResourceWithConfigure   = &networkIsolatedResource{}
	_ resource.ResourceWithImportState = &networkIsolatedResource{}
	_ resource.ResourceWithModifyPlan  = &networkIsolatedResource{}
	_ vcdNetworkIsolatedOrRouted       = &networkIsolatedResource{}
)

// NewNetworkIsolatedResource is a helper function to simplify the provider implementation.
func NewNetworkIsolatedResource() resource.Resource {
	return &networkIsolatedResource{}
}

// networkIsolatedResource is the resource implementation.
type networkIsolatedResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	org    org.Org
}

type networkIsolatedResourceModel struct {
	ID           types.String `tfsdk:"id"`
	VDC          types.String `tfsdk:"vdc"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Gateway      types.String `tfsdk:"gateway"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	DNS1         types.String `tfsdk:"dns1"`
	DNS2         types.String `tfsdk:"dns2"`
	DNSSuffix    types.String `tfsdk:"dns_suffix"`
	StaticIPPool types.Set    `tfsdk:"static_ip_pool"`
}

func (r *networkIsolatedResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	configVDC := &types.String{}
	req.Config.GetAttribute(ctx, path.Root("vdc"), configVDC)
	stateVDC := &types.String{}
	req.State.GetAttribute(ctx, path.Root("vdc"), stateVDC)
	if (configVDC.IsNull() || configVDC.IsUnknown()) && !stateVDC.IsNull() {
		if r.client.GetDefaultVDC() != stateVDC.ValueString() {
			x := &path.Paths{}
			resp.RequiresReplace = x.Append(path.Root("vdc"))
			resp.Plan.SetAttribute(ctx, path.Root("vdc"), r.client.GetDefaultVDC())
		}
	}
}

// Metadata returns the resource type name.
func (r *networkIsolatedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_isolated"
}

// Init resource used to initialize the resource.
func (r *networkIsolatedResource) Init(_ context.Context, rm *networkIsolatedResourceModel) (diags diag.Diagnostics) {
	// Init Org
	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}
	// Init Vdc
	r.vdc, diags = vdc.Init(r.client, rm.VDC)
	return
}

// Schema defines the schema for the resource.
func (r *networkIsolatedResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetIsolated()).GetResource()
}

func (r *networkIsolatedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *networkIsolatedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan := &networkIsolatedResourceModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Define VDC or VDCGroup
	vdcOrVDCGroup, err := r.client.GetVDCOrVDCGroup(plan.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or VDCGroup", err.Error())
		return
	}

	// Lock VDC or VDCGroup
	vcdMutexKV := mutex.NewKV()
	vcdMutexKV.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer vcdMutexKV.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	// Set network type
	networkType, diag := r.SetVCDNetwork(ctx, vdcOrVDCGroup.GetID(), *plan)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create network
	orgNetwork, err := r.org.CreateOpenApiOrgVdcNetwork(networkType)
	if err != nil {
		resp.Diagnostics.AddError("[CREATE] Error creating isolated network", err.Error())
		return
	}

	// Set Plan only for compute values
	plan.ID = types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID)
	plan.VDC = types.StringValue(vdcOrVDCGroup.GetName())

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *networkIsolatedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	state := &networkIsolatedResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Define VDC or VDCGroup
	vdcOrVDCGroup, err := r.client.GetVDCOrVDCGroup(state.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or VDCGroup", err.Error())
		return
	}

	// Get network
	orgNetwork, err := vdcOrVDCGroup.GetOpenApiOrgVdcNetworkByName(state.Name.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// Network not found, so remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[READ] Error retrieving isolated network", err.Error())
		return
	}

	// Get network static IP pools
	ipPools := []staticIPPool{}
	if len(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		for _, ipRange := range orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipPools = append(ipPools, staticIPPool{
				StartAddress: types.StringValue(ipRange.StartAddress),
				EndAddress:   types.StringValue(ipRange.EndAddress),
			})
		}
	}

	// Set Plan updated
	plan := &networkIsolatedResourceModel{
		ID:           types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID),
		Name:         types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Name),
		Description:  types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Description),
		VDC:          types.StringValue(vdcOrVDCGroup.GetName()),
		Gateway:      types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].Gateway),
		PrefixLength: types.Int64Value(int64(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].PrefixLength)),
		DNS1:         types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1),
		DNS2:         types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2),
		DNSSuffix:    types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix),
	}

	// Set static IP pools
	var diags diag.Diagnostics
	plan.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolAttrTypes}, ipPools)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkIsolatedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	plan := &networkIsolatedResourceModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	state := &networkIsolatedResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Define VDC or VDCGroup
	vdcOrVDCGroup, err := r.client.GetVDCOrVDCGroup(state.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or VDCGroup", err.Error())
		return
	}

	// Lock VDC or VDCGroup
	vcdMutexKV := mutex.NewKV()
	vcdMutexKV.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer vcdMutexKV.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	// Get network
	orgNetwork, err := vdcOrVDCGroup.GetOpenApiOrgVdcNetworkByName(state.Name.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// Network not found, so remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[READ] Error retrieving isolated network", err.Error())
		return
	}

	// Set network type
	ipPool := []staticIPPool{}
	resp.Diagnostics.Append(plan.StaticIPPool.ElementsAs(ctx, &ipPool, true)...)
	if resp.Diagnostics.HasError() {
		return
	}
	myshared := false // Cloudavenue does not support shared networks
	networkType := &govdctypes.OpenApiOrgVdcNetwork{
		ID:          plan.ID.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Shared:      &myshared,
		NetworkType: govdctypes.OrgVdcNetworkTypeIsolated,
		OwnerRef:    &govdctypes.OpenApiReference{ID: vdcOrVDCGroup.GetID()},
		Subnets: govdctypes.OrgVdcNetworkSubnets{
			Values: []govdctypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      plan.Gateway.ValueString(),
					PrefixLength: int(plan.PrefixLength.ValueInt64()),
					IPRanges: govdctypes.OrgVdcNetworkSubnetIPRanges{
						Values: processIPRanges(ipPool),
					},
					DNSServer1: plan.DNS1.ValueString(),
					DNSServer2: plan.DNS2.ValueString(),
					DNSSuffix:  plan.DNSSuffix.ValueString(),
				},
			},
		},
	}

	// Update network
	_, err = orgNetwork.Update(networkType)
	if err != nil {
		resp.Diagnostics.AddError("Error updating isolated network", err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkIsolatedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	state := &networkIsolatedResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Define VDC or VDCGroup
	vdcOrVDCGroup, err := r.client.GetVDCOrVDCGroup(state.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or VDCGroup", err.Error())
		return
	}

	// Get network
	orgNetwork, err := vdcOrVDCGroup.GetOpenApiOrgVdcNetworkByName(state.Name.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// Network not found, so remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[READ] Error retrieving isolated network", err.Error())
		return
	}

	// Lock VDC or VDCGroup
	vcdMutexKV := mutex.NewKV()
	vcdMutexKV.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer vcdMutexKV.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	// Delete network
	err = orgNetwork.Delete()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting isolated network", err.Error())
		return
	}
}

func (r *networkIsolatedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Get URI from import ID
	resourceURI := strings.Split(req.ID, ".")
	if len(resourceURI) != 2 {
		resp.Diagnostics.AddError("Error importing network_routed", "Resource name must be specified as vdc-name.network-name or vdc-group-name.network-name")
		return
	}
	vdcOrVDCGroupName, networkName := resourceURI[0], resourceURI[1]

	// Get VDC or VDCGroup
	vdcOrVDCGroup, err := r.client.GetVDCOrVDCGroup(vdcOrVDCGroupName)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or VDCGroup", err.Error())
		return
	}

	// Get network
	orgNetwork, err := vdcOrVDCGroup.GetOpenApiOrgVdcNetworkByName(networkName)
	if err != nil { // If network is not found, return error
		resp.Diagnostics.AddError("Error retrieving org network by name", err.Error())
		return
	}

	// Check if network is isolated
	if !orgNetwork.IsIsolated() {
		resp.Diagnostics.AddError("Error importing network_isolated", "Network is not isolated")
		return
	}

	// Get network static IP pools
	ipPools := []staticIPPool{}
	if len(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		for _, ipRange := range orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipPools = append(ipPools, staticIPPool{
				StartAddress: types.StringValue(ipRange.StartAddress),
				EndAddress:   types.StringValue(ipRange.EndAddress),
			})
		}
	}

	// Set state to fully populated data
	// Set Plan updated
	plan := &networkIsolatedResourceModel{
		ID:           types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID),
		Name:         types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Name),
		Description:  types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Description),
		VDC:          types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.OwnerRef.Name),
		Gateway:      types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].Gateway),
		PrefixLength: types.Int64Value(int64(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].PrefixLength)),
		DNS1:         types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1),
		DNS2:         types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2),
		DNSSuffix:    types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix),
	}
	// Set static IP pools
	var diags diag.Diagnostics
	plan.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolAttrTypes}, ipPools)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
