/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vm provides a Terraform resource.
package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

// Ensure the implementation satisfies the expected interfaces.VAppName.
var (
	_ resource.Resource              = &insertedMediaResource{}
	_ resource.ResourceWithConfigure = &insertedMediaResource{}
)

// NewInsertedMediaResource is a helper function to simplify the provider implementation.
func NewInsertedMediaResource() resource.Resource {
	return &insertedMediaResource{}
}

// Metadata returns the resource type name.
func (r *insertedMediaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_inserted_media"
}

// Schema defines the schema for the resource.
func (r *insertedMediaResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vmInsertedMediaSuperSchema().GetResource(ctx)
}

func (r *insertedMediaResource) Init(_ context.Context, rm *insertedMediaResourceModel) (diags diag.Diagnostics) {
	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.vdc, diags = vdc.Init(r.client, rm.VDC)
	if diags.HasError() {
		return
	}

	r.vapp, diags = vapp.Init(r.client, r.vdc, rm.VAppID, rm.VAppName)
	if diags.HasError() {
		return
	}

	r.vm, diags = vm.Get(r.vapp, vm.GetVMOpts{
		ID:   types.StringNull(),
		Name: rm.VMName,
	})

	return
}

func (r *insertedMediaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *insertedMediaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vm_inserted_media", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	var (
		plan *insertedMediaResourceModel
		err  error
	)

	// Read the plan
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
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	// Insert media
	task, err := r.vm.HandleInsertMedia(r.org.Org.Org, plan.Catalog.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error inserting media", err.Error())
		return
	}
	if err = task.WaitTaskCompletion(); err != nil {
		resp.Diagnostics.AddError("Error during inserting media", err.Error())
		return
	}

	plan.ID = types.StringValue(r.vm.GetID())
	plan.VDC = types.StringValue(r.vdc.GetName())

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *insertedMediaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vm_inserted_media", r.client.GetOrgName(), metrics.Read)()

	var state *insertedMediaResourceModel

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

	// Check if media is mounted
	var isIsoMounted bool

	for _, hardwareItem := range r.vm.GetVirtualHardwareSectionItems() {
		if hardwareItem.ResourceType == 15 { // 15 = CD/DVD Drive
			isIsoMounted = true
			break
		}
	}

	if !isIsoMounted {
		resp.Diagnostics.AddError("Media not mounted", "Media is not mounted on the VM")
		resp.State.RemoveResource(ctx)
		return
	}

	// Set Plan state
	plan := &insertedMediaResourceModel{
		ID:       types.StringValue(r.vm.GetID()),
		VDC:      types.StringValue(r.vdc.GetName()),
		Catalog:  state.Catalog,
		Name:     state.Name,
		VAppName: state.VAppName,
		VAppID:   state.VAppID,
		VMName:   state.VMName,
		// EjectForce: state.EjectForce,
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *insertedMediaResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	/* linked with issue - Disable attributes - Issue referrer: vmware/go-vcloud-director#552
	var plan, state *insertedMediaResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	// Check if VM exists
	vm, err := r.vapp.GetVMByName(state.VMName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	plan = &insertedMediaResourceModel{
		ID:       types.StringValue(vm.VM.ID),
		VDC:      plan.VDC,
		Catalog:  plan.Catalog,
		Name:     plan.Name,
		VAppName: plan.VAppName,
		VMName:   plan.VMName,
		// EjectForce: plan.EjectForce,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	*/
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *insertedMediaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vm_inserted_media", r.client.GetOrgName(), metrics.Delete)()

	var state *insertedMediaResourceModel

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
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	// Eject media
	if _, err := r.vm.HandleEjectMediaAndAnswer(r.org.Org.Org, state.Catalog.ValueString(), state.Name.ValueString(), true); err != nil {
		resp.Diagnostics.AddError("Error ejecting media", err.Error())
	}
}
