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
	return
}

func (r *EdgeGatewayFirewallDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
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
	}
}

func TestAccEdgeGatewayFirewallDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayFirewallDataSource{}),
	})
}
