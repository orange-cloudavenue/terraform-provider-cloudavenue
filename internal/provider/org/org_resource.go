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

	"github.com/orange-cloudavenue/common-go/regex"
	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/organization/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
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

	// oClient is the Organization client from the SDK V2
	oClient *organization.Client
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

	// Get the provider client from the request data.
	client, ok := req.ProviderData.(*client.CloudAvenue)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client

	// Create the Organisation client from the SDK V2
	oC, err := organization.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Organization client, got error: %s", err),
		)
		return
	}
	r.oClient = oC
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

	// Refresh the state
	stateRefreshed, d := r.read(ctx)
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

	/*
		Implement the resource update here
	*/

	// Prepare the update request
	reqP := types.ParamsUpdateOrganization{}

	if !plan.FullName.IsNull() && !plan.FullName.Equal(state.FullName) {
		reqP.FullName = plan.FullName.Get()
	}
	if !plan.Description.Equal(state.Description) {
		reqP.Description = plan.Description.GetPtr()
	}
	if !plan.Email.IsNull() && !plan.Email.Equal(state.Email) {
		reqP.Email = plan.Email.Get()
	}
	if !plan.InternetBillingMode.IsNull() && !plan.InternetBillingMode.Equal(state.InternetBillingMode) {
		reqP.InternetBillingMode = plan.InternetBillingMode.Get()
	}

	// Update the organization
	orgUpdated, err := r.oClient.UpdateOrganization(ctx, reqP)
	if err != nil {
		resp.Diagnostics.AddError("Error updating organization properties", err.Error())
		return
	}

	// Refresh state after update
	stateRefreshed := state.Copy()
	diags := stateRefreshed.fromModel(ctx, orgUpdated)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
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

	/*
		Implement the resource deletion here
	*/

	resp.Diagnostics.AddWarning("Resource does not support delete", "The resource is not deletable. It will be removed from the state file but will still exist in Cloud Avenue.")
	resp.State.RemoveResource(ctx)
}

func (r *OrgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_org", r.client.GetOrgName(), metrics.Import)()

	// * ID (urn format) or Name

	// ID format is urn:vcloud:org:<org_uuid>
	// Name format is the organization name
	// If ID is provided, it will be used
	// If Name is provided, it will be used
	// if wrong format, return error

	// No properties is needed for the import, but we force to precise the name or id for clarity
	if urn.IsOrg(req.ID) {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	} else if regex.OrganizationNameRegex().MatchString(req.ID) {
		resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	} else {
		resp.Diagnostics.AddError("Invalid import identifier", fmt.Sprintf("The import identifier '%s' is not valid. It should be either the organization URN (urn:vcloud:org:<org_uuid>) or the organization name.", req.ID))
	}
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *OrgResource) read(ctx context.Context) (stateRefreshed *OrgModel, diags diag.Diagnostics) {
	// Get the organization details
	org, err := r.oClient.GetOrganization(ctx)
	if err != nil {
		diags.AddError("Error getting organization", err.Error())
		return nil, diags
	}

	// Set Refresh state with the latest data
	stateRefreshed = &OrgModel{}
	diags = stateRefreshed.fromModel(ctx, org)
	if diags.HasError() {
		return nil, diags
	}

	// Return the refreshed state
	return stateRefreshed, nil
}
