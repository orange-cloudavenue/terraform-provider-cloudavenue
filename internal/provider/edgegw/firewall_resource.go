/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"
	"slices"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgegateway"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &firewallResource{}
	_ resource.ResourceWithConfigure   = &firewallResource{}
	_ resource.ResourceWithImportState = &firewallResource{}
	_ resource.ResourceWithModifyPlan  = &firewallResource{}
)

// NewFirewallResource is a helper function to simplify the provider implementation.
func NewFirewallResource() resource.Resource {
	return &firewallResource{}
}

// firewallResource is the resource implementation.
type firewallResource struct {
	client *client.CloudAvenue
	edgegw *edgegateway.EdgeGateway
}

// Init Initializes the resource.
func (r *firewallResource) Init(ctx context.Context, rm *firewallModel) (diags diag.Diagnostics) {
	var err error

	edgegw, err := edgegateway.NewClient()
	if err != nil {
		diags.AddError("Error creating Edge Gateway client", err.Error())
		return
	}

	nameOrID := rm.EdgeGatewayID.Get()
	if nameOrID == "" {
		nameOrID = rm.EdgeGatewayName.Get()
	}

	r.edgegw, err = edgegw.GetEdgeGateway(ctx, nameOrID)
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
	}

	return
}

// Metadata returns the resource type name.
func (r *firewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_firewall"
}

// Schema defines the schema for the resource.
func (r *firewallResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = firewallSchema(ctx).GetResource(ctx)
}

