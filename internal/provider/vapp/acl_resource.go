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

// Package vapp provides a Terraform resource.
package vapp

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/acl"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &aclResource{}
	_ resource.ResourceWithConfigure   = &aclResource{}
	_ resource.ResourceWithImportState = &aclResource{}
)

// NewaclResource is a helper function to simplify the provider implementation.
func NewACLResource() resource.Resource {
	return &aclResource{}
}

// aclResource is the resource implementation.
type aclResource struct {
	client       *client.CloudAvenue
	vappClientV2 vAppClientV2
	vdc          resolvedVDC
	vapp         *sdktypes.ModelVApp
}

// sharedWithObjectType is the type of the shared_with set elements.
var sharedWithObjectType = types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes}

// Metadata returns the resource type name.
func (r *aclResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "acl"
}

// Schema defines the schema for the resource.
func (r *aclResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = aclSchema().GetResource(ctx)
}

// Init resolves the VDC and the parent vApp. The SDK client is built in Configure.
func (r *aclResource) Init(ctx context.Context, rm *aclResourceModel) (diags diag.Diagnostics) {
	vdc, diags := r.vappClientV2.resolveVDC(ctx, rm.VDC.ValueString())
	if diags.HasError() {
		return diags
	}
	r.vdc = vdc

	app, diags := r.vappClientV2.resolveVApp(ctx, r.vdc, rm.VAppID.ValueString(), rm.VAppName.ValueString())
	if diags.HasError() {
		return diags
	}
	r.vapp = app

	return diags
}

func (r *aclResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Expected *client.CloudAvenue, got %T. Report this to provider maintainers.", req.ProviderData),
		)

		return
	}

	r.client = client

	v2, diags := InitVAppClientV2(client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.vappClientV2 = v2
}

// Create creates the resource and sets the initial Terraform state.
func (r *aclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vapp_acl", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	var (
		plan *aclResourceModel
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

	// Create resource
	plan, diags := r.createOrUpdateACL(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *aclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vapp_acl", r.client.GetOrgName(), metrics.Read)()

	var state *aclResourceModel

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

	// Request acl
	vappACL, err := r.vappClientV2.vappClient.GetVAppACL(ctx, sdktypes.ParamsGetVAppACL{ID: r.vapp.ID})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, fmt.Sprintf("vApp %s(%s) access control", r.vapp.Name, r.vapp.ID), err)
		return
	}

	plan := &aclResourceModel{
		ID:                  types.StringValue(r.vapp.ID),
		VDC:                 types.StringValue(r.vdc.Name),
		VAppID:              state.VAppID,
		VAppName:            state.VAppName,
		EveryoneAccessLevel: types.StringNull(),
		SharedWith:          types.SetNull(sharedWithObjectType),
	}

	// SharedWithEveryone
	if vappACL.EveryoneAccessLevel != nil && *vappACL.EveryoneAccessLevel != "" {
		plan.EveryoneAccessLevel = types.StringValue(*vappACL.EveryoneAccessLevel)
	}

	// SharedWith
	if len(vappACL.SharedWith) > 0 {
		sharedWith, err := sharedWithFromSDK(vappACL.SharedWith)
		if err != nil {
			resp.Diagnostics.AddError("Error converting slice AccessSetting into set", err.Error())
			return
		}

		var diags diag.Diagnostics
		plan.SharedWith, diags = types.SetValueFrom(ctx, sharedWithObjectType, sharedWith)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Remove resource if all rights are null
	if plan.EveryoneAccessLevel.IsNull() && plan.SharedWith.IsNull() {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *aclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vapp_acl", r.client.GetOrgName(), metrics.Update)()

	var plan *aclResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update resource
	plan, diags := r.createOrUpdateACL(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *aclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vapp_acl", r.client.GetOrgName(), metrics.Delete)()

	var state *aclResourceModel

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
	unlock, d := lockVApp(ctx, r.vdc.Name, r.vapp.Name)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer unlock()

	// Delete vApp access control
	if _, err := r.vappClientV2.vappClient.RemoveVAppACL(ctx, sdktypes.ParamsRemoveVAppACL{ID: r.vapp.ID}); err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionDelete, fmt.Sprintf("access control of vApp %s(%s)", r.vapp.Name, r.vapp.ID), err)
		return
	}
}

func (r *aclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vapp_acl", r.client.GetOrgName(), metrics.Import)()

	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 1 && len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: vdc.vapp_name or vapp_name. Got: %q", req.ID),
		)
		return
	}

	if len(idParts) == 1 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), req.ID)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vdc"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), idParts[1])...)
	}
}

func (r *aclResource) createOrUpdateACL(ctx context.Context, plan *aclResourceModel) (*aclResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	sharedList := []acl.SharedWithModel{}
	diags.Append(plan.SharedWith.ElementsAs(ctx, &sharedList, true)...)
	if diags.HasError() {
		return nil, diags
	}

	// Lock vApp
	unlock, d := lockVApp(ctx, r.vdc.Name, r.vapp.Name)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	defer unlock()

	params := sdktypes.ParamsSetVAppACL{ID: r.vapp.ID}

	isSharedWithEveryone := !plan.EveryoneAccessLevel.IsNull() && !plan.EveryoneAccessLevel.IsUnknown()
	if isSharedWithEveryone {
		everyoneAccessLevel := plan.EveryoneAccessLevel.ValueString()
		params.IsSharedToEveryone = true
		params.EveryoneAccessLevel = &everyoneAccessLevel
	} else {
		subjects, err := sharedWithToSDK(sharedList)
		if err != nil {
			diags.AddError("Error when reading shared_with from schema.", err.Error())
			return nil, diags
		}
		params.SharedWith = subjects
	}

	vappACL, err := r.vappClientV2.vappClient.SetVAppACL(ctx, params)
	if err != nil {
		cerrs.AddError(&diags, cerrs.ActionUpdate, fmt.Sprintf("access control of vApp %s(%s)", r.vapp.Name, r.vapp.ID), err)
		return nil, diags
	}

	sharedWith := plan.SharedWith
	if !isSharedWithEveryone {
		// Use the subjects returned by the API (they carry the subject name)
		sharedWith = types.SetNull(sharedWithObjectType)
		if len(vappACL.SharedWith) > 0 {
			sharedWithModels, err := sharedWithFromSDK(vappACL.SharedWith)
			if err != nil {
				diags.AddError("Error converting slice AccessSetting into set", err.Error())
				return nil, diags
			}

			var d diag.Diagnostics
			sharedWith, d = types.SetValueFrom(ctx, sharedWithObjectType, sharedWithModels)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}
		}
	}

	return &aclResourceModel{
		ID:                  types.StringValue(r.vapp.ID),
		VDC:                 types.StringValue(r.vdc.Name),
		VAppID:              plan.VAppID,
		VAppName:            plan.VAppName,
		EveryoneAccessLevel: plan.EveryoneAccessLevel,
		SharedWith:          sharedWith,
	}, diags
}
