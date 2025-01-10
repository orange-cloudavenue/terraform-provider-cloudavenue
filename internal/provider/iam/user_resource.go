/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package org provides a Terraform resource to manage org users.
package iam

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// NewuserResource is a helper function to simplify the provider implementation.
func NewIAMUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = userSchema().GetResource(ctx)
}

func (r *userResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userResource) Init(_ context.Context, rm *userResourceModel) (diags diag.Diagnostics) {
	r.adminOrg, diags = adminorg.Init(r.client)
	return
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_iam_user", r.client.GetOrgName(), metrics.Create)()

	plan := &userResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userData := govcd.OrgUserConfiguration{
		ProviderType:    govcd.OrgUserProviderIntegrated,
		Name:            plan.Name.Get(),
		RoleName:        plan.RoleName.Get(),
		FullName:        plan.FullName.Get(),
		EmailAddress:    plan.Email.Get(),
		Telephone:       plan.Telephone.Get(),
		IsEnabled:       plan.Enabled.Get(),
		Password:        plan.Password.Get(),
		DeployedVmQuota: plan.DeployedVMQuota.GetInt(),
		StoredVmQuota:   plan.StoredVMQuota.GetInt(),
	}

	user, err := r.adminOrg.CreateUserSimple(userData)
	if err != nil {
		// Here bypass the error where the API return user not found after creation
		// This is a known issue in the API
		if !govcd.ContainsNotFound(err) {
			resp.Diagnostics.AddError("Error creating user", err.Error())
			return
		}
	}

	// Catch if user is not nil (nil means user not found after creation)
	if user != nil {
		plan.ID.Set(user.User.ID)
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("User not found after create", fmt.Sprintf("User with name %s not found after create", plan.Name.Get()))
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_iam_user", r.client.GetOrgName(), metrics.Read)()

	state := &userResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("User not found", fmt.Sprintf("User with name %s(%s) not found after update.", state.Name.Get(), state.ID.Get()))
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_iam_user", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &userResourceModel{}
		state = &userResourceModel{}
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

	user, err := r.adminOrg.GetUserById(state.ID.Get(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	userData := govcd.OrgUserConfiguration{
		ProviderType:    govcd.OrgUserProviderIntegrated,
		Name:            plan.Name.Get(),
		RoleName:        plan.RoleName.Get(),
		FullName:        plan.FullName.Get(),
		EmailAddress:    plan.Email.Get(),
		Telephone:       plan.Telephone.Get(),
		IsEnabled:       plan.Enabled.Get(),
		Password:        plan.Password.Get(),
		DeployedVmQuota: plan.DeployedVMQuota.GetInt(),
		StoredVmQuota:   plan.StoredVMQuota.GetInt(),
	}

	if err = user.UpdateSimple(userData); err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_iam_user", r.client.GetOrgName(), metrics.Delete)()

	state := &userResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.adminOrg.GetUserById(state.ID.Get(), true)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	if err = user.Delete(state.TakeOwnership.Get()); err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_iam_user", r.client.GetOrgName(), metrics.Import)()

	userData := &userResourceModel{}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, userData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.adminOrg.GetUserByName(req.ID, true)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.Diagnostics.AddError("User not found", fmt.Sprintf("User with name %s not found", req.ID))
			return
		}
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	if user.User.ProviderType != govcd.OrgUserProviderIntegrated {
		resp.Diagnostics.AddError("User is not an local user", fmt.Sprintf("User with name %s is %s type and not local user", req.ID, user.User.ProviderType))
		return
	}

	userData.ID.Set(user.User.ID)
	userData.Name.Set(user.User.Name)

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, userData)
	if !found {
		resp.Diagnostics.AddError("User not found", fmt.Sprintf("User with name %s not found", req.ID))
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *userResource) read(_ context.Context, planOrState *userResourceModel) (stateRefreshed *userResourceModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	var (
		user *govcd.OrgUser
		err  error
	)

	if stateRefreshed.ID.IsKnown() {
		user, err = r.adminOrg.GetUserByNameOrId(stateRefreshed.ID.Get(), true)
	} else {
		user, err = r.adminOrg.GetUserByNameOrId(stateRefreshed.Name.Get(), true)
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving user", err.Error())
		return
	}

	stateRefreshed.ID.Set(user.User.ID)
	stateRefreshed.Name.Set(user.User.Name)
	stateRefreshed.RoleName.Set(user.User.Role.Name)
	stateRefreshed.FullName.Set(user.User.FullName)
	stateRefreshed.Email.Set(user.User.EmailAddress)
	stateRefreshed.Telephone.Set(user.User.Telephone)
	stateRefreshed.Enabled.Set(user.User.IsEnabled)
	stateRefreshed.DeployedVMQuota.SetInt(user.User.DeployedVmQuota)
	stateRefreshed.StoredVMQuota.SetInt(user.User.StoredVmQuota)

	return stateRefreshed, true, diags
}
