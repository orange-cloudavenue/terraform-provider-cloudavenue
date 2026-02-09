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
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	BucketLifecycleConfigurationDatasourceModel struct {
		Timeouts timeoutsD.Value                                                           `tfsdk:"timeouts"`
		Bucket   supertypes.StringValue                                                    `tfsdk:"bucket"`
		ID       supertypes.StringValue                                                    `tfsdk:"id"`
		Rules    supertypes.ListNestedObjectValueOf[BucketLifecycleConfigurationModelRule] `tfsdk:"rules"`
	}
	BucketLifecycleConfigurationModel struct {
		Timeouts timeoutsR.Value                                                           `tfsdk:"timeouts"`
		Bucket   supertypes.StringValue                                                    `tfsdk:"bucket"`
		ID       supertypes.StringValue                                                    `tfsdk:"id"`
		Rules    supertypes.ListNestedObjectValueOf[BucketLifecycleConfigurationModelRule] `tfsdk:"rules"`
	}

	BucketLifecycleConfigurationModelRules []*BucketLifecycleConfigurationModelRule
	BucketLifecycleConfigurationModelRule  struct {
		AbortIncompleteMultipartUpload supertypes.SingleNestedObjectValueOf[BucketLifecycleConfigurationModelAbortIncompleteMultipartUpload] `tfsdk:"abort_incomplete_multipart_upload"`
		Expiration                     supertypes.SingleNestedObjectValueOf[BucketLifecycleConfigurationModelExpiration]                     `tfsdk:"expiration"`
		Filter                         supertypes.SingleNestedObjectValueOf[BucketLifecycleConfigurationModelFilter]                         `tfsdk:"filter"`
		ID                             supertypes.StringValue                                                                                `tfsdk:"id"`
		NoncurrentVersionExpiration    supertypes.SingleNestedObjectValueOf[BucketLifecycleConfigurationModelNoncurrentVersionExpiration]    `tfsdk:"noncurrent_version_expiration"`
		Status                         supertypes.StringValue                                                                                `tfsdk:"status"`
	}
	BucketLifecycleConfigurationModelAbortIncompleteMultipartUpload struct {
		DaysAfterInitiation supertypes.Int64Value `tfsdk:"days_after_initiation"`
	}
	BucketLifecycleConfigurationModelExpiration struct {
		Date                      supertypes.StringValue `tfsdk:"date"`
		Days                      supertypes.Int64Value  `tfsdk:"days"`
		ExpiredObjectDeleteMarker supertypes.BoolValue   `tfsdk:"expired_object_delete_marker"`
	}
	BucketLifecycleConfigurationModelFilter struct {
		And    supertypes.SingleNestedObjectValueOf[BucketLifecycleConfigurationModelAnd] `tfsdk:"and"`
		Tag    supertypes.SingleNestedObjectValueOf[BucketLifecycleConfigurationModelTag] `tfsdk:"tag"`
		Prefix supertypes.StringValue                                                     `tfsdk:"prefix"`
	}
	BucketLifecycleConfigurationModelAnd struct {
		Tags   supertypes.ListNestedObjectValueOf[BucketLifecycleConfigurationModelTag] `tfsdk:"tags"`
		Prefix supertypes.StringValue                                                   `tfsdk:"prefix"`
	}
	BucketLifecycleConfigurationModelTag struct {
		Key   supertypes.StringValue `tfsdk:"key"`
		Value supertypes.StringValue `tfsdk:"value"`
	}
	BucketLifecycleConfigurationModelNoncurrentVersionExpiration struct {
		NewerNoncurrentVersions supertypes.Int64Value `tfsdk:"newer_noncurrent_versions"`
		NoncurrentDays          supertypes.Int64Value `tfsdk:"noncurrent_days"`
	}
)

