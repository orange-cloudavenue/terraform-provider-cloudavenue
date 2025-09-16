/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &EdgeGatewayDataSource{}

const (
	EdgeGatewayDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway")
)

type EdgeGatewayDataSource struct{}

func NewEdgeGatewayDataSourceTest() testsacc.TestACC {
	return &EdgeGatewayDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayDataSource) GetResourceName() string {
	return EdgeGatewayDataSourceName.String()
}

func (r *EdgeGatewayDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *EdgeGatewayDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway" "example" {
						name = cloudavenue_edgegateway.example.name
					}`,
					Checks: GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultChecks(),
				},
			}
		},
		"example_with_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway" "example_with_id" {
						id = cloudavenue_edgegateway.example.id
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
					},
				},
			}
		},
		"example_for_elb": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway" "example_for_elb" {
						name = "tn01e02ocb0006205spt101"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
					},
				},
			}
		},
	}
}

func TestAccEdgeGatewayDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayDataSource{}),
	})
}
