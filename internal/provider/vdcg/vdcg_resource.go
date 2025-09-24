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

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdcgroup/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vdcgResource{}
	_ resource.ResourceWithConfigure   = &vdcgResource{}
	_ resource.ResourceWithImportState = &vdcgResource{}
)

// NewVDCGResource is a helper function to simplify the provider implementation.
func NewVDCGResource() resource.Resource {
	return &vdcgResource{}
}

// vdcgResource is the resource implementation.
type vdcgResource struct {
	// Client is a terraform Client
	client *client.CloudAvenue

	// vgClient is the VDC Group client from the SDK V2
	vgClient *vdcgroup.Client
}

// Metadata returns the resource type name.
func (r *vdcgResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *vdcgResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vdcgSchema(ctx).GetResource(ctx)
}

func (r *vdcgResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	vgC, err := vdcgroup.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create VDC Group client, got error: %s", err),
		)
		return
	}

	r.client = client
	r.vgClient = vgC
}

// Create creates the resource and sets the initial Terraform state.
func (r *vdcgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Create)()

	plan := &vdcgModel{}

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource creation logic here.
	*/

	vdcIDs, d := plan.VDCIDs.Get(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	vdcGroup, err := r.vgClient.CreateVdcGroup(ctx, types.ParamsCreateVdcGroup{
		Name:        plan.Name.Get(),
		Description: plan.Description.Get(),
		Vdcs: func() (vdcs []types.ParamsCreateVdcGroupVdc) {
			for _, id := range vdcIDs {
				vdcs = append(vdcs, types.ParamsCreateVdcGroupVdc{
					ID: id,
				})
			}
			return
		}(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating VDC Group", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.fromSDK(ctx, vdcGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vdcgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Read)()

	state := &vdcgModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
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
func (r *vdcgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &vdcgModel{}
		state = &vdcgModel{}
	)

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource update here
	*/

	vdcIDs, d := plan.VDCIDs.Get(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	vdcGroup, err := r.vgClient.UpdateVdcGroup(ctx, types.ParamsUpdateVdcGroup{
		ID:   state.ID.Get(),
		Name: plan.Name.Get(),
		Description: func() *string {
			if plan.Description.IsNull() {
				return nil
			}
			return plan.Description.GetPtr()
		}(),
		Vdcs: func() (vdcs []types.ParamsCreateVdcGroupVdc) {
			for _, id := range vdcIDs {
				vdcs = append(vdcs, types.ParamsCreateVdcGroupVdc{
					ID: id,
				})
			}
			return
		}(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating VDC Group", err.Error())
		return
	}

	resp.Diagnostics.Append(state.fromSDK(ctx, vdcGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vdcgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Delete)()

	state := &vdcgModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.vgClient.DeleteVdcGroup(ctx, types.ParamsDeleteVdcGroup{
		ID:    state.ID.Get(),
		Name:  state.Name.Get(),
		Force: false,
	}); err != nil {
		resp.Diagnostics.AddError("Error deleting VDC Group", err.Error())
		return
	}
}

func (r *vdcgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdcg", r.client.GetOrgName(), metrics.Import)()

	// id format is vdcGroupIDOrName

	param := types.ParamsGetVdcGroup{}
	if urn.IsVDCGroup(req.ID) {
		param.ID = req.ID
	} else {
		param.Name = req.ID
	}

	vdcGroup, err := r.vgClient.GetVdcGroup(ctx, param)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC Group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), vdcGroup.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), vdcGroup.Name)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *vdcgResource) read(ctx context.Context, planOrState *vdcgModel) (stateRefreshed *vdcgModel, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	vdcGroup, err := r.vgClient.GetVdcGroup(ctx, types.ParamsGetVdcGroup{
		ID:   stateRefreshed.ID.Get(),
		Name: stateRefreshed.Name.Get(),
	})
	if err != nil {
		diags.AddError("Error reading VDC Group", err.Error())
		return nil, diags
	}

	diags.Append(stateRefreshed.fromSDK(ctx, vdcGroup)...)

	return stateRefreshed, diags
}
