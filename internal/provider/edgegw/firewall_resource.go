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

// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"
	"time"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	edge "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	vdcgroup "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdcgroup/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &firewallResource{}
	_ resource.ResourceWithConfigure   = &firewallResource{}
	_ resource.ResourceWithImportState = &firewallResource{}
)

// readTimeout is the timeout for reading the firewall resource. It is set to 30 seconds to allow for the backend to process the request and return the updated state.
const readTimeout = 30 * time.Second

// NewFirewallResource is a helper function to simplify the provider implementation.
func NewFirewallResource() resource.Resource {
	return &firewallResource{}
}

// firewallResource is the resource implementation.
type firewallResource struct {
	client     *client.CloudAvenue
	edgeClient *edge.Client
	edgegw     *types.ModelEdgeGateway
}

// Init Initializes the resource.
func (r *firewallResource) Init(ctx context.Context, rm *firewallModel) (diags diag.Diagnostics) {
	var err error
	var v2 cav.Client
	var vdcgroupClient *vdcgroup.Client

	if v2, err = r.client.V2(); err == nil {
		r.edgeClient, err = edge.New(v2)
	}
	if err != nil {
		diags.AddError("Error initializing EdgeGateway client", err.Error())
		return diags
	}
	vdcgroupClient, err = vdcgroup.New(v2)
	if err != nil {
		diags.AddError("Error initializing VdcGroup client", err.Error())
		return diags
	}
	r.edgegw, err = r.edgeClient.GetEdgeGateway(ctx, types.ParamsEdgeGateway{ID: rm.EdgeGatewayID.Get(), Name: rm.EdgeGatewayName.Get()})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return diags
	}

	_, err = vdcgroupClient.GetVDCGroup(ctx, types.ParamsGetVDCGroup{ID: r.edgegw.OwnerRef.ID, Name: r.edgegw.OwnerRef.Name})
	if err != nil {
		diags.AddError(
			"Edge Gateway belongs to a VDC Group",
			"Edge Gateway belongs to a VDC Group, please use the VDC Group resource instead `cloudavenue_vdcg_firewall`.",
		)
	}

	return diags
}

// Metadata returns the resource type name.
func (r *firewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_firewall"
}

// Schema defines the schema for the resource.
func (r *firewallResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = firewallSchema(ctx).GetResource(ctx)
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
		resp.Diagnostics.AddWarning("Resource not found", fmt.Sprintf("Unable to find firewall on edge %s", plan.EdgeGatewayName.Get()))
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

	edgeGateway, err := r.edgeClient.GetEdgeGateway(ctx, types.ParamsEdgeGateway{ID: state.EdgeGatewayID.Get(), Name: state.EdgeGatewayName.Get()})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	mutex.GlobalMutex.KvLock(ctx, edgeGateway.ID)
	defer mutex.GlobalMutex.KvUnlock(ctx, edgeGateway.ID)

	if err := r.edgeClient.DeleteFirewall(ctx, types.ParamsDeleteEdgeGatewayFirewall{EdgeGatewayID: state.EdgeGatewayID.Get(), EdgeGatewayName: state.EdgeGatewayName.Get()}); err != nil {
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
		d   diag.Diagnostics
		err error
	)

	state := &firewallModel{}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	edgeGateway, err := r.edgeClient.GetEdgeGateway(ctx, types.ParamsEdgeGateway{ID: req.ID})
	if err != nil {
		resp.Diagnostics.AddError("Failed to import firewall.", err.Error())
		return
	}

	state.EdgeGatewayID.Set(edgeGateway.ID)
	state.EdgeGatewayName.Set(edgeGateway.Name)

	// Refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.Diagnostics.AddError("Failed to import firewall.", fmt.Sprintf("Unable to find firewall on edge %s", edgeGateway.Name))
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
	// Lock object VDC or VDC Group
	edgeGateway, err := r.edgeClient.GetEdgeGateway(ctx, types.ParamsEdgeGateway{ID: plan.EdgeGatewayID.Get(), Name: plan.EdgeGatewayName.Get()})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway parent", err.Error())
		return diags
	}

	mutex.GlobalMutex.KvLock(ctx, edgeGateway.ID)
	defer mutex.GlobalMutex.KvUnlock(ctx, edgeGateway.ID)

	// Set the rules
	fwRules, d := plan.rulesToFirewallRule(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	_, err = r.edgeClient.UpdateFirewall(ctx, types.ParamsUpdateEdgeGatewayFirewall{EdgeGatewayID: plan.EdgeGatewayID.Get(), EdgeGatewayName: plan.EdgeGatewayName.Get(), Rules: fwRules})
	if err != nil {
		diags.AddError("Error to Update Firewall", err.Error())
		return diags
	}

	return diags
}

