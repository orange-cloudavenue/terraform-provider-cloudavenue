// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
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

	loadRemainingBandwidth := func() (int, error) {
		edgegws, err := r.client.CAVSDK.V1.EdgeGateway.List()
		if err != nil {
			return 0, err
		}

		return edgegws.GetBandwidthCapacityRemaining(plan.Tier0VrfID.Get())
	}

	allowedValuesFunc := func() {
		allowedValues, err := r.client.CAVSDK.V1.EdgeGateway.GetAllowedBandwidthValues(plan.Tier0VrfID.Get())
		if err != nil {
			resp.Diagnostics.AddError("Error on calculating allowed Bandwidth values", err.Error())
			return
		}

		if !slices.Contains(allowedValues, plan.Bandwidth.GetInt()) {
			resp.Diagnostics.AddError("Invalid Bandwidth value", fmt.Sprintf("Bandwidth value must be one of %v", allowedValues))
			return
		}
	}

	// determine la valeur autorisé la plus proche de la valeur demandé
	calculBestValue := func(value int, allowedValues []int) int {
		var bestValue int
		for _, v := range allowedValues {
			if v <= value && v > bestValue {
				bestValue = v
			}
		}
		return bestValue
	}

	// If the plan is nil, then this is a delete operation.
	if plan == nil {
		return
	}

	switch {
	// Create case with value is known
	case plan.Bandwidth.IsKnown() && (state == nil || !state.Bandwidth.IsKnown()):
		allowedValuesFunc()
		remaining, err := loadRemainingBandwidth()
		if err != nil {
			if errors.Is(err, fmt.Errorf("no bandwidth capacity remaining")) {
				resp.Diagnostics.AddError("Error on calculating remaining bandwidth", "Not enough bandwidth available")
				return
			}
			resp.Diagnostics.AddError("Error on calculating remaining bandwidth", err.Error())
			return
		}

		if plan.Bandwidth.GetInt() > remaining {
			resp.Diagnostics.AddAttributeError(path.Root("bandwidth"), "Overcommitting bandwidth", fmt.Sprintf("Not enough bandwidth available, requested: %dMbps, available: %dMbps", plan.Bandwidth.GetInt(), remaining))
		}
		goto END

	// Create case with value is unknown
	case !plan.Bandwidth.IsKnown():
		remaining, err := loadRemainingBandwidth()
		if err != nil {
			if errors.Is(err, fmt.Errorf("no bandwidth capacity remaining")) {
				resp.Diagnostics.AddError("Error on calculating remaining bandwidth", "Not enough bandwidth available")
				return
			}
			resp.Diagnostics.AddError("Error on calculating remaining bandwidth", err.Error())
			return
		}

		allowedValues, err := r.client.CAVSDK.V1.EdgeGateway.GetAllowedBandwidthValues(plan.Tier0VrfID.Get())
		if err != nil {
			resp.Diagnostics.AddError("Error on calculating allowed Bandwidth values", err.Error())
			return
		}

		remaining = calculBestValue(remaining, allowedValues)

		resp.Diagnostics.AddAttributeWarning(path.Root("bandwidth"), "Bandwidth value is unknown, will be set to remaining bandwidth", fmt.Sprintf("Bandwidth defined to %dMbps", remaining))
		plan.Bandwidth.SetInt(remaining)
		goto END

	// Update case
	case !plan.Bandwidth.Equal(state.Bandwidth):
		allowedValues, err := r.client.CAVSDK.V1.EdgeGateway.GetAllowedBandwidthValues(plan.Tier0VrfID.Get())
		if err != nil {
			resp.Diagnostics.AddError("Error on calculating allowed Bandwidth values", err.Error())
			return
		}

		// Ignore error because recalculating remaining bandwidth with bandwidth released by the update
		remaining, _ := loadRemainingBandwidth()
		remainingOnUpdate := calculBestValue(remaining+state.Bandwidth.GetInt(), allowedValues)

		if plan.Bandwidth.IsUnknown() && remainingOnUpdate > 0 {
			plan.Bandwidth.SetInt(remainingOnUpdate)
		} else if plan.Bandwidth.GetInt() > remainingOnUpdate {
			resp.Diagnostics.AddAttributeError(path.Root("bandwidth"), "Overcommitting bandwidth", fmt.Sprintf("Not enough bandwidth available, requested: %dMbps, available: %dMbps", plan.Bandwidth.GetInt(), remainingOnUpdate))
			return
		}

		allowedValuesFunc()
		goto END
	}

