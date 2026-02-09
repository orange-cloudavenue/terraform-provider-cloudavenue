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
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BucketLifecycleConfigurationResource{}
	_ resource.ResourceWithConfigure   = &BucketLifecycleConfigurationResource{}
	_ resource.ResourceWithImportState = &BucketLifecycleConfigurationResource{}
)

// NewBucketLifecycleConfigurationResource is a helper function to simplify the provider implementation.
func NewBucketLifecycleConfigurationResource() resource.Resource {
	return &BucketLifecycleConfigurationResource{}
}

// BucketLifecycleConfigurationResource is the resource implementation.
type BucketLifecycleConfigurationResource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the resource.
func (r *BucketLifecycleConfigurationResource) Init(_ context.Context, _ *BucketLifecycleConfigurationModel) (diags diag.Diagnostics) {
	r.s3Client = r.client.CAVSDK.V1.S3()
	return diags
}

// Metadata returns the resource type name.
func (r *BucketLifecycleConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_lifecycle_configuration"
}

// Schema defines the schema for the resource.
func (r *BucketLifecycleConfigurationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bucketLifecycleConfigurationSchema(ctx).GetResource(ctx)
}

func (r *BucketLifecycleConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BucketLifecycleConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_lifecycle_configuration", r.client.GetOrgName(), metrics.Create)()

	plan := &BucketLifecycleConfigurationModel{}

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

	// Set default timeouts
	createTimeout, diags := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	// use generic createOrUpdate function
	resp.Diagnostics.Append(r.createOrUpdateLifeCycle(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	state, _, diags := r.read(ctx, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *BucketLifecycleConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_s3_bucket_lifecycle_configuration", r.client.GetOrgName(), metrics.Read)()

	state := &BucketLifecycleConfigurationModel{}

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
func (r *BucketLifecycleConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_lifecycle_configuration", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &BucketLifecycleConfigurationModel{}
		state = &BucketLifecycleConfigurationModel{}
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

	updateTimeout, diags := plan.Timeouts.Update(ctx, defaultUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	// use generic createOrUpdate function
	resp.Diagnostics.Append(r.createOrUpdateLifeCycle(ctx, plan)...)
	if resp.Diagnostics.HasError() {
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
func (r *BucketLifecycleConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_s3_bucket_lifecycle_configuration", r.client.GetOrgName(), metrics.Delete)()

	state := &BucketLifecycleConfigurationModel{}

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

	deleteTimeout, diags := state.Timeouts.Update(ctx, defaultUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	if _, err := r.s3Client.DeleteBucketLifecycleWithContext(ctx, &s3.DeleteBucketLifecycleInput{
		Bucket: state.Bucket.GetPtr(),
	}); err != nil {
		resp.Diagnostics.AddError("Error deleting S3 Bucket Lifecycle Configuration", err.Error())
	}
}

func (r *BucketLifecycleConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_lifecycle_configuration", r.client.GetOrgName(), metrics.Import)()

	// * Import basic
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *BucketLifecycleConfigurationResource) read(ctx context.Context, planOrState *BucketLifecycleConfigurationModel) (stateRefreshed *BucketLifecycleConfigurationModel, found bool, diags diag.Diagnostics) {
	return genericReadLifeCycleConfiguration(ctx, &readLifeCycleConfigurationConfig[*BucketLifecycleConfigurationModel]{
		Client: r.s3Client.S3,
		Timeout: func() (time.Duration, diag.Diagnostics) {
			return planOrState.Timeouts.Read(ctx, defaultReadTimeout)
		},
		BucketName: func() *string {
			return planOrState.Bucket.GetPtr()
		},
	}, planOrState)
}

// createOrUpdate is a generic create or update function that can be used by the resource Create and Update functions.
func (r *BucketLifecycleConfigurationResource) createOrUpdateLifeCycle(ctx context.Context, planOrState *BucketLifecycleConfigurationModel) (diags diag.Diagnostics) {
	/*
		Implement the resource create or update here
	*/

	lifeCycleRules, d := planOrState.Rules.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	rules := make([]*s3.LifecycleRule, 0)

	for _, rule := range lifeCycleRules {
		s3Rule := &s3.LifecycleRule{
			ID:     rule.ID.GetPtr(),
			Status: rule.Status.GetPtr(),
		}

		// * AbortIncompleteMultipartUpload
		if rule.AbortIncompleteMultipartUpload.IsKnown() {
			abortIncompleteMultipartUpload, d := rule.AbortIncompleteMultipartUpload.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			s3Rule.AbortIncompleteMultipartUpload = &s3.AbortIncompleteMultipartUpload{
				DaysAfterInitiation: abortIncompleteMultipartUpload.DaysAfterInitiation.GetPtr(),
			}
		}

		// * Expiration
		if rule.Expiration.IsKnown() {
			expiration, d := rule.Expiration.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			s3Rule.Expiration = &s3.LifecycleExpiration{}

			if expiration.Days.IsKnown() {
				s3Rule.Expiration.Days = expiration.Days.GetPtr()
			}

			if expiration.ExpiredObjectDeleteMarker.IsKnown() {
				s3Rule.Expiration.ExpiredObjectDeleteMarker = expiration.ExpiredObjectDeleteMarker.GetPtr()
			}

			if expiration.Date.IsKnown() {
				t, err := time.Parse(time.RFC3339, expiration.Date.Get())
				if err != nil {
					diags.AddError("Error parsing S3 Bucket Lifecycle Rule Expiration date", err.Error())
					return diags
				}
				s3Rule.Expiration.Date = utils.TakePointer(t)
			}
		}

		// * Filter
		if rule.Filter.IsKnown() {
			filter, d := rule.Filter.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			s3Rule.Filter = &s3.LifecycleRuleFilter{
				Prefix: filter.Prefix.GetPtr(),
			}

			// ? Tag
			if filter.Tag.IsKnown() {
				tag, d := filter.Tag.Get(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				s3Rule.Filter.Tag = &s3.Tag{
					Key:   tag.Key.GetPtr(),
					Value: tag.Value.GetPtr(),
				}
			}

			// ? And
			if !filter.And.IsNull() {
				and, d := filter.And.Get(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				s3Rule.Filter.And = &s3.LifecycleRuleAndOperator{
					Prefix: and.Prefix.GetPtr(),
				}

				// ? And.Tags
				if and.Tags.IsKnown() {
					tags, d := and.Tags.Get(ctx)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					s3Rule.Filter.And.Tags = make([]*s3.Tag, 0)

					for _, t := range tags {
						s3Rule.Filter.And.Tags = append(s3Rule.Filter.And.Tags, &s3.Tag{
							Key:   t.Key.GetPtr(),
							Value: t.Value.GetPtr(),
						})
					}
				}
			}
		}

		// * NoncurrentVersionExpiration
		if rule.NoncurrentVersionExpiration.IsKnown() {
			noncurrentVersionExpiration, d := rule.NoncurrentVersionExpiration.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			s3Rule.NoncurrentVersionExpiration = &s3.NoncurrentVersionExpiration{
				NoncurrentDays:          noncurrentVersionExpiration.NoncurrentDays.GetPtr(),
				NewerNoncurrentVersions: noncurrentVersionExpiration.NewerNoncurrentVersions.GetPtr(),
			}
		}

		rules = append(rules, s3Rule)
	}

	input := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: planOrState.Bucket.GetPtr(),
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: rules,
		},
	}

	createTimeout, d := planOrState.Timeouts.Create(ctx, defaultCreateTimeout)
	if d.HasError() {
		diags.Append(d...)
		return diags
	}

	if _, err := retryWhenAWSErrCodeEquals(ctx, []string{s3.ErrCodeNoSuchBucket}, &RetryWhenConfig[*s3.PutBucketLifecycleConfigurationOutput]{
		Timeout:  createTimeout,
		Interval: 15 * time.Second,
		Function: func() (*s3.PutBucketLifecycleConfigurationOutput, error) {
			return r.s3Client.PutBucketLifecycleConfigurationWithContext(ctx, input)
		},
	}); err != nil {
		diags.AddError("Error putting S3 Bucket Lifecycle Configuration", err.Error())
		return diags
	}

	if err := waitForLifecycleConfigurationRulesStatus(ctx, r.s3Client.S3, planOrState.Bucket.Get(), rules); err != nil {
		diags.AddError("Error waiting for S3 Lifecycle Configuration for bucket to reach expected rules status", err.Error())
	}

	return diags
}

const (
	lifecycleConfigurationExtraRetryDelay         = 5 * time.Second
	lifecycleConfigurationRulesPropagationTimeout = 10 * time.Minute
	lifecycleConfigurationRulesSteadyTimeout      = 2 * time.Minute

	// LifecycleConfigurationRulesStatusReady occurs when all configured rules reach their desired state (Enabled or Disabled).
	LifecycleConfigurationRulesStatusReady = "READY"
	// LifecycleConfigurationRulesStatusNotReady occurs when all configured rules have not reached their desired state (Enabled or Disabled).
	LifecycleConfigurationRulesStatusNotReady = "NOT_READY"
)

func waitForLifecycleConfigurationRulesStatus(ctx context.Context, conn *s3.S3, bucket string, rules []*s3.LifecycleRule) error {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{"", LifecycleConfigurationRulesStatusNotReady},
		Target:                    []string{LifecycleConfigurationRulesStatusReady},
		Refresh:                   lifecycleConfigurationRulesStatus(ctx, conn, bucket, rules),
		Timeout:                   lifecycleConfigurationRulesPropagationTimeout,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 3,
		NotFoundChecks:            20,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func lifecycleConfigurationRulesStatus(ctx context.Context, conn *s3.S3, bucket string, rules []*s3.LifecycleRule) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucket),
		}

		output, err := conn.GetBucketLifecycleConfigurationWithContext(ctx, input)

		if tfawserr.ErrCodeEquals(err, ErrCodeNoSuchLifecycleConfiguration, s3.ErrCodeNoSuchBucket) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		if output == nil {
			return nil, "", &retry.NotFoundError{
				Message:     "Empty result",
				LastRequest: input,
			}
		}

		for _, expectedRule := range rules {
			found := false

			for _, actualRule := range output.Rules {
				if aws.StringValue(actualRule.ID) != aws.StringValue(expectedRule.ID) {
					continue
				}
				found = true
				if aws.StringValue(actualRule.Status) != aws.StringValue(expectedRule.Status) {
					return output, LifecycleConfigurationRulesStatusNotReady, nil
				}
			}

			if !found {
				return output, LifecycleConfigurationRulesStatusNotReady, nil
			}
		}

		return output, LifecycleConfigurationRulesStatusReady, nil
	}
}