// read is a generic read function for the resource.
func (r *firewallResource) read(ctx context.Context, planOrState *firewallModel) (stateRefreshed *firewallModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	// When state already has rules, the backend can briefly return an empty list
	// right after apply. Retry a few times before accepting the empty response.
	previousRulesCount := 0
	if !planOrState.Rules.IsNull() && !planOrState.Rules.IsUnknown() {
		previousRules, d := planOrState.Rules.Get(ctx)
		diags.Append(d...)
		if d.HasError() {
			return stateRefreshed, true, diags
		}
		previousRulesCount = len(previousRules)
	}

	// Use GetFirewallExtended to also retrieve NetworkContextProfiles per rule.
	var (
		firewall *types.ModelGetFirewall
		err      error
	)

	err = retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		firewall, err = r.edgeClient.GetFirewall(ctx, types.ParamsGetEdgeGatewayFirewall{EdgeGatewayID: planOrState.EdgeGatewayID.Get(), EdgeGatewayName: planOrState.EdgeGatewayName.Get()})
		if err != nil {
			return retry.NonRetryableError(err)
		}

		if previousRulesCount > 0 && (firewall == nil || len(firewall.Rules) == 0) {
			return retry.RetryableError(fmt.Errorf("firewall rules not yet visible after apply"))
		}

		return nil
	})
	if err != nil {
		if cerrs.IsNotFound(err) {
			return stateRefreshed, false, nil
		}
		if previousRulesCount == 0 || (firewall != nil && len(firewall.Rules) != 0) {
			cerrs.AddError(&diags, cerrs.ActionRead, "Edge Gateway Firewall", err)
			return stateRefreshed, true, diags
		}
	}

	stateRefreshed.ID.Set(planOrState.ID.Get())
	stateRefreshed.EdgeGatewayID.Set(planOrState.EdgeGatewayID.Get())
	stateRefreshed.EdgeGatewayName.Set(planOrState.EdgeGatewayName.Get())

	rules := make([]*firewallModelRule, 0)

	if firewall == nil || len(firewall.Rules) == 0 {
		// Explicitly set an empty rules list so that out-of-band deletions
		// are reflected in state rather than leaving stale rules behind.
		diags.Append(stateRefreshed.Rules.Set(ctx, rules)...)
		return stateRefreshed, true, diags
	}

	for _, rule := range firewall.Rules {
		fwRule := &firewallModelRule{
			ID:                       supertypes.NewStringNull(),
			Name:                     supertypes.NewStringNull(),
			Enabled:                  supertypes.NewBoolNull(),
			Direction:                supertypes.NewStringNull(),
			IPProtocol:               supertypes.NewStringNull(),
			Action:                   supertypes.NewStringNull(),
			Logging:                  supertypes.NewBoolNull(),
			SourceIDs:                supertypes.NewSetValueOfNull[string](ctx),
			DestinationIDs:           supertypes.NewSetValueOfNull[string](ctx),
			AppPortProfileIDs:        supertypes.NewSetValueOfNull[string](ctx),
			NetworkContextProfileIDs: supertypes.NewSetValueOfNull[string](ctx),
		}
		fwRule.ID.Set(rule.ID)
		fwRule.Name.Set(rule.Name)
		fwRule.Enabled.Set(rule.Enabled)
		fwRule.Direction.Set(rule.Direction)
		fwRule.IPProtocol.Set(rule.IPProtocol)
		fwRule.Action.Set(rule.Action)
		fwRule.Logging.Set(rule.Logging)
		diags.Append(setIDsFromGroups(ctx, &fwRule.SourceIDs, rule.SourceFirewallGroups)...)
		diags.Append(setIDsFromGroups(ctx, &fwRule.DestinationIDs, rule.DestinationFirewallGroups)...)
		diags.Append(setIDsFromGroups(ctx, &fwRule.AppPortProfileIDs, rule.ApplicationPortProfiles)...)
		diags.Append(setIDsFromGroups(ctx, &fwRule.NetworkContextProfileIDs, rule.NetworkContextProfiles)...)
		rules = append(rules, fwRule)
	}

	diags.Append(stateRefreshed.Rules.Set(ctx, rules)...)
	return stateRefreshed, true, diags
}

// setIDsFromGroups sets the set attribute with the IDs of the given members.
// An empty list of members leaves the attribute null (equivalent to "any").
func setIDsFromGroups(ctx context.Context, attr *supertypes.SetValueOf[string], members []types.ModelGetFirewallGroupMember) diag.Diagnostics {
	if len(members) == 0 {
		*attr = supertypes.NewSetValueOfNull[string](ctx)
		return nil
	}

	ids := make([]string, 0, len(members))
	for _, m := range members {
		ids = append(ids, m.ID)
	}

	return attr.Set(ctx, ids)
}
