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

var _ testsacc.TestACC = &S3BucketWebsiteConfigurationResource{}

const (
	S3BucketWebsiteConfigurationResourceName = testsacc.ResourceName("cloudavenue_s3_bucket_website_configuration")
)

type S3BucketWebsiteConfigurationResource struct{}

func NewS3BucketWebsiteConfigurationResourceTest() testsacc.TestACC {
	return &S3BucketWebsiteConfigurationResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketWebsiteConfigurationResource) GetResourceName() string {
	return S3BucketWebsiteConfigurationResourceName.String()
}

func (r *S3BucketWebsiteConfigurationResource) DependenciesConfig() (deps testsacc.DependenciesConfigResponse) {
	deps.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketWebsiteConfigurationResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
					resource.TestCheckResourceAttrSet(resourceName, "website_endpoint"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_s3_bucket_website_configuration" "example" {
						bucket = cloudavenue_s3_bucket.example.name
						index_document = {
						  suffix = "index.html"
						}
						
						error_document = {
						  key = "error.html"
						}
						
						routing_rules = [{
							condition = {
							  key_prefix_equals = "docs/"
							}
							redirect = {
							  replace_key_prefix_with = "documents/"
							}
						}]
						
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "index_document.suffix", "index.html"),
						resource.TestCheckResourceAttr(resourceName, "error_document.key", "error.html"),
						resource.TestCheckResourceAttr(resourceName, "routing_rules.0.condition.key_prefix_equals", "docs/"),
						resource.TestCheckResourceAttr(resourceName, "routing_rules.0.redirect.replace_key_prefix_with", "documents/"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_website_configuration" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							index_document = {
							  suffix = "home.html"
							}
							
							error_document = {
							  key = "errors.html"
							}
							
							routing_rules = [{
								condition = {
								  key_prefix_equals = "img/"
								}
								redirect = {
								  replace_key_prefix_with = "imgs/"
								}
							}]
							
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "index_document.suffix", "home.html"),
							resource.TestCheckResourceAttr(resourceName, "error_document.key", "errors.html"),
							resource.TestCheckResourceAttr(resourceName, "routing_rules.0.condition.key_prefix_equals", "img/"),
							resource.TestCheckResourceAttr(resourceName, "routing_rules.0.redirect.replace_key_prefix_with", "imgs/"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_website_configuration" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							index_document = {
							  suffix = "home.html"
							}
							
							error_document = {
							  key = "errors.html"
							}
							
							routing_rules = [{
								condition = {
								  key_prefix_equals = "img/"
								}
								redirect = {
								  replace_key_prefix_with = "imgs/"
								  hostname = "www.example.com"
								  http_redirect_code = "302"
								  protocol = "https"
								}
							}]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "index_document.suffix", "home.html"),
							resource.TestCheckResourceAttr(resourceName, "error_document.key", "errors.html"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "routing_rules.*", map[string]string{
								"condition.key_prefix_equals":      "img/",
								"redirect.replace_key_prefix_with": "imgs/",
								"redirect.hostname":                "www.example.com",
								"redirect.http_redirect_code":      "302",
								"redirect.protocol":                "https",
							}),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_website_configuration" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							index_document = {
							  suffix = "home.html"
							}
							
							error_document = {
							  key = "errors.html"
							}
							
							routing_rules = [
							{
								condition = {
								  key_prefix_equals = "img/"
								}
								redirect = {
								  replace_key_prefix_with = "imgs/"
								  hostname = "www.example.com"
								  http_redirect_code = "302"
								  protocol = "https"
								}
							},
							{
							  condition = {
								http_error_code_returned_equals = "404"
							  }
							  redirect = {
								replace_key_with = "errors.html"
								hostname = "www.example.com"
								http_redirect_code = "301"
								protocol = "https"
							  }
						  	}]
							
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "index_document.suffix", "home.html"),
							resource.TestCheckResourceAttr(resourceName, "error_document.key", "errors.html"),
							resource.TestCheckResourceAttr(resourceName, "routing_rules.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "routing_rules.*", map[string]string{
								"condition.key_prefix_equals":      "img/",
								"redirect.replace_key_prefix_with": "imgs/",
								"redirect.hostname":                "www.example.com",
								"redirect.http_redirect_code":      "302",
								"redirect.protocol":                "https",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "routing_rules.*", map[string]string{
								"condition.http_error_code_returned_equals": "404",
								"redirect.replace_key_with":                 "errors.html",
								"redirect.hostname":                         "www.example.com",
								"redirect.http_redirect_code":               "301",
								"redirect.protocol":                         "https",
							}),
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
				Destroy: true,
			}
		},
		"example_redirect_all_requests_to": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
					resource.TestCheckResourceAttrSet(resourceName, "website_endpoint"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_s3_bucket_website_configuration" "example_redirect_all_requests_to" {
						bucket = cloudavenue_s3_bucket.example.name
						redirect_all_requests_to = {
							hostname = "example.com"
						}
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "redirect_all_requests_to.hostname", "example.com"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_website_configuration" "example_redirect_all_requests_to" {
							bucket = cloudavenue_s3_bucket.example.name
							redirect_all_requests_to = {
								hostname = "www.example.com"
								protocol = "https"
							}
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "redirect_all_requests_to.hostname", "www.example.com"),
							resource.TestCheckResourceAttr(resourceName, "redirect_all_requests_to.protocol", "https"),
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
				Destroy: true,
			}
		},
	}
}

func TestAccS3BucketWebsiteConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketWebsiteConfigurationResource{}),
	})
}
