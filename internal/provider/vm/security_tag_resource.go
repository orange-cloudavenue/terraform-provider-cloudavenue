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

// Package vm provides a Terraform resource.
package vm

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &securityTagResource{}
	_ resource.ResourceWithConfigure   = &securityTagResource{}
	_ resource.ResourceWithImportState = &securityTagResource{}
)

// NewSecurityTagResource is a helper function to simplify the provider implementation.
func NewSecurityTagResource() resource.Resource {
	return &securityTagResource{}
}

// securityTagResource is the resource implementation.
type securityTagResource struct {
	client *client.CloudAvenue
	org    org.Org
}

// Metadata returns the resource type name.
func (r *securityTagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "security_tag"
}

// Init resource used to initialize the resource.
func (r *securityTagResource) Init(_ context.Context, _ *securityTagResourceModel) (diags diag.Diagnostics) {
	// Init Org
	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return diags
	}
	return diags
}

// Schema defines the schema for the resource.
func (r *securityTagResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = securityTagSchema().GetResource(ctx)
}

func (r *securityTagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *securityTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *securityTagResourceModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert List into []string
	var listvmids []string
	for _, vmid := range plan.VMIDs.Elements() {
		listvmids = append(listvmids, strings.ReplaceAll(vmid.String(), `"`, ``))
	}

	// Create the type SecurityTag
	securityTag := &govcdtypes.SecurityTag{
		Entities: listvmids,
		Tag:      plan.Name.ValueString(),
	}

	// Update the security tag in vCD
	_, err := r.org.UpdateSecurityTag(securityTag)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Tag", err.Error())
		return
	}

	plan = &securityTagResourceModel{
		Name:  plan.Name,
		VMIDs: plan.VMIDs,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *securityTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *securityTagResourceModel

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

	// Get all VM tagged in struct taggedEntities
	taggedEntities, err := r.org.GetAllSecurityTaggedEntitiesByName(state.Name.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// Tag not found, so remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to get Tagged Entities", err.Error())
		return
	}

	// Convert taggedEntities into []string
	readEntities := make([]string, len(taggedEntities))
	for i, entity := range taggedEntities {
		readEntities[i] = entity.ID
	}

	// Convert []string into List to refesh state

	VMIDs, d := types.SetValueFrom(ctx, types.StringType, readEntities)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := &securityTagResourceModel{
		Name:  state.Name,
		VMIDs: VMIDs,
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securityTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *securityTagResourceModel

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

	// Convert List into []string
	var listvmids []string
	for _, vmid := range plan.VMIDs.Elements() {
		listvmids = append(listvmids, strings.ReplaceAll(vmid.String(), `"`, ``))
	}

	// Create the type SecurityTag
	securityTag := &govcdtypes.SecurityTag{
		Entities: listvmids,
		Tag:      plan.Name.ValueString(),
	}

	// Update the security tag in vCD
	_, err := r.org.UpdateSecurityTag(securityTag)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Tag", err.Error())
		return
	}

	plan = &securityTagResourceModel{
		Name:  plan.Name,
		VMIDs: plan.VMIDs,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securityTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &securityTagResourceModel{}

	// Get current state and plan
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the type SecurityTag
	securityTag := &govcdtypes.SecurityTag{
		Entities: []string{},
		Tag:      state.Name.ValueString(),
	}

	// Update the security tag in vCD
	if _, err := r.org.UpdateSecurityTag(securityTag); err != nil {
		resp.Diagnostics.AddError("Unable to Delete Tag", err.Error())
		return
	}
}

func (r *securityTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all VM tagged in struct taggedEntities
	taggedEntities, err := r.org.GetAllSecurityTaggedEntitiesByName(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing security_tag", "name not found:"+err.Error())
		return
	}

	// Convert taggedEntities into []string
	readEntities := make([]string, 0)
	for _, entity := range taggedEntities {
		readEntities = append(readEntities, entity.ID)
	}

	VMIDs, d := types.SetValueFrom(ctx, types.StringType, readEntities)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := &securityTagResourceModel{
		Name:  types.StringValue(req.ID),
		VMIDs: VMIDs,
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