END:
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

	// List all edge gateways for determining the ID of the new edge gateway
	edgegws, err := r.client.CAVSDK.V1.EdgeGateway.List()
	if err != nil {
		resp.Diagnostics.AddError("Error listing edge gateways", err.Error())
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	var job *commoncloudavenue.JobStatus

	switch plan.OwnerType.ValueString() {
	case "vdc":
		job, err = r.client.CAVSDK.V1.EdgeGateway.New(plan.OwnerName.Get(), plan.Tier0VrfID.Get())
	case "vdc-group":
		job, err = r.client.CAVSDK.V1.EdgeGateway.NewFromVDCGroup(plan.OwnerName.Get(), plan.Tier0VrfID.Get())
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

	var edgegwNew v1.EdgeGw

	// Find the new edge gateway in the list of all edge gateways and set the ID. New edge gateway is in the list refreshed but not in the old list.
	for _, edgegw := range *edgegwsRefreshed {
		var found bool
		for _, edgegwOld := range *edgegws {
			if edgegw.GetID() == edgegwOld.GetID() {
				found = true
			}
		}
		if !found {
			plan.ID.Set(uuid.Normalize(uuid.Gateway, edgegw.GetID()).String())
			plan.Name.Set(edgegw.GetName())
			edgegwNew = edgegw
			break
		}
	}

	if edgegwNew == (v1.EdgeGw{}) {
		resp.Diagnostics.AddError("Error retrieving new edge gateway", "New edge gateway not found")
		return
	}

	if int(edgegwNew.GetBandwidth()) != plan.Bandwidth.GetInt() {
		job, err = edgegwNew.UpdateBandwidth(plan.Bandwidth.GetInt())
		if err != nil {
			resp.Diagnostics.AddError("Error setting Bandwidth", err.Error())
			return
		}

		if err := job.Wait(1, int(createTimeout.Seconds())); err != nil {
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

// Read refreshes the Terraform state with the latest data.
func (r *edgeGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Read)()

	state := &edgeGatewayResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read timeout
	readTimeout, errTO := state.Timeouts.Read(ctx, 8*time.Minute)
	if errTO != nil {
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

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	edgegw, err := r.client.CAVSDK.V1.EdgeGateway.GetByID(plan.ID.Get())
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

	deleteTimeout, errTO := state.Timeouts.Update(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	edgegw, err := r.client.CAVSDK.V1.EdgeGateway.GetByID(state.ID.Get())
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

	var (
		edgegw *v1.EdgeGw
		err    error
	)

	switch {
	case planOrState.ID.IsKnown():
		edgegw, err = r.client.CAVSDK.V1.EdgeGateway.GetByID(planOrState.ID.Get())
		if err != nil {
			if commoncloudavenue.IsNotFound(err) {
				return nil, false, nil
			}
			diags.AddError("Error retrieving edge gateway", err.Error())
			return nil, true, diags
		}
	case planOrState.Name.IsKnown():
		edgegw, err = r.client.CAVSDK.V1.EdgeGateway.GetByName(planOrState.Name.Get())
		if err != nil {
			if commoncloudavenue.IsNotFound(err) {
				return nil, false, nil
			}
			diags.AddError("Error retrieving edge gateway", err.Error())
			return nil, true, diags
		}
	default:
		diags.AddError("Error retrieving edge gateway", "Either name or ID must be set")
		return nil, true, diags
	}

	if !planOrState.ID.IsKnown() {
		stateRefreshed.ID.Set(uuid.Normalize(uuid.Gateway, edgegw.GetID()).String())
	}

	stateRefreshed.Tier0VrfID.Set(edgegw.GetTier0VrfID())
	stateRefreshed.OwnerName.Set(edgegw.GetOwnerName())
	stateRefreshed.OwnerType.Set(string(edgegw.GetOwnerType()))
	stateRefreshed.Description.Set(edgegw.GetDescription())
	stateRefreshed.Bandwidth.SetInt(int(edgegw.GetBandwidth()))

	// EnableLoadBalancing is now deprecated, but we still need to set it to false if it is unknown
	if !stateRefreshed.EnableLoadBalancing.IsKnown() {
		stateRefreshed.EnableLoadBalancing.Set(false)
	}

	return stateRefreshed, true, nil
}
