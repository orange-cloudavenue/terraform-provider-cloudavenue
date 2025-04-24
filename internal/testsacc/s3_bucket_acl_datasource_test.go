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

var _ testsacc.TestACC = &S3BucketACLDataSource{}

const (
	S3BucketACLDataSourceName = testsacc.ResourceName("data.cloudavenue_s3_bucket_acl")
)

type S3BucketACLDataSource struct{}

func NewS3BucketACLDataSourceTest() testsacc.TestACC {
	return &S3BucketACLDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketACLDataSource) GetResourceName() string {
	return S3BucketACLDataSourceName.String()
}

func (r *S3BucketACLDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketACLDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_bucket_acl" "example" {
						bucket = cloudavenue_s3_bucket.example.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.0.permission", "FULL_CONTROL"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.0.grantee.type", "CanonicalUser"),
						resource.TestCheckResourceAttrSet(resourceName, "access_control_policy.grants.0.grantee.id"),
						resource.TestCheckResourceAttrSet(resourceName, "access_control_policy.owner.id"),
					},
				},
			}
		},
	}
}

func TestAccS3BucketACLDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketACLDataSource{}),
	})
}
