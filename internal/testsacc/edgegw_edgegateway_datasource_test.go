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
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway" "example" {
						name = cloudavenue_edgegateway.example.name
					}`,
					Checks: GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultChecks(),
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
