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
	_ resource.Resource                = &BucketCorsConfigurationResource{}
	_ resource.ResourceWithConfigure   = &BucketCorsConfigurationResource{}
	_ resource.ResourceWithImportState = &BucketCorsConfigurationResource{}
)

// NewBucketCorsConfigurationResource is a helper function to simplify the provider implementation.
func NewBucketCorsConfigurationResource() resource.Resource {
	return &BucketCorsConfigurationResource{}
}

// BucketCorsConfigurationResource is the resource implementation.
type BucketCorsConfigurationResource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the resource.
func (r *BucketCorsConfigurationResource) Init(_ context.Context, _ *BucketCorsConfigurationModel) (diags diag.Diagnostics) {
	r.s3Client = r.client.CAVSDK.V1.S3()
	return diags
}

// Metadata returns the resource type name.
func (r *BucketCorsConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_cors_configuration"
}

// Schema defines the schema for the resource.
func (r *BucketCorsConfigurationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bucketCorsConfigurationSchema(ctx).GetResource(ctx)
}

func (r *BucketCorsConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BucketCorsConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_cors_configuration", r.client.GetOrgName(), metrics.Create)()

	plan := &BucketCorsConfigurationModel{}

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

	// Set timeouts
	createTimeout, diags := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	// Create the CORS policy
	resp.Diagnostics.Append(r.createOrUpdateCorsPolicy(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	state, found, diags := r.read(ctx, plan)
	if !found {
		diags.AddWarning("Resource not found", "The resource was not found after creation")
		resp.State.RemoveResource(ctx)
		return
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *BucketCorsConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_s3_bucket_cors_configuration", r.client.GetOrgName(), metrics.Read)()

	state := &BucketCorsConfigurationModel{}

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

	// Set timeouts
	readTimeout, diags := state.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	// Refresh the state
	stateRefreshed, found, diags := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *BucketCorsConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_cors_configuration", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &BucketCorsConfigurationModel{}
		state = &BucketCorsConfigurationModel{}
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

	// Set default timeouts
	updateTimeout, diags := plan.Timeouts.Update(ctx, defaultUpdateTimeout)
	diags.Append(diags...)
	if diags.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	// Update the CORS policy
	resp.Diagnostics.Append(r.createOrUpdateCorsPolicy(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, diags := r.read(ctx, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BucketCorsConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_s3_bucket_cors_configuration", r.client.GetOrgName(), metrics.Delete)()

	state := &BucketCorsConfigurationModel{}

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

	// Set timeouts
	deleteTimeout, diags := state.Timeouts.Delete(ctx, defaultDeleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	if _, err := r.s3Client.DeleteBucketCorsWithContext(ctx, &s3.DeleteBucketCorsInput{
		Bucket: state.Bucket.GetPtr(),
	}); err != nil {
		resp.Diagnostics.AddError("Error deleting CORS policy", err.Error())
		return
	}
}

func (r *BucketCorsConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_cors_configuration", r.client.GetOrgName(), metrics.Import)()
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *BucketCorsConfigurationResource) read(ctx context.Context, planOrState *BucketCorsConfigurationModel) (stateRefreshed *BucketCorsConfigurationModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	// Read the CORS
	corsResponse, err := r.s3Client.GetBucketCorsWithContext(ctx, &s3.GetBucketCorsInput{
		Bucket: planOrState.Bucket.GetPtr(),
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), ErrCodeNoSuchCORSConfiguration):
			diags.AddError("CORS policy not found", err.Error())
			return nil, false, diags
		default:
			diags.AddError("Error retrieving CORS policy", err.Error())
			return nil, true, diags
		}
	}
	corsPolicy := new(BucketCorsConfigurationModelCorsRules)
	for _, corsRule := range corsResponse.CORSRules {
		corsRuleModel := NewBucketCorsConfigurationModelCorsRule()

		corsRuleModel.MaxAgeSeconds.SetPtr(corsRule.MaxAgeSeconds)
		corsRuleModel.ID.SetPtr(corsRule.ID)

		corsRuleModel.AllowedMethods.Set(ctx, utils.SlicePointerToSlice(corsRule.AllowedMethods))
		corsRuleModel.AllowedOrigins.Set(ctx, utils.SlicePointerToSlice(corsRule.AllowedOrigins))

		// AllowedHeaders and ExposeHeaders are optional
		if len(corsRule.AllowedHeaders) > 0 {
			corsRuleModel.AllowedHeaders.Set(ctx, utils.SlicePointerToSlice(corsRule.AllowedHeaders))
		}
		if len(corsRule.ExposeHeaders) > 0 {
			corsRuleModel.ExposeHeaders.Set(ctx, utils.SlicePointerToSlice(corsRule.ExposeHeaders))
		}

		*corsPolicy = append(*corsPolicy, corsRuleModel)
	}
	diags.Append(stateRefreshed.CorsRules.Set(ctx, corsPolicy)...)

	if !stateRefreshed.ID.IsKnown() {
		stateRefreshed.ID.Set(planOrState.Bucket.Get())
	}

	return stateRefreshed, true, diags
}

// createOrUpdateCorsPolicy creates or updates the CORS policy.
func (r *BucketCorsConfigurationResource) createOrUpdateCorsPolicy(ctx context.Context, planOrState *BucketCorsConfigurationModel) (diags diag.Diagnostics) {
	corsRules, d := planOrState.GetCorsRules(ctx)
	diags.Append(d...)
	if d.HasError() {
		return diags
	}

	corsRulesS3, d := corsRules.CorsRulesToS3CorsRules(ctx)
	diags.Append(d...)
	if d.HasError() {
		return diags
	}

	corsInput := &s3.PutBucketCorsInput{
		Bucket: planOrState.Bucket.GetPtr(),
		CORSConfiguration: &s3.CORSConfiguration{
			CORSRules: corsRulesS3,
		},
	}

	if _, err := r.s3Client.PutBucketCorsWithContext(ctx, corsInput); err != nil {
		diags.AddError("Error on PUT CORS policy", err.Error())
	}

	return diags
}
