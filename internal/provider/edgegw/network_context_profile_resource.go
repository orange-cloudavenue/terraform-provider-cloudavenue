/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	"context"
	"errors"
	"fmt"
	"strings"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	sdkv1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

var (
	_ resource.Resource                = &networkContextProfileResource{}
	_ resource.ResourceWithConfigure   = &networkContextProfileResource{}
	_ resource.ResourceWithImportState = &networkContextProfileResource{}
)

// NewNetworkContextProfileResource returns a new context profile resource.
func NewNetworkContextProfileResource() resource.Resource {
	return &networkContextProfileResource{}
}

type networkContextProfileResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init initializes the resource.
func (r *networkContextProfileResource) Init(_ context.Context, rm *networkContextProfileModel) (diags diag.Diagnostics) {
	var err error

	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return diags
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(rm.EdgeGatewayID.Get()),
		Name: types.StringValue(rm.EdgeGatewayName.Get()),
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
	}

	return diags
}

func (r *networkContextProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_network_context_profile"
}

func (r *networkContextProfileResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = networkContextProfileSchema(ctx).GetResource(ctx)
}

func (r *networkContextProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *networkContextProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_context_profile", r.client.GetOrgName(), metrics.Create)()

	plan := &networkContextProfileModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())

	profile, d := plan.toSDKProfile(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.edgegw.EdgeClient.CreateNetworkContextProfile(profile)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Network Context Profile", err.Error())
		return
	}

	plan.ID.Set(created.ID)

	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Error refreshing state", "Could not find the created Network Context Profile")
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

func (r *networkContextProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_context_profile", r.client.GetOrgName(), metrics.Read)()

	state := &networkContextProfileModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

func (r *networkContextProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_context_profile", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &networkContextProfileModel{}
		state = &networkContextProfileModel{}
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())

	// Carry ID from state to plan for the update call.
	plan.ID.Set(state.ID.Get())

	profile, d := plan.toSDKProfile(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.edgegw.EdgeClient.UpdateNetworkContextProfile(profile); err != nil {
		resp.Diagnostics.AddError("Error updating Network Context Profile", err.Error())
		return
	}

	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Error refreshing state", "Could not find the updated Network Context Profile")
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

func (r *networkContextProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_context_profile", r.client.GetOrgName(), metrics.Delete)()

	state := &networkContextProfileModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())

	if err := r.edgegw.EdgeClient.DeleteNetworkContextProfile(state.ID.Get()); err != nil {
		resp.Diagnostics.AddError("Error deleting Network Context Profile", err.Error())
	}
}

func (r *networkContextProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway_network_context_profile", r.client.GetOrgName(), metrics.Import)()

	// Import format: <edge_gateway_name_or_id>.<profile_id_or_name>
	parts := strings.Split(req.ID, ".")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format <edge_gateway_name_or_id>.<profile_id_or_name>",
		)
		return
	}
	edgeIDOrName, profileIDOrName := parts[0], parts[1]

	x := &networkContextProfileModel{
		ID:              supertypes.NewStringNull(),
		Name:            supertypes.NewStringNull(),
		EdgeGatewayID:   supertypes.NewStringNull(),
		EdgeGatewayName: supertypes.NewStringNull(),
	}

	if urn.IsEdgeGateway(edgeIDOrName) {
		x.EdgeGatewayID.Set(edgeIDOrName)
	} else {
		x.EdgeGatewayName.Set(edgeIDOrName)
	}

	if urn.IsNetworkContextProfile(profileIDOrName) {
		x.ID.Set(profileIDOrName)
	} else {
		x.Name.Set(profileIDOrName)
	}

	resp.Diagnostics.Append(r.Init(ctx, x)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, d := r.read(ctx, x)
	if !found {
		resp.Diagnostics.AddError("Import failed", fmt.Sprintf("Network Context Profile %q not found", profileIDOrName))
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// read is the generic read function.
func (r *networkContextProfileResource) read(ctx context.Context, planOrState *networkContextProfileModel) (stateRefreshed *networkContextProfileModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	var (
		profile *sdkv1.NetworkContextProfile
		err     error
	)

	if planOrState.ID.IsKnown() && planOrState.ID.Get() != "" {
		profile, err = r.edgegw.EdgeClient.GetNetworkContextProfileByID(planOrState.ID.Get())
	} else {
		profile, err = r.edgegw.EdgeClient.GetNetworkContextProfileByName(planOrState.Name.Get())
	}

	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			return stateRefreshed, false, nil
		}
		diags.AddError("Error reading Network Context Profile", err.Error())
		return stateRefreshed, true, diags
	}

	diags.Append(stateRefreshed.fromSDKProfile(ctx, profile)...)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())

	return stateRefreshed, true, diags
}
