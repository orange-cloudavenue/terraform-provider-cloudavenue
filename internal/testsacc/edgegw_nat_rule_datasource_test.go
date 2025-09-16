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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &EdgeGatewayNATRuleDataSource{}

const (
	EdgeGatewayNATRuleDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_nat_rule")
)

type EdgeGatewayNATRuleDataSource struct{}

func NewEdgeGatewayNATRuleDataSourceTest() testsacc.TestACC {
	return &EdgeGatewayNATRuleDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayNATRuleDataSource) GetResourceName() string {
	return EdgeGatewayNATRuleDataSourceName.String()
}

func (r *EdgeGatewayNATRuleDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayNATRuleResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayNATRuleDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (with edge_gateway_id)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_nat_rule" "example" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						name = cloudavenue_edgegateway_nat_rule.example.name
					}`,
					Checks: GetResourceConfig()[EdgeGatewayNATRuleResourceName]().GetDefaultChecks(),
				},
			}
		},
		// * Test Two (with edge_gateway_name)
		"example_with_edge_gateway_name": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_nat_rule" "example_with_edge_gateway_name" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
						name = cloudavenue_edgegateway_nat_rule.example.name
					}`,
					Checks: GetResourceConfig()[EdgeGatewayNATRuleResourceName]().GetDefaultChecks(),
				},
			}
		},
		// * Test nat rule with same name
		"example_with_same_name": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayNATRuleResourceName]().GetSpecificConfig("example_two_rules_with_same_name"))
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_nat_rule" "example_with_same_name" {
						depends_on = [
							cloudavenue_edgegateway_nat_rule.example_two_rules_with_same_name,
							cloudavenue_edgegateway_nat_rule.example_two_rules_with_same_name_2	
						]
						edge_gateway_name = cloudavenue_edgegateway.example.name
						name = cloudavenue_edgegateway_nat_rule.example_two_rules_with_same_name.name
					}`,
					Checks: []resource.TestCheckFunc{},
					TFAdvanced: testsacc.TFAdvanced{
						ExpectError: regexp.MustCompile(`Multiple NAT Rules found with the same name`),
					},
				},
			}
		},
	}
}

func TestAccEdgeGatewayNATRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayNATRuleDataSource{}),
	})
}
