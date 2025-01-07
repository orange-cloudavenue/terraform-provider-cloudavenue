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
	"strings"

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
	_ resource.Resource                = &BucketResource{}
	_ resource.ResourceWithConfigure   = &BucketResource{}
	_ resource.ResourceWithImportState = &BucketResource{}
)

// NewBucketResource is a helper function to simplify the provider implementation.
func NewBucketResource() resource.Resource {
	return &BucketResource{}
}

// BucketResource is the resource implementation.
type BucketResource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the resource.
func (r *BucketResource) Init(ctx context.Context, rm *BucketModel) (diags diag.Diagnostics) {
	r.s3Client = r.client.CAVSDK.V1.S3()
	return
}

// Metadata returns the resource type name.
func (r *BucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket"
}

// Schema defines the schema for the resource.
func (r *BucketResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = s3BucketSchema(ctx).GetResource(ctx)
}

func (r *BucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_s3_bucket", r.client.GetOrgName(), metrics.Create)()

	plan := &BucketModel{}

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

	if _, err := r.s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket:                     plan.Name.GetPtr(),
		ObjectLockEnabledForBucket: utils.TakeBoolPointer(plan.ObjectLock.Get()),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating bucket", fmt.Sprintf("Error creating bucket (name=%s): %s", plan.Name.Get(), err.Error()))
		return
	}

	if err := r.s3Client.SyncBucket(plan.Name.Get()); err != nil {
		resp.Diagnostics.AddError("Error syncing bucket", fmt.Sprintf("Error syncing bucket (name=%s): %s", plan.Name.Get(), err.Error()))
		return
	}

	// Use generic read function to refresh the state
	state, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *BucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_s3_bucket", r.client.GetOrgName(), metrics.Read)()

	state := &BucketModel{}

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
func (r *BucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_s3_bucket", r.client.GetOrgName(), metrics.Update)()
	// All attributes are immutable
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_s3_bucket", r.client.GetOrgName(), metrics.Delete)()

	state := &BucketModel{}

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

	_, err := r.s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: state.Name.GetPtr(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting bucket", fmt.Sprintf("Error deleting bucket: %s", err.Error()))
		return
	}
}

func (r *BucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_s3_bucket", r.client.GetOrgName(), metrics.Import)()
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *BucketResource) read(ctx context.Context, planOrState *BucketModel) (stateRefreshed *BucketModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	// Get object_lock_enabled
	objectLockConfiguration, err := r.s3Client.GetObjectLockConfigurationWithContext(ctx, &s3.GetObjectLockConfigurationInput{
		Bucket: planOrState.Name.GetPtr(),
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), ErrCodeObjectLockConfigurationNotFoundError):
			stateRefreshed.ObjectLock.Set(false)
		case strings.Contains(err.Error(), s3.ErrCodeNoSuchBucket):
			// Not Found
			return nil, false, nil
		default:
			diags.AddError("Error retrieving bucket", fmt.Sprintf("Error retrieving bucket: %s", err.Error()))
			return nil, true, diags
		}
	} else if objectLockConfiguration.ObjectLockConfiguration != nil {
		stateRefreshed.ObjectLock.Set(true)
	}

	stateRefreshed.Endpoint.Set(fmt.Sprintf("https://%s.s3-region01.cloudavenue.orange-business.com", planOrState.Name.Get()))

	if !stateRefreshed.ID.IsKnown() {
		stateRefreshed.ID.Set(stateRefreshed.Name.Get())
	}

	return stateRefreshed, true, nil
}
