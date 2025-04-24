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

var _ testsacc.TestACC = &S3BucketPolicyResource{}

const (
	S3BucketPolicyResourceName = testsacc.ResourceName("cloudavenue_s3_bucket_policy")
)

type S3BucketPolicyResource struct{}

func NewS3BucketPolicyResourceTest() testsacc.TestACC {
	return &S3BucketPolicyResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketPolicyResource) GetResourceName() string {
	return S3BucketPolicyResourceName.String()
}

func (r *S3BucketPolicyResource) DependenciesConfig() (deps testsacc.DependenciesConfigResponse) {
	deps.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketPolicyResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_s3_bucket_policy" "example" {
					bucket = cloudavenue_s3_bucket.example.name
					policy = jsonencode({
						Version = "2012-10-17"
						Statement = [
						  {
						    Effect = "Allow"
						    Principal = "*"
						    Action = [
						      "s3:*"
						    ]
						    Resource = [
						      "arn:aws:s3:::${cloudavenue_s3_bucket.example.name}/*"
						    ]
						  }
						]
					})
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "policy"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_policy" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							policy = jsonencode({
								Version = "2012-10-17"
								Statement = [
								  {
									Effect = "Allow"
									Principal = "*"
									Action = [
									  "s3:*"
									]
									Resource = [
									  "arn:aws:s3:::${cloudavenue_s3_bucket.example.name}/*"
									]
								  }
								]
							  })
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "policy"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_policy" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							policy = jsonencode({
								Version = "2012-10-17"
								Statement = [
								  {
									Effect = "Allow"
									Principal = "*"
									Action = [ 
									  "s3:DeleteBucket",
									  "s3:GetObject",
									  "s3:ListBucketVersions",
									]
									Resource = [
									  "arn:aws:s3:::${cloudavenue_s3_bucket.example.name}",
									  "arn:aws:s3:::${cloudavenue_s3_bucket.example.name}/*",
									]
								  }
								]
							  })
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "policy"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"bucket"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
				// ! Destroy testing
				Destroy: true,
			}
		},
	}
}

func TestAccS3BucketPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketPolicyResource{}),
	})
}
