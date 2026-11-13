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

// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	edge "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &edgeGatewayResource{}
	_ resource.ResourceWithConfigure   = &edgeGatewayResource{}
	_ resource.ResourceWithImportState = &edgeGatewayResource{}
	_ resource.ResourceWithModifyPlan  = &edgeGatewayResource{}
)

// NewEdgeGatewayResource returns a new resource implementing the edge_gateway data source.
func NewEdgeGatewayResource() resource.Resource {
	return &edgeGatewayResource{}
}

// edgeGatewayResource is the resource implementation.
type edgeGatewayResource struct {
	client     *client.CloudAvenue
	edgeClient *edge.Client
}

// modifyPlanHelper holds shared state for ModifyPlan sub-operations.
type modifyPlanHelper struct {
	r    *edgeGatewayResource
	plan *edgeGatewayResourceModel
	diag *diag.Diagnostics
}

type modifyPlanBandwidthAction int

const (
	modifyPlanBandwidthActionNone modifyPlanBandwidthAction = iota
	modifyPlanBandwidthActionCreateKnown
	modifyPlanBandwidthActionCreateUnknown
	modifyPlanBandwidthActionUpdate
)

func determineModifyPlanBandwidthAction(plan, state *edgeGatewayResourceModel) modifyPlanBandwidthAction {
	if state == nil {
		if plan.Bandwidth.IsKnown() {
			return modifyPlanBandwidthActionCreateKnown
		}
		return modifyPlanBandwidthActionCreateUnknown
	}

	planKnown := plan.Bandwidth.IsKnown()
	stateKnown := state.Bandwidth.IsKnown()

	switch {
	case planKnown && !stateKnown:
		return modifyPlanBandwidthActionCreateKnown

	case !planKnown && stateKnown:
		return modifyPlanBandwidthActionUpdate

	case !planKnown && !stateKnown:
		return modifyPlanBandwidthActionCreateUnknown

	case !plan.Bandwidth.Equal(state.Bandwidth):
		return modifyPlanBandwidthActionUpdate

	default:
		return modifyPlanBandwidthActionNone
	}
}

// validateAllowedBandwidth checks that the plan bandwidth is an allowed value.
func (h *modifyPlanHelper) validateAllowedBandwidth(t0 *types.ModelT0) {
	if !slices.Contains(t0.Bandwidth.AllowedBandwidthValues, h.plan.Bandwidth.GetInt()) {
		h.diag.AddError("Invalid Bandwidth value", fmt.Sprintf("Bandwidth value must be one of %v", t0.Bandwidth.AllowedBandwidthValues))
	}
}

// maxAllowedBandwidth returns highest bandwidth from allowedValues.
func maxAllowedBandwidth(allowedValues []int) (int, bool) {
	if len(allowedValues) == 0 {
		return 0, false
	}

	maxVal := allowedValues[0]
	for _, value := range allowedValues[1:] {
		if value > maxVal {
			maxVal = value
		}
	}

	return maxVal, true
}

// bestValueAtMost returns the largest value in allowedValues that is <= value.
func bestValueAtMost(value int, allowedValues []int) int {
	var best int
	for _, v := range allowedValues {
		if v <= value && v > best {
			best = v
		}
	}

	return best
}

func bestValueAtMostOrError(value int, allowedValues []int) (int, error) {
	best := bestValueAtMost(value, allowedValues)
	if best == 0 {
		return 0, fmt.Errorf("no allowed bandwidth value fits current available capacity")
	}

	return best, nil
}

func dedicatedT0UpdateBandwidth(stateBandwidth int, allowedValues []int) (int, error) {
	if stateBandwidth > 0 && slices.Contains(allowedValues, stateBandwidth) {
		return stateBandwidth, nil
	}

	maxAllowed, ok := maxAllowedBandwidth(allowedValues)
	if !ok {
		return 0, fmt.Errorf("no allowed bandwidth values returned for dedicated T0")
	}

	return maxAllowed, nil
}

// allowedBandwidthAtMostUpdate returns best allowed bandwidth for current availability on update.
func allowedBandwidthAtMostUpdate(remaining, currentBandwidth int, allowedValues []int) int {
	return bestValueAtMost(remaining+currentBandwidth, allowedValues)
}

// modifyPlanCreateKnown handles the create case when bandwidth is already known.
func (h *modifyPlanHelper) modifyPlanCreateKnown(ctx context.Context) bool {
	t0, err := h.r.edgeClient.GetT0(ctx, types.ParamsGetT0{T0Name: h.plan.Tier0VRFName.Get(), EdgegatewayID: h.plan.ID.Get(), EdgegatewayName: h.plan.Name.Get()})
	if err != nil {
		h.diag.AddError("Error on retrieving associated TierO", err.Error())
		return false
	}

	h.validateAllowedBandwidth(t0)
	if h.diag.HasError() {
		return false
	}

	if h.plan.Bandwidth.GetInt() > t0.Bandwidth.Remaining {
		h.diag.AddAttributeError(path.Root("bandwidth"), "Overcommitting bandwidth", fmt.Sprintf("Not enough bandwidth available, requested: %dMbps, available: %dMbps", h.plan.Bandwidth.GetInt(), t0.Bandwidth.Remaining))
	}

	return true
}

