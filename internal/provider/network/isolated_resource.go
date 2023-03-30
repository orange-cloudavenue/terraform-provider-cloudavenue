// Package network provides a Terraform resource.
package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govdctypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"

	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkIsolatedResource{}
	_ resource.ResourceWithConfigure   = &networkIsolatedResource{}
	_ resource.ResourceWithImportState = &networkIsolatedResource{}
	_ resource.ResourceWithModifyPlan  = &networkIsolatedResource{}
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
	PrimaryDNS   types.String `tfsdk:"dns1"`
	SecondaryDNS types.String `tfsdk:"dns2"`
	SuffixDNS    types.String `tfsdk:"dns_suffix"`
	StaticIPPool types.Set    `tfsdk:"static_ip_pool"`
}

type staticIPPoolResourceModel struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

var staticIPPoolResourceModelAttrTypes = map[string]attr.Type{
	"start_address": types.StringType,
	"end_address":   types.StringType,
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
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a VMware Cloud Director Org VDC isolated Network. This can be used to create, modify, and delete isolated VDC networks",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the network. This is a generated value and cannot be specified during creation. This value is used to identify the network in other resources.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vdc": vdc.Schema(),
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the network. This value must be unique within the `VDC` or `VDC Group` that owns the network.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the network.",
			},
			"gateway": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "(Force replacement) The gateway IP address for the network. This value define also the network IP range with the prefix length.",
				Validators: []validator.String{
					fstringvalidator.IsValidIP(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"prefix_length": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "(Force replacement) The prefix length for the network. This value must be a valid prefix length for the network IP range.(e.g. 24 for netmask 255.255.255.0)",
				Validators: []validator.Int64{
					int64validator.Between(1, 32),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"dns1": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The primary DNS server IP address for the network.",
				Validators: []validator.String{
					fstringvalidator.IsValidIP(),
				},
			},
			"dns2": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The secondary DNS server IP address for the network.",
				Validators: []validator.String{
					fstringvalidator.IsValidIP(),
				},
			},
			"dns_suffix": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The DNS suffix for the network.",
			},
			"static_ip_pool": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "A set of static IP pools to be used for this network.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_address": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The start address of the IP pool. This value must be a valid IP address in the network IP range.",
							Validators: []validator.String{
								fstringvalidator.IsValidIP(),
							},
						},
						"end_address": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The end address of the IP pool. This value must be a valid IP address in the network IP range.",
							Validators: []validator.String{
								fstringvalidator.IsValidIP(),
							},
						},
					},
				},
			},
		},
	}
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
	ipPool := []staticIPPoolResourceModel{}
	resp.Diagnostics.Append(plan.StaticIPPool.ElementsAs(ctx, &ipPool, true)...)
	if resp.Diagnostics.HasError() {
		return
	}
	myshared := false // Cloudavenue does not support shared networks
	networkType := &govdctypes.OpenApiOrgVdcNetwork{
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
						Values: myProcessIPRanges(ipPool),
					},
					DNSServer1: plan.PrimaryDNS.ValueString(),
					DNSServer2: plan.SecondaryDNS.ValueString(),
					DNSSuffix:  plan.SuffixDNS.ValueString(),
				},
			},
		},
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
	ipPools := []staticIPPoolResourceModel{}
	if len(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		for _, ipRange := range orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipPools = append(ipPools, staticIPPoolResourceModel{
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
		PrimaryDNS:   types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1),
		SecondaryDNS: types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2),
		SuffixDNS:    types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix),
	}

	// Set static IP pools
	var diags diag.Diagnostics
	plan.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolResourceModelAttrTypes}, ipPools)
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
	ipPool := []staticIPPoolResourceModel{}
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
						Values: myProcessIPRanges(ipPool),
					},
					DNSServer1: plan.PrimaryDNS.ValueString(),
					DNSServer2: plan.SecondaryDNS.ValueString(),
					DNSSuffix:  plan.SuffixDNS.ValueString(),
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
	ipPools := []staticIPPoolResourceModel{}
	if len(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		for _, ipRange := range orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipPools = append(ipPools, staticIPPoolResourceModel{
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
		PrimaryDNS:   types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1),
		SecondaryDNS: types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2),
		SuffixDNS:    types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix),
	}
	// Set static IP pools
	var diags diag.Diagnostics
	plan.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolResourceModelAttrTypes}, ipPools)
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

// StaticIPPool is a helper function to get the static IP pool from the resource data.
func myProcessIPRanges(mystaticIPPool []staticIPPoolResourceModel) []govdctypes.ExternalNetworkV2IPRange {
	subnetRng := make([]govdctypes.ExternalNetworkV2IPRange, len(mystaticIPPool))
	for i, ipRange := range mystaticIPPool {
		subnetRng[i].StartAddress = ipRange.StartAddress.ValueString()
		subnetRng[i].EndAddress = ipRange.EndAddress.ValueString()
	}
	return subnetRng
}
