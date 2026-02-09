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

package iam

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/iam"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &UserSAMLResource{}
	_ resource.ResourceWithConfigure   = &UserSAMLResource{}
	_ resource.ResourceWithImportState = &UserSAMLResource{}
)

// NewUserSAMLResource is a helper function to simplify the provider implementation.
func NewUserSAMLResource() resource.Resource {
	return &UserSAMLResource{}
}

// UserSAMLResource is the resource implementation.
type UserSAMLResource struct {
	client    *client.CloudAvenue
	iamClient *iam.Client
}

// Init Initializes the resource.
func (r *UserSAMLResource) Init(_ context.Context, _ *UserSAMLModel) (diags diag.Diagnostics) {
	var err error

	r.iamClient, err = r.client.CAVSDK.V1.IAM()
	if err != nil {
		diags.AddError("Error initializing IAM client", err.Error())
	}

	return diags
}

// Metadata returns the resource type name.
func (r *UserSAMLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_user_saml"
}

// Schema defines the schema for the resource.
func (r *UserSAMLResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = userSAMLSchema(ctx).GetResource(ctx)
}

func (r *UserSAMLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *UserSAMLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Create)()

	plan := &UserSAMLModel{}

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

	userCreated, err := r.iamClient.CreateSAMLUser(iam.SAMLUser{
		User: iam.User{
			Name:            plan.UserName.Get(),
			RoleName:        plan.RoleName.Get(),
			Enabled:         plan.Enabled.Get(),
			DeployedVMQuota: plan.DeployedVMQuota.GetInt(),
			StoredVMQuota:   plan.StoredVMQuota.GetInt(),
		},
	})
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.Diagnostics.AddError("User not found after create", fmt.Sprintf("User with name %s not found after create", plan.UserName.Get()))
			return
		}
		resp.Diagnostics.AddError("Error creating user SAML", err.Error())
		return
	}

	plan.ID.Set(userCreated.User.ID)

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("User not found after create", fmt.Sprintf("User with name %s not found after create", plan.UserName.Get()))
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
func (r *UserSAMLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Read)()

	state := &UserSAMLModel{}

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
func (r *UserSAMLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &UserSAMLModel{}
		state = &UserSAMLModel{}
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

	user, err := r.iamClient.GetUser(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	user.User.ID = plan.ID.Get()
	user.User.Name = plan.UserName.Get()
	user.User.RoleName = plan.RoleName.Get()
	user.User.Enabled = plan.Enabled.Get()
	user.User.DeployedVMQuota = plan.DeployedVMQuota.GetInt()
	user.User.StoredVMQuota = plan.StoredVMQuota.GetInt()

	if err := user.Update(); err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	// Special case to inject TakeOwnership value
	state.TakeOwnership.Set(plan.TakeOwnership.Get())

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("User not found", fmt.Sprintf("User with name %s not found", plan.UserName.Get()))
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
func (r *UserSAMLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Delete)()

	state := &UserSAMLModel{}

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

	user, err := r.iamClient.GetUser(state.ID.Get())
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

func (r *UserSAMLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Import)()

	userData := &UserSAMLModel{}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, userData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// req.ID is the user name
	user, err := r.iamClient.GetUser(req.ID)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.Diagnostics.AddError("User not found", fmt.Sprintf("User with name %s not found", req.ID))
			return
		}
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	if user.User.Type != iam.UserTypeSAML {
		resp.Diagnostics.AddError("User is not SAML user", fmt.Sprintf("User with name %s is not a SAML user.", req.ID))
		return
	}

	userData.ID.Set(user.User.ID)
	userData.UserName.Set(user.User.Name)

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
func (r *UserSAMLResource) read(_ context.Context, planOrState *UserSAMLModel) (stateRefreshed *UserSAMLModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	var (
		user *iam.UserClient
		err  error
	)

	// Get user by ID is more efficient
	if stateRefreshed.ID.IsKnown() {
		user, err = r.iamClient.GetUser(stateRefreshed.ID.Get())
	} else {
		user, err = r.iamClient.GetUser(stateRefreshed.UserName.Get())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving user", err.Error())
		return nil, true, diags
	}

	stateRefreshed.ID.Set(user.User.ID)
	stateRefreshed.UserName.Set(user.User.Name)
	stateRefreshed.RoleName.Set(user.User.RoleName)
	stateRefreshed.Enabled.Set(user.User.Enabled)
	stateRefreshed.DeployedVMQuota.SetInt(user.User.DeployedVMQuota)
	stateRefreshed.StoredVMQuota.SetInt(user.User.StoredVMQuota)

	return stateRefreshed, true, nil
}
