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

// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
)

const (
	defaultCheckJobDelayEdgeGateway = 10 * time.Second
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &edgeGatewayResource{}
	_ resource.ResourceWithConfigure   = &edgeGatewayResource{}
	_ resource.ResourceWithImportState = &edgeGatewayResource{}
	_ resource.ResourceWithModifyPlan  = &edgeGatewayResource{}

	// ConfigEdgeGateway is the default configuration for edge gateway.
	ConfigEdgeGateway setDefaultEdgeGateway = func() EdgeGatewayConfig {
		return EdgeGatewayConfig{
			CheckJobDelay: defaultCheckJobDelayEdgeGateway,
		}
	}
)

// NewEdgeGatewayResource returns a new resource implementing the edge_gateway data source.
func NewEdgeGatewayResource() resource.Resource {
	return &edgeGatewayResource{}
}

type setDefaultEdgeGateway func() EdgeGatewayConfig

// EdgeGatewayConfig is the configuration for edge gateway.
type EdgeGatewayConfig struct {
	CheckJobDelay time.Duration
}

// edgeGatewayResource is the resource implementation.
type edgeGatewayResource struct {
	client *client.CloudAvenue
	EdgeGatewayConfig
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

// loadRemainingBandwidth returns the remaining bandwidth capacity for the T0.
func (h *modifyPlanHelper) loadRemainingBandwidth() (int, error) {
	edgegws, err := h.r.client.CAVSDK.V1.EdgeGateway.List()
	if err != nil {
		return 0, err
	}

	return edgegws.GetBandwidthCapacityRemaining(h.plan.Tier0VRFName.Get())
}

// validateAllowedBandwidth checks that the plan bandwidth is an allowed value.
func (h *modifyPlanHelper) validateAllowedBandwidth() {
	allowedValues, err := h.r.client.CAVSDK.V1.EdgeGateway.GetAllowedBandwidthValues(h.plan.Tier0VRFName.Get())
	if err != nil {
		h.diag.AddError("Error on calculating allowed Bandwidth values", err.Error())
		return
	}

	if !slices.Contains(allowedValues, h.plan.Bandwidth.GetInt()) { //nolint:govet
		h.diag.AddError("Invalid Bandwidth value", fmt.Sprintf("Bandwidth value must be one of %v", allowedValues))
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
	if stateBandwidth > 0 && slices.Contains(allowedValues, stateBandwidth) { //nolint: govet
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
func (h *modifyPlanHelper) modifyPlanCreateKnown() bool {
	h.validateAllowedBandwidth()
	if h.diag.HasError() {
		return false
	}

	remaining, err := h.loadRemainingBandwidth()
	if err != nil {
		if errors.Is(err, v1.ErrDedicatedT0BandwidthNotComputable) {
			// API returns rateLimit=0 for dedicated T0s — capacity check is unreliable, skip it.
			// Bandwidth will be validated by GetAllowedBandwidthValues only.
			// See: https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1229
			return true
		}
		if errors.Is(err, v1.ErrNoBandwidthCapacityRemaining) {
			h.diag.AddError("Error on calculating remaining bandwidth", "Not enough bandwidth available")
			return false
		}
		h.diag.AddError("Error on calculating remaining bandwidth", err.Error())

		return false
	}

	if h.plan.Bandwidth.GetInt() > remaining {
		h.diag.AddAttributeError(path.Root("bandwidth"), "Overcommitting bandwidth", fmt.Sprintf("Not enough bandwidth available, requested: %dMbps, available: %dMbps", h.plan.Bandwidth.GetInt(), remaining))
	}

	return true
}

// modifyPlanCreateUnknown handles the create case when bandwidth is not yet known.
func (h *modifyPlanHelper) modifyPlanCreateUnknown() bool {
	remaining, err := h.loadRemainingBandwidth()

	allowedValues, errAllowed := h.r.client.CAVSDK.V1.EdgeGateway.GetAllowedBandwidthValues(h.plan.Tier0VRFName.Get())
	if errAllowed != nil {
		h.diag.AddError("Error on calculating allowed Bandwidth values", errAllowed.Error())
		return false
	}

	if err != nil {
		return h.modifyPlanCreateUnknownWithError(err, allowedValues)
	}

	remaining, err = bestValueAtMostOrError(remaining, allowedValues)
	if err != nil {
		h.diag.AddError("Error on calculating remaining bandwidth", err.Error())
		return false
	}
	h.diag.AddAttributeWarning(path.Root("bandwidth"), "Bandwidth value is unknown, will be set to remaining bandwidth", fmt.Sprintf("Bandwidth defined to %dMbps", remaining))
	h.plan.Bandwidth.SetInt(remaining)

	return true
}

// modifyPlanCreateUnknownWithError handles the dedicated-T0 and other error paths for the unknown-bandwidth create case.
func (h *modifyPlanHelper) modifyPlanCreateUnknownWithError(err error, allowedValues []int) bool {
	if errors.Is(err, v1.ErrDedicatedT0BandwidthNotComputable) {
		// API returns rateLimit=0 for dedicated T0s — use the max allowed value as best effort.
		// See: https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1229
		maxAllowed, ok := maxAllowedBandwidth(allowedValues)
		if !ok {
			h.diag.AddError("Error on calculating allowed Bandwidth values", "No allowed bandwidth values returned for dedicated T0")
			return false
		}

		h.plan.Bandwidth.SetInt(maxAllowed)
		h.diag.AddAttributeWarning(path.Root("bandwidth"), "Bandwidth value is unknown, will be set to max allowed value (dedicated T0)", fmt.Sprintf("Bandwidth defined to %dMbps", maxAllowed))

		return true
	}

	if errors.Is(err, v1.ErrNoBandwidthCapacityRemaining) {
		h.diag.AddError("Error on calculating remaining bandwidth", "Not enough bandwidth available")
		return false
	}

	h.diag.AddError("Error on calculating remaining bandwidth", err.Error())

	return false
}

// modifyPlanUpdate handles the bandwidth-update case.
func (h *modifyPlanHelper) modifyPlanUpdate(state *edgeGatewayResourceModel) bool {
	allowedValues, err := h.r.client.CAVSDK.V1.EdgeGateway.GetAllowedBandwidthValues(h.plan.Tier0VRFName.Get())
	if err != nil {
		h.diag.AddError("Error on calculating allowed Bandwidth values", err.Error())
		return false
	}

	remaining, remainingErr := h.loadRemainingBandwidth()
	if remainingErr != nil {
		if errors.Is(remainingErr, v1.ErrDedicatedT0BandwidthNotComputable) {
			// API returns rateLimit=0 for dedicated T0s — skip overcommit check.
			// See: https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1229
			selectedBandwidth, bwErr := dedicatedT0UpdateBandwidth(state.Bandwidth.GetInt(), allowedValues)
			if bwErr != nil {
				h.diag.AddError("Error on calculating allowed Bandwidth values", bwErr.Error())
				return false
			}

			h.plan.Bandwidth.SetInt(selectedBandwidth)
			h.validateAllowedBandwidth()
			return !h.diag.HasError()
		}
		if errors.Is(remainingErr, v1.ErrNoBandwidthCapacityRemaining) {
			h.diag.AddError("Error on calculating remaining bandwidth", "Not enough bandwidth available")
			return false
		}
		h.diag.AddError("Error on calculating remaining bandwidth", remainingErr.Error())

		return false
	}

	remainingOnUpdate := allowedBandwidthAtMostUpdate(remaining, state.Bandwidth.GetInt(), allowedValues)

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

	h.validateAllowedBandwidth()

	return true
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
		t0s, err := r.client.CAVSDK.V1.T0.GetT0s()
		if err != nil {
			resp.Diagnostics.AddError("Error listing T0s", err.Error())
			return
		}

		switch len(*t0s) {
		case 0:
			resp.Diagnostics.AddError("Error listing T0s", "No T0s found")
			return
		case 1:
			plan.Tier0VRFName.Set((*t0s)[0].GetName())
		default:
			resp.Diagnostics.AddError("Error listing T0s", "Multiple T0s found, please specify the T0 name")
			return
		}
	}

	h := &modifyPlanHelper{r: r, plan: plan, diag: &resp.Diagnostics}

	switch determineModifyPlanBandwidthAction(plan, state) {
	case modifyPlanBandwidthActionCreateKnown:
		// Create case: bandwidth value is already known.
		h.modifyPlanCreateKnown()
	case modifyPlanBandwidthActionCreateUnknown:
		// Create case: bandwidth value is not yet known (computed).
		h.modifyPlanCreateUnknown()
	case modifyPlanBandwidthActionUpdate:
		// Update case: bandwidth value changed.
		h.modifyPlanUpdate(state)
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

func (r *edgeGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.EdgeGatewayConfig = ConfigEdgeGateway()
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

	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, d := plan.Timeouts.Create(ctx, 8*time.Minute)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(ctx, createTimeout)
	defer cancel()

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// List all edge gateways for determining the ID of the new edge gateway
	edgegws, err := r.client.CAVSDK.V1.EdgeGateway.List()
	if err != nil {
		resp.Diagnostics.AddError("Error listing edge gateways", err.Error())
		return
	}

	var job *commoncloudavenue.JobStatus

	vdcOrVDCGroup, err := r.client.CAVSDK.V1.VDC().GetVDCOrVDCGroup(plan.OwnerName.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or VDC Group", err.Error())
		return
	}

	if vdcOrVDCGroup.IsVDCGroup() {
		job, err = r.client.CAVSDK.V1.EdgeGateway.NewFromVDCGroup(plan.OwnerName.Get(), plan.Tier0VRFName.Get())
	} else {
		job, err = r.client.CAVSDK.V1.EdgeGateway.New(plan.OwnerName.Get(), plan.Tier0VRFName.Get())
	}

	if err != nil {
		resp.Diagnostics.AddError("Error creating edge gateway", err.Error())
		return
	}
	if err := job.Wait(1, int(createTimeout.Seconds())); err != nil {
		resp.Diagnostics.AddError("Error waiting for edge gateway creation", err.Error())
		return
	}

	// Find the new edge gateway
	edgegwsRefreshed, err := r.client.CAVSDK.V1.EdgeGateway.List()
	if err != nil {
		resp.Diagnostics.AddError("Error listing edge gateways", err.Error())
		return
	}

	var edgegwNew v1.EdgeGatewayType

	// Find the new edge gateway in the list of all edge gateways and set the ID. New edge gateway is in the list refreshed but not in the old list.
	for _, edgegw := range *edgegwsRefreshed {
		var found bool
		for _, edgegwOld := range *edgegws {
			if edgegw.GetID() == edgegwOld.GetID() {
				found = true
			}
		}
		if !found {
			plan.ID.Set(urn.Normalize(urn.Gateway, edgegw.GetID()).String())
			plan.Name.Set(edgegw.GetName())
			edgegwNew = edgegw
			break
		}
	}

	if edgegwNew == (v1.EdgeGatewayType{}) {
		resp.Diagnostics.AddError("Error retrieving new edge gateway", "New edge gateway not found")
		return
	}

	if edgegwNew.GetBandwidth() != plan.Bandwidth.GetInt() {
		job, err = edgegwNew.UpdateBandwidth(plan.Bandwidth.GetInt())
		if err != nil {
			resp.Diagnostics.AddError("Error setting Bandwidth", err.Error())
			return
		}

		if job != nil {
			if err := job.Wait(1, int(createTimeout.Seconds())); err != nil {
				resp.Diagnostics.AddError("Error waiting for Bandwidth update", err.Error())
				return
			}
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

// Read refreshes the Terraform state with the latest data.
func (r *edgeGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Read)()

	state := &edgeGatewayResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read timeout
	readTimeout, err := state.Timeouts.Read(ctx, 8*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(ctx, readTimeout)
	defer cancel()

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

	// Update() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	updateTimeout, errTO := plan.Timeouts.Update(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	if !plan.Bandwidth.Equal(state.Bandwidth) {
		edgegw, err := r.client.CAVSDK.V1.EdgeGateway.Get(urn.ExtractUUID(plan.ID.Get()))
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving edge gateway", err.Error())
			return
		}

		job, err := edgegw.UpdateBandwidth(plan.Bandwidth.GetInt())
		if err != nil {
			resp.Diagnostics.AddError("Error setting Bandwidth", err.Error())
			return
		}

		if err := job.Wait(1, int(updateTimeout.Seconds())); err != nil {
			resp.Diagnostics.AddError("Error waiting for Bandwidth update", err.Error())
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

	deleteTimeout, errTO := state.Timeouts.Delete(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	edgegw, err := r.client.CAVSDK.V1.EdgeGateway.Get(urn.ExtractUUID(state.ID.Get()))
	if err != nil {
		if commoncloudavenue.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving edge gateway", err.Error())
		return
	}

	job, err := edgegw.Delete()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting edge gateway", err.Error())
		return
	}

	if err := job.Wait(1, int(deleteTimeout.Seconds())); err != nil {
		resp.Diagnostics.AddError("Error waiting for edge gateway deletion", err.Error())
		return
	}
}

func (r *edgeGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Import)()

	// Retrieve import Name and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// * Custom funcs.
func (r *edgeGatewayResource) read(_ context.Context, planOrState *edgeGatewayResourceModel) (stateRefreshed *edgeGatewayResourceModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	nameOrID := planOrState.ID.Get()
	if nameOrID == "" {
		nameOrID = planOrState.Name.Get()
	}

	edgegw, err := r.client.CAVSDK.V1.EdgeGateway.Get(nameOrID)
	if err != nil {
		if commoncloudavenue.IsNotFound(err) || govcd.IsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving edge gateway", err.Error())
		return nil, true, diags
	}

	if !planOrState.ID.IsKnown() {
		stateRefreshed.ID.Set(urn.Normalize(urn.Gateway, edgegw.GetID()).String())
	}

	stateRefreshed.Tier0VRFName.Set(edgegw.GetTier0VrfID())
	stateRefreshed.OwnerName.Set(edgegw.GetOwnerName())
	stateRefreshed.Description.Set(edgegw.GetDescription())
	stateRefreshed.Bandwidth.SetInt(edgegw.GetBandwidth())

	return stateRefreshed, true, nil
}
