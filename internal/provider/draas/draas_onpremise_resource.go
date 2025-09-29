/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package draas provides a Terraform resource to manage draas.
package draas

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/draas/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vcda"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &draasIPResource{}
	_ resource.ResourceWithConfigure   = &draasIPResource{}
	_ resource.ResourceWithImportState = &draasIPResource{}
	_ resource.ResourceWithMoveState   = &draasIPResource{}
)

// NewDraasIPResource is a helper function to simplify the provider implementation.
func NewDraasIPResource() resource.Resource {
	return &draasIPResource{}
}

// draasIPResource is the resource implementation.
type draasIPResource struct {
	client  *client.CloudAvenue
	dClient *draas.Client
}

func (r *draasIPResource) MoveState(ctx context.Context) []resource.StateMover {
	sc := vcda.VcdaIPSchema().GetResource(ctx)
	return []resource.StateMover{
		{
			SourceSchema: &sc,
			StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
				if req.SourceTypeName != "cloudavenue_vcda_ip" {
					return
				}

				var sourceStateData vcda.VcdaIPResourceModel

				resp.Diagnostics.Append(req.SourceState.Get(ctx, &sourceStateData)...)
				if resp.Diagnostics.HasError() {
					return
				}

				targetStateData := &draasIPResourceModel{
					ID:        sourceStateData.ID,
					IPAddress: sourceStateData.IPAddress,
				}

				resp.Diagnostics.Append(resp.TargetState.Set(ctx, targetStateData)...)
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *draasIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "onpremise"
}

// Schema defines the schema for the resource.
func (r *draasIPResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = draasIPSchema().GetResource(ctx)
}

// Configure configures the resource.
func (r *draasIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	dClient, err := draas.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Failed to create draas client: %v", err),
		)

		return
	}

	r.client = client
	r.dClient = dClient
}

// Create creates the resource and sets the initial Terraform state.
func (r *draasIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_draas_ip", r.client.GetOrgName(), metrics.Create)()

	plan := new(draasIPResourceModel)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ! Lock/Unlock is kept while all resources migrate to sdk v2
	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	err := r.dClient.AddOnPremiseIp(ctx, types.ParamsAddDraasOnPremiseIP{
		IP: plan.IPAddress.Get(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error on add Draas onpremise IP", err.Error())
		return
	}

	stateRefreshed, diags := r.read(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *draasIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_draas_ip", r.client.GetOrgName(), metrics.Read)()

	state := new(draasIPResourceModel)

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, diags := r.read(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *draasIPResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *draasIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_draas_ip", r.client.GetOrgName(), metrics.Delete)()

	state := new(draasIPResourceModel)

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ! Lock/Unlock is kept while all resources migrate to sdk v2
	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	err := r.dClient.RemoveOnPremiseIp(ctx, types.ParamsRemoveDraasOnPremiseIP{
		IP: state.IPAddress.Get(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error on delete VDCA IP", err.Error())
		return
	}
}

func (r *draasIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip_address"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), tftypes.StringValue(urn.Normalize(
		urn.VCDA,
		utils.GenerateUUID(
			req.ID,
		).ValueString(),
	).String()))...)
}

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *draasIPResource) read(ctx context.Context, planOrState *draasIPResourceModel) (stateRefreshed *draasIPResourceModel, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	listOfIps, err := r.dClient.ListOnPremiseIp(ctx)
	if err != nil {
		diags.AddError("Error on list VDCA IPs", err.Error())
		return nil, diags
	}

	if !slices.Contains(listOfIps.IPs, planOrState.IPAddress.Get()) {
		diags.AddError("Draas OnPremise IP not found", fmt.Sprintf("The Draas OnPremise IP '%s' was not found", planOrState.IPAddress.Get()))
		return nil, diags
	}

	stateRefreshed.ID.Set(urn.Normalize(
		urn.VCDA,
		utils.GenerateUUID(planOrState.IPAddress.Get()).ValueString(),
	).String())

	return stateRefreshed, nil
}