// modifyPlanCreateUnknown handles the create case when bandwidth is not yet known.
func (h *modifyPlanHelper) modifyPlanCreateUnknown(ctx context.Context) bool {
	t0, err := h.r.edgeClient.GetT0(ctx, types.ParamsGetT0{T0Name: h.plan.Tier0VRFName.Get(), EdgegatewayID: h.plan.ID.Get(), EdgegatewayName: h.plan.Name.Get()})
	if err != nil {
		h.diag.AddError("Error on retrieving associated TierO", err.Error())
		return false
	}

	remaining, err := bestValueAtMostOrError(t0.Bandwidth.Remaining, t0.Bandwidth.AllowedBandwidthValues)
	if err != nil {
		h.diag.AddError("Error on calculating remaining bandwidth", err.Error())
		return false
	}
	h.diag.AddAttributeWarning(path.Root("bandwidth"), "Bandwidth value is unknown, will be set to remaining bandwidth", fmt.Sprintf("Bandwidth defined to %dMbps", remaining))
	h.plan.Bandwidth.SetInt(remaining)

	return true
}

// modifyPlanUpdate handles the bandwidth-update case.
func (h *modifyPlanHelper) modifyPlanUpdate(ctx context.Context, state *edgeGatewayResourceModel) bool {
	t0, err := h.r.edgeClient.GetT0(ctx, types.ParamsGetT0{T0Name: h.plan.Tier0VRFName.Get(), EdgegatewayID: h.plan.ID.Get(), EdgegatewayName: h.plan.Name.Get()})
	if err != nil {
		h.diag.AddError("Error on retrieving associated Tier0", err.Error())
		return false
	}

	remainingOnUpdate := allowedBandwidthAtMostUpdate(t0.Bandwidth.Remaining, state.Bandwidth.GetInt(), t0.Bandwidth.AllowedBandwidthValues)

	if h.plan.Bandwidth.IsUnknown() {
		if remainingOnUpdate == 0 {
			h.diag.AddError("Error on calculating remaining bandwidth", "No allowed bandwidth value fits current available capacity")
			return false
		}

		h.plan.Bandwidth.SetInt(remainingOnUpdate)
		return true
	}

	if h.plan.Bandwidth.GetInt() > remainingOnUpdate {
		h.diag.AddAttributeError(path.Root("bandwidth"), "Overcommitting bandwidth", fmt.Sprintf("Not enough bandwidth available, requested: %dMbps, available: %dMbps", h.plan.Bandwidth.GetInt(), remainingOnUpdate))
		return false
	}

	h.validateAllowedBandwidth(t0)

	return true
}

// Init Initializes the resource.
func (r *edgeGatewayResource) Init(_ context.Context) (diags diag.Diagnostics) {
	var err error

	var v2 cav.Client
	if v2, err = r.client.V2(); err == nil {
		r.edgeClient, err = edge.New(v2)
	}
	if err != nil {
		diags.AddError("Error initializing EdgeGateway client", err.Error())
		return diags
	}

	return diags
}