// ModifyPlan modifies the plan before it is applied.
func (r *firewallResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Get the current plan
	plan := &firewallModel{}
	if d := req.Plan.Get(ctx, plan); d.HasError() {
		// return because plan is empty
		return
	}

	// Apply the default values to the plan
	// The schema is not used here because a bug in the Terraform framework crashes when using `SetDefault` on a SET nested object.
	// https://github.com/hashicorp/terraform-plugin-framework/issues/783

	rules, d := plan.Rules.Get(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	for _, rule := range rules {
		// Priority attribute
		if rule.Priority.IsUnknown() {
			rule.Priority.SetInt64(1)
		}

		// IPProtocol attribute
		if rule.IPProtocol.IsUnknown() {
			rule.IPProtocol.Set("IPV4")
		}

		// Enabled attribute
		if rule.Enabled.IsUnknown() {
			rule.Enabled.Set(true)
		}

		// Logging attribute
		if rule.Logging.IsUnknown() {
			rule.Logging.Set(false)
		}
	}

	plan.Rules.DiagsSet(ctx, resp.Diagnostics, rules)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the modified plan back to the response
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	return
}

func (r *firewallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *firewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { //nolint:dupl
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Create)()

	plan := &firewallModel{}

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

	// Create or update the resource
	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning("Resource not found", fmt.Sprintf("Unable to find firewall on edgegateway %s", plan.EdgeGatewayName.Get()))
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *firewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = new(firewallModel)
		state = new(firewallModel)
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	// Identify the rules deleted from the state
	rulesPlan, d := plan.Rules.Get(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	rulesState, d := state.Rules.Get(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Find rules that are in the state but not in the plan (= deleted)
	deletedRules := make([]string, 0)
	for _, ruleState := range rulesState {
		if !slices.ContainsFunc(rulesPlan, func(rulePlan *firewallModelRule) bool {
			return rulePlan.ID.Get() == ruleState.ID.Get()
		}) {
			deletedRules = append(deletedRules, ruleState.ID.Get())
		}
	}
	if len(deletedRules) > 0 {
		if err := r.edgegw.DeleteFirewallRules(ctx, deletedRules); err != nil {
			resp.Diagnostics.AddError("Error deleting Edge Gateway Firewall rules", err.Error())
			return
		}
	}

	// Use generic createOrUpdate function to update the resource
	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *firewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Delete)()

	state := &firewallModel{}

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

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.ID)
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.ID)

	rules, d := state.Rules.Get(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	rulesID := make([]string, len(rules))
	for i, rule := range rules {
		rulesID[i] = rule.ID.Get()
	}

	if err := r.edgegw.DeleteFirewallRules(ctx, rulesID); err != nil {
		resp.Diagnostics.AddError("Error deleting Edge Gateway Firewall", err.Error())
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *firewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Read)()

	state := &firewallModel{}

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

func (r *firewallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Import)()

	var (
		edgegwID   string
		edgegwName string
	)

	if urn.IsValid(req.ID) {
		edgegwID = urn.Normalize(urn.Gateway, req.ID).String()
	} else {
		edgegwName = req.ID
	}

	state := &firewallModel{}
	state.EdgeGatewayID.Set(edgegwID)
	state.EdgeGatewayName.Set(edgegwName)

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID.Set(r.edgegw.ID)

	// Refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.Diagnostics.AddError("Failed to import firewall.", fmt.Sprintf("Unable to find firewall on edgegateway %s", r.edgegw.Name))
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// * custom functions

// createOrUpdate creates or updates the resource and sets the Terraform state.
func (r *firewallResource) createOrUpdate(ctx context.Context, plan *firewallModel) (diags diag.Diagnostics) {
	// Set the rules
	fwRules, d := plan.ToSDK(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.ID)
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.ID)

	if err := r.edgegw.UpdateFirewallRules(ctx, fwRules); err != nil {
		diags.AddError("Error creating or updating Edge Gateway Firewall rules", err.Error())
		return diags
	}

	return
}

// read is a generic read function for the resource.
func (r *firewallResource) read(ctx context.Context, planOrState *firewallModel) (stateRefreshed *firewallModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	fwRules, err := r.edgegw.GetFirewallRules(ctx)
	if err != nil {
		if govcd.IsNotFound(err) {
			return stateRefreshed, false, nil
		}
		diags.AddError("Error retrieving Edge Gateway Firewall", err.Error())
		return stateRefreshed, true, diags
	}

	// ID is a generated URN used to identify the resource in Terraform state
	stateRefreshed.ID.Set(r.edgegw.ID)

	stateRefreshed.EdgeGatewayID.Set(r.edgegw.ID)
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.Name)

	rules := make([]*firewallModelRule, 0)

	if fwRules.Rules == nil {
		return stateRefreshed, true, nil
	}

	existingRules, d := planOrState.Rules.Get(ctx)
	if d.HasError() {
		diags.Append(d...)
		return stateRefreshed, true, diags
	}

	for _, rule := range fwRules.Rules {
		// Only rules already in the plan or state should be processed
		// This is to avoid conflicts with the other rules managed by `cloudavenue_edgegateway_firewall_rule` or manually created rules.
		for _, existingRule := range existingRules {
			if (existingRule.ID.Get() == rule.ID) || (edgegateway.FirewallHashRule(existingRule.Name.Get(), existingRule.Action.Get(), existingRule.Direction.Get()) == rule.Hash) {
				fwRule := &firewallModelRule{
					ID:                     supertypes.NewStringNull(),
					Name:                   supertypes.NewStringNull(),
					Enabled:                supertypes.NewBoolNull(),
					Direction:              supertypes.NewStringNull(),
					IPProtocol:             supertypes.NewStringNull(),
					Action:                 supertypes.NewStringNull(),
					Logging:                supertypes.NewBoolNull(),
					Priority:               supertypes.NewInt64Null(),
					SourceIDs:              supertypes.NewSetValueOfNull[string](ctx),
					SourceIPAddresses:      supertypes.NewSetValueOfNull[string](ctx),
					DestinationIDs:         supertypes.NewSetValueOfNull[string](ctx),
					DestinationIPAddresses: supertypes.NewSetValueOfNull[string](ctx),
					AppPortProfileIDs:      supertypes.NewSetValueOfNull[string](ctx),
				}

				fwRule.ID.Set(rule.ID)
				fwRule.Name.Set(rule.Name)
				fwRule.Enabled.Set(rule.Enabled)
				fwRule.Direction.Set(rule.Direction)
				fwRule.IPProtocol.Set(rule.IPProtocol)
				fwRule.Action.Set(rule.Action)
				fwRule.Logging.Set(rule.Logging)
				fwRule.Priority.SetIntPtr(rule.Priority)
				fwRule.AppPortProfileIDs.DiagsSet(ctx, diags, common.FromOpenAPIReferenceID(ctx, rule.ApplicationPortProfiles))
				// * Sources
				fwRule.SourceIDs.DiagsSet(ctx, diags, common.FromOpenAPIReferenceID(ctx, rule.SourceFirewallGroups))
				fwRule.SourceIPAddresses.DiagsSet(ctx, diags, rule.SourceIPAddresses)
				// * Destinations
				fwRule.DestinationIDs.DiagsSet(ctx, diags, common.FromOpenAPIReferenceID(ctx, rule.DestinationFirewallGroups))
				fwRule.DestinationIPAddresses.DiagsSet(ctx, diags, rule.DestinationIPAddresses)
				if diags.HasError() {
					return stateRefreshed, true, diags
				}
				rules = append(rules, fwRule)
			}
		}
	}

	if diags.HasError() {
		return stateRefreshed, true, diags
	}

	diags.Append(stateRefreshed.Rules.Set(ctx, rules)...)
	return stateRefreshed, true, diags
}
