/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketPolicyDataSource{}

const (
	S3BucketPolicyDataSourceName = testsacc.ResourceName("data.cloudavenue_s3_bucket_policy")
)

type S3BucketPolicyDataSource struct{}

func NewS3BucketPolicyDataSourceTest() testsacc.TestACC {
	return &S3BucketPolicyDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketPolicyDataSource) GetResourceName() string {
	return S3BucketPolicyDataSourceName.String()
}

func (r *S3BucketPolicyDataSource) DependenciesConfig() (deps testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	deps.Append(GetResourceConfig()[S3BucketPolicyResourceName]().GetDefaultConfig)
	return deps
}

func (r *S3BucketPolicyDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_bucket_policy" "example" {
						bucket = cloudavenue_s3_bucket.example.name
					}`,
					Checks: GetResourceConfig()[S3BucketPolicyResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccS3BucketPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketPolicyDataSource{}),
	})
}
