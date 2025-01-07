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

var _ testsacc.TestACC = &S3BucketCorsConfigurationDataSource{}

const (
	S3BucketCorsConfigurationDataSourceName = testsacc.ResourceName("data.cloudavenue_s3_bucket_cors_configuration")
)

type S3BucketCorsConfigurationDataSource struct{}

func NewS3BucketCorsConfigurationDataSourceTest() testsacc.TestACC {
	return &S3BucketCorsConfigurationDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketCorsConfigurationDataSource) GetResourceName() string {
	return S3BucketCorsConfigurationDataSourceName.String()
}

func (r *S3BucketCorsConfigurationDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[S3BucketCorsConfigurationResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketCorsConfigurationDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_bucket_cors_configuration" "example" {
						bucket = cloudavenue_s3_bucket_cors_configuration.example.id
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[S3BucketCorsConfigurationResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccS3BucketCorsConfigurationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketCorsConfigurationDataSource{}),
	})
}
