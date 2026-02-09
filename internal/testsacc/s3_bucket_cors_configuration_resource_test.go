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

var _ testsacc.TestACC = &S3BucketCorsConfigurationResource{}

const (
	S3BucketCorsConfigurationResourceName = testsacc.ResourceName("cloudavenue_s3_bucket_cors_configuration")
)

type S3BucketCorsConfigurationResource struct{}

func NewS3BucketCorsConfigurationResourceTest() testsacc.TestACC {
	return &S3BucketCorsConfigurationResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketCorsConfigurationResource) GetResourceName() string {
	return S3BucketCorsConfigurationResourceName.String()
}

func (r *S3BucketCorsConfigurationResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return resp
}

func (r *S3BucketCorsConfigurationResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_s3_bucket_cors_configuration" "example" {
						bucket = cloudavenue_s3_bucket.example.name
						cors_rules = [{
							allowed_headers = ["*"]
							allowed_methods = ["GET", "PUT"]
							allowed_origins = ["*"]
							expose_headers = ["ETag"]
							max_age_seconds = 3000
						}]
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_headers.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_headers.0", "*"),
						resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_methods.#", "2"),
						resource.TestCheckTypeSetElemAttr(resourceName, "cors_rules.0.allowed_methods.*", "GET"),
						resource.TestCheckTypeSetElemAttr(resourceName, "cors_rules.0.allowed_methods.*", "PUT"),
						resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_origins.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_origins.0", "*"),
						resource.TestCheckResourceAttr(resourceName, "cors_rules.0.expose_headers.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "cors_rules.0.expose_headers.0", "ETag"),
						resource.TestCheckResourceAttr(resourceName, "cors_rules.0.max_age_seconds", "3000"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_cors_configuration" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							cors_rules = [{
								allowed_headers = ["Content-Type"]
								allowed_methods = ["GET", "DELETE"]
								allowed_origins = ["https://www.example.com"]
								expose_headers = ["X-Custom-Header"]
								max_age_seconds = 3600
							}]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_headers.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_headers.0", "Content-Type"),
							resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_methods.#", "2"),
							resource.TestCheckTypeSetElemAttr(resourceName, "cors_rules.0.allowed_methods.*", "GET"),
							resource.TestCheckTypeSetElemAttr(resourceName, "cors_rules.0.allowed_methods.*", "DELETE"),
							resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_origins.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_origins.0", "https://www.example.com"),
							resource.TestCheckResourceAttr(resourceName, "cors_rules.0.expose_headers.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "cors_rules.0.expose_headers.0", "X-Custom-Header"),
							resource.TestCheckResourceAttr(resourceName, "cors_rules.0.max_age_seconds", "3600"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_cors_configuration" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							cors_rules = [{
								allowed_headers = ["Content-Type"]
								allowed_methods = ["GET", "DELETE"]
								allowed_origins = ["https://www.example.com"]
								expose_headers = ["X-Custom-Header"]
								max_age_seconds = 3600
							},
							{
								allowed_headers = ["Accept"]
								allowed_methods = ["GET"]
								allowed_origins = ["https://www.example.com"]
							}]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "cors_rules.#", "2"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"bucket"},
						ImportState:          true,
					},
				},
			}
		},
	}
}

func TestAccS3BucketCorsConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketCorsConfigurationResource{}),
	})
}
