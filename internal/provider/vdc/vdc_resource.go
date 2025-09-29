/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"context"
	"fmt"
	"slices"

	"golang.org/x/sync/errgroup"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vdcResource{}
	_ resource.ResourceWithConfigure   = &vdcResource{}
	_ resource.ResourceWithImportState = &vdcResource{}
)

// NewVDCResource is a helper function to simplify the provider implementation.
func NewVDCResource() resource.Resource {
	return &vdcResource{}
}

// vdcResource is the resource implementation.
type vdcResource struct {
	// Client is a terraform Client
	client *client.CloudAvenue

	// vClient is the VDC client from the SDK V2
	vClient *vdc.Client
}

// Metadata returns the resource type name.
func (r *vdcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *vdcResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vdcSchema().GetResource(ctx)
}

// Configure configures the resource.
func (r *vdcResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	vC, err := vdc.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create VDC client, got error: %s", err),
		)
		return
	}

	r.client = client
	r.vClient = vC
}

// Create creates the resource and sets the initial Terraform state.
func (r *vdcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	plan := new(vdcModel)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ! Lock/Unlock is kept while all resources migrate to sdk v2
	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	paramsCreate := types.ParamsCreateVDC{
		Name:                plan.Name.Get(),
		Description:         plan.Description.Get(),
		ServiceClass:        plan.ServiceClass.Get(),
		DisponibilityClass:  plan.DisponibilityClass.Get(),
		BillingModel:        plan.BillingModel.Get(),
		StorageBillingModel: plan.StorageBillingModel.Get(),
		Vcpu: func() int {
			if plan.VCPU.IsKnown() {
				return plan.VCPU.GetInt()
			}
			return plan.CPUAllocated.GetInt() / plan.VCPUInMhz.GetInt() // ! Deprecated fields - Maintain for backward compatibility
		}(),
		Memory: func() int {
			if plan.Memory.IsKnown() {
				return plan.Memory.GetInt()
			}
			return plan.MemoryAllocated.GetInt() // ! Deprecated fields - Maintain for backward compatibility
		}(),
		// * StorageProfiles
		StorageProfiles: func() (sps []types.ParamsCreateVDCStorageProfile) {
			if plan.StorageProfiles.IsNull() || plan.StorageProfiles.IsUnknown() {
				return sps
			}
			storageProfiles := plan.StorageProfiles.DiagsGet(ctx, resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return sps
			}

			for _, sp := range storageProfiles {
				sps = append(sps, types.ParamsCreateVDCStorageProfile{
					Class:   sp.Class.Get(),
					Limit:   sp.Limit.GetInt(),
					Default: *sp.Default.GetPtr(),
				})
			}
			return sps
		}(),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	vdcCreated, err := r.vClient.CreateVDC(ctx, paramsCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating VDC",
			"Could not create VDC, unexpected error: "+err.Error(),
		)
		return
	}

	storageProfiles, err := r.vClient.ListStorageProfile(ctx, types.ParamsListStorageProfile{
		VdcID: vdcCreated.ID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Listing Storage Profiles",
			"Could not list storage profiles, unexpected error: "+err.Error(),
		)
		return
	}

	plan.fromSDK(ctx, vdcCreated, storageProfiles, &resp.Diagnostics)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vdcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Read)()

	state := new(vdcModel)
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
func (r *vdcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Update)()

	// Retrieve values from plan
	plan := &vdcModel{}
	state := &vdcModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ! Lock/Unlock is kept while all resources migrate to sdk v2
	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// VDC update are split in two calls:
	// - VDC properties/resources update (description, vcpu, memory)
	// - VDC storage profiles update (add/remove storage profile, change limit/default)

	var requireUpdateVDC bool

	paramsUpdateVDC := types.ParamsUpdateVDC{
		ID:   state.ID.Get(),
		Name: state.Name.Get(),
	}

	if !plan.Description.Equal(state.Description) {
		requireUpdateVDC = true
		paramsUpdateVDC.Description = plan.Description.GetPtr()
	}

	if !plan.VCPU.Equal(state.VCPU) {
		requireUpdateVDC = true
		paramsUpdateVDC.Vcpu = plan.VCPU.GetIntPtr()
	}

	if !plan.Memory.Equal(state.Memory) {
		requireUpdateVDC = true
		paramsUpdateVDC.Memory = plan.Memory.GetIntPtr()
	}

	if requireUpdateVDC {
		_, err := r.vClient.UpdateVDC(ctx, paramsUpdateVDC)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating VDC",
				"Could not update VDC, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if !plan.StorageProfiles.Equal(state.StorageProfiles) {
		storageProfilesPlan := plan.StorageProfiles.DiagsGet(ctx, resp.Diagnostics)
		storageProfilesState := state.StorageProfiles.DiagsGet(ctx, resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		listOfChanges := map[string][]*vdcModelStorageProfile{}

		for _, sPlan := range storageProfilesPlan {
			// If a difference is detected, the profile is added to the to_update list.
			if slices.ContainsFunc(storageProfilesState, func(sp *vdcModelStorageProfile) bool {
				return sp.Class.Equal(sPlan.Class) && (!sp.Limit.Equal(sPlan.Limit) &&
					sp.Default.Equal(sPlan.Default))
			}) {
				listOfChanges["to_update"] = append(listOfChanges["to_update"], sPlan)
			}

			// If Class is not find in state, the profile is added to the to_add list.
			if !slices.ContainsFunc(storageProfilesState, func(sp *vdcModelStorageProfile) bool {
				return sp.Class.Equal(sPlan.Class)
			}) {
				listOfChanges["to_add"] = append(listOfChanges["to_add"], sPlan)
			}
		}

		for _, sState := range storageProfilesState {
			// If Class is not find in plan, the profile is add to the to_remove list.
			if !slices.ContainsFunc(storageProfilesPlan, func(sp *vdcModelStorageProfile) bool {
				return sp.Class.Equal(sState.Class)
			}) {
				listOfChanges["to_remove"] = append(listOfChanges["to_remove"], sState)
			}
		}

		if len(listOfChanges["to_update"]) > 0 {
			params := types.ParamsUpdateStorageProfile{
				VdcID:           state.ID.Get(),
				StorageProfiles: make([]types.ParamsUpdateVDCStorageProfile, 0),
			}
			for _, sp := range listOfChanges["to_update"] {
				params.StorageProfiles = append(params.StorageProfiles, types.ParamsUpdateVDCStorageProfile{
					Class:   sp.Class.Get(),
					Limit:   sp.Limit.GetInt(),
					Default: sp.Default.GetPtr(),
				})
			}

			_, err := r.vClient.UpdateStorageProfile(ctx, params)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Updating VDC Storage Profiles",
					"Could not update VDC storage profiles, unexpected error: "+err.Error(),
				)
				return
			}
		}
		if len(listOfChanges["to_add"]) > 0 {
			params := types.ParamsAddStorageProfile{
				VdcID:           state.ID.Get(),
				StorageProfiles: make([]types.ParamsCreateVDCStorageProfile, 0),
			}
			for _, sp := range listOfChanges["to_add"] {
				params.StorageProfiles = append(params.StorageProfiles, types.ParamsCreateVDCStorageProfile{
					Class:   sp.Class.Get(),
					Limit:   sp.Limit.GetInt(),
					Default: sp.Default.Get(),
				})
			}

			err := r.vClient.AddStorageProfile(ctx, params)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Adding VDC Storage Profiles",
					"Could not add VDC storage profiles, unexpected error: "+err.Error(),
				)
				return
			}
		}
		if len(listOfChanges["to_remove"]) > 0 {
			params := types.ParamsDeleteStorageProfile{
				VdcID:           state.ID.Get(),
				StorageProfiles: make([]types.ParamsDeleteVDCStorageProfile, 0),
			}
			for _, sp := range listOfChanges["to_remove"] {
				params.StorageProfiles = append(params.StorageProfiles, types.ParamsDeleteVDCStorageProfile{
					Class: sp.Class.Get(),
				})
			}

			err := r.vClient.DeleteStorageProfile(ctx, params)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Deleting VDC Storage Profiles",
					"Could not delete VDC storage profiles, unexpected error: "+err.Error(),
				)
				return
			}
		}
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

