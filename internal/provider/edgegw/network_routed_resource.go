/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	"context"
	"fmt"
	"strings"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/network"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &NetworkRoutedResource{}
	_ resource.ResourceWithConfigure   = &NetworkRoutedResource{}
	_ resource.ResourceWithImportState = &NetworkRoutedResource{}
)

// NewNetworkRoutedResource is a helper function to simplify the provider implementation.
func NewNetworkRoutedResource() resource.Resource {
	return &NetworkRoutedResource{}
}

// NetworkRoutedResource is the resource implementation.
type NetworkRoutedResource struct {
	client *client.CloudAvenue
	edgegw *v1.EdgeClient
	vdc    v1.VDCOrVDCGroupInterface
}

// Init Initializes the resource.
func (r *NetworkRoutedResource) Init(_ context.Context, rm *NetworkRoutedModel) (diags diag.Diagnostics) {
	var err error

	idOrName := rm.EdgeGatewayID.Get()
	if idOrName == "" {
		idOrName = rm.EdgeGatewayName.Get()
	}

	r.edgegw, err = r.client.CAVSDK.V1.EdgeGateway.Get(idOrName)
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	rm.EdgeGatewayID.Set(r.edgegw.GetURN())
	rm.EdgeGatewayName.Set(r.edgegw.GetName())

	r.vdc, err = r.client.CAVSDK.V1.VDC().GetVDCOrVDCGroup(r.edgegw.OwnerName)
	if err != nil {
		diags.AddError("Error retrieving VDC from EdgeGateway", fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	if r.vdc.IsVDCGroup() {
		diags.AddError("EdgeGateway is connected to a VDCGroup", "EdgeGateway is connected to a VDCGroup, please use the VDCGroup resource instead. `cloudavenue_vdcg_network_routed`")
	}

	return
}

// Metadata returns the resource type name.
func (r *NetworkRoutedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_network_routed"
}

// Schema defines the schema for the resource.
func (r *NetworkRoutedResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = networkRoutedSchema(ctx).GetResource(ctx)
}

func (r *NetworkRoutedResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ResourceWithMoveState interface implementation.
func (r *NetworkRoutedResource) MoveState(_ context.Context) []resource.StateMover {
	return []resource.StateMover{
		{
			SourceSchema: func() *schema.Schema {
				ctx := context.Background()
				schemaRequest := resource.SchemaRequest{}
				schemaResponse := &resource.SchemaResponse{}

				network.NewNetworkRoutedResource().Schema(ctx, schemaRequest, schemaResponse)
				if schemaResponse.Diagnostics.HasError() {
					return nil
				}

				return &schemaResponse.Schema
			}(),
			StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
				if req.SourceTypeName != "cloudavenue_network_routed" {
					return
				}

				if req.SourceSchemaVersion != 0 {
					return
				}

				if !strings.HasSuffix(req.SourceProviderAddress, "orange-cloudavenue/cloudavenue") {
					return
				}

				source := &network.RoutedModel{}
				resp.Diagnostics.Append(req.SourceState.Get(ctx, source)...)
				if resp.Diagnostics.HasError() {
					return
				}

				dest := &NetworkRoutedModel{
					ID:              source.ID,
					Name:            source.Name,
					EdgeGatewayID:   source.EdgeGatewayID,
					EdgeGatewayName: source.EdgeGatewayName,
					Description:     source.Description,
					Gateway:         source.Gateway,
					PrefixLength:    source.PrefixLength,
					DNS1:            source.DNS1,
					DNS2:            source.DNS2,
					DNSSuffix:       source.DNSSuffix,
					StaticIPPool:    supertypes.NewSetNestedObjectValueOfNull[NetworkRoutedModelStaticIPPool](ctx),

					// GuestVLANAllowed is a new field
					GuestVLANAllowed: supertypes.NewBoolNull(),
				}

				sIPPools, d := source.StaticIPPool.Get(ctx)
				if d.HasError() {
					resp.Diagnostics.Append(d...)
					return
				}

				dIPPools := []*NetworkRoutedModelStaticIPPool{}
				for _, ipPool := range sIPPools {
					dIPPools = append(dIPPools, &NetworkRoutedModelStaticIPPool{
						StartAddress: ipPool.StartAddress,
						EndAddress:   ipPool.EndAddress,
					})
				}

				resp.Diagnostics.Append(dest.StaticIPPool.Set(ctx, dIPPools)...)
				if resp.Diagnostics.HasError() {
					return
				}

				resp.Diagnostics.Append(resp.TargetState.Set(ctx, &dest)...)
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *NetworkRoutedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_routed", r.client.GetOrgName(), metrics.Create)()

	plan := &NetworkRoutedModel{}

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource creation logic here.
	*/

	sdkValues, diags := plan.ToSDKNetworkRoutedModel(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())

	netRouted, err := r.vdc.CreateNetworkRouted(sdkValues)
	if err != nil {
		resp.Diagnostics.AddError("Error creating network routed", fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	// Set the ID
	plan.ID.Set(netRouted.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", "The resource was not found after creation")
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *NetworkRoutedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_routed", r.client.GetOrgName(), metrics.Read)()

	state := &NetworkRoutedModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.Diagnostics.AddError("Resource not found", "The resource was not found after refresh")
		resp.State.RemoveResource(ctx)
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *NetworkRoutedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_routed", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &NetworkRoutedModel{}
		state = &NetworkRoutedModel{}
	)

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource update here
	*/

	sdkValues, diags := plan.ToSDKNetworkRoutedModel(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())

	netRouted, err := r.vdc.GetNetworkRouted(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving network routed", fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	if err = netRouted.Update(sdkValues); err != nil {
		resp.Diagnostics.AddError("Error updating network routed", fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", "The resource was not found after update")
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *NetworkRoutedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_routed", r.client.GetOrgName(), metrics.Delete)()

	state := &NetworkRoutedModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource deletion here
	*/

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())

	netRouted, err := r.vdc.GetNetworkRouted(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving network routed", fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	if err = netRouted.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting network routed", fmt.Sprintf("Error: %s", err.Error()))
	}
}

func (r *NetworkRoutedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_routed", r.client.GetOrgName(), metrics.Import)()

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: edgeGatewayNameOrID.networkNameOrID Got: %q", req.ID),
		)
		return
	}

	x := &NetworkRoutedModel{
		ID:              supertypes.NewStringNull(),
		Name:            supertypes.NewStringNull(),
		EdgeGatewayName: supertypes.NewStringNull(),
		EdgeGatewayID:   supertypes.NewStringNull(),
	}

	if urn.IsEdgeGateway(idParts[0]) {
		x.EdgeGatewayID.Set(idParts[0])
	} else {
		x.EdgeGatewayName.Set(idParts[0])
	}

	if urn.IsNetwork(idParts[1]) {
		x.ID.Set(idParts[1])
	} else {
		x.Name.Set(idParts[1])
	}

	resp.Diagnostics.Append(r.Init(ctx, x)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, d := r.read(ctx, x)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *NetworkRoutedResource) read(ctx context.Context, planOrState *NetworkRoutedModel) (stateRefreshed *NetworkRoutedModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	idOrName := stateRefreshed.ID.Get()
	if idOrName == "" {
		idOrName = stateRefreshed.Name.Get()
	}

	net, err := r.vdc.GetNetworkRouted(idOrName)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error getting routed network", err.Error())
		return nil, true, diags
	}

	// Populate the state with the network data
	stateRefreshed.ID.Set(net.ID)
	stateRefreshed.Name.Set(net.Name)
	stateRefreshed.Description.Set(net.Description)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetURN())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())
	stateRefreshed.Gateway.Set(net.Subnet.Gateway)
	stateRefreshed.PrefixLength.SetInt(net.Subnet.PrefixLength)
	stateRefreshed.DNS1.Set(net.Subnet.DNSServer1)
	stateRefreshed.DNS2.Set(net.Subnet.DNSServer2)
	stateRefreshed.DNSSuffix.Set(net.Subnet.DNSSuffix)
	stateRefreshed.GuestVLANAllowed.SetPtr(net.GuestVLANTaggingAllowed)

	if len(net.Subnet.IPRanges) == 0 {
		stateRefreshed.StaticIPPool.SetNull(ctx)
	} else {
		x := []*NetworkRoutedModelStaticIPPool{}
		for _, ipRange := range net.Subnet.IPRanges {
			x = append(x, &NetworkRoutedModelStaticIPPool{
				StartAddress: supertypes.NewStringValueOrNull(ipRange.StartAddress),
				EndAddress:   supertypes.NewStringValueOrNull(ipRange.EndAddress),
			})
		}

		diags.Append(stateRefreshed.StaticIPPool.Set(ctx, x)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	return stateRefreshed, true, nil
}
