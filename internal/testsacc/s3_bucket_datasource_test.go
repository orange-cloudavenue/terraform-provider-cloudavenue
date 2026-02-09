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

var _ testsacc.TestACC = &S3BucketDatasource{}

const (
	S3BucketDatasourceName = testsacc.ResourceName("data.cloudavenue_s3_bucket")
)

type S3BucketDatasource struct{}

func NewS3BucketDatasourceTest() testsacc.TestACC {
	return &S3BucketDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketDatasource) GetResourceName() string {
	return S3BucketDatasourceName.String()
}

func (r *S3BucketDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return resp
}

func (r *S3BucketDatasource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_bucket" "example" {
						name = cloudavenue_s3_bucket.example.name
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[S3BucketResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccS3BucketDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketDatasource{}),
	})
}
