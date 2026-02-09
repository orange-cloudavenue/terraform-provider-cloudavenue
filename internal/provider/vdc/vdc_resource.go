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

// Package vdc provides a resource to manage VDCs.
package vdc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi/rules"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &vdcResource{}
	_ resource.ResourceWithConfigure      = &vdcResource{}
	_ resource.ResourceWithImportState    = &vdcResource{}
	_ resource.ResourceWithValidateConfig = &vdcResource{}
	_ resource.ResourceWithModifyPlan     = &vdcResource{}
)

var validationDisabled = os.Getenv(envVarValidation) == "false"

const (
	envVarValidation = "CLOUDAVENUE_VDC_VALIDATION"
)

// NewVDCResource is a helper function to simplify the provider implementation.
func NewVDCResource() resource.Resource {
	return &vdcResource{}
}

// vdcResource is the resource implementation.
type vdcResource struct {
	client *client.CloudAvenue
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

	r.client = client
}

func (r *vdcResource) validateConfig(ctx context.Context, config *vdcResourceModel) (diags diag.Diagnostics) {
	StorageProfiles, d := config.StorageProfiles.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	if err := rules.Validate(rules.ValidateData{
		ServiceClass:        rules.ServiceClass(config.ServiceClass.Get()),
		DisponibilityClass:  rules.DisponibilityClass(config.DisponibilityClass.Get()),
		BillingModel:        rules.BillingModel(config.BillingModel.Get()),
		VCPUInMhz:           config.VCPUInMhz.GetInt(),
		CPUAllocated:        config.CPUAllocated.GetInt(),
		MemoryAllocated:     config.MemoryAllocated.GetInt(),
		StorageBillingModel: rules.BillingModel(config.StorageBillingModel.Get()),
		StorageProfiles: func() map[rules.StorageProfileClass]struct {
			Limit   int
			Default bool
		} {
			storageProfiles := make(map[rules.StorageProfileClass]struct {
				Limit   int
				Default bool
			})
			for _, sP := range StorageProfiles {
				storageProfiles[rules.StorageProfileClass(sP.Class.Get())] = struct {
					Limit   int
					Default bool
				}{Limit: sP.Limit.GetInt(), Default: sP.Default.Get()}
			}
			return storageProfiles
		}(),
	}, false); err != nil {
		switch {
		case errors.Is(err, rules.ErrBillingModelNotAvailable):
			diags.AddAttributeError(path.Root("billing_model"), "Billing model attribute is not valid", err.Error())
		case errors.Is(err, rules.ErrServiceClassNotFound):
			diags.AddAttributeError(path.Root("service_class"), "Service class attribute is not valid", err.Error())
		case errors.Is(err, rules.ErrDisponibilityClassNotFound):
			diags.AddAttributeError(path.Root("disponibility_class"), "Disponibility class attribute is not valid", err.Error())
		case errors.Is(err, rules.ErrStorageBillingModelNotFound):
			diags.AddAttributeError(path.Root("storage_billing_model"), "Storage billing model attribute is not valid", err.Error())
		case errors.Is(err, rules.ErrStorageProfileClassNotFound):
			diags.AddAttributeError(path.Root("storage_profiles"), "Storage profile class attribute is not valid", err.Error())
		case errors.Is(err, rules.ErrCPUAllocatedInvalid):
			diags.AddAttributeError(path.Root("cpu_allocated"), "CPU allocated attribute is not valid", err.Error())
		case errors.Is(err, rules.ErrMemoryAllocatedInvalid):
			diags.AddAttributeError(path.Root("memory_allocated"), "Memory allocated attribute is not valid", err.Error())
		case errors.Is(err, rules.ErrVCPUInMhzInvalid):
			diags.AddAttributeError(path.Root("cpu_speed_in_mhz"), "CPU speed in MHz attribute is not valid", err.Error())
		case errors.Is(err, rules.ErrStorageProfileLimitInvalid) || errors.Is(err, rules.ErrStorageProfileLimitNotIntegrer):
			diags.AddAttributeError(path.Root("storage_profiles").AtName("limit"), "Storage profile limit attribute is not valid", err.Error())
		case errors.Is(err, rules.ErrStorageProfileDefault):
			diags.AddAttributeError(path.Root("storage_profiles").AtName("default"), "Storage profile default attribute is not valid", err.Error())
		default:
			diags.AddError("Error validating VDC", err.Error())
		}
	}

	return diags
}

func (r *vdcResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if validationDisabled {
		return
	}

	config := new(vdcResourceModel)

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the configuration
	resp.Diagnostics.Append(r.validateConfig(ctx, config)...)
}