// ToBucketLifecycleConfigurationModel converts a slice of *s3.GetBucketLifecycleConfigurationOutput to a slice of BucketLifecycleConfigurationModelRule.
func ToBucketLifecycleConfigurationModel(ctx context.Context, lifecycle *s3.GetBucketLifecycleConfigurationOutput) (values BucketLifecycleConfigurationModelRules, diags diag.Diagnostics) {
	rulesRefreshed := BucketLifecycleConfigurationModelRules{}

	if len(lifecycle.Rules) > 0 {
		for _, rule := range lifecycle.Rules {
			ruleRefreshed := &BucketLifecycleConfigurationModelRule{
				ID:                             supertypes.NewStringNull(),
				Status:                         supertypes.NewStringNull(),
				AbortIncompleteMultipartUpload: supertypes.NewSingleNestedObjectValueOfNull[BucketLifecycleConfigurationModelAbortIncompleteMultipartUpload](ctx),
				Expiration:                     supertypes.NewSingleNestedObjectValueOfNull[BucketLifecycleConfigurationModelExpiration](ctx),
				Filter:                         supertypes.NewSingleNestedObjectValueOfNull[BucketLifecycleConfigurationModelFilter](ctx),
				NoncurrentVersionExpiration:    supertypes.NewSingleNestedObjectValueOfNull[BucketLifecycleConfigurationModelNoncurrentVersionExpiration](ctx),
			}

			ruleRefreshed.ID.SetPtr(rule.ID)
			ruleRefreshed.Status.SetPtr(rule.Status)

			if rule.AbortIncompleteMultipartUpload != nil {
				abortIncompleteMultipartUpload := &BucketLifecycleConfigurationModelAbortIncompleteMultipartUpload{
					DaysAfterInitiation: supertypes.NewInt64Null(),
				}
				abortIncompleteMultipartUpload.DaysAfterInitiation.SetPtr(rule.AbortIncompleteMultipartUpload.DaysAfterInitiation)
				diags.Append(ruleRefreshed.AbortIncompleteMultipartUpload.Set(ctx, abortIncompleteMultipartUpload)...)
			}

			if rule.Expiration != nil {
				expiration := &BucketLifecycleConfigurationModelExpiration{
					Date:                      supertypes.NewStringNull(),
					Days:                      supertypes.NewInt64Null(),
					ExpiredObjectDeleteMarker: supertypes.NewBoolNull(),
				}
				if rule.Expiration.Date != nil {
					expiration.Date.Set(rule.Expiration.Date.Format(time.RFC3339))
				}
				expiration.Days.SetPtr(rule.Expiration.Days)
				expiration.ExpiredObjectDeleteMarker.SetPtr(rule.Expiration.ExpiredObjectDeleteMarker)

				diags.Append(ruleRefreshed.Expiration.Set(ctx, expiration)...)
			}

			if rule.Filter != nil {
				filter := &BucketLifecycleConfigurationModelFilter{
					And:    supertypes.NewSingleNestedObjectValueOfNull[BucketLifecycleConfigurationModelAnd](ctx),
					Tag:    supertypes.NewSingleNestedObjectValueOfNull[BucketLifecycleConfigurationModelTag](ctx),
					Prefix: supertypes.NewStringNull(),
				}
				filter.Prefix.SetPtr(rule.Filter.Prefix)

				if rule.Filter.Tag != nil {
					tag := &BucketLifecycleConfigurationModelTag{
						Key:   supertypes.NewStringNull(),
						Value: supertypes.NewStringNull(),
					}

					tag.Key.SetPtr(rule.Filter.Tag.Key)
					tag.Value.SetPtr(rule.Filter.Tag.Value)
					diags.Append(filter.Tag.Set(ctx, tag)...)
				}

				if rule.Filter.And != nil {
					and := &BucketLifecycleConfigurationModelAnd{
						Tags:   supertypes.NewListNestedObjectValueOfNull[BucketLifecycleConfigurationModelTag](ctx),
						Prefix: supertypes.NewStringNull(),
					}
					and.Prefix.SetPtr(rule.Filter.And.Prefix)
					tags := make([]*BucketLifecycleConfigurationModelTag, 0)
					if len(rule.Filter.And.Tags) > 0 {
						for _, tag := range rule.Filter.And.Tags {
							tagRefreshed := &BucketLifecycleConfigurationModelTag{
								Key:   supertypes.NewStringNull(),
								Value: supertypes.NewStringNull(),
							}

							tagRefreshed.Key.SetPtr(tag.Key)
							tagRefreshed.Value.SetPtr(tag.Value)
							tags = append(tags, tagRefreshed)
						}
					}
					diags.Append(and.Tags.Set(ctx, tags)...)
					diags.Append(filter.And.Set(ctx, and)...)
				}
				diags.Append(ruleRefreshed.Filter.Set(ctx, filter)...)
			}

			if rule.NoncurrentVersionExpiration != nil {
				noncurrentVersionExpiration := &BucketLifecycleConfigurationModelNoncurrentVersionExpiration{
					NewerNoncurrentVersions: supertypes.NewInt64Null(),
					NoncurrentDays:          supertypes.NewInt64Null(),
				}
				noncurrentVersionExpiration.NewerNoncurrentVersions.SetPtr(rule.NoncurrentVersionExpiration.NewerNoncurrentVersions)
				noncurrentVersionExpiration.NoncurrentDays.SetPtr(rule.NoncurrentVersionExpiration.NoncurrentDays)

				diags.Append(ruleRefreshed.NoncurrentVersionExpiration.Set(ctx, noncurrentVersionExpiration)...)
			}

			if diags.HasError() {
				return rulesRefreshed, diags
			}

			rulesRefreshed = append(rulesRefreshed, ruleRefreshed)
		}
	}

	return rulesRefreshed, diags
}

