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
	"slices"
	"strings"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &orgNetworkResource{}
	_ resource.ResourceWithConfigure   = &orgNetworkResource{}
	_ resource.ResourceWithImportState = &orgNetworkResource{}
)

// NewOrgNetworkResource is a helper function to simplify the provider implementation.
func NewOrgNetworkResource() resource.Resource {
	return &orgNetworkResource{}
}

// orgNetworkResource is the resource implementation.
type orgNetworkResource struct {
	client       *client.CloudAvenue
	vappClientV2 vAppClientV2
	vdc          resolvedVDC
	vapp         *sdktypes.ModelVApp
}

// Metadata returns the resource type name.
func (r *orgNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_org_network"
}

// Schema defines the schema for the resource.
func (r *orgNetworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetRoutedVapp()).GetResource(ctx)
}

// Init resolves the VDC and the parent vApp. The SDK client is built in Configure.
func (r *orgNetworkResource) Init(ctx context.Context, rm *orgNetworkModel) (diags diag.Diagnostics) {
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
func (r *orgNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *orgNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	var (
		plan *orgNetworkModel
	)

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

	orgNetworkName := plan.NetworkName.ValueString()
	vdcNetwork, err := r.vappClientV2.vdcClient.GetVDCNetwork(ctx, sdktypes.ParamsGetVDCNetwork{Name: orgNetworkName, OwnerID: r.vdc.ID})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, fmt.Sprintf("Vdc Network %s", orgNetworkName), err)
		return
	}

	retain := false
	section, err := r.vappClientV2.vappClient.AddOrgNetworkToVApp(
		ctx,
		sdktypes.ParamsAddOrgNetworkToVApp{
			VAppID: r.vapp.ID,
			Network: sdktypes.ModelVAppRuntimeNetworkConfig{
				NetworkName: vdcNetwork.Name,
				Configuration: sdktypes.ModelVAppRuntimeNetworkConfiguration{
					FenceMode:                      "bridged",
					RetainNetInfoAcrossDeployments: &retain,
					ParentNetwork:                  vdcNetwork.VAppParentNetworkReference(),
				},
			},
		},
	)
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionCreate, fmt.Sprintf("Add Network %s to vApp %s(%s)", orgNetworkName, r.vapp.Name, r.vapp.ID), err)
		return
	}

	vAppNetwork, diags := plan.findOrgNetwork(section.NetworkConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := &orgNetworkModel{
		ID:          types.StringValue(urn.Normalize(urn.Network, urn.ExtractUUID(vAppNetwork.ID)).String()),
		VAppName:    plan.VAppName,
		VAppID:      plan.VAppID,
		VDC:         types.StringValue(r.vdc.Name),
		NetworkName: plan.NetworkName,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *orgNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Read)()

	var state *orgNetworkModel

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

	vAppNetworkConfig, err := r.vappClientV2.vappClient.GetVAppNetworkConfig(ctx, sdktypes.ParamsGetVAppNetworkConfig{ID: r.vapp.ID})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, fmt.Sprintf("vApp %s(%s)", r.vapp.Name, r.vapp.ID), err)
		return
	}

	vAppNetwork, d := state.findOrgNetwork(vAppNetworkConfig.NetworkConfig)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := &orgNetworkModel{
		ID:          types.StringValue(urn.Normalize(urn.Network, urn.ExtractUUID(vAppNetwork.ID)).String()),
		VAppName:    state.VAppName,
		VAppID:      state.VAppID,
		VDC:         types.StringValue(r.vdc.Name),
		NetworkName: state.NetworkName,
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *orgNetworkResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Update)()
	// No update for this resource
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *orgNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Delete)()

	var state *orgNetworkModel

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

	// A deployed vApp (poweredON, poweredOFF, suspended, partiallyDeployed, ...)
	// cannot have a network detached (API error 409). Only an undeployed
	// (resolved) vApp accepts network detach, so the vApp is undeployed first.
	vapp, err := r.vappClientV2.vappClient.GetVApp(ctx, sdktypes.ParamsGetVApp{ID: r.vapp.ID})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, fmt.Sprintf("vApp %s(%s)", r.vapp.Name, r.vapp.ID), err)
		return
	}

	vAppRequiredStatuses := []string{resolved}
	vAppStatusBeforeAction := vapp.Status
	if !slices.Contains(vAppRequiredStatuses, vapp.Status) {
		// A suspended vApp cannot be undeployed directly: discard its state first.
		if vapp.Status == suspended {
			if err := r.vappClientV2.vappClient.DiscardSuspendedStateVApp(ctx, vapp.ID); err != nil {
				cerrs.AddError(&resp.Diagnostics, cerrs.ActionDelete, fmt.Sprintf("Error discarding suspended state of vApp %s(%s)", r.vapp.Name, r.vapp.ID), err)
				return
			}
		}

		err := r.vappClientV2.vappClient.UndeployVApp(ctx, sdktypes.ParamsUndeployVApp{ID: vapp.ID, Mode: sdktypes.VAppUndeployModePowerOff})
		if err != nil {
			cerrs.AddError(&resp.Diagnostics, cerrs.ActionDelete, fmt.Sprintf("Error undeploying vApp %s(%s)", r.vapp.Name, r.vapp.ID), err)
			return
		}
	}

	if _, err := r.vappClientV2.vappClient.RemoveVAppNetwork(ctx, sdktypes.ParamsRemoveVAppNetwork{VAppID: vapp.ID, NetworkName: state.NetworkName.ValueString()}); err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionDelete, fmt.Sprintf("Error deleting vApp network from vApp %s(%s)", r.vapp.Name, r.vapp.ID), err)
	}

	// Restore power if the vApp was powered on before the network was detached.
	if slices.Contains([]string{poweredON, partiallyDeployed, partiallyPoweredOff, partiallySuspended}, vAppStatusBeforeAction) {
		// Power On vApp
		err := r.vappClientV2.vappClient.PowerOnVApp(ctx, vapp.ID)
		if err != nil {
			resp.Diagnostics.AddError("Error powering on vApp", err.Error())
			return
		}
	}
}

func (r *orgNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Import)()

	var state *orgNetworkModel
	resourceURI := strings.Split(req.ID, ".")

	if len(resourceURI) != 3 && len(resourceURI) != 2 {
		resp.Diagnostics.AddError("Error importing org_network", "Wrong resource URI format. Expected vdc.vapp.org_network_name or vapp.org_network_name")
		return
	}

	state = &orgNetworkModel{
		VAppName:    types.StringValue(resourceURI[0]),
		NetworkName: types.StringValue(resourceURI[1]),
	}

	if len(resourceURI) == 3 {
		state = &orgNetworkModel{
			VDC:         types.StringValue(resourceURI[0]),
			VAppName:    types.StringValue(resourceURI[1]),
			NetworkName: types.StringValue(resourceURI[2]),
		}
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
