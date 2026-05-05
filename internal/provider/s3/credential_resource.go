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

package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource               = &CredentialResource{}
	_ resource.ResourceWithConfigure  = &CredentialResource{}
	_ resource.ResourceWithModifyPlan = &CredentialResource{}
)

// NewCredentialResource is a helper function to simplify the provider implementation.
func NewCredentialResource() resource.Resource {
	return &CredentialResource{}
}

// CredentialResource is the resource implementation.
type CredentialResource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the resource.
func (r *CredentialResource) Init(_ context.Context, _ *CredentialModel) (diags diag.Diagnostics) {
	r.s3Client = r.client.CAVSDK.V1.S3()
	return diags
}

// Metadata returns the resource type name.
func (r *CredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_credential"
}

// Schema defines the schema for the resource.
func (r *CredentialResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = credentialSchema(ctx).GetResource(ctx)
}

func (r *CredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CredentialResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	plan := &CredentialModel{}
	state := &CredentialModel{}

	// Retrieve values from plan
	d := req.Plan.Get(ctx, plan)
	if d.HasError() {
		// If there is an error in the plan, we don't need to continue
		return
	}

	d = req.State.Get(ctx, state)
	// If error in state will be is in create mode
	if !d.HasError() {
		return
	}

	// If save_in_tfstate is true print warning security risk
	if plan.SaveInTFState.Get() {
		resp.Diagnostics.AddWarning(
			"save_in_tfstate is true",
			"SaveInTFState is true. This is a security risk and should only be used for testing purposes.",
		)
	}

	// if print_token is true print warning security risk
	if plan.PrintToken.Get() {
		resp.Diagnostics.AddWarning(
			"print_token is true",
			"PrintToken is true. This is a security risk and should only be used for testing purposes.",
		)
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *CredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_s3_credential", r.client.GetOrgName(), metrics.Create)()

	plan := &CredentialModel{}

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

	user, oseErr := r.s3Client.GetUser(r.client.GetUserName())
	if oseErr != nil {
		if oseErr.IsNotFountError() {
			resp.Diagnostics.AddError("User not found", fmt.Sprintf("The user %s is not found", r.client.GetUserName()))
			return
		}
		resp.Diagnostics.AddError("Error getting user", oseErr.Error())
		return
	}

	cred, err := user.NewCredential()
	if err != nil {
		resp.Diagnostics.AddError("Error creating credential", err.Error())
		return
	}

	plan.AccessKey.Set(cred.GetAccessKey())
	plan.Username.Set(user.GetName())

	if !plan.SaveInTFState.Get() {
		// Token is computed attribute so we need to set it to empty string if SaveInTfstate is false
		plan.SecretKey.Set("")
	}

	if plan.PrintToken.Get() {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("New credential created for the user %s ", plan.Username.Get()),
			fmt.Sprintf("Access key: %s\nSecret Key: %s", cred.GetAccessKey(), cred.GetSecretKey()),
		)
	}

	if plan.SaveInTFState.Get() {
		plan.SecretKey.Set(cred.GetSecretKey())
	}

	if plan.SaveInFile.Get() {
		type credentialFile struct {
			AK string `json:"accesskey"`
			SK string `json:"secretkey"`
		}

		credFile := credentialFile{
			AK: cred.GetAccessKey(),
			SK: cred.GetSecretKey(),
		}

		b, err := json.Marshal(credFile)
		if err != nil {
			resp.Diagnostics.AddError("Error marshalling credential", err.Error())
			return
		}

		if err := os.WriteFile(plan.FileName.Get(), b, 0o600); err != nil {
			resp.Diagnostics.AddError("Error saving credential", err.Error())
			return
		}
	}

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
func (r *CredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_s3_credential", r.client.GetOrgName(), metrics.Read)()

	state := &CredentialModel{}

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
func (r *CredentialResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_s3_credential", r.client.GetOrgName(), metrics.Update)()
	// No update for this resource
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *CredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_s3_credential", r.client.GetOrgName(), metrics.Delete)()

	state := &CredentialModel{}

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

	user, oseErr := r.s3Client.GetUser(state.Username.Get())
	if oseErr != nil {
		resp.Diagnostics.AddError("Error getting user", oseErr.Error())
		return
	}

	cred, err := user.GetCredential(state.AccessKey.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error getting credential", err.Error())
		return
	}

	if err := cred.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting credential", err.Error())
		return
	}
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *CredentialResource) read(_ context.Context, planOrState *CredentialModel) (stateRefreshed *CredentialModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	user, err := r.s3Client.GetUser(stateRefreshed.Username.Get())
	if err != nil {
		diags.AddError("Error getting user", err.Error())
		return stateRefreshed, found, diags
	}

	if _, err := user.GetCredential(stateRefreshed.AccessKey.Get()); err != nil {
		diags.AddError("Error getting credential", err.Error())
		return stateRefreshed, found, diags
	}

	// ID is a username and 4 first characters of the access key. (e.g. `username-1234`)
	stateRefreshed.ID.Set(fmt.Sprintf("%s-%s", stateRefreshed.Username.Get(), stateRefreshed.AccessKey.Get()[:4]))

	return stateRefreshed, true, nil
}
