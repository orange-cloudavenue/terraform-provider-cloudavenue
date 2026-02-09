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

package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketACLResource{}

const (
	S3BucketACLResourceName = testsacc.ResourceName("cloudavenue_s3_bucket_acl")
)

type S3BucketACLResource struct{}

func NewS3BucketACLResourceTest() testsacc.TestACC {
	return &S3BucketACLResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketACLResource) GetResourceName() string {
	return S3BucketACLResourceName.String()
}

func (r *S3BucketACLResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return resp
}

func (r *S3BucketACLResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example" with an ACL canned policy
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
							resource "cloudavenue_s3_bucket_acl" "example" {
								bucket = cloudavenue_s3_bucket.example.name
								acl = "public-read"
							}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.1.permission", "READ"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.1.grantee.uri", "http://acs.amazonaws.com/groups/global/AllUsers"),
						resource.TestCheckResourceAttrSet(resourceName, "access_control_policy.owner.id"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
							resource "cloudavenue_s3_bucket_acl" "example" {
								bucket = cloudavenue_s3_bucket.example.name
								acl = "public-read"
							}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.1.permission", "READ"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.1.grantee.uri", "http://acs.amazonaws.com/groups/global/AllUsers"),
							resource.TestCheckResourceAttrSet(resourceName, "access_control_policy.owner.id"),
						},
					},
					{
						TFConfig: `
							resource "cloudavenue_s3_bucket_acl" "example" {
								bucket = cloudavenue_s3_bucket.example.name
								acl = "private"
							}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "acl", "private"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.0.permission", "FULL_CONTROL"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.0.grantee.type", "CanonicalUser"),
							resource.TestCheckResourceAttrSet(resourceName, "access_control_policy.owner.id"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder:    []string{"bucket"},
						ImportStateVerifyIgnore: []string{"acl"},
						ImportState:             true,
						ImportStateVerify:       true,
					},
				},
				// ! Destroy
				Destroy: true,
			}
		},
		// *  Second test named "example_with_custom_policy"
		"example_with_custom_policy": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
						resource "cloudavenue_s3_bucket_acl" "example_with_custom_policy" {
							bucket = cloudavenue_s3_bucket.example.name
							access_control_policy = {
              					grants = [{
                    	  			grantee    = {
                    	      			type = "Group"
                    	      			uri  = "http://acs.amazonaws.com/groups/global/AllUsers"
                    	    		},
                    	  			permission = "READ"
                    			}]
              					owner  = {
                  					id           = "1f680dc8d84ca778f885628a39e1980850d408dbc3ca3a706bc182a9672f95ce"
                				}
            				}
						}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.0.permission", "READ"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.0.grantee.uri", "http://acs.amazonaws.com/groups/global/AllUsers"),
						resource.TestCheckResourceAttrSet(resourceName, "access_control_policy.owner.id"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
							resource "cloudavenue_s3_bucket_acl" "example_with_custom_policy" {
								bucket = cloudavenue_s3_bucket.example.name
								access_control_policy = {
									grants = [{
										grantee    = {
											type = "Group"
											uri  = "http://acs.amazonaws.com/groups/global/AllUsers"
										},
										permission = "READ"
									}]
									owner  = {
										id           = "1f680dc8d84ca778f885628a39e1980850d408dbc3ca3a706bc182a9672f95ce"
									}
								}
							}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.0.permission", "READ"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.0.grantee.uri", "http://acs.amazonaws.com/groups/global/AllUsers"),
							resource.TestCheckResourceAttrSet(resourceName, "access_control_policy.owner.id"),
						},
					},
					{
						TFConfig: `
							resource "cloudavenue_s3_bucket_acl" "example_with_custom_policy" {
								bucket = cloudavenue_s3_bucket.example.name
								access_control_policy = {
									grants = [{
										grantee    = {
											type = "CanonicalUser"
											id  = "1f680dc8d84ca778f885628a39e1980850d408dbc3ca3a706bc182a9672f95ce"
										}
										permission = "FULL_CONTROL"
									}],
									owner  = {
										id           = "1f680dc8d84ca778f885628a39e1980850d408dbc3ca3a706bc182a9672f95ce"
									}
								}
							}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policy.grants.0.permission", "FULL_CONTROL"),
							resource.TestCheckResourceAttrSet(resourceName, "access_control_policy.grants.0.grantee.id"),
							resource.TestCheckResourceAttrSet(resourceName, "access_control_policy.owner.id"),
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
				// ! Destroy
				Destroy: true,
			}
		},
	}
}

func TestAccS3BucketACLResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketACLResource{}),
	})
}
