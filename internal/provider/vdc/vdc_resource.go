// Package vdc provides a resource to manage VDCs.
package vdc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	apiclient "github.com/orange-cloudavenue/infrapi-sdk-go"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
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
func (r *vdcResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vdcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	var plan *vdcResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

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

	ctxTO, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Prepare the body to create a VDC.
	body := apiclient.CreateOrgVdcV2{
		Vdc: &apiclient.OrgVdcV2{
			Name:                   plan.Name.Get(),
			Description:            plan.Description.Get(),
			VdcServiceClass:        plan.VDCServiceClass.Get(),
			VdcDisponibilityClass:  plan.VDCDisponibilityClass.Get(),
			VdcBillingModel:        plan.VDCBillingModel.Get(),
			VcpuInMhz2:             float64(plan.VcpuInMhz2.Get()),
			CpuAllocated:           float64(plan.CPUAllocated.Get()),
			MemoryAllocated:        float64(plan.MemoryAllocated.Get()),
			VdcStorageBillingModel: plan.VDCStorageBillingModel.Get(),
		},
	}

	// Get the storage profiles
	storageProfiles, d := plan.GetVDCStorageProfiles(ctx)
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	// Iterate over the storage profiles and add them to the body.
	for _, storageProfile := range storageProfiles {
		body.Vdc.VdcStorageProfiles = append(body.Vdc.VdcStorageProfiles, apiclient.VdcStorageProfilesV2{
			Class:    storageProfile.Class.Get(),
			Limit:    storageProfile.Limit.GetInt32(),
			Default_: storageProfile.Default.Get(),
		})
	}

	var err error
	var job apiclient.Jobcreated
	var httpR *http.Response

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// Call API to create the resource and test for errors.
	job, httpR, err = r.client.APIClient.VDCApi.CreateOrgVdc(auth, body)

	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	// Wait for job to complete
	errRetry := retry.RetryContext(ctxTO, createTimeout, func() *retry.RetryError {
		jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
		if errGetJob != nil {
			retry.NonRetryableError(errGetJob)
		}
		if !slices.Contains(helpers.JobStateDone(), jobStatus.String()) {
			return retry.RetryableError(fmt.Errorf("expected job done but was %s - %w", jobStatus, errGetJob))
		}

		return nil
	})

	if errRetry != nil {
		resp.Diagnostics.AddError("Error waiting job to complete", errRetry.Error())
		return
	}

	// Get vDC UUID by parsing vDCs list and set URN ID
	var ID string
	vdcs, httpR, err := r.client.APIClient.VDCApi.GetOrgVdcs(auth)
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vdcs detail, got error: %s", err))
		return
	}

	for _, v := range vdcs {
		if plan.Name.ValueString() == v.VdcName {
			ID = uuid.Normalize(uuid.VDC, v.VdcUuid).String()
			break
		}
	}

	plan.ID.Set(ID)

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vdcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Read)()

	// Get current state
	var state *vdcResourceModel

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

	ctxTO, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Get vDC info
	vdc, httpR, err := r.client.APIClient.VDCApi.GetOrgVdcByName(auth, state.Name.ValueString())

	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		if !apiErr.IsNotFound() {
			resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
			return
		}
		// 404, resource was not found, remove it from state
		resp.State.RemoveResource(ctx)

		return
	}

	// Get vDC UUID by parsing vDCs list and set URN ID
	var ID string
	vdcs, httpR, err := r.client.APIClient.VDCApi.GetOrgVdcs(auth)

	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vdcs detail, got error: %s", err))
		return
	}

	for _, v := range vdcs {
		if state.Name.Get() == v.VdcName {
			ID = uuid.Normalize(uuid.VDC, v.VdcUuid).String()
			break
		}
	}

	// Get storageProfile
	profiles := make(vdcResourceModelVDCStorageProfiles, 0)
	for _, profile := range vdc.Vdc.VdcStorageProfiles {
		p := vdcResourceModelVDCStorageProfile{}
		p.Class.Set(profile.Class)
		p.Limit.SetInt32(profile.Limit)
		p.Default.Set(profile.Default_)
		profiles = append(profiles, p)
	}

	// Convert from the API data model to the Terraform data model
	// and refresh any attribute values.
	state.ID.Set(ID)
	state.Name.Set(vdc.Vdc.Name)
	state.Description.Set(vdc.Vdc.Description)
	state.VDCServiceClass.Set(vdc.Vdc.VdcServiceClass)
	state.VDCDisponibilityClass.Set(vdc.Vdc.VdcDisponibilityClass)
	state.VDCBillingModel.Set(vdc.Vdc.VdcBillingModel)
	state.VcpuInMhz2.Set(int64(vdc.Vdc.VcpuInMhz2))
	state.CPUAllocated.Set(int64(vdc.Vdc.CpuAllocated))
	state.MemoryAllocated.Set(int64(vdc.Vdc.MemoryAllocated))
	state.VDCStorageBillingModel.Set(vdc.Vdc.VdcStorageBillingModel)
	resp.Diagnostics.Append(state.VDCStorageProfiles.Set(ctx, profiles)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vdcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Update)()
	var (
		plan  *vdcResourceModel
		state *vdcResourceModel
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	ctxTO, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Get vDC info
	var httpR *http.Response
	var err error
	// Due a bug in CloudAvenue the field VdcGroup is mandatory in the body
	vdc, httpR, err := r.client.APIClient.VDCApi.GetOrgVdcByName(auth, state.Name.Get())
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	// Convert from Terraform data model into API data model
	body := apiclient.UpdateOrgVdcV2{
		VdcGroup: vdc.VdcGroup,
		Vdc: &apiclient.OrgVdcV2{
			Name:                   plan.Name.Get(),
			Description:            plan.Description.Get(),
			VdcServiceClass:        plan.VDCServiceClass.Get(),
			VdcDisponibilityClass:  plan.VDCDisponibilityClass.Get(),
			VdcBillingModel:        plan.VDCBillingModel.Get(),
			VcpuInMhz2:             float64(plan.VcpuInMhz2.Get()),
			CpuAllocated:           float64(plan.CPUAllocated.Get()),
			MemoryAllocated:        float64(plan.MemoryAllocated.Get()),
			VdcStorageBillingModel: plan.VDCStorageBillingModel.Get(),
		},
	}

	// Get the storage profiles
	storageProfiles, d := plan.GetVDCStorageProfiles(ctx)
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}

	// Iterate over the storage profiles and add them to the body.
	for _, storageProfile := range storageProfiles {
		body.Vdc.VdcStorageProfiles = append(body.Vdc.VdcStorageProfiles, apiclient.VdcStorageProfilesV2{
			Class:    storageProfile.Class.Get(),
			Limit:    storageProfile.Limit.GetInt32(),
			Default_: storageProfile.Default.Get(),
		})
	}

	var job apiclient.Jobcreated

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// Call API to update the resource and test for errors.
	job, httpR, err = r.client.APIClient.VDCApi.UpdateOrgVdc(auth, body, body.Vdc.Name)
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	// Wait for job to complete
	errRetry := retry.RetryContext(ctxTO, updateTimeout, func() *retry.RetryError {
		jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
		if errGetJob != nil {
			retry.NonRetryableError(err)
		}
		if !slices.Contains(helpers.JobStateDone(), jobStatus.String()) {
			return retry.RetryableError(fmt.Errorf("expected job done but was %s", jobStatus))
		}

		return nil
	})

	if errRetry != nil {
		resp.Diagnostics.AddError("Error waiting job to complete", errRetry.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vdcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Delete)()

	var state *vdcResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	deleteTimeout, errTO := state.Timeouts.Delete(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctxTO, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// Delete the VDC
	job, httpR, err := r.client.APIClient.VDCApi.DeleteOrgVdc(auth, state.Name.ValueString())

	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	// Wait for job to complete
	errRetry := retry.RetryContext(ctxTO, deleteTimeout, func() *retry.RetryError {
		jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
		if errGetJob != nil {
			retry.NonRetryableError(err)
		}
		if !slices.Contains(helpers.JobStateDone(), jobStatus.String()) {
			return retry.RetryableError(fmt.Errorf("expected job done but was %s", jobStatus))
		}

		return nil
	})

	if errRetry != nil {
		resp.Diagnostics.AddError("Error waiting job to complete", errRetry.Error())
		return
	}
}

func (r *vdcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdc", r.client.GetOrgName(), metrics.Import)()

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
