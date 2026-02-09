/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vdcgResource{}
	_ resource.ResourceWithConfigure   = &vdcgResource{}
	_ resource.ResourceWithImportState = &vdcgResource{}
)

// NewvdcgResource is a helper function to simplify the provider implementation.
func NewVDCGResource() resource.Resource {
	return &vdcgResource{}
}

// vdcgResource is the resource implementation.
type vdcgResource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
}

// Init Initializes the resource.
func (r *vdcgResource) Init(_ context.Context, _ *vdcgModel) (diags diag.Diagnostics) {
	r.adminOrg, diags = adminorg.Init(r.client)
	return diags
}

// Metadata returns the resource type name.
func (r *vdcgResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *vdcgResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vdcgSchema(ctx).GetResource(ctx)
}

func (r *vdcgResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vdcgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Create)()

	plan := &vdcgModel{}

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

	vdcIDs, d := plan.VDCIDs.Get(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	vdcGroup, err := r.adminOrg.CreateNsxtVdcGroup(plan.Name.Get(), plan.Description.Get(), vdcIDs[0], vdcIDs)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VDC Group", err.Error())
		return
	}

	plan.ID.Set(vdcGroup.VdcGroup.Id)

	// Use generic read function to refresh the state
	state, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vdcgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Read)()

	state := &vdcgModel{}

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
func (r *vdcgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &vdcgModel{}
		state = &vdcgModel{}
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

	vdcIDs, d := plan.VDCIDs.Get(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Here use GetVdcGroupById instead of GetVdcGroupByNameOrID because we want to update the name of VDC Group
	vdcGroup, err := r.adminOrg.GetVdcGroupById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VDC Group", err.Error())
		return
	}

	if _, err := vdcGroup.Update(plan.Name.Get(), plan.Description.Get(), vdcIDs); err != nil {
		resp.Diagnostics.AddError("Error updating VDC Group", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("Resource with name %s(%s) not found after update.", plan.Name.Get(), plan.ID.Get()))
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
func (r *vdcgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Delete)()

	state := &vdcgModel{}

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

	vdcGroup, err := r.adminOrg.GetVdcGroupById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VDC Group", err.Error())
		return
	}

	if err = vdcGroup.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting VDC Group", err.Error())
		return
	}
}

func (r *vdcgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Import)()

	// id format is vdcGroupIDOrName

	var d diag.Diagnostics

	r.adminOrg, d = adminorg.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	var (
		vdcGroup *govcd.VdcGroup
		err      error
	)

	if urn.IsVDCGroup(req.ID) {
		vdcGroup, err = r.adminOrg.GetVdcGroupById(req.ID)
	} else {
		vdcGroup, err = r.adminOrg.GetVdcGroupByName(req.ID)
	}
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC Group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), vdcGroup.VdcGroup.Id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), vdcGroup.VdcGroup.Name)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *vdcgResource) read(ctx context.Context, planOrState *vdcgModel) (stateRefreshed *vdcgModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	var (
		vdcGroup *govcd.VdcGroup
		err      error
	)

	if planOrState.ID.IsKnown() {
		vdcGroup, err = r.adminOrg.GetVdcGroupById(planOrState.ID.Get())
	} else {
		vdcGroup, err = r.adminOrg.GetVdcGroupByName(planOrState.Name.Get())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error reading VDC Group", err.Error())
		return nil, true, diags
	}

	var vdcIDs []string
	for _, vdc := range vdcGroup.VdcGroup.ParticipatingOrgVdcs {
		vdcIDs = append(vdcIDs, vdc.VdcRef.ID)
	}

	stateRefreshed.ID.Set(vdcGroup.VdcGroup.Id)
	stateRefreshed.Name.Set(vdcGroup.VdcGroup.Name)
	stateRefreshed.Description.Set(vdcGroup.VdcGroup.Description)
	stateRefreshed.Status.Set(vdcGroup.VdcGroup.Status)
	stateRefreshed.Type.Set(vdcGroup.VdcGroup.Type)
	diags.Append(stateRefreshed.VDCIDs.Set(ctx, vdcIDs)...)

	return stateRefreshed, true, diags
}
