/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &OrgResource{}
	_ resource.ResourceWithConfigure   = &OrgResource{}
	_ resource.ResourceWithImportState = &OrgResource{}
)

// NewOrgResource is a helper function to simplify the provider implementation.
func NewOrgResource() resource.Resource {
	return &OrgResource{}
}

// OrgResource is the resource implementation.
type OrgResource struct { //nolint:revive
	client *client.CloudAvenue
	org    org.Client
}

// Init Initializes the resource.
func (r *OrgResource) Init(_ context.Context, _ *OrgModel) (diags diag.Diagnostics) {
	var err error

	r.org, err = org.NewClient()
	if err != nil {
		diags.AddError("Error creating org client", err.Error())
		return diags
	}
	return diags
}

// Metadata returns the resource type name.
func (r *OrgResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *OrgResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = orgSchema(ctx).GetResource(ctx)
}

func (r *OrgResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *OrgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_org", r.client.GetOrgName(), metrics.Create)()

	plan := &OrgModel{}

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

	resp.Diagnostics.AddError("Resource does not support creation", "The resource does not support creation. Import the resource instead.")
}

// Read refreshes the Terraform state with the latest data.
func (r *OrgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_org", r.client.GetOrgName(), metrics.Read)()

	state := &OrgModel{}

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
		resp.Diagnostics.AddError("Resource not found", "The resource was not found after refresh")
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
func (r *OrgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_org", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &OrgModel{}
		state = &OrgModel{}
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

	reqP := &org.PropertiesRequest{}

	if !state.Description.Equal(plan.Description) {
		reqP.Description = plan.Description.Get()
	}

	if !state.Email.Equal(plan.Email) {
		reqP.Email = plan.Email.Get()
	}

	if !state.InternetBillingModel.Equal(plan.InternetBillingModel) {
		reqP.BillingModel = plan.InternetBillingModel.Get()
	}

	if !state.Name.Equal(plan.Name) {
		reqP.FullName = plan.Name.Get()
	}

	job, err := r.org.UpdateProperties(ctx, reqP)
	if err != nil {
		resp.Diagnostics.AddError("Error updating properties", err.Error())
		return
	}

	// Wait for the job to complete
	jobStatus, err := job.GetJobStatus()
	if err != nil {
		resp.Diagnostics.AddError("Error getting job status", err.Error())
	}

	if err := jobStatus.WaitWithContext(ctx, 2); err != nil {
		resp.Diagnostics.AddError("Error waiting for job to complete", err.Error())
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", "The resource was not found after update")
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
func (r *OrgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_org", r.client.GetOrgName(), metrics.Delete)()

	state := &OrgModel{}

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

	resp.State.RemoveResource(ctx)
}

func (r *OrgResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_org", r.client.GetOrgName(), metrics.Import)()

	// No properties is needed for the import

	x := &OrgModel{
		ID: supertypes.NewStringNull(),
	}

	resp.Diagnostics.Append(r.Init(ctx, x)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, d := r.read(ctx, x)
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

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *OrgResource) read(ctx context.Context, planOrState *OrgModel) (stateRefreshed *OrgModel, found bool, diags diag.Diagnostics) { //nolint:unparam
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	properties, err := r.org.GetProperties(ctx)
	if err != nil {
		diags.AddError("Error getting properties", err.Error())
		// GetProperties never return not found error
		return nil, true, diags
	}

	stateRefreshed.ID.Set(r.client.GetOrgName())
	stateRefreshed.Name.Set(properties.FullName)
	stateRefreshed.Description.Set(properties.Description)
	stateRefreshed.Email.Set(properties.Email)
	stateRefreshed.InternetBillingModel.Set(properties.BillingModel)

	return stateRefreshed, true, nil
}
