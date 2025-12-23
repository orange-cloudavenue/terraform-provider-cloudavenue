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
	"strings"

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
	_ resource.Resource                = &AppPortProfileResource{}
	_ resource.ResourceWithConfigure   = &AppPortProfileResource{}
	_ resource.ResourceWithImportState = &AppPortProfileResource{}
)

// NewAppPortProfileResource is a helper function to simplify the provider implementation.
func NewAppPortProfileResource() resource.Resource {
	return &AppPortProfileResource{}
}

// AppPortProfileResource is the resource implementation.
type AppPortProfileResource struct {
	client   *client.CloudAvenue
	vdcGroup *v1.VDCGroup
}

// Init Initializes the resource.
func (r *AppPortProfileResource) Init(_ context.Context, rm *AppPortProfileModel) (diags diag.Diagnostics) {
	var err error

	idOrName := rm.VDCGroupName.Get()
	if rm.VDCGroupID.IsKnown() {
		// Use the ID
		idOrName = rm.VDCGroupID.Get()
	}

	r.vdcGroup, err = r.client.CAVSDK.V1.VDC().GetVDCGroup(idOrName)
	if err != nil {
		diags.AddError("Error retrieving VDC Group", err.Error())
		return diags
	}
	return diags
}

// Metadata returns the resource type name.
func (r *AppPortProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_app_port_profile"
}

// Schema defines the schema for the resource.
func (r *AppPortProfileResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = appPortProfileSchema(ctx).GetResource(ctx)
}

func (r *AppPortProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *AppPortProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdcg_app_port_profile", r.client.GetOrgName(), metrics.Create)()

	plan := &AppPortProfileModel{}

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

	appPortProfileModel, d := plan.toSDKAppPortProfile(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Create the application port profile
	appPortProfile, err := r.vdcGroup.CreateFirewallAppPortProfile(appPortProfileModel)
	if err != nil {
		resp.Diagnostics.AddError("Error creating application port profile", err.Error())
		return
	}

	// Set the ID
	plan.ID.Set(appPortProfile.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Error refreshing state", "Could not find the created application port profile")
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
func (r *AppPortProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdcg_app_port_profile", r.client.GetOrgName(), metrics.Read)()

	state := &AppPortProfileModel{}

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
func (r *AppPortProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdcg_app_port_profile", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &AppPortProfileModel{}
		state = &AppPortProfileModel{}
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

	appPortProfile, err := r.vdcGroup.GetFirewallAppPortProfile(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application port profile", err.Error())
		return
	}

	appPortProfileModel, d := plan.toSDKAppPortProfile(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Update the application port profile
	if err := appPortProfile.Update(appPortProfileModel); err != nil {
		resp.Diagnostics.AddError("Error updating application port profile", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Error refreshing state", "Could not find the updated application port profile")
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
func (r *AppPortProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdcg_app_port_profile", r.client.GetOrgName(), metrics.Delete)()

	state := &AppPortProfileModel{}

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

	appPortProfile, err := r.vdcGroup.GetFirewallAppPortProfile(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application port profile", err.Error())
		return
	}

	// Delete the application port profile
	if err := appPortProfile.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting application port profile", err.Error())
		return
	}
}

func (r *AppPortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdcg_app_port_profile", r.client.GetOrgName(), metrics.Import)()

	// split req.ID into edge gateway ID and app port profile ID/name
	split := strings.Split(req.ID, ".")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be in the format <edge_gateway_id_or_name>.<app_port_profile_id_or_name>")
		return
	}
	vdcGroupIDOrName, appPortProfileIDOrName := split[0], split[1]

	x := &AppPortProfileModel{
		ID:           supertypes.NewStringNull(),
		Name:         supertypes.NewStringNull(),
		VDCGroupID:   supertypes.NewStringNull(),
		VDCGroupName: supertypes.NewStringNull(),
	}

	if urn.IsVDCGroup(vdcGroupIDOrName) {
		x.VDCGroupID.Set(vdcGroupIDOrName)
	} else {
		x.VDCGroupName.Set(vdcGroupIDOrName)
	}

	if urn.IsAppPortProfile(appPortProfileIDOrName) {
		x.ID.Set(appPortProfileIDOrName)
	} else {
		x.Name.Set(appPortProfileIDOrName)
	}

	resp.Diagnostics.Append(r.Init(ctx, x)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, d := r.read(ctx, x)
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

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *AppPortProfileResource) read(ctx context.Context, planOrState *AppPortProfileModel) (stateRefreshed *AppPortProfileModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	var (
		appPortProfile *v1.FirewallGroupAppPortProfile
		err            error
		nameOrID       = planOrState.Name.Get()
	)

	if planOrState.ID.IsKnown() {
		// Use the ID
		nameOrID = planOrState.ID.Get()
	}

	appPortProfile, err = r.vdcGroup.GetFirewallAppPortProfile(nameOrID)
	if err != nil {
		if govcd.IsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error reading App Port Profile", err.Error())
		return stateRefreshed, found, diags
	}

	appPorts := make([]*AppPortProfileModelAppPort, len(appPortProfile.ApplicationPorts))
	for index, singlePort := range appPortProfile.ApplicationPorts {
		ap := &AppPortProfileModelAppPort{
			Protocol: supertypes.NewStringNull(),
			Ports:    supertypes.NewSetValueOfNull[string](ctx),
		}

		ap.Protocol.Set(string(singlePort.Protocol))
		if singlePort.Protocol == v1.FirewallGroupAppPortProfileModelPortProtocolTCP || singlePort.Protocol == v1.FirewallGroupAppPortProfileModelPortProtocolUDP {
			// DestinationPorts is optional
			if len(singlePort.DestinationPorts) > 0 {
				diags.Append(ap.Ports.Set(ctx, singlePort.DestinationPorts)...)
				if diags.HasError() {
					return stateRefreshed, found, diags
				}
			}
		}
		appPorts[index] = ap
	}

	stateRefreshed.ID.Set(appPortProfile.ID)
	stateRefreshed.Name.Set(appPortProfile.Name)
	stateRefreshed.Description.Set(appPortProfile.Description)
	stateRefreshed.AppPorts.Set(ctx, appPorts)
	stateRefreshed.VDCGroupID.Set(r.vdcGroup.GetID())
	stateRefreshed.VDCGroupName.Set(r.vdcGroup.GetName())

	return stateRefreshed, true, nil
}
