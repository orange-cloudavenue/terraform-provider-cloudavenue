// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// Package vapp provides a Terraform resource.
package vapp

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &isolatedNetworkResource{}
	_ resource.ResourceWithConfigure   = &isolatedNetworkResource{}
	_ resource.ResourceWithImportState = &isolatedNetworkResource{}
)

// NewIsolatedNetworkResource returns isolated network resource.
func NewIsolatedNetworkResource() resource.Resource {
	return &isolatedNetworkResource{}
}

// isolatedNetworkResource is isolated network resource implementation.
type isolatedNetworkResource struct {
	client       *client.CloudAvenue
	vappClientV2 vAppClientV2
	vdc          resolvedVDC
	vapp         *sdktypes.ModelVApp
}

// Metadata returns the resource type name.
func (r *isolatedNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "isolated_network"
}

// Schema defines the schema for the resource.
func (r *isolatedNetworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = isolatedNetworkSchema().GetResource(ctx)
}

// Init resolves the VDC and the parent vApp. The SDK client is built in Configure.
func (r *isolatedNetworkResource) Init(ctx context.Context, rm *isolatedNetworkModel) (diags diag.Diagnostics) {
	vdc, diags := r.vappClientV2.resolveVDC(ctx, rm.VDC.ValueString())
	if diags.HasError() {
		return diags
	}
	r.vdc = vdc

	app, diags := r.vappClientV2.resolveVApp(ctx, r.vdc, rm.VAppID.ValueString(), rm.VAppName.ValueString())
	if diags.HasError() {
		return diags
	}
	r.vapp = app

	return diags
}

// Configure clients
func (r *isolatedNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Expected *client.CloudAvenue, got %T. Report this to provider maintainers.", req.ProviderData),
		)

		return
	}

	r.client = client

	v2, diags := InitVAppClientV2(client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.vappClientV2 = v2
}

// Create creates the resource and sets the initial Terraform state.
func (r *isolatedNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Create)()

	plan := &isolatedNetworkModel{}

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	unlock, d := lockVApp(ctx, r.vdc.Name, r.vapp.Name)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer unlock()

	network, d := plan.toSDK(ctx, nil)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create network
	_, err := r.vappClientV2.vappClient.CreateVAppNetwork(ctx, sdktypes.ParamsCreateVAppNetwork{VAppID: r.vapp.ID, Network: network})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionCreate, fmt.Sprintf("vApp isolated network %s", plan.Name.ValueString()), err)
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("Error creating vApp isolated network", fmt.Sprintf("network %s not found in vApp %s(%s) right after creation", plan.Name.ValueString(), r.vapp.Name, r.vapp.ID))
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *isolatedNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Read)()

	state := &isolatedNetworkModel{}

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

	// Refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *isolatedNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &isolatedNetworkModel{}
		state = &isolatedNetworkModel{}
	)

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	unlock, d := lockVApp(ctx, r.vdc.Name, r.vapp.Name)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer unlock()

	// Current network is the base of the update to preserve unmanaged runtime fields
	current, d := r.findNetwork(ctx, plan.Name.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	if current == nil {
		resp.Diagnostics.AddError("Error updating vApp isolated network", fmt.Sprintf("network %s not found in vApp %s(%s)", plan.Name.ValueString(), r.vapp.Name, r.vapp.ID))
		return
	}

	network, d := plan.toSDK(ctx, current)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update network
	_, err := r.vappClientV2.vappClient.UpdateVAppNetwork(ctx, sdktypes.ParamsUpdateVAppNetwork{VAppID: r.vapp.ID, NetworkName: plan.Name.ValueString(), Network: network})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionUpdate, fmt.Sprintf("vApp isolated network %s", plan.Name.ValueString()), err)
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("Error updating vApp isolated network", fmt.Sprintf("network %s not found in vApp %s(%s) right after update", plan.Name.ValueString(), r.vapp.Name, r.vapp.ID))
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *isolatedNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Delete)()

	state := &isolatedNetworkModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	unlock, d := lockVApp(ctx, r.vdc.Name, r.vapp.Name)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer unlock()

	_, err := r.vappClientV2.vappClient.RemoveVAppNetwork(ctx, sdktypes.ParamsRemoveVAppNetwork{VAppID: r.vapp.ID, NetworkName: state.Name.ValueString()})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionDelete, fmt.Sprintf("vApp %s(%s) network %s", r.vapp.Name, r.vapp.ID, state.Name.ValueString()), err)
		return
	}
}

func (r *isolatedNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Import)()

	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 3 && len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: vdc.vapp_name.network_name or vapp_name.network_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)

	if len(idParts) == 3 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vdc"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), idParts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[2])...)
	}
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
// found is false (without error) when the network does not exist in the vApp.
func (r *isolatedNetworkResource) read(ctx context.Context, planOrState *isolatedNetworkModel) (stateRefreshed *isolatedNetworkModel, found bool, diags diag.Diagnostics) {
	network, d := r.findNetwork(ctx, planOrState.Name.ValueString())
	diags.Append(d...)
	if diags.HasError() {
		return nil, true, diags
	}
	if network == nil {
		return nil, false, diags
	}

	stateRefreshed = planOrState.Copy()
	diags.Append(stateRefreshed.fromSDK(ctx, network)...)
	if diags.HasError() {
		return nil, true, diags
	}

	stateRefreshed.VDC.Set(r.vdc.Name)
	stateRefreshed.VAppName.Set(r.vapp.Name)
	stateRefreshed.VAppID.Set(r.vapp.ID)

	return stateRefreshed, true, diags
}

// findNetwork returns the vApp network matching the name, or nil if there is none.
func (r *isolatedNetworkResource) findNetwork(ctx context.Context, networkName string) (*sdktypes.ModelVAppRuntimeNetworkConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	section, err := r.vappClientV2.vappClient.GetVAppNetworkConfig(ctx, sdktypes.ParamsGetVAppNetworkConfig{ID: r.vapp.ID})
	if err != nil {
		cerrs.AddError(&diags, cerrs.ActionRead, fmt.Sprintf("vApp %s(%s) network config", r.vapp.Name, r.vapp.ID), err)
		return nil, diags
	}

	for i := range section.NetworkConfig {
		if section.NetworkConfig[i].NetworkName == networkName {
			return &section.NetworkConfig[i], diags
		}
	}

	return nil, diags
}
