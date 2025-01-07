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

var _ testsacc.TestACC = &S3BucketWebsiteConfigurationDataSource{}

const (
	S3BucketWebsiteConfigurationDataSourceName = testsacc.ResourceName("data.cloudavenue_s3_bucket_website_configuration")
)

type S3BucketWebsiteConfigurationDataSource struct{}

func NewS3BucketWebsiteConfigurationDataSourceTest() testsacc.TestACC {
	return &S3BucketWebsiteConfigurationDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketWebsiteConfigurationDataSource) GetResourceName() string {
	return S3BucketWebsiteConfigurationDataSourceName.String()
}

func (r *S3BucketWebsiteConfigurationDataSource) DependenciesConfig() (deps testsacc.DependenciesConfigResponse) {
	deps.Append(GetResourceConfig()[S3BucketWebsiteConfigurationResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketWebsiteConfigurationDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_bucket_website_configuration" "example" {
						bucket = cloudavenue_s3_bucket_website_configuration.example.bucket
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[S3BucketWebsiteConfigurationResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccS3BucketWebsiteConfigurationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketWebsiteConfigurationDataSource{}),
	})
}
