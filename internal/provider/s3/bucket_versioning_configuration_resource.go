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

package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BucketVersioningConfigurationResource{}
	_ resource.ResourceWithConfigure   = &BucketVersioningConfigurationResource{}
	_ resource.ResourceWithImportState = &BucketVersioningConfigurationResource{}
)

// NewBucketVersioningConfigurationResource is a helper function to simplify the provider implementation.
func NewBucketVersioningConfigurationResource() resource.Resource {
	return &BucketVersioningConfigurationResource{}
}

// BucketVersioningConfigurationResource is the resource implementation.
type BucketVersioningConfigurationResource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the resource.
func (r *BucketVersioningConfigurationResource) Init(_ context.Context, _ *BucketVersioningConfigurationModel) (diags diag.Diagnostics) {
	r.s3Client = r.client.CAVSDK.V1.S3()
	return diags
}

// Metadata returns the resource type name.
func (r *BucketVersioningConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_versioning_configuration"
}

// Schema defines the schema for the resource.
func (r *BucketVersioningConfigurationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bucketVersioningConfigurationSchema(ctx).GetResource(ctx)
}

func (r *BucketVersioningConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BucketVersioningConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_versioning_configuration", r.client.GetOrgName(), metrics.Create)()

	plan := &BucketVersioningConfigurationModel{}

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

	if _, err := r.s3Client.PutBucketVersioningWithContext(ctx, &s3.PutBucketVersioningInput{
		Bucket: plan.Bucket.GetPtr(),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: plan.Status.GetPtr(),
		},
	}); err != nil {
		resp.Diagnostics.AddError("Error suspending bucket versioning", err.Error())
	}

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
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
func (r *BucketVersioningConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_s3_bucket_versioning_configuration", r.client.GetOrgName(), metrics.Read)()

	state := &BucketVersioningConfigurationModel{}

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
func (r *BucketVersioningConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_versioning_configuration", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &BucketVersioningConfigurationModel{}
		state = &BucketVersioningConfigurationModel{}
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

	if _, err := r.s3Client.PutBucketVersioningWithContext(ctx, &s3.PutBucketVersioningInput{
		Bucket: state.Bucket.GetPtr(),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: plan.Status.GetPtr(),
		},
	}); err != nil {
		resp.Diagnostics.AddError("Error updating bucket versioning", err.Error())
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
func (r *BucketVersioningConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_s3_bucket_versioning_configuration", r.client.GetOrgName(), metrics.Delete)()

	state := &BucketVersioningConfigurationModel{}

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

	if state.Status.IsKnown() && state.Status.Get() != s3.BucketVersioningStatusSuspended {
		if _, err := r.s3Client.PutBucketVersioningWithContext(ctx, &s3.PutBucketVersioningInput{
			Bucket: state.Bucket.GetPtr(),
			VersioningConfiguration: &s3.VersioningConfiguration{
				Status: utils.TakePointer(s3.BucketVersioningStatusSuspended),
			},
		}); err != nil {
			resp.Diagnostics.AddError("Error suspending bucket versioning", err.Error())
		}
	}
}

func (r *BucketVersioningConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_versioning_configuration", r.client.GetOrgName(), metrics.Import)()

	// * Import basic
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *BucketVersioningConfigurationResource) read(ctx context.Context, planOrState *BucketVersioningConfigurationModel) (stateRefreshed *BucketVersioningConfigurationModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	// Read the versioning configuration
	versioningResponse, err := r.s3Client.GetBucketVersioningWithContext(ctx, &s3.GetBucketVersioningInput{
		Bucket: planOrState.Bucket.GetPtr(),
	})
	if err != nil {
		diags.AddError("Error retrieving bucket versioning configuration", err.Error())
		return stateRefreshed, found, diags
	}

	stateRefreshed.Status.SetPtr(versioningResponse.Status)

	if !stateRefreshed.ID.IsKnown() {
		stateRefreshed.ID.Set(planOrState.Bucket.Get())
	}

	return stateRefreshed, true, nil
}
