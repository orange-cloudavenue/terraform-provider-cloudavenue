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

package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
)

const roleSubsystem = "iam.role"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithConfigure   = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
)

// NewroleResource is a helper function to simplify the provider implementation.
func NewRoleResource() resource.Resource {
	return &roleResource{}
}

// roleResource is the resource implementation.
type roleResource struct {
	client     *client.CloudAvenue
	roleClient roleClientV2
}

func (r *roleResource) Init(_ context.Context, _ *RoleResourceModel) (diags diag.Diagnostics) {
	r.roleClient, diags = InitRoleClientV2(r.client)
	return diags
}

// Metadata returns the resource type name.
func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "role"
}

// Schema defines the schema for the resource.
func (r *roleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = roleSchema().GetResource(ctx)
}

func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// findMissingImpliedRights computes implied rights that need to be added.
// This is a provider-side reimplementation of the govcd.FindMissingImpliedRights logic:
// every requested right is seeded first so an implied right that is itself requested
// is not reported as missing, and each requested right's implied set is examined once.
func findMissingImpliedRights(rightsMap map[string]*sdktypes.ModelRight, requestedRights []string) []string {
	// Track all rights we've seen (both directly requested and implied)
	seen := make(map[string]bool)
	for _, rightName := range requestedRights {
		seen[rightName] = true
	}

	var missing []string

	for _, rightName := range requestedRights {
		// Check if this right has implied rights
		right, exists := rightsMap[rightName]
		if !exists {
			continue
		}

		for _, impliedRight := range right.ImpliedRights {
			if seen[impliedRight.Name] {
				continue
			}
			seen[impliedRight.Name] = true
			missing = append(missing, impliedRight.Name)
		}
	}

	return missing
}