type readLifeCycleConfigurationResourceDatasource interface {
	*BucketLifecycleConfigurationModel | *BucketLifecycleConfigurationDatasourceModel
	SetID(*string)
	SetRules(context.Context, BucketLifecycleConfigurationModelRules) diag.Diagnostics
	Copy() any
}

type readLifeCycleConfigurationConfig[T readLifeCycleConfigurationResourceDatasource] struct {
	Timeout    func() (time.Duration, diag.Diagnostics)
	Client     *s3.S3
	BucketName func() *string
}

// genericReadLifeCycleConfiguration is a generic function that reads the lifecycle configuration of a bucket.
func genericReadLifeCycleConfiguration[T readLifeCycleConfigurationResourceDatasource](ctx context.Context, config *readLifeCycleConfigurationConfig[T], planOrState T) (stateRefreshed T, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy().(T)

	readTimeout, d := config.Timeout()
	if d.HasError() {
		diags.Append(d...)
		return stateRefreshed, found, diags
	}

	var lastOutput, lifecycle *s3.GetBucketLifecycleConfigurationOutput

	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		var err error

		time.Sleep(lifecycleConfigurationExtraRetryDelay)

		lifecycle, err = config.Client.GetBucketLifecycleConfigurationWithContext(ctx, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: config.BucketName(),
		})

		if tfawserr.ErrCodeEquals(err, ErrCodeNoSuchLifecycleConfiguration, s3.ErrCodeNoSuchBucket) {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		if lastOutput == nil || !reflect.DeepEqual(*lastOutput, *lifecycle) {
			lastOutput = lifecycle
			return retry.RetryableError(fmt.Errorf("bucket lifecycle configuration has not stablized; trying again"))
		}

		return nil
	})

	if err != nil && !tfawserr.ErrMessageContains(err, ErrCodeNoSuchLifecycleConfiguration, "") {
		diags.AddError("Error retrieving bucket lifecycle configuration", err.Error())
		return stateRefreshed, false, diags
	}

	rulesRefreshed, d := ToBucketLifecycleConfigurationModel(ctx, lifecycle)
	if d.HasError() {
		diags.Append(d...)
		return stateRefreshed, true, diags
	}

	if len(rulesRefreshed) == 0 {
		diags.AddError("Error retrieving bucket lifecycle configuration", "no rules found")
		return stateRefreshed, false, diags
	}

	stateRefreshed.SetID(config.BucketName())
	diags.Append(stateRefreshed.SetRules(ctx, rulesRefreshed)...)

	return stateRefreshed, true, diags
}

// BucketLifecycleConfigurationModel SetID.
func (rm *BucketLifecycleConfigurationModel) SetID(id *string) {
	rm.ID.SetPtr(id)
}

// BucketLifecycleConfigurationModel SetRules.
func (rm *BucketLifecycleConfigurationModel) SetRules(ctx context.Context, rules BucketLifecycleConfigurationModelRules) (diags diag.Diagnostics) {
	return rm.Rules.Set(ctx, rules)
}

// BucketLifecycleConfigurationDatasourceModel SetID.
func (rm *BucketLifecycleConfigurationDatasourceModel) SetID(id *string) {
	rm.ID.SetPtr(id)
}

// BucketLifecycleConfigurationDatasourceModel SetRules.
func (rm *BucketLifecycleConfigurationDatasourceModel) SetRules(ctx context.Context, rules BucketLifecycleConfigurationModelRules) (diags diag.Diagnostics) {
	return rm.Rules.Set(ctx, rules)
}

func (rm *BucketLifecycleConfigurationModel) Copy() any {
	x := &BucketLifecycleConfigurationModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *BucketLifecycleConfigurationDatasourceModel) Copy() any {
	x := &BucketLifecycleConfigurationDatasourceModel{}
	utils.ModelCopy(rm, x)
	return x
}
