/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

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
	cnetwork "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/network"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &NetworkIsolatedResource{}
	_ resource.ResourceWithConfigure   = &NetworkIsolatedResource{}
	_ resource.ResourceWithImportState = &NetworkIsolatedResource{}
	_ resource.ResourceWithMoveState   = &NetworkIsolatedResource{}
)

// NewNetworkIsolatedResource is a helper function to simplify the provider implementation.
func NewNetworkIsolatedResource() resource.Resource {
	return &NetworkIsolatedResource{}
}

// NetworkIsolatedResource is the resource implementation.
type NetworkIsolatedResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
}

// Init Initializes the resource.
func (r *NetworkIsolatedResource) Init(ctx context.Context, rm *networkIsolatedModel) (diags diag.Diagnostics) {
	r.vdc, diags = vdc.Init(r.client, rm.VDC.StringValue)
	return
}

// Metadata returns the resource type name.
func (r *NetworkIsolatedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_network_isolated"
}

// Schema defines the schema for the resource.
func (r *NetworkIsolatedResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = networkIsolatedSchema(ctx).GetResource(ctx)
}

func (r *NetworkIsolatedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NetworkIsolatedResource) MoveState(ctx context.Context) []resource.StateMover {
	return []resource.StateMover{
		{
			SourceSchema: func() *schema.Schema {
				ctx := context.Background()
				schemaRequest := resource.SchemaRequest{}
				schemaResponse := &resource.SchemaResponse{}

				network.NewNetworkIsolatedResource().Schema(ctx, schemaRequest, schemaResponse)
				if schemaResponse.Diagnostics.HasError() {
					return nil
				}

				return &schemaResponse.Schema
			}(),
			StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
				if req.SourceTypeName != "cloudavenue_network_isolated" {
					return
				}

				if req.SourceSchemaVersion != 0 {
					return
				}

				if !strings.HasSuffix(req.SourceProviderAddress, "orange-cloudavenue/cloudavenue") {
					return
				}

				source := network.IsolatedModel{}
				dest := &networkIsolatedModel{
					ID:               supertypes.NewStringNull(),
					Name:             supertypes.NewStringNull(),
					VDC:              supertypes.NewStringNull(),
					Description:      supertypes.NewStringNull(),
					Gateway:          supertypes.NewStringNull(),
					PrefixLength:     supertypes.NewInt64Null(),
					DNS1:             supertypes.NewStringNull(),
					DNS2:             supertypes.NewStringNull(),
					DNSSuffix:        supertypes.NewStringNull(),
					StaticIPPool:     supertypes.NewSetNestedObjectValueOfNull[networkIsolatedModelStaticIPPool](ctx),
					GuestVLANAllowed: supertypes.NewBoolValue(false),
				}

				resp.Diagnostics.Append(req.SourceState.Get(ctx, &source)...)
				if resp.Diagnostics.HasError() {
					return
				}

				dest.ID.Set(source.ID.ValueString())
				dest.VDC.Set(source.VDC.ValueString())
				dest.Description.Set(source.Description.ValueString())
				dest.Gateway.Set(source.Gateway.ValueString())
				dest.PrefixLength.SetInt64(source.PrefixLength.ValueInt64())
				dest.DNS1.Set(source.DNS1.ValueString())
				dest.DNS2.Set(source.DNS2.ValueString())
				dest.DNSSuffix.Set(source.DNSSuffix.ValueString())
				dIPPools := []*networkIsolatedModelStaticIPPool{}
				sIPPools := []cnetwork.StaticIPPool{}

				resp.Diagnostics.Append(source.StaticIPPool.ElementsAs(ctx, &sIPPools, true)...)
				if resp.Diagnostics.HasError() {
					return
				}

				for _, ipPool := range sIPPools {
					dIPPools = append(dIPPools, &networkIsolatedModelStaticIPPool{
						StartAddress: supertypes.NewStringValue(ipPool.StartAddress.ValueString()),
						EndAddress:   supertypes.NewStringValue(ipPool.EndAddress.ValueString()),
					})
				}

				resp.Diagnostics.Append(dest.StaticIPPool.Set(ctx, dIPPools)...)
				resp.Diagnostics.Append(resp.TargetState.Set(ctx, &dest)...)
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *NetworkIsolatedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdc_network_isolated", r.client.GetOrgName(), metrics.Create)()

	plan := &networkIsolatedModel{}

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

	values, d := plan.ToSDK(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	networkIsolated, err := r.vdc.CreateNetworkIsolated(values)
	if err != nil {
		resp.Diagnostics.AddError("Error creating isolated network", err.Error())
		return
	}

	plan.ID.Set(networkIsolated.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("Resource not found after creation", "The resource was not found after creation.")
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
func (r *NetworkIsolatedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdc_network_isolated", r.client.GetOrgName(), metrics.Read)()

	state := &networkIsolatedModel{}

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
func (r *NetworkIsolatedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdc_network_isolated", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &networkIsolatedModel{}
		state = &networkIsolatedModel{}
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

	values, d := plan.ToSDK(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	net, err := r.vdc.GetNetworkIsolated(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error getting isolated network", err.Error())
		return
	}

	values.ID = state.ID.Get()

	// Update the network
	if err := net.Update(values); err != nil {
		resp.Diagnostics.AddError("Error updating isolated network", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("Resource not found after update", "The resource was not found after update. Please refresh the state.")
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
func (r *NetworkIsolatedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdc_network_isolated", r.client.GetOrgName(), metrics.Delete)()

	state := &networkIsolatedModel{}

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

	net, err := r.vdc.GetNetworkIsolated(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error getting isolated network", err.Error())
		return
	}

	// Delete the network
	if err := net.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting isolated network", err.Error())
		return
	}
}

func (r *NetworkIsolatedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdc_network_isolated", r.client.GetOrgName(), metrics.Import)()

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: vdc.networkNameOrID Got: %q", req.ID),
		)
		return
	}

	x := &networkIsolatedModel{
		ID:   supertypes.NewStringNull(),
		Name: supertypes.NewStringNull(),
		VDC:  supertypes.NewStringNull(),
	}

	x.VDC.Set(idParts[0])

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
func (r *NetworkIsolatedResource) read(ctx context.Context, planOrState *networkIsolatedModel) (stateRefreshed *networkIsolatedModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	var (
		net *v1.VDCNetworkIsolated
		err error
	)

	if urn.IsNetwork(planOrState.ID.Get()) {
		net, err = r.vdc.GetNetworkIsolated(planOrState.ID.Get())
	} else {
		net, err = r.vdc.GetNetworkIsolated(planOrState.Name.Get())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error getting isolated network", err.Error())
		return
	}

	// Populate the state with the network data
	stateRefreshed.ID.Set(net.ID)
	stateRefreshed.Name.Set(net.Name)
	stateRefreshed.Description.Set(net.Description)
	stateRefreshed.VDC.Set(r.vdc.GetName())
	stateRefreshed.Gateway.Set(net.Subnet.Gateway)
	stateRefreshed.PrefixLength.SetInt(net.Subnet.PrefixLength)
	stateRefreshed.DNS1.Set(net.Subnet.DNSServer1)
	stateRefreshed.DNS2.Set(net.Subnet.DNSServer2)
	stateRefreshed.DNSSuffix.Set(net.Subnet.DNSSuffix)
	stateRefreshed.GuestVLANAllowed.SetPtr(net.GuestVLANTaggingAllowed)

	x := []*networkIsolatedModelStaticIPPool{}
	for _, ipRange := range net.Subnet.IPRanges {
		n := &networkIsolatedModelStaticIPPool{
			StartAddress: supertypes.NewStringNull(),
			EndAddress:   supertypes.NewStringNull(),
		}
		n.StartAddress.Set(ipRange.StartAddress)
		n.EndAddress.Set(ipRange.EndAddress)
		x = append(x, n)
	}

	diags.Append(stateRefreshed.StaticIPPool.Set(ctx, x)...)
	if diags.HasError() {
		return nil, true, diags
	}

	return stateRefreshed, true, nil
}