// Delete deletes the resource and removes the Terraform state on success.
func (r *vdcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Delete)()

	state := &vdcModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ! Lock/Unlock is kept while all resources migrate to sdk v2
	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	err := r.vClient.DeleteVDC(ctx, types.ParamsDeleteVDC{
		ID:   state.ID.Get(),
		Name: state.Name.Get(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting VDC", err.Error())
	}
}

func (r *vdcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Import)()

	// Retrieve import ID and save to id attribute
	if urn.IsVDC(req.ID) {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	} else {
		resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	}
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *vdcResource) read(ctx context.Context, planOrState *vdcModel) (stateRefreshed *vdcModel, diags diag.Diagnostics) {
	errGroup, ctxE := errgroup.WithContext(ctx)

	var (
		dataVDC             *types.ModelGetVDC
		dataStorageProfiles *types.ModelListStorageProfiles
	)

	errGroup.Go(func() (err error) {
		dataVDC, err = r.vClient.GetVDC(ctxE, types.ParamsGetVDC{
			ID:   planOrState.ID.Get(),
			Name: planOrState.Name.Get(),
		})
		return err
	})

	errGroup.Go(func() (err error) {
		dataStorageProfiles, err = r.vClient.ListStorageProfile(ctxE, types.ParamsListStorageProfile{
			VdcID:   planOrState.ID.Get(),
			VdcName: planOrState.Name.Get(),
		})
		return err
	})

	if err := errGroup.Wait(); err != nil {
		diags.AddError(
			"Error Reading VDC",
			"Could not read VDC, unexpected error: "+err.Error(),
		)
		return stateRefreshed, diags
	}

	planOrState.fromSDK(ctx, dataVDC, dataStorageProfiles, &diags)

	return planOrState, diags
}
