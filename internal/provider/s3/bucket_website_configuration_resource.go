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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BucketWebsiteConfigurationResource{}
	_ resource.ResourceWithConfigure   = &BucketWebsiteConfigurationResource{}
	_ resource.ResourceWithImportState = &BucketWebsiteConfigurationResource{}
)

// NewBucketWebsiteConfigurationResource is a helper function to simplify the provider implementation.
func NewBucketWebsiteConfigurationResource() resource.Resource {
	return &BucketWebsiteConfigurationResource{}
}

// BucketWebsiteConfigurationResource is the resource implementation.
type BucketWebsiteConfigurationResource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the resource.
func (r *BucketWebsiteConfigurationResource) Init(_ context.Context, _ *BucketWebsiteConfigurationModel) (diags diag.Diagnostics) {
	r.s3Client = r.client.CAVSDK.V1.S3()

	return diags
}

// Metadata returns the resource type name.
func (r *BucketWebsiteConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_website_configuration"
}

// Schema defines the schema for the resource.
func (r *BucketWebsiteConfigurationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bucketWebsiteConfigurationSchema(ctx).GetResource(ctx)
}

func (r *BucketWebsiteConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BucketWebsiteConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_website_configuration", r.client.GetOrgName(), metrics.Create)()

	plan := &BucketWebsiteConfigurationModel{}

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

	websiteConfig, d := plan.CreateS3WebsiteConfigurationAPIObject(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &s3.PutBucketWebsiteInput{
		Bucket:               plan.Bucket.GetPtr(),
		WebsiteConfiguration: websiteConfig,
	}

	createTimeout, d := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := retryWhenAWSErrCodeEquals(ctx, []string{ErrCodeNoSuchBucket}, &RetryWhenConfig[*s3.PutBucketWebsiteOutput]{
		Timeout:  createTimeout,
		Interval: 15 * time.Second,
		Function: func() (*s3.PutBucketWebsiteOutput, error) {
			return r.s3Client.PutBucketWebsite(input)
		},
	}); err != nil {
		resp.Diagnostics.AddError("Error putting website configuration", err.Error())
		return
	}

	if _, err := retryWhenNotFound(ctx, s3BucketPropagationTimeout, func() (interface{}, error) {
		return findBucketWebsite(ctx, r.s3Client.S3, plan.Bucket.Get())
	}); err != nil {
		resp.Diagnostics.AddError("Waiting for S3 Bucket Website Configuration", err.Error())
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
func (r *BucketWebsiteConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_s3_bucket_website_configuration", r.client.GetOrgName(), metrics.Read)()

	state := &BucketWebsiteConfigurationModel{}

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
func (r *BucketWebsiteConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_website_configuration", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &BucketWebsiteConfigurationModel{}
		state = &BucketWebsiteConfigurationModel{}
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

	websiteConfig, d := plan.CreateS3WebsiteConfigurationAPIObject(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &s3.PutBucketWebsiteInput{
		Bucket:               plan.Bucket.GetPtr(),
		WebsiteConfiguration: websiteConfig,
	}

	updateTimeout, d := plan.Timeouts.Update(ctx, defaultUpdateTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	if _, err := r.s3Client.PutBucketWebsiteWithContext(ctx, input); err != nil {
		resp.Diagnostics.AddError("Error putting website configuration", err.Error())
		return
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
func (r *BucketWebsiteConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_s3_bucket_website_configuration", r.client.GetOrgName(), metrics.Delete)()

	state := &BucketWebsiteConfigurationModel{}

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

	deleteTimeout, d := state.Timeouts.Delete(ctx, defaultUpdateTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	input := &s3.DeleteBucketWebsiteInput{
		Bucket: state.Bucket.GetPtr(),
	}

	if _, err := r.s3Client.DeleteBucketWebsiteWithContext(ctx, input); err != nil {
		if tfawserr.ErrCodeEquals(err, ErrCodeNoSuchBucket, ErrCodeNoSuchWebsiteConfiguration) {
			return
		}

		resp.Diagnostics.AddError("Error deleting website configuration", err.Error())
		return
	}
}

func (r *BucketWebsiteConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_website_configuration", r.client.GetOrgName(), metrics.Import)()

	// * Import basic
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *BucketWebsiteConfigurationResource) read(ctx context.Context, planOrState *BucketWebsiteConfigurationModel) (stateRefreshed *BucketWebsiteConfigurationModel, found bool, diags diag.Diagnostics) {
	return genericReadWebsiteConfiguration(ctx, &readWebsiteConfigurationConfig[*BucketWebsiteConfigurationModel]{
		Client:     r.s3Client.S3,
		BucketName: planOrState.Bucket.GetPtr(),
	}, planOrState)
}

func findBucketWebsite(ctx context.Context, conn *s3.S3, bucket string) (*s3.GetBucketWebsiteOutput, error) {
	input := &s3.GetBucketWebsiteInput{
		Bucket: aws.String(bucket),
	}

	output, err := conn.GetBucketWebsiteWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, ErrCodeNoSuchBucket, ErrCodeNoSuchWebsiteConfiguration) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, fmt.Errorf("unexpected nil output from GetBucketWebsite")
	}

	return output, nil
}
