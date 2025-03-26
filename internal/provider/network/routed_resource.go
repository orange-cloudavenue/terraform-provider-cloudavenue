/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	sdkv1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkRoutedResource{}
	_ resource.ResourceWithConfigure   = &networkRoutedResource{}
	_ resource.ResourceWithImportState = &networkRoutedResource{}
)

// NewNetworkRoutedResource is a helper function to simplify the provider implementation.
func NewNetworkRoutedResource() resource.Resource {
	return &networkRoutedResource{}
}

// networkRoutedResource is the resource implementation.
type networkRoutedResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Metadata returns the resource type name.
func (r *networkRoutedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "routed"
}

// Schema defines the schema for the resource.
func (r *networkRoutedResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = routedSchema(ctx).GetResource(ctx)
}

// Init resource used to initialize the resource.
func (r *networkRoutedResource) Init(_ context.Context, rm *RoutedModel) (diags diag.Diagnostics) {
	// Init Org
	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	var err error
	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   rm.EdgeGatewayID.StringValue,
		Name: rm.EdgeGatewayName.StringValue,
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

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
	defer metrics.New("cloudavenue_network_routed", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	plan := &RoutedModel{}

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

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}

	// Set Network
	orgVDCNetworkConfig, diag := r.setNetworkAPIObject(ctx, plan, vdcOrVDCGroup)
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
	plan.ID.Set(orgNetwork.OpenApiOrgVdcNetwork.ID)
	plan.EdgeGatewayID.Set(r.edgegw.GetID())
	plan.EdgeGatewayName.Set(r.edgegw.GetName())

	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Error creating routing network", "Network not found after creation")
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *networkRoutedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_network_routed", r.client.GetOrgName(), metrics.Read)()

	state := &RoutedModel{}
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

	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.Diagnostics.AddError("Error reading routing network", "Network not found")
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkRoutedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_network_routed", r.client.GetOrgName(), metrics.Update)()

	plan := &RoutedModel{}

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

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}

	// Get current network
	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(plan.ID.Get())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving routing network", err.Error())
		return
	}

	// Set Network
	orgVDCNetworkConfig, diag := r.setNetworkAPIObject(ctx, plan, vdcOrVDCGroup)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	// Update network
	_, err = orgNetwork.Update(orgVDCNetworkConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error updating routing network", err.Error())
		return
	}

	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Error reading routing network", "Network not found")
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkRoutedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_network_routed", r.client.GetOrgName(), metrics.Delete)()

	state := &RoutedModel{}

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

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}

	// Get current network
	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(state.ID.Get())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving routing network", err.Error())
		return
	}

	if err := orgNetwork.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting routing network", err.Error())
	}
}

func (r *networkRoutedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_network_routed", r.client.GetOrgName(), metrics.Import)()

	resourceURI := strings.Split(req.ID, ".")

	if len(resourceURI) != 2 {
		resp.Diagnostics.AddError("Error importing network_routed", "Resource name must be specified as vdc-name.network-name or vdc-group-name.network-name")
		return
	}

	edgeGatewayNameOrEdgeGatewayID, networkName := resourceURI[0], resourceURI[1]

	// Get Edge Gateway
	var (
		edgeGWName string
		edgeGWID   string
	)
	if urn.IsEdgeGateway(edgeGatewayNameOrEdgeGatewayID) {
		edgeGWID = edgeGatewayNameOrEdgeGatewayID
	} else {
		edgeGWName = edgeGatewayNameOrEdgeGatewayID
	}

	newPlan := &RoutedModel{}
	newPlan.EdgeGatewayID.Set(edgeGWID)
	newPlan.EdgeGatewayName.Set(edgeGWName)

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, newPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get vdc or vdc group
	v, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_id"), r.edgegw.GetID())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), r.edgegw.GetName())...)
}

func (r *networkRoutedResource) read(ctx context.Context, planOrState *RoutedModel) (stateRefreshed *RoutedModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	// Get Parent Edge Gateway ID to define the owner (VDC or VDC Group)
	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(stateRefreshed.ID.Get())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving routing network", err.Error())
		return
	}

	stateRefreshed.ID.Set(orgNetwork.OpenApiOrgVdcNetwork.ID)
	stateRefreshed.Name.Set(orgNetwork.OpenApiOrgVdcNetwork.Name)
	stateRefreshed.Description.Set(orgNetwork.OpenApiOrgVdcNetwork.Description)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())
	stateRefreshed.InterfaceType.Set(orgNetwork.OpenApiOrgVdcNetwork.Connection.ConnectionType)
	if len(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values) == 0 {
		diags.AddError("Error retrieving subnet", "No subnet found")
		return
	}
	stateRefreshed.Gateway.Set(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].Gateway)
	stateRefreshed.PrefixLength.Set(int64(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].PrefixLength))
	stateRefreshed.DNS1.Set(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1)
	stateRefreshed.DNS2.Set(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2)
	stateRefreshed.DNSSuffix.Set(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix)

	ipPools := []*RoutedModelStaticIPPool{}
	for _, ipRange := range orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
		ipPool := &RoutedModelStaticIPPool{}
		ipPool.StartAddress.Set(ipRange.StartAddress)
		ipPool.EndAddress.Set(ipRange.EndAddress)
		ipPools = append(ipPools, ipPool)
	}
	diags.Append(stateRefreshed.StaticIPPool.Set(ctx, ipPools)...)

	return stateRefreshed, true, diags
}

func (r *networkRoutedResource) setNetworkAPIObject(ctx context.Context, plan *RoutedModel, vdcOrVDCGroup sdkv1.VDCOrVDCGroupInterface) (orgVDCNetwork *govcdtypes.OpenApiOrgVdcNetwork, diags diag.Diagnostics) {
	ipPools, d := plan.StaticIPPool.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	ipRange := make([]govcdtypes.ExternalNetworkV2IPRange, len(ipPools))
	for rangeIndex, subnetRange := range ipPools {
		ipRange[rangeIndex] = govcdtypes.ExternalNetworkV2IPRange{
			StartAddress: subnetRange.StartAddress.Get(),
			EndAddress:   subnetRange.EndAddress.Get(),
		}
	}

	return &govcdtypes.OpenApiOrgVdcNetwork{
		ID:          plan.ID.Get(),
		Name:        plan.Name.Get(),
		Description: plan.Description.Get(),
		OwnerRef:    &govcdtypes.OpenApiReference{ID: vdcOrVDCGroup.GetID()},
		NetworkType: govcdtypes.OrgVdcNetworkTypeRouted,
		Connection: &govcdtypes.Connection{
			RouterRef: govcdtypes.OpenApiReference{
				ID:   plan.EdgeGatewayID.Get(),
				Name: plan.EdgeGatewayName.Get(),
			},
			ConnectionType: plan.InterfaceType.Get(),
		},
		Subnets: govcdtypes.OrgVdcNetworkSubnets{
			Values: []govcdtypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      plan.Gateway.Get(),
					PrefixLength: plan.PrefixLength.GetInt(),
					IPRanges: govcdtypes.OrgVdcNetworkSubnetIPRanges{
						Values: ipRange,
					},
					DNSServer1: plan.DNS1.Get(),
					DNSServer2: plan.DNS2.Get(),
					DNSSuffix:  plan.DNSSuffix.Get(),
				},
			},
		},
	}, nil
}
