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
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	awspolicy "github.com/hashicorp/awspolicyequivalence"
	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type BucketPolicyModel struct {
	Bucket   supertypes.StringValue `tfsdk:"bucket"`
	ID       supertypes.StringValue `tfsdk:"id"`
	Policy   supertypes.StringValue `tfsdk:"policy"`
	Timeouts timeoutsR.Value        `tfsdk:"timeouts"`
}

type BucketPolicyModelDatasource struct {
	Bucket   supertypes.StringValue `tfsdk:"bucket"`
	ID       supertypes.StringValue `tfsdk:"id"`
	Policy   supertypes.StringValue `tfsdk:"policy"`
	Timeouts timeoutsD.Value        `tfsdk:"timeouts"`
}

func (rm *BucketPolicyModel) Copy() any {
	x := &BucketPolicyModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *BucketPolicyModelDatasource) Copy() any {
	x := &BucketPolicyModelDatasource{}
	utils.ModelCopy(rm, x)
	return x
}

// SetID sets the value of the ID field.
func (rm *BucketPolicyModel) SetID(id *string) {
	rm.ID.SetPtr(id)
}

func (rm *BucketPolicyModelDatasource) SetID(id *string) {
	rm.ID.SetPtr(id)
}

func (rm *BucketPolicyModel) GetPolicy(_ context.Context) string {
	return rm.Policy.Get()
}

func (rm *BucketPolicyModelDatasource) GetPolicy(_ context.Context) string {
	return rm.Policy.Get()
}

// SetPolicy sets the value of the Policy field.
func (rm *BucketPolicyModel) SetPolicy(policy string) {
	rm.Policy.Set(policy)
}

// SetPolicy sets the value of the Policy field.
func (rm *BucketPolicyModelDatasource) SetPolicy(policy string) {
	rm.Policy.Set(policy)
}

type readPolicyConfig[T readPolicyResourceDatasource] struct {
	Timeout    func() (time.Duration, diag.Diagnostics)
	Client     *s3.S3
	BucketName func() *string
}

type readPolicyResourceDatasource interface {
	*BucketPolicyModel | *BucketPolicyModelDatasource
	SetID(*string)
	GetPolicy(context.Context) string
	SetPolicy(string)
	Copy() any
}

func SecondJSONUnlessEquivalent(existingPolicy, newPolicy string) (string, error) {
	// valid empty JSON is "{}" not "" so handle special case to avoid
	// Error unmarshaling policy: unexpected end of JSON input
	if strings.TrimSpace(newPolicy) == "" {
		return "", nil
	}

	if strings.TrimSpace(newPolicy) == "{}" {
		return "{}", nil
	}

	if strings.TrimSpace(existingPolicy) == "" || strings.TrimSpace(existingPolicy) == "{}" {
		return newPolicy, nil
	}

	equivalent, err := awspolicy.PoliciesAreEquivalent(existingPolicy, newPolicy)
	if err != nil {
		return "", err
	}

	if equivalent {
		return existingPolicy, nil
	}

	return newPolicy, nil
}

// findBucketPolicy returns the bucket policy if found.
func findBucketPolicy[T readPolicyResourceDatasource](ctx context.Context, config *readPolicyConfig[T]) (string, error) {
	var (
		output *s3.GetBucketPolicyOutput
		err    error
	)

	input := &s3.GetBucketPolicyInput{
		Bucket: config.BucketName(),
	}

	timeout, _ := config.Timeout()

	if _, err := retryWhenAWSErrCodeEquals(ctx, []string{ErrCodeNoSuchBucket, ErrCodeNoSuchBucketPolicy}, &RetryWhenConfig[*s3.GetBucketPolicyOutput]{
		Timeout:  timeout,
		Interval: 15 * time.Second,
		Function: func() (*s3.GetBucketPolicyOutput, error) {
			output, err = config.Client.GetBucketPolicy(input)
			return output, err
		},
	}); err != nil {
		return "", err
	}

	return *output.Policy, nil
}

// PolicyToSet returns the existing policy if the new policy is equivalent.
// Otherwise, it returns the new policy. Either policy is normalized.
func policyToSet(existingPolicy, newPolicy string) (string, error) {
	policyToSet, err := SecondJSONUnlessEquivalent(existingPolicy, newPolicy)
	if err != nil {
		return "", fmt.Errorf("while checking equivalency of existing policy (%s) and new policy (%s), encountered: %w", existingPolicy, newPolicy, err)
	}

	policyToSet, err = structure.NormalizeJsonString(policyToSet)
	if err != nil {
		return "", fmt.Errorf("policy (%s) is invalid JSON: %w", policyToSet, err)
	}

	return policyToSet, nil
}
