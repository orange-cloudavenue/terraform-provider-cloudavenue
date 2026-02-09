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

var _ testsacc.TestACC = &EdgeGatewayServicesDataSource{}

const (
	EdgeGatewayServicesDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_services")
)

type EdgeGatewayServicesDataSource struct{}

func NewEdgeGatewayServicesDataSourceTest() testsacc.TestACC {
	return &EdgeGatewayServicesDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayServicesDataSource) GetResourceName() string {
	return EdgeGatewayServicesDataSourceName.String()
}

func (r *EdgeGatewayServicesDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[EdgeGatewayServicesResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayServicesDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_services" "example" {
						edge_gateway_name = cloudavenue_edgegateway_services.example.edge_gateway_name
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[EdgeGatewayServicesResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccEdgeGatewayServicesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayServicesDataSource{}),
	})
}
