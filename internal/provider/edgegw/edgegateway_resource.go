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
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
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

type edgeGatewayResource struct {
	// Client is a terraform Client
	client *client.CloudAvenue

	// eClient is the Edge Gateway client from the SDK V2
	eClient *edgegateway.Client
}

// ModifyPlan modifies the plan to add the default values.
func (r *edgeGatewayResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) { //nolint:gocyclo
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

	if plan != nil && state == nil {
		// This is a create operation.
		// Disallow deprecated fields for new resources
		if plan.Tier0VRFName.IsKnown() {
			resp.Diagnostics.AddAttributeError(path.Root("tier0_vrf_name"), "Field is deprecated", "Please use 't0_name' instead")
		}
	}

	// TO

	t0s, err := r.eClient.ListT0(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing T0s", err.Error())
		return
	}

	if len(t0s.T0s) == 0 {
		resp.Diagnostics.AddError("Error listing T0s", "No T0s found")
		return
	}

	var t0 types.ModelT0

	// If Tier0VRFName is not known, we need to find the T0
	// IF multiple T0s are available, return an error
	if !plan.Tier0VRFName.IsKnown() && !plan.T0Name.IsKnown() {
		if len(t0s.T0s) > 1 {
			resp.Diagnostics.AddAttributeError(path.Root("t0_name"), "Error listing T0s", "Multiple T0s found, please specify the T0 name")
			return
		}

		t0 = t0s.T0s[0]
	} else {
		for _, t0x := range t0s.T0s {
			if t0x.Name == plan.Tier0VRFName.Get() || t0x.Name == plan.T0Name.Get() {
				t0 = t0x
				break
			}
		}
	}

	if t0.Name == "" {
		resp.Diagnostics.AddError("Error retrieving T0", "T0 not found")
		return
	}

	// This is a create operation.
	if plan != nil && state == nil {
		if t0.MaxEdgeGateways == len(t0.EdgeGateways) {
			resp.Diagnostics.AddError("Error creating edge gateway", "Maximum number of edge gateways reached for T0 "+t0.Name)
			return
		}
	}

	plan.Tier0VRFName.Set(t0.Name)
	plan.T0Name.Set(t0.Name)

	switch {
	// Create case with value is known
	case plan.Bandwidth.IsKnown() && (state == nil || !state.Bandwidth.IsKnown()):
		// In this case edgegateway is not already created or the T0 allow unlimited bandwidth
		// t0.Bandwidth.AllowedBandwidthValues AllowedBandwidthValues returns the allowed bandwidth values for the T0 router. It's used to determine the available bandwidth options for the new edge gateway.

		if len(t0.Bandwidth.AllowedBandwidthValues) == 0 {
			resp.Diagnostics.AddAttributeError(path.Root("bandwidth"), "Error on calculating remaining bandwidth", "Not enough bandwidth available")
			return
		}

		if slices.Contains(t0.Bandwidth.AllowedBandwidthValues, plan.Bandwidth.GetInt()) {
			// Value defined match with AllowedBandwidthValues
			goto END
		}

		// If we reach this point, the value is not allowed
		resp.Diagnostics.AddAttributeError(path.Root("bandwidth"), "Invalid bandwidth value", fmt.Sprintf("Bandwidth value %dMbps is not allowed. (Allowed values: %v)", plan.Bandwidth.GetInt(), t0.Bandwidth.AllowedBandwidthValues))
		goto END

	// Create case with value is unknown
	case !plan.Bandwidth.IsKnown():
		if t0.Bandwidth.AllowUnlimited {
			goto END
		}

		if len(t0.Bandwidth.AllowedBandwidthValues) == 0 {
			resp.Diagnostics.AddAttributeError(path.Root("bandwidth"), "Error on calculating remaining bandwidth", "Not enough bandwidth available")
			return
		}

		// The best values is a last element in a slice
		bestAllowedValues := t0.Bandwidth.AllowedBandwidthValues[len(t0.Bandwidth.AllowedBandwidthValues)-1]

		resp.Diagnostics.AddAttributeWarning(path.Root("bandwidth"), "Bandwidth value is unknown, will be set to remaining bandwidth.", fmt.Sprintf("Bandwidth defined to %dMbps. (Allowed values : %v)", bestAllowedValues, t0.Bandwidth.AllowedBandwidthValues))
		plan.Bandwidth.SetInt(bestAllowedValues)
		goto END

	// Update case
	case !plan.Bandwidth.Equal(state.Bandwidth):
		if (plan.Bandwidth.IsUnknown() || plan.Bandwidth.GetInt() == 0) && t0.Bandwidth.AllowUnlimited {
			// This case is allowed return without error
			goto END
		}

		// Find edgegateway
		var edgegateway types.ModelT0EdgeGateway

		for _, gw := range t0.EdgeGateways {
			if gw.ID == plan.ID.Get() {
				edgegateway = gw
				break
			}
		}

		if edgegateway.ID == "" {
			resp.Diagnostics.AddError("Error retrieving edge gateway", "Edge gateway not found")
			return
		}

		if slices.Contains(edgegateway.AllowedBandwidthValues, plan.Bandwidth.GetInt()) {
			// Value defined match with AllowedBandwidthValues
			goto END
		}

		resp.Diagnostics.AddAttributeError(path.Root("bandwidth"), "Invalid bandwidth value", fmt.Sprintf("Bandwidth value %dMbps is not allowed. (Allowed values: %v)", plan.Bandwidth.GetInt(), edgegateway.AllowedBandwidthValues))
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

	eC, err := edgegateway.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Edge Gateway client, got error: %s", err),
		)
		return
	}

	r.client = client
	r.eClient = eC
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

	// ! Lock/Unlock is kept while all resources migrate to sdk v2
	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	edgegatewayCreated, err := r.eClient.CreateEdgeGateway(ctx, types.ParamsCreateEdgeGateway{
		OwnerName: plan.OwnerName.Get(),
		T0Name:    plan.T0Name.Get(),
		Bandwidth: plan.Bandwidth.GetInt(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating edge gateway", err.Error())
		return
	}

	plan.fromSDK(edgegatewayCreated)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	stateRefreshed, d := r.read(ctx, state)
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

	// ! Lock/Unlock is kept while all resources migrate to sdk v2
	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	if !plan.Bandwidth.Equal(state.Bandwidth) {
		_, err := r.eClient.UpdateEdgeGateway(ctx, types.ParamsUpdateEdgeGateway{
			ID:        plan.ID.Get(),
			Name:      plan.Name.Get(),
			Bandwidth: plan.Bandwidth.GetInt(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating edge gateway", err.Error())
			return
		}
	}

	// Use generic read function to refresh the state
	stateRefreshed, d := r.read(ctx, plan)
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

	// ! Lock/Unlock is kept while all resources migrate to sdk v2
	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	err := r.eClient.DeleteEdgeGateway(ctx, types.ParamsEdgeGateway{
		ID:   state.ID.Get(),
		Name: state.Name.Get(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting edge gateway", err.Error())
		return
	}
}

func (r *edgeGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Import)()

	// Retrieve import Name and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// * Custom funcs.
func (r *edgeGatewayResource) read(ctx context.Context, planOrState *edgeGatewayResourceModel) (stateRefreshed *edgeGatewayResourceModel, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	edgegateway, err := r.eClient.GetEdgeGateway(ctx, types.ParamsEdgeGateway{
		ID:   planOrState.ID.Get(),
		Name: planOrState.Name.Get(),
	})
	if err != nil {
		diags.AddError("Error retrieving edge gateway", err.Error())
		return nil, diags
	}

	stateRefreshed.fromSDK(edgegateway)

	bandwidth, err := r.eClient.GetBandwidth(ctx, types.ParamsEdgeGateway{
		ID:   edgegateway.ID,
		Name: edgegateway.Name,
	})
	if err != nil {
		diags.AddError("Error retrieving bandwidth", err.Error())
		return nil, diags
	}

	stateRefreshed.Bandwidth.SetInt(bandwidth.Bandwidth)

	return stateRefreshed, nil
}