// ModifyPlan modifies the plan to add the default values.
func (r *edgeGatewayResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var (
		plan  = &edgeGatewayResourceModel{}
		state = &edgeGatewayResourceModel{}
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the plan is nil, then this is a delete operation.
	if plan == nil {
		return
	}

	// If Tier0VRFName is not known, we need to find the T0.
	// If multiple T0s are available, return an error.
	if !plan.Tier0VRFName.IsKnown() {
		t0s, err := r.edgeClient.ListT0(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Error listing T0s", err.Error())
			return
		}
		switch len(t0s.T0s) {
		case 0:
			resp.Diagnostics.AddError("Error listing T0s", "No T0s found")
			return
		case 1:
			plan.Tier0VRFName.Set(t0s.T0s[0].Name)
		default:
			resp.Diagnostics.AddError("Error listing T0s", "Multiple T0s found, please specify the T0 name")
			return
		}
	}

	h := &modifyPlanHelper{r: r, plan: plan, diag: &resp.Diagnostics}

	switch determineModifyPlanBandwidthAction(plan, state) {
	case modifyPlanBandwidthActionCreateKnown:
		// Create case: bandwidth value is already known.
		h.modifyPlanCreateKnown(ctx)
	case modifyPlanBandwidthActionCreateUnknown:
		// Create case: bandwidth value is not yet known (computed).
		h.modifyPlanCreateUnknown(ctx)
	case modifyPlanBandwidthActionUpdate:
		// Update case: bandwidth value changed.
		h.modifyPlanUpdate(ctx, state)
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

// Metadata returns the resource type name.
func (r *edgeGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *edgeGatewayResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = edgegwSchema().GetResource(ctx)
}

func (r *edgeGatewayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = client

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *edgeGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Create)()

	plan := &edgeGatewayResourceModel{}

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	edgeGw, err := r.edgeClient.CreateEdgeGateway(ctx, types.ParamsCreateEdgeGateway{OwnerName: plan.OwnerName.Get(), T0Name: plan.Tier0VRFName.Get()})
	if err != nil {
		resp.Diagnostics.AddError("Error creating edge gateway", fmt.Sprintf("error creating edge gateway %s: %s", plan.Name.Get(), err.Error()))
		return
	}

	plan.ID.Set(edgeGw.ID)

	edgeBandwidth, err := r.edgeClient.GetBandwidth(ctx, types.ParamsEdgeGateway{ID: plan.ID.Get(), Name: plan.Name.Get()})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving edge gateway bandwidth", fmt.Sprintf("error retrieving edge gateway bandwidth %s: %s", plan.Name.Get(), err.Error()))
		return
	}

	if edgeBandwidth.Bandwidth != plan.Bandwidth.GetInt() {
		_, err = r.edgeClient.UpdateEdgeGateway(ctx, types.ParamsUpdateEdgeGateway{ID: edgeGw.ID, Bandwidth: plan.Bandwidth.GetInt()})
		if err != nil {
			resp.Diagnostics.AddError("Error updating bandwidth", fmt.Sprintf("error updating bandwidth for edge gateway %s during create: %s", plan.Name.Get(), err.Error()))
			return
		}
	}

	updatedEdgeBandwidth, err := r.edgeClient.GetBandwidth(ctx, types.ParamsEdgeGateway{ID: plan.ID.Get(), Name: plan.Name.Get()})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving edge gateway bandwidth", fmt.Sprintf("error retrieving edge gateway bandwidth %s: %s", plan.Name.Get(), err.Error()))
		return
	}

	plan.Bandwidth.Set(int64(updatedEdgeBandwidth.Bandwidth))

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The EdgeGateway %s(%s) was not found after creation.", plan.Name.Get(), plan.ID.Get()))
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *edgeGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Read)()

	state := &edgeGatewayResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *edgeGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Update)()

	plan := &edgeGatewayResourceModel{}
	state := &edgeGatewayResourceModel{}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	if !plan.Bandwidth.Equal(state.Bandwidth) {
		_, err := r.edgeClient.UpdateEdgeGateway(ctx, types.ParamsUpdateEdgeGateway{ID: plan.ID.Get(), Bandwidth: plan.Bandwidth.GetInt()})
		if err != nil {
			resp.Diagnostics.AddError("Error updating bandwidth", fmt.Sprintf("error updating bandwidth for edge gateway %s: %s", plan.Name.Get(), err.Error()))
			return
		}
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *edgeGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Delete)()

	state := &edgeGatewayResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	err := r.edgeClient.DeleteEdgeGateway(ctx, types.ParamsEdgeGateway{ID: state.ID.Get()})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting edge gateway", fmt.Sprintf("error deleting edge gateway %s: %s", state.Name.Get(), err.Error()))
		return
	}
}

func (r *edgeGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Import)()

	// Retrieve import Name and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// * Custom funcs.
func (r *edgeGatewayResource) read(ctx context.Context, planOrState *edgeGatewayResourceModel) (stateRefreshed *edgeGatewayResourceModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	edgegw, err := r.edgeClient.GetEdgeGateway(ctx, types.ParamsEdgeGateway{ID: planOrState.ID.Get(), Name: planOrState.Name.Get()})
	if err != nil {
		if cerrs.IsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving edge gateway", fmt.Sprintf("error retrieving edge gateway %s(%s) during read: %s", planOrState.Name.Get(), planOrState.ID.Get(), err.Error()))
		return nil, true, diags
	}

	edgeBandwidth, err := r.edgeClient.GetBandwidth(ctx, types.ParamsEdgeGateway{ID: planOrState.ID.Get(), Name: planOrState.Name.Get()})
	if err != nil {
		diags.AddError("Error retrieving edge gateway bandwidth", fmt.Sprintf("error retrieving edge gateway bandwidth %s: %s", planOrState.Name.Get(), err.Error()))
		return stateRefreshed, true, diags
	}

	if !planOrState.ID.IsKnown() {
		stateRefreshed.ID.Set(edgegw.ID)
	}

	stateRefreshed.Name.Set(edgegw.Name)
	if edgegw.UplinkT0 != nil && edgegw.UplinkT0.Name != "" {
		stateRefreshed.Tier0VRFName.Set(edgegw.UplinkT0.Name)
	} else if t0, err := r.edgeClient.GetT0(ctx, types.ParamsGetT0{EdgegatewayID: planOrState.ID.Get(), EdgegatewayName: planOrState.Name.Get()}); err == nil && t0 != nil && t0.Name != "" {
		stateRefreshed.Tier0VRFName.Set(t0.Name)
	}
	if edgegw.OwnerRef != nil {
		stateRefreshed.OwnerName.Set(edgegw.OwnerRef.Name)
	}
	stateRefreshed.Description.Set(edgegw.Description)
	stateRefreshed.Bandwidth.SetInt(edgeBandwidth.Bandwidth)

	return stateRefreshed, true, nil
}
