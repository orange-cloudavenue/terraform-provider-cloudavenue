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

var _ testsacc.TestACC = &EdgeGatewayStaticRouteDataSource{}

const (
	EdgeGatewayStaticRouteDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_static_route")
)

type EdgeGatewayStaticRouteDataSource struct{}

func NewEdgeGatewayStaticRouteDataSourceTest() testsacc.TestACC {
	return &EdgeGatewayStaticRouteDataSource{}
}

func (r *EdgeGatewayStaticRouteDataSource) GetResourceName() string {
	return EdgeGatewayStaticRouteDataSourceName.String()
}

func (r *EdgeGatewayStaticRouteDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayStaticRouteResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayStaticRouteDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"basic": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_static_route" "example" {
					  edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
					  name = cloudavenue_edgegateway_static_route.example.name
					}
					`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet("data.cloudavenue_edgegateway_static_route.example", "id"),
						resource.TestCheckResourceAttrSet("data.cloudavenue_edgegateway_static_route.example", "network"),
						resource.TestCheckResourceAttrSet("data.cloudavenue_edgegateway_static_route.example", "next_hop"),
					},
				},
			}
		},
	}
}

func TestAccEdgeGatewayStaticRouteDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayStaticRouteDataSource{}),
	})
}