func (r *vdcResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var (
		plan  = new(vdcResourceModel)
		state = new(vdcResourceModel)
	)

	// Retrieve values from plan
	d := req.Plan.Get(ctx, plan)
	if d.HasError() {
		// If there is an error in the plan, we don't need to continue
		return
	}

	d = req.State.Get(ctx, state)
	// If error in state will be is in create mode
	if !d.HasError() {
		return
	}

	// If there is no error in the state, we are in update mode

	// "Force replacement attributes, however you can change the `cpu_speed_in_mhz` attribute only if the `billing_model` is set to **RESERVED**."
	if plan.VCPUInMhz.Equal(state.VCPUInMhz) && plan.BillingModel.Get() != string(rules.BillingModelReserved) {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("cpu_speed_in_mhz"))
		resp.Diagnostics.AddAttributeWarning(path.Root("cpu_speed_in_mhz"), "CPU speed in MHz attribute require replacement", "You can change the `cpu_speed_in_mhz` attribute only if the `billing_model` is set to **RESERVED**.")
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *vdcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	plan := new(vdcResourceModel)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if validationDisabled {
		// Validate the configuration
		resp.Diagnostics.Append(r.validateConfig(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, errTO := plan.Timeouts.Create(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	body, d := plan.ToCAVVirtualDataCenter(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CAVSDK.V1.VDC().New(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VDC", err.Error())
		return
	}

	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning("Resource not found", fmt.Sprintf("Unable to find resource %s after creation", plan.Name.Get()))
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vdcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Read)()

	// Get current state
	state := new(vdcResourceModel)

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call API to read the resource and test for errors.
	readTimeout, errTO := state.Timeouts.Read(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning("Resource not found", fmt.Sprintf("Unable to find resource %s", state.Name.Get()))
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
func (r *vdcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Update)()
	var (
		plan  = new(vdcResourceModel)
		state = new(vdcResourceModel)
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if validationDisabled {
		// Validate the configuration
		resp.Diagnostics.Append(r.validateConfig(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
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

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	vdc, err := r.client.CAVSDK.V1.VDC().GetVDC(plan.Name.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VDC", err.Error())
		return
	}

	vdc.SetDescription(plan.Description.Get())
	vdc.SetVCPUInMhz(plan.VCPUInMhz.GetInt())
	vdc.SetCPUAllocated(plan.CPUAllocated.GetInt())
	vdc.SetMemoryAllocated(plan.MemoryAllocated.GetInt())

	storageProfiles, d := plan.StorageProfiles.Get(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	vdcStorageProfiles := make([]infrapi.StorageProfile, 0)

	for _, storageProfile := range storageProfiles {
		vdcStorageProfiles = append(vdcStorageProfiles, infrapi.StorageProfile{
			Class:   infrapi.StorageProfileClass(storageProfile.Class.Get()),
			Limit:   storageProfile.Limit.GetInt(),
			Default: storageProfile.Default.Get(),
		})
	}

	vdc.SetStorageProfiles(vdcStorageProfiles)

	if err := vdc.Update(ctx); err != nil {
		resp.Diagnostics.AddError("Error updating VDC", err.Error())
		return
	}

	if !plan.Timeouts.Equal(state.Timeouts) {
		state.Timeouts = plan.Timeouts
	}

	stateRefreshed, _, d := r.read(ctx, state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vdcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Delete)()

	state := new(vdcResourceModel)

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// Update() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	deleteTimeout, errTO := state.Timeouts.Delete(ctx, 5*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	vdc, err := r.client.CAVSDK.V1.VDC().GetVDC(state.Name.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VDC", err.Error())
		return
	}

	if err := vdc.Delete(ctx); err != nil {
		resp.Diagnostics.AddError("Error deleting VDC", err.Error())
		return
	}
}

func (r *vdcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Import)()

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// * Custom Functions.
// read is a generic function to read a resource.
func (r *vdcResource) read(ctx context.Context, planOrState *vdcResourceModel) (stateRefreshed *vdcResourceModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	vdc, err := r.client.CAVSDK.V1.VDC().GetVDC(planOrState.Name.Get())
	if err != nil {
		if govcd.IsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error reading VDC", err.Error())
		return nil, true, diags
	}

	stateRefreshed.ID.Set(vdc.GetID())
	stateRefreshed.Name.Set(vdc.GetName())
	stateRefreshed.Description.Set(vdc.GetDescription())
	stateRefreshed.ServiceClass.Set(string(vdc.GetServiceClass()))
	stateRefreshed.StorageBillingModel.Set(string(vdc.GetStorageBillingModel()))
	stateRefreshed.DisponibilityClass.Set(string(vdc.GetDisponibilityClass()))
	stateRefreshed.BillingModel.Set(string(vdc.GetBillingModel()))
	stateRefreshed.VCPUInMhz.SetInt(vdc.GetVCPUInMhz())
	stateRefreshed.CPUAllocated.SetInt(vdc.GetCPUAllocated())
	stateRefreshed.MemoryAllocated.SetInt(vdc.GetMemoryAllocated())

	storageProfiles := make([]*vdcResourceModelVDCStorageProfile, 0)
	for _, storageProfile := range vdc.GetStorageProfiles() {
		p := new(vdcResourceModelVDCStorageProfile)
		p.Class.Set(string(storageProfile.Class))
		p.Limit.SetInt(storageProfile.Limit)
		p.Default.Set(storageProfile.Default)
		storageProfiles = append(storageProfiles, p)
	}

	diags.Append(stateRefreshed.StorageProfiles.Set(ctx, storageProfiles)...)
	if diags.HasError() {
		return stateRefreshed, found, diags
	}

	return stateRefreshed, true, diags
}
