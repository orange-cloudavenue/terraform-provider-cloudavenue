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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
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
		VdcGroup: plan.VDCGroup.ValueString(),
		Vdc: &apiclient.OrgVdcV2{
			Name:                   plan.Name.ValueString(),
			Description:            plan.Description.ValueString(),
			VdcServiceClass:        plan.VDCServiceClass.ValueString(),
			VdcDisponibilityClass:  plan.VDCDisponibilityClass.ValueString(),
			VdcBillingModel:        plan.VDCBillingModel.ValueString(),
			VcpuInMhz2:             plan.VcpuInMhz2.ValueFloat64(),
			CpuAllocated:           plan.CPUAllocated.ValueFloat64(),
			MemoryAllocated:        plan.MemoryAllocated.ValueFloat64(),
			VdcStorageBillingModel: plan.VDCStorageBillingModel.ValueString(),
		},
	}

	// Iterate over the storage profiles and add them to the body.
	for _, storageProfile := range plan.VDCStorageProfiles {
		body.Vdc.VdcStorageProfiles = append(body.Vdc.VdcStorageProfiles, apiclient.VdcStorageProfilesV2{
			Class:    storageProfile.Class.ValueString(),
			Limit:    int32(storageProfile.Limit.ValueInt64()),
			Default_: storageProfile.Default.ValueBool(),
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
			ID = common.NormalizeID("urn:vcloud:vdc:", v.VdcUuid)
			break
		}
	}
	plan.ID = types.StringValue(ID)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "VDC created")

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vdcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
		if state.Name.ValueString() == v.VdcName {
			ID = common.NormalizeID("urn:vcloud:vdc:", v.VdcUuid)
			break
		}
	}

	// Get storageProfile
	var profiles []vdcStorageProfileModel
	for _, profile := range vdc.Vdc.VdcStorageProfiles {
		p := vdcStorageProfileModel{
			Class:   types.StringValue(profile.Class),
			Limit:   types.Int64Value(int64(profile.Limit)),
			Default: types.BoolValue(profile.Default_),
		}
		profiles = append(profiles, p)
	}

	// Convert from the API data model to the Terraform data model
	// and refresh any attribute values.
	state = &vdcResourceModel{
		Timeouts:               state.Timeouts,
		ID:                     types.StringValue(ID),
		Name:                   types.StringValue(vdc.Vdc.Name),
		Description:            types.StringValue(vdc.Vdc.Description),
		VDCGroup:               types.StringValue(vdc.VdcGroup),
		VDCServiceClass:        types.StringValue(vdc.Vdc.VdcServiceClass),
		VDCDisponibilityClass:  types.StringValue(vdc.Vdc.VdcDisponibilityClass),
		VDCBillingModel:        types.StringValue(vdc.Vdc.VdcBillingModel),
		VcpuInMhz2:             types.Float64Value(vdc.Vdc.VcpuInMhz2),
		CPUAllocated:           types.Float64Value(vdc.Vdc.CpuAllocated),
		MemoryAllocated:        types.Float64Value(vdc.Vdc.MemoryAllocated),
		VDCStorageBillingModel: types.StringValue(vdc.Vdc.VdcStorageBillingModel),
		VDCStorageProfiles:     profiles,
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vdcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *vdcResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

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

	// Convert from Terraform data model into API data model
	body := apiclient.UpdateOrgVdcV2{
		VdcGroup: plan.VDCGroup.ValueString(),
		Vdc: &apiclient.OrgVdcV2{
			Name:                   plan.Name.ValueString(),
			Description:            plan.Description.ValueString(),
			VdcServiceClass:        plan.VDCServiceClass.ValueString(),
			VdcDisponibilityClass:  plan.VDCDisponibilityClass.ValueString(),
			VdcBillingModel:        plan.VDCBillingModel.ValueString(),
			VcpuInMhz2:             plan.VcpuInMhz2.ValueFloat64(),
			CpuAllocated:           plan.CPUAllocated.ValueFloat64(),
			MemoryAllocated:        plan.MemoryAllocated.ValueFloat64(),
			VdcStorageBillingModel: plan.VDCStorageBillingModel.ValueString(),
		},
	}

	// Iterate over the storage profiles and add them to the body.
	for _, storageProfile := range plan.VDCStorageProfiles {
		body.Vdc.VdcStorageProfiles = append(body.Vdc.VdcStorageProfiles, apiclient.VdcStorageProfilesV2{
			Class:    storageProfile.Class.ValueString(),
			Limit:    int32(storageProfile.Limit.ValueInt64()),
			Default_: storageProfile.Default.ValueBool(),
		})
	}

	var err error
	var job apiclient.Jobcreated
	var httpR *http.Response

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

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "VDC updated")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vdcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
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

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "VDC deleted")
}

func (r *vdcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
