/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vapp provides a Terraform resource.
package vapp

import (
	"context"
	"fmt"
	"strings"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/acl"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
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
	client *client.CloudAvenue
	vdc    vdc.VDC
	vapp   vapp.VAPP
}

// Metadata returns the resource type name.
func (r *aclResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "acl"
}

// Schema defines the schema for the resource.
func (r *aclResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = aclSchema().GetResource(ctx)
}

func (r *aclResource) Init(_ context.Context, rm *aclResourceModel) (diags diag.Diagnostics) {
	r.vdc, diags = vdc.Init(r.client, rm.VDC)
	if diags.HasError() {
		return
	}

	r.vapp, diags = vapp.Init(r.client, r.vdc, rm.VAppID, rm.VAppName)

	return
}

func (r *aclResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	if diags.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
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
	accessControl, err := r.vapp.GetAccessControl(false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving access control properties", err.Error())
		return
	}

	// SharedWithEveryone
	everyoneAccessLevel := ""
	if accessControl.EveryoneAccessLevel != nil {
		everyoneAccessLevel = *accessControl.EveryoneAccessLevel
	}

	plan := &aclResourceModel{
		ID:                  types.StringValue(r.vapp.GetID()),
		VDC:                 types.StringValue(r.vdc.GetName()),
		VAppID:              state.VAppID,
		VAppName:            state.VAppName,
		EveryoneAccessLevel: types.StringValue(everyoneAccessLevel),
	}

	if plan.EveryoneAccessLevel.ValueString() == "" {
		plan.EveryoneAccessLevel = types.StringNull()
	}

	plan.SharedWith = types.SetNull(types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes})
	if accessControl.AccessSettings != nil {
		accessControlListSet, err := acl.AccessControlListToSharedSet(accessControl.AccessSettings.AccessSetting)
		if err != nil {
			resp.Diagnostics.AddError("Error converting slice AccessSetting into set", err.Error())
			return
		}
		var diags diag.Diagnostics
		plan.SharedWith, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes}, accessControlListSet)
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
	if resp.Diagnostics.HasError() {
		return
	}
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
	if diags.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
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
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	// Delete vApp access control
	if err := r.vapp.RemoveAccessControl(false); err != nil {
		resp.Diagnostics.AddError("Error removing vApp", err.Error())
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
	var accessControl govcdtypes.ControlAccessParams

	everyoneAccessLevel := plan.EveryoneAccessLevel.ValueString()

	sharedList := []acl.SharedWithModel{}
	diags.Append(plan.SharedWith.ElementsAs(ctx, &sharedList, true)...)
	if diags.HasError() {
		return nil, diags
	}

	// Lock vApp
	diags.Append(r.vapp.LockVAPP(ctx)...)
	if diags.HasError() {
		return nil, diags
	}
	defer r.vapp.UnlockVAPP(ctx)

	var accessSettings []*govcdtypes.AccessSetting

	isSharedWithEveryone := !plan.EveryoneAccessLevel.IsNull() && !plan.EveryoneAccessLevel.IsUnknown()
	if isSharedWithEveryone {
		accessControl.IsSharedToEveryone = true
		accessControl.EveryoneAccessLevel = &everyoneAccessLevel
	} else {
		var sharedListOutput []*acl.SharedWithModel

		// Get admin Org
		adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
		if err != nil {
			diags.AddError("Error retrieving Org", err.Error())
			return nil, diags
		}

		accessSettings, sharedListOutput, err = acl.SharedSetToAccessControl(r.client.Vmware, adminOrg, sharedList)
		if err != nil {
			diags.AddError("Error when reading shared_with from schema.", err.Error())
			return nil, diags
		}

		plan.SharedWith = types.SetNull(types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes})
		if len(sharedListOutput) != 0 {
			accessControl.AccessSettings = &govcdtypes.AccessSettingList{
				AccessSetting: accessSettings,
			}
			plan.SharedWith, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes}, sharedListOutput)

			diags.Append(diags...)
			if diags.HasError() {
				return nil, diags
			}
		}
	}

	err := r.vapp.SetAccessControl(&accessControl, false)
	if err != nil {
		diags.AddError("Error setting access control for vApp", err.Error())
		return nil, diags
	}

	plan = &aclResourceModel{
		ID:                  types.StringValue(r.vapp.GetID()),
		VDC:                 types.StringValue(r.vdc.GetName()),
		VAppID:              plan.VAppID,
		VAppName:            plan.VAppName,
		EveryoneAccessLevel: plan.EveryoneAccessLevel,
		SharedWith:          plan.SharedWith,
	}

	return plan, diags
}
