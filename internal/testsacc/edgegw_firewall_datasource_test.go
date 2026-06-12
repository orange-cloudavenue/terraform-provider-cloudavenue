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

var _ testsacc.TestACC = &EdgeGatewayFirewallDataSource{}

const (
	EdgeGatewayFirewallDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_firewall")
)

type EdgeGatewayFirewallDataSource struct{}

func NewEdgeGatewayFirewallDataSourceTest() testsacc.TestACC {
	return &EdgeGatewayFirewallDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayFirewallDataSource) GetResourceName() string {
	return EdgeGatewayFirewallDataSourceName.String()
}

func (r *EdgeGatewayFirewallDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayFirewallResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayFirewallDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		testNameExample: func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_firewall" "example" {
						edge_gateway_id = cloudavenue_edgegateway_firewall.example.edge_gateway_id
					}`,
					Checks: GetResourceConfig()[EdgeGatewayFirewallResourceName]().GetDefaultChecks(),
				},
			}
		},
		"example_with_context_profile": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetDataSourceConfig()[EdgeGatewayNetworkContextProfileDatasourceName]().GetDefaultConfig)
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
            resource "cloudavenue_edgegateway_firewall" "example_with_context_profile" {
              edge_gateway_id = cloudavenue_edgegateway.example.id
              rules = [
                {
                  action      = "ALLOW"
                  name        = "allow outbound SSL"
                  direction   = "OUT"
                  ip_protocol = "IPV4"
                  network_context_profile_ids = [data.cloudavenue_edgegateway_network_context_profile.example.id]
                }
              ]
            }

            data "cloudavenue_edgegateway_firewall" "example_with_context_profile" {
              edge_gateway_id = cloudavenue_edgegateway.example.id
              depends_on      = [cloudavenue_edgegateway_firewall.example_with_context_profile]
            }`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.name", "allow outbound SSL"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.network_context_profile_ids.#", "1"),
						resource.TestCheckResourceAttrSet(resourceName, "rules.0.network_context_profile_ids.0"),
					},
				},
				Destroy: true,
			}
		},
	}
}

func TestAccEdgeGatewayFirewallDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayFirewallDataSource{}),
	})
}
