/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &FirewallResource{}
	_ resource.ResourceWithConfigure   = &FirewallResource{}
	_ resource.ResourceWithImportState = &FirewallResource{}
)

// NewFirewallResource is a helper function to simplify the provider implementation.
func NewFirewallResource() resource.Resource {
	return &FirewallResource{}
}

// FirewallResource is the resource implementation.
type FirewallResource struct {
	client   *client.CloudAvenue
	vdcGroup *v1.VDCGroup
}

// Init Initializes the resource.
func (r *FirewallResource) Init(_ context.Context, rm *FirewallModel) (diags diag.Diagnostics) {
	var err error

	idOrName := rm.VDCGroupName.Get()
	if rm.VDCGroupID.IsKnown() {
		// Use the ID
		idOrName = rm.VDCGroupID.Get()
	}

	r.vdcGroup, err = r.client.CAVSDK.V1.VDC().GetVDCGroup(idOrName)
	if err != nil {
		diags.AddError("Error retrieving VDC Group", err.Error())
		return
	}
	return
}

// Metadata returns the resource type name.
func (r *FirewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_firewall"
}

// Schema defines the schema for the resource.
func (r *FirewallResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = firewallSchema(ctx).GetResource(ctx)
}

func (r *FirewallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *FirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdcg_firewall", r.client.GetOrgName(), metrics.Create)()

	plan := &FirewallModel{}

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

	mutex.GlobalMutex.KvLock(ctx, r.vdcGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.vdcGroup.GetID())

	rules, d := plan.rulesToSDKRules(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	_, err := r.vdcGroup.CreateFirewall(rules)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VDC Group Firewall", err.Error())
		return
	}

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Fail to retrieve VDC Group Firewall after creation", "VDC Group Firewall not found")
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *FirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdcg_firewall", r.client.GetOrgName(), metrics.Read)()

	state := &FirewallModel{}

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
func (r *FirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdcg_firewall", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &FirewallModel{}
		state = &FirewallModel{}
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

	mutex.GlobalMutex.KvLock(ctx, r.vdcGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.vdcGroup.GetID())

	vdcgfw, err := r.vdcGroup.GetFirewall()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC Group Firewall", err.Error())
		return
	}

	rules, d := plan.rulesToSDKRules(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if err := vdcgfw.UpdateFirewall(rules); err != nil {
		resp.Diagnostics.AddError("Error updating VDC Group Firewall rules", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Fail to retrieve VDC Group Firewall after update", "VDC Group Firewall not found")
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
func (r *FirewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdcg_firewall", r.client.GetOrgName(), metrics.Delete)()

	state := &FirewallModel{}

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

	mutex.GlobalMutex.KvLock(ctx, r.vdcGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.vdcGroup.GetID())

	vdcgfw, err := r.vdcGroup.GetFirewall()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC Group Firewall", err.Error())
		return
	}

	if err := vdcgfw.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting VDC Group Firewall", err.Error())
		return
	}
}

func (r *FirewallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdcg_firewall", r.client.GetOrgName(), metrics.Import)()

	state := &FirewallModel{
		VDCGroupID:   supertypes.NewStringNull(),
		VDCGroupName: supertypes.NewStringNull(),
		Enabled:      supertypes.NewBoolNull(),
		Rules:        supertypes.NewListNestedObjectValueOfNull[FirewallModelRule](ctx),
	}

	if urn.IsVDCGroup(req.ID) {
		state.VDCGroupID.Set(req.ID)
	} else {
		state.VDCGroupName.Set(req.ID)
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.Diagnostics.AddError("Fail to retrieve VDC Group Firewall after import", "VDC Group Firewall not found")
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *FirewallResource) read(ctx context.Context, planOrState *FirewallModel) (stateRefreshed *FirewallModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	vdcgfw, err := r.vdcGroup.GetFirewall()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving VDC Group Firewall", err.Error())
		return
	}

	if !stateRefreshed.ID.IsKnown() {
		// firewall don't have an ID, use the VDC Group ID instead
		stateRefreshed.ID.Set(r.vdcGroup.GetID())
	}

	// * Enabled
	isEnabled, err := vdcgfw.IsEnabled()
	if err != nil {
		diags.AddError("Error retrieving VDC Group Firewall enabled status", err.Error())
		return
	}
	stateRefreshed.Enabled.Set(isEnabled)

	stateRefreshed.VDCGroupID.Set(r.vdcGroup.GetID())
	stateRefreshed.VDCGroupName.Set(r.vdcGroup.GetName())

	// * Rules
	rules, d := stateRefreshed.sdkRulesToRules(ctx, vdcgfw.GetRules())
	if d.HasError() {
		diags.Append(d...)
		return stateRefreshed, true, diags
	}

	diags.Append(stateRefreshed.Rules.Set(ctx, rules)...)

	return stateRefreshed, true, diags
}
