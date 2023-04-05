// Package network provides a Terraform resource.
package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkRoutedResource{}
	_ resource.ResourceWithConfigure   = &networkRoutedResource{}
	_ resource.ResourceWithImportState = &networkRoutedResource{}
	_ network.Network                  = &networkRoutedResource{}
)

// NewNetworkRoutedResource is a helper function to simplify the provider implementation.
func NewNetworkRoutedResource() resource.Resource {
	return &networkRoutedResource{}
}

// networkRoutedResource is the resource implementation.
type networkRoutedResource struct {
	client  *client.CloudAvenue
	org     org.Org
	network network.Kind
}

// Metadata returns the resource type name.
func (r *networkRoutedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "routed"
}

// Schema defines the schema for the resource.
func (r *networkRoutedResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetRouted()).GetResource(ctx)
}

// Init resource used to initialize the resource.
func (r *networkRoutedResource) Init(_ context.Context, rm *networkRoutedModel) (diags diag.Diagnostics) {
	// Set Network Type
	r.network.TypeOfNetwork = network.NAT_ROUTED
	// Init Org
	r.org, diags = org.Init(r.client)
	return
}

func (r *networkRoutedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *networkRoutedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan := &networkRoutedModel{}

	// Get Plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	edgegw, err := r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   plan.EdgeGatewayID,
		Name: plan.EdgeGatewayName,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	vdcOrVDCGroup, err := edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	networkMutexKV.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer networkMutexKV.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	// Set Network
	orgVDCNetworkConfig, diag := r.SetNetworkAPIObject(ctx, plan)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	// Create Network
	orgNetwork, err := r.org.CreateOpenApiOrgVdcNetwork(orgVDCNetworkConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error creating routing network", err.Error())
		return
	}

	// Set ID
	plan.ID = types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID)
	plan.EdgeGatewayID = types.StringValue(edgegw.EdgeGateway.ID)
	plan.EdgeGatewayName = types.StringValue(edgegw.EdgeGateway.Name)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *networkRoutedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &networkRoutedModel{}
	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Parent Edge Gateway ID to define the owner (VDC or VDC Group)
	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(state.ID.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Debug(ctx, "Network not found, removing resource from state")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving routing network", err.Error())
		return
	}

	// Set data into the network model
	plan := SetDataToNetworkRoutedModel(orgNetwork)

	// Set Static IP Pool
	ipPools := []staticIPPool{}
	if len(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		for _, ipRange := range orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipPool := staticIPPool{
				StartAddress: types.StringValue(ipRange.StartAddress),
				EndAddress:   types.StringValue(ipRange.EndAddress),
			}
			ipPools = append(ipPools, ipPool)
		}
	}
	var diags diag.Diagnostics
	plan.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolAttrTypes}, ipPools)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkRoutedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &networkRoutedModel{}

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Parent Edge Gateway ID to define the owner (VDC or VDC Group)
	edgegw, err := r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   plan.EdgeGatewayID,
		Name: plan.EdgeGatewayName,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	vdcOrVDCGroup, err := edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	networkMutexKV.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer networkMutexKV.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	// Get current network
	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(plan.ID.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Debug(ctx, "Network not found, removing resource from state")
			resp.State.RemoveResource(ctx)
		}
		resp.Diagnostics.AddError("Error retrieving routing network", err.Error())
		return
	}

	// Set Network
	newOrgNetwork, diag := r.SetNetworkAPIObject(ctx, plan)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	// Update network
	_, err = orgNetwork.Update(newOrgNetwork)
	if err != nil {
		resp.Diagnostics.AddError("Error updating routing network", err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkRoutedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &networkRoutedModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Parent Edge Gateway ID to define the owner (VDC or VDC Group)
	edgegw, err := r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   state.EdgeGatewayID,
		Name: state.EdgeGatewayName,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	vdcOrVDCGroup, err := edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	networkMutexKV.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer networkMutexKV.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	// Get current network
	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(state.ID.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Debug(ctx, "Network not found, removing resource from state")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving routing network", err.Error())
		return
	}

	err = orgNetwork.Delete()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting routing network", err.Error())
	}
}

func (r *networkRoutedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceURI := strings.Split(req.ID, ".")

	if len(resourceURI) != 2 {
		resp.Diagnostics.AddError("Error importing network_routed", "Resource name must be specified as vdc-name.network-name or vdc-group-name.network-name")
		return
	}

	vdcOrVDCGroupName, networkName := resourceURI[0], resourceURI[1]

	v, err := r.client.GetVDCOrVDCGroup(vdcOrVDCGroupName)
	if err != nil && govcd.ContainsNotFound(err) {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	orgNetwork, err := v.GetOpenApiOrgVdcNetworkByName(networkName)
	if err != nil && !govcd.ContainsNotFound(err) || orgNetwork == nil {
		resp.Diagnostics.AddError("Error retrieving org vdc network by name", err.Error())
		return
	}

	if !orgNetwork.IsRouted() {
		resp.Diagnostics.AddError("Error importing routed network", fmt.Sprintf("Org network with name '%s' found, but is not of type Routed (type is '%s')", networkName, orgNetwork.GetType()))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), orgNetwork.OpenApiOrgVdcNetwork.ID)...)
}

func (r *networkRoutedResource) SetNetworkAPIObject(ctx context.Context, plan any) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	d := diag.Diagnostics{}

	p, ok := plan.(*networkRoutedModel)
	if !ok {
		d.AddError("Error", "Error converting plan to network routed resource model")
		return nil, d
	}

	// Get Parent Edge Gateway ID to define the owner (VDC or VDC Group)
	edgegw, err := r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   p.EdgeGatewayID,
		Name: p.EdgeGatewayName,
	})
	if err != nil {
		d.AddError("Error retrieving Edge Gateway", err.Error())
		return nil, d
	}

	vdcOrVDCGroup, err := edgegw.GetParent()
	if err != nil {
		d.AddError("Error retrieving Edge Gateway parent", err.Error())
		return nil, d
	}

	// Set global resource model
	return r.network.SetNetworkAPIObject(ctx, network.GlobalResourceModel{
		ID:                p.ID,
		Name:              p.Name,
		Description:       p.Description,
		Gateway:           p.Gateway,
		PrefixLength:      p.PrefixLength,
		DNS1:              p.DNS1,
		DNS2:              p.DNS2,
		DNSSuffix:         p.DNSSuffix,
		StaticIPPool:      p.StaticIPPool,
		VDCIDOrVDCGroupID: types.StringValue(vdcOrVDCGroup.GetID()),
		EdgeGatewayID:     p.EdgeGatewayID,
		EdgegatewayName:   p.EdgeGatewayName,
		InterfaceType:     p.InterfaceType,
	})
}
