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

var _ testsacc.TestACC = &EdgeGatewayNetworkRoutedDataSource{}

const (
	EdgeGatewayNetworkRoutedDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_network_routed")
)

type EdgeGatewayNetworkRoutedDataSource struct{}

func NewEdgeGatewayNetworkRoutedDataSourceTest() testsacc.TestACC {
	return &EdgeGatewayNetworkRoutedDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayNetworkRoutedDataSource) GetResourceName() string {
	return EdgeGatewayNetworkRoutedDataSourceName.String()
}

func (r *EdgeGatewayNetworkRoutedDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[EdgeGatewayNetworkRoutedResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayNetworkRoutedDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_network_routed" "example" {
						name = cloudavenue_edgegateway_network_routed.example.name
						edge_gateway_name = cloudavenue_edgegateway.example.name
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[EdgeGatewayNetworkRoutedResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccEdgeGatewayNetworkRoutedDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayNetworkRoutedDataSource{}),
	})
}
