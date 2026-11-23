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

// Package vapp provides a Terraform resource to manage vApps.
package vapp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
)

const vappSubsystem = "vapp"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vappResource{}
	_ resource.ResourceWithConfigure   = &vappResource{}
	_ resource.ResourceWithImportState = &vappResource{}
)

// NewVappResource is a helper function to simplify the provider implementation.
func NewVappResource() resource.Resource {
	return &vappResource{}
}

// vappResource is the resource implementation.
type vappResource struct {
	client       *client.CloudAvenue
	vappClientV2 vAppClientV2
	vdc          resolvedVDC
}

// Metadata returns the resource type name.
func (r *vappResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *vappResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vappSchema().GetResource(ctx)
}

// Init resolves the VDC of the vApp. The SDK client is built in Configure.
func (r *vappResource) Init(ctx context.Context, rm *vappResourceModel) (diags diag.Diagnostics) {
	resolved, diags := r.vappClientV2.resolveVDC(ctx, rm.VDC.ValueString())
	if diags.HasError() {
		return diags
	}
	r.vdc = resolved

	return diags
}

func (r *vappResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vappResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	var plan *vappResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if tflog.IsDebug(ctx) {
		tflog.SubsystemDebug(ctx, vappSubsystem, "Creating vApp", map[string]any{attrVappName: plan.VAppName.ValueString(), attrVDC: plan.VDC.ValueString()})
	}

	// Create vApp
	app, err := r.vappClientV2.vappClient.ComposeVApp(ctx, types.ParamsComposeVApp{VDCID: r.vdc.ID, Name: plan.VAppName.Get(), Description: plan.Description.Get()})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionCreate, fmt.Sprintf("vApp %s", plan.VAppName.Get()), err)
		return
	}

	// Wait for the vApp to leave the UNRESOLVED status
	if err := r.waitVAppResolved(ctx, app.ID); err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionCreate, fmt.Sprintf("vApp %s(%s) status", app.Name, app.ID), err)
		return
	}

	// Persist the ID right away so the vApp is not orphaned if a later step fails
	plan.VAppID.Set(app.ID)
	plan.VDC.Set(r.vdc.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("Error creating vApp", fmt.Sprintf("vApp %s(%s) not found right after creation", app.Name, app.ID))
		return
	}

	// Update vApp
	resp.Diagnostics.Append(r.updateVapp(ctx, plan, stateRefreshed)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, _, d = r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

const (
	vappStatusUnresolved     = "UNRESOLVED"
	vappResolvedTimeout      = 90 * time.Second
	vappResolvedPollInterval = 2 * time.Second
)

// waitVAppResolved polls GetVApp until the vApp status is no longer UNRESOLVED.
// ComposeVApp only waits for the job, not for the status transition.
func (r *vappResource) waitVAppResolved(ctx context.Context, vappID string) error {
	ctx, cancel := context.WithTimeout(ctx, vappResolvedTimeout)
	defer cancel()

	ticker := time.NewTicker(vappResolvedPollInterval)
	defer ticker.Stop()

	for {
		app, err := r.vappClientV2.vappClient.GetVApp(ctx, types.ParamsGetVApp{ID: vappID})
		if err != nil {
			return err
		}
		if tflog.IsDebug(ctx) {
			tflog.SubsystemDebug(ctx, vappSubsystem, "vApp status", map[string]any{"id": vappID, "status": app.Status})
		}
		if app.Status != vappStatusUnresolved {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for status != %s (last status %s): %w", vappStatusUnresolved, app.Status, ctx.Err())
		case <-ticker.C:
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vappResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Read)()

	state := new(vappResourceModel)

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
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed plan
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vappResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = new(vappResourceModel)
		state = new(vappResourceModel)
	)

	// Get current state
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
	if tflog.IsDebug(ctx) {
		tflog.SubsystemDebug(ctx, vappSubsystem, "Updating vApp", map[string]any{attrVappName: plan.VAppName.ValueString(), attrVDC: plan.VDC.ValueString()})
	}

	// Update vApp
	resp.Diagnostics.Append(r.updateVapp(ctx, plan, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read from the plan (not the old state) so removed guest properties are not kept
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
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
func (r *vappResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Delete)()

	// Get current state
	state := new(vappResourceModel)

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Request vApp
	err := r.vappClientV2.vappClient.DeleteVApp(ctx, types.ParamsDeleteVApp{ID: state.VAppID.ValueString()})
	if err != nil {
		if cerrs.IsNotFound(err) {
			cerrs.AddError(&resp.Diagnostics, cerrs.ActionDelete, fmt.Sprintf("vApp not found %s(%s)", state.VAppName.ValueString(), state.VAppID.ValueString()), err)
			return
		}
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionDelete, fmt.Sprintf("vApp %s(%s)", state.VAppName.ValueString(), state.VAppID.ValueString()), err)
		return
	}
}

func (r *vappResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Import)()

	var (
		diags        diag.Diagnostics
		resolved     resolvedVDC
		vAppIDOrName string
	)

	// Split req.ID with dot. ID format is  VDCName.VAppIDOrName or VAppIDOrName
	idParts := strings.Split(req.ID, ".")

	switch len(idParts) {
	case 1:
		vAppIDOrName = idParts[0]
		resolved, diags = r.vappClientV2.resolveVDC(ctx, "")
		resp.Diagnostics.Append(diags...)
	case 2:
		vAppIDOrName = idParts[1]
		resolved, diags = r.vappClientV2.resolveVDC(ctx, idParts[0])
		resp.Diagnostics.Append(diags...)
	default:
		resp.Diagnostics.AddError("Invalid ID format", fmt.Sprintf("ID format is VDCName.VAppIDOrName or VAppIDOrName, got: %s", req.ID))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var p types.ParamsGetVApp
	if urn.IsValid(vAppIDOrName) && urn.IsVAPP(vAppIDOrName) {
		p = types.ParamsGetVApp{ID: vAppIDOrName}
	} else {
		p = types.ParamsGetVApp{Name: vAppIDOrName}
	}
	app, err := r.vappClientV2.vappClient.GetVApp(ctx, p)
	if err != nil {
		if cerrs.IsNotFound(err) {
			cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, fmt.Sprintf("vApp not found %s", vAppIDOrName), err)
			return
		}
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, fmt.Sprintf("vApp %s", vAppIDOrName), err)
		return
	}

	if !vappInVDC(app, resolved) {
		resp.Diagnostics.AddError("Vdc not found", fmt.Sprintf("App vdc id %s different from resolved vdc id %s", app.VDC, resolved.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), app.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), app.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vdc"), resolved.Name)...)
}

// updateVapp make updates only on elements that must be updated.
func (r *vappResource) updateVapp(ctx context.Context, plan, state *vappResourceModel) (diags diag.Diagnostics) {
	var runtimeLease, storageLease *int
	var planLease *vappResourceModelLease
	// Get lease config
	if !plan.Lease.IsNull() {
		lease, d := plan.Lease.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		// Lease is unknown: nothing to send
		if lease != nil {
			planLease = lease
			runtimeLease = lease.RuntimeLeaseInSec.GetIntPtr()
			storageLease = lease.StorageLeaseInSec.GetIntPtr()
		}
	}

	p := types.ParamsUpdateVApp{ID: plan.VAppID.ValueString()}
	stateLease, d := state.Lease.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	var modify bool
	switch {
	case planLease == nil:
		// No lease in plan: leave the lease untouched
	case stateLease == nil:
		p.DeploymentLeaseInSeconds = runtimeLease
		p.StorageLeaseInSeconds = storageLease
		modify = true
	default:
		if !planLease.RuntimeLeaseInSec.Equal(stateLease.RuntimeLeaseInSec) {
			p.DeploymentLeaseInSeconds = runtimeLease
			modify = true
		}
		if !planLease.StorageLeaseInSec.Equal(stateLease.StorageLeaseInSec) {
			p.StorageLeaseInSeconds = storageLease
			modify = true
		}
	}

	// Check updates
	if !plan.Description.Equal(state.Description) {
		description := plan.Description.ValueString()
		p.Description = &description
		modify = true
	}

	// Update if needed
	if modify {
		if _, err := r.vappClientV2.vappClient.UpdateVApp(ctx, p); err != nil {
			cerrs.AddError(&diags, cerrs.ActionUpdate, fmt.Sprintf("vApp %s(%s)", plan.VAppName.ValueString(), plan.VAppID.ValueString()), err)
			return diags
		}
	}

	// Update GuestProperties if needed
	if !plan.GuestProperties.Equal(state.GuestProperties) {
		gp, d := plan.guestPropertiesToSDK(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		_, err := r.vappClientV2.vappClient.UpdateVAppGuestProperties(ctx, types.ParamsUpdateVAppGuestProperties{ID: plan.VAppID.ValueString(), GuestProperties: gp})
		if err != nil {
			cerrs.AddError(&diags, cerrs.ActionUpdate, fmt.Sprintf("vApp %s(%s)", plan.VAppName.ValueString(), plan.VAppID.ValueString()), err)
			return diags
		}
	}

	return diags
}

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *vappResource) read(ctx context.Context, planOrState *vappResourceModel) (stateRefreshed *vappResourceModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	app, err := r.vappClientV2.vappClient.GetVApp(ctx, types.ParamsGetVApp{ID: planOrState.VAppID.ValueString()})
	if err != nil {
		if cerrs.IsNotFound(err) {
			tflog.SubsystemDebug(ctx, vappSubsystem, "vApp not found, removing from state", map[string]any{attrVappName: planOrState.VAppName.ValueString(), attrVDC: planOrState.VDC.ValueString()})
			return nil, false, diags
		}
		tflog.SubsystemError(ctx, vappSubsystem, "vApp read failed", map[string]any{attrVappName: planOrState.VAppName.ValueString(), attrVDC: planOrState.VDC.ValueString()})
		cerrs.AddError(&diags, cerrs.ActionRead, fmt.Sprintf("vApp %s", planOrState.VAppName.ValueString()), err)
		return nil, true, diags
	}

	if !vappInVDC(app, r.vdc) {
		diags.AddError("Vdc not found", fmt.Sprintf("App vdc id %s different from resolved vdc id %s", app.VDC, r.vdc.ID))
		return nil, true, diags
	}

	// * Guest properties
	guestProperties, err := r.vappClientV2.vappClient.GetVAppGuestProperties(ctx, types.ParamsGetVAppGuestProperties{ID: app.ID})
	if err != nil {
		tflog.SubsystemError(ctx, vappSubsystem, "vApp guest properties read failed", map[string]any{attrVappName: planOrState.VAppName.ValueString()})
		cerrs.AddError(&diags, cerrs.ActionRead, fmt.Sprintf("vApp %s(%s) guest properties", app.Name, app.ID), err)
		return nil, true, diags
	}

	refreshDiag := stateRefreshed.guestPropertiesFromSDK(ctx, guestProperties)
	diags.Append(refreshDiag...)
	if diags.HasError() {
		return nil, true, diags
	}
	refreshDiag = stateRefreshed.fromSDK(ctx, app, r.vdc.Name)
	diags.Append(refreshDiag...)
	if diags.HasError() {
		return nil, true, diags
	}

	return stateRefreshed, true, diags
}
