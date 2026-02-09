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

// Package iam provides a Terraform resource.
package iam

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &tokenResource{}
	_ resource.ResourceWithConfigure = &tokenResource{}
)

// NewTokenResource is a helper function to simplify the provider implementation.
func NewTokenResource() resource.Resource {
	return &tokenResource{}
}

// tokenResource is the resource implementation.
type tokenResource struct {
	client *client.CloudAvenue
	org    org.Org
}

// Init Initializes the resource.
func (r *tokenResource) Init(_ context.Context, _ *TokenModel) (diags diag.Diagnostics) {
	r.org, diags = org.Init(r.client)
	return diags
}

// Metadata returns the resource type name.
func (r *tokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_token"
}

// Schema defines the schema for the resource.
func (r *tokenResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = tokenSchema(ctx).GetResource(ctx)
}

func (r *tokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *tokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_iam_token", r.client.GetOrgName(), metrics.Create)()

	plan := &TokenModel{}

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

	token, err := r.client.Vmware.CreateToken(r.org.GetName(), plan.Name.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error creating token", err.Error())
		return
	}

	tokenString, err := token.GetInitialApiToken()
	if err != nil {
		resp.Diagnostics.AddError("Error getting token", err.Error())
		return
	}

	state := plan.Copy()
	if !plan.SaveInTfstate.Get() {
		// Token is computed attribute so we need to set it to empty string if SaveInTfstate is false
		state.Token.Set("")
	}

	if plan.PrintToken.Get() {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%s token generated", plan.Name.Get()), tokenString.AccessToken)
	}

	if plan.SaveInTfstate.Get() {
		state.Token.Set(tokenString.AccessToken)
	}

	if plan.SaveInFile.Get() {
		if err := govcd.SaveApiTokenToFile(plan.FileName.Get(), r.client.Vmware.Client.UserAgent, tokenString); err != nil {
			resp.Diagnostics.AddError("Error saving token", err.Error())
			return
		}
	}

	if !plan.PrintToken.Get() && !plan.SaveInTfstate.Get() && !plan.SaveInFile.Get() {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%s token generated", plan.Name.Get()), "Token not saved")
	}

	state.ID.Set(token.Token.ID)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *tokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_iam_token", r.client.GetOrgName(), metrics.Read)()

	state := &TokenModel{}

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
		Implement the resource read here
	*/

	token, err := r.client.Vmware.GetTokenById(state.ID.Get())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error getting token", err.Error())
		return
	}

	state.ID.Set(token.Token.ID)
	state.Name.Set(token.Token.Name)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tokenResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *tokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_iam_token", r.client.GetOrgName(), metrics.Delete)()

	state := &TokenModel{}

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

	token, err := r.client.Vmware.GetTokenById(state.ID.Get())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error getting token", err.Error())
		return
	}

	if err := token.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting token", err.Error())
		return
	}
}