// Create creates the resource and sets the initial Terraform state.
func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_iam_role", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	plan := &RoleResourceModel{}

	// Read the plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if tflog.IsDebug(ctx) {
		tflog.SubsystemDebug(ctx, roleSubsystem, "Creating role", map[string]interface{}{attrName: plan.Name.ValueString()})
	}

	// Get all rights to build a name->ModelRight map
	allRights, err := r.roleClient.ListRights(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving rights", err.Error())
		return
	}

	rightsMap := make(map[string]*sdktypes.ModelRight)
	for _, right := range allRights {
		rightsMap[right.Name] = right
	}

	// Validate and collect requested rights
	planRights, d := plan.Rights.Get(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build list of rights with their IDs
	rightsEntities := make([]sdktypes.ParamsEntityReference, 0)
	for _, rightName := range planRights {
		right, exists := rightsMap[rightName]
		if !exists {
			resp.Diagnostics.AddError("This right does not exist", rightName)
			return
		}
		rightsEntities = append(rightsEntities, sdktypes.ParamsEntityReference{
			ID:   right.ID,
			Name: right.Name,
		})
	}

	// Find missing implied rights
	missingImpliedRights := findMissingImpliedRights(rightsMap, planRights)

	// Print missing implied rights
	if len(missingImpliedRights) > 0 {
		message := "The rights set for this role require the following implied rights to be added:"
		rightsList := ""
		for _, right := range missingImpliedRights {
			rightsList += fmt.Sprintf("\"%s\",\n", right)
		}
		resp.Diagnostics.AddError(message, rightsList)
		return
	}

	// Create the role
	role, err := r.roleClient.CreateOrgRole(ctx, plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		tflog.SubsystemError(ctx, roleSubsystem, "Role create failed", map[string]interface{}{attrName: plan.Name.ValueString()})
		resp.Diagnostics.AddError("Error creating role", err.Error())
		return
	}

	// Assign rights if any
	if len(rightsEntities) > 0 {
		_, err = r.roleClient.ReplaceOrgRoleRights(ctx, role.ID, rightsEntities)
		if err != nil {
			tflog.SubsystemError(ctx, roleSubsystem, "Role rights update failed", map[string]interface{}{attrName: plan.Name.ValueString(), attrRights: len(rightsEntities)})
			resp.Diagnostics.AddError("Error adding rights to role", err.Error())
			return
		}
	}

	// Set Plan state
	plan.ID.Set(role.ID)
	plan.Name.Set(role.Name)
	plan.Description.Set(role.Description)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_iam_role", r.client.GetOrgName(), metrics.Read)()

	var state *RoleResourceModel

	// Read state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if tflog.IsDebug(ctx) {
		tflog.SubsystemDebug(ctx, roleSubsystem, "Reading role", map[string]interface{}{attrName: state.Name.ValueString()})
	}

	// Get Role - use ID from state if available, otherwise use Name
	lookupKey := state.ID.ValueString()
	if lookupKey == "" {
		lookupKey = state.Name.ValueString()
	}

	role, err := r.roleClient.GetOrgRole(ctx, lookupKey)
	if err != nil {
		if cerrs.IsNotFound(err) {
			tflog.SubsystemDebug(ctx, roleSubsystem, "Role not found, removing from state", map[string]interface{}{attrName: state.Name.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		tflog.SubsystemError(ctx, roleSubsystem, "Role read failed", map[string]interface{}{attrName: state.Name.ValueString()})
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, "role", err)
		return
	}

	// Get rights
	rights, err := r.roleClient.ListOrgRoleRights(ctx, role.ID)
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, "role", err)
		return
	}

	assignedRights := []string{}
	for _, right := range rights {
		assignedRights = append(assignedRights, right.Name)
	}

	resp.Diagnostics.Append(state.Rights.Set(ctx, assignedRights)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID.Set(role.ID)
	state.Name.Set(role.Name)
	state.Description.Set(role.Description)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_iam_role", r.client.GetOrgName(), metrics.Delete)()

	var state *RoleResourceModel

	// Read state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if tflog.IsDebug(ctx) {
		tflog.SubsystemDebug(ctx, roleSubsystem, "Deleting role", map[string]interface{}{attrName: state.Name.ValueString()})
	}

	// Get the role - use ID from state if available, otherwise use Name
	lookupKey := state.ID.ValueString()
	if lookupKey == "" {
		lookupKey = state.Name.ValueString()
	}

	role, err := r.roleClient.GetOrgRole(ctx, lookupKey)
	if err != nil {
		if cerrs.IsNotFound(err) {
			tflog.SubsystemDebug(ctx, roleSubsystem, "Role not found, removing from state", map[string]interface{}{attrName: state.Name.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		tflog.SubsystemError(ctx, roleSubsystem, "Role delete lookup failed", map[string]interface{}{attrName: state.Name.ValueString()})
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionDelete, "role", err)
		return
	}

	err = r.roleClient.DeleteOrgRole(ctx, role.ID)
	if err != nil {
		tflog.SubsystemError(ctx, roleSubsystem, "Role delete failed", map[string]interface{}{attrName: state.Name.ValueString()})
		resp.Diagnostics.AddError("Error deleting role", err.Error())
		return
	}
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_iam_role", r.client.GetOrgName(), metrics.Update)()

	var (
		plan, state *RoleResourceModel
		err         error
	)

	// Get state and plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if tflog.IsDebug(ctx) {
		tflog.SubsystemDebug(ctx, roleSubsystem, "Updating role", map[string]interface{}{attrName: state.Name.ValueString()})
	}

	// Get the role - use ID from state if available, otherwise use Name
	lookupKey := state.ID.ValueString()
	if lookupKey == "" {
		lookupKey = state.Name.ValueString()
	}

	role, err := r.roleClient.GetOrgRole(ctx, lookupKey)
	if err != nil {
		tflog.SubsystemError(ctx, roleSubsystem, "Role update lookup failed", map[string]interface{}{attrName: state.Name.ValueString()})
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionUpdate, "role", err)
		return
	}

	// Update the role Name or Description
	if (!plan.Name.Equal(state.Name)) || (!plan.Description.Equal(state.Description)) {
		role, err = r.roleClient.UpdateOrgRole(ctx, role.ID, plan.Name.ValueString(), plan.Description.ValueString())
		if err != nil {
			tflog.SubsystemError(ctx, roleSubsystem, "Role attribute update failed", map[string]interface{}{attrName: plan.Name.ValueString()})
			resp.Diagnostics.AddError("Error updating role", err.Error())
			return
		}
	}

	// Get all rights to build a name->ModelRight map
	allRights, err := r.roleClient.ListRights(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving rights", err.Error())
		return
	}

	rightsMap := make(map[string]*sdktypes.ModelRight)
	for _, right := range allRights {
		rightsMap[right.Name] = right
	}

	// Check rights are valid
	planRights, d := plan.Rights.Get(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build list of rights with their IDs
	rightsEntities := make([]sdktypes.ParamsEntityReference, 0)
	for _, rightName := range planRights {
		right, exists := rightsMap[rightName]
		if !exists {
			resp.Diagnostics.AddError("This right does not exist", rightName)
			return
		}
		rightsEntities = append(rightsEntities, sdktypes.ParamsEntityReference{
			ID:   right.ID,
			Name: right.Name,
		})
	}

	// Find missing implied rights
	missingImpliedRights := findMissingImpliedRights(rightsMap, planRights)

	// Print missing implied rights
	if len(missingImpliedRights) > 0 {
		message := "The rights set for this role require the following implied rights to be added:"
		rightsList := ""
		for _, right := range missingImpliedRights {
			rightsList += fmt.Sprintf("\"%s\",\n", right)
		}
		resp.Diagnostics.AddError(message, rightsList)
		return
	}

	// Update the role rights
	if len(rightsEntities) > 0 {
		_, err = r.roleClient.ReplaceOrgRoleRights(ctx, role.ID, rightsEntities)
		if err != nil {
			resp.Diagnostics.AddError("Error updating role rights", err.Error())
			return
		}
	} else {
		// Check if there are current rights to remove
		currentRights, err := r.roleClient.ListOrgRoleRights(ctx, role.ID)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving role rights", err.Error())
			return
		}
		if len(currentRights) > 0 {
			err = r.roleClient.RemoveAllOrgRoleRights(ctx, role.ID)
			if err != nil {
				resp.Diagnostics.AddError("Error removing role rights", err.Error())
				return
			}
		}
	}

	// Set Plan state
	plan.ID.Set(role.ID)
	plan.Name.Set(role.Name)
	plan.Description.Set(role.Description)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_iam_role", r.client.GetOrgName(), metrics.Import)()

	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
