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

var _ testsacc.TestACC = &EdgeGatewayIPSetDataSource{}

const (
	EdgeGatewayIPSetDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_ip_set")
)

type EdgeGatewayIPSetDataSource struct{}

func NewEdgeGatewayIPSetDataSourceTest() testsacc.TestACC {
	return &EdgeGatewayIPSetDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayIPSetDataSource) GetResourceName() string {
	return EdgeGatewayIPSetDataSourceName.String()
}

func (r *EdgeGatewayIPSetDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *EdgeGatewayIPSetDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayIPSetResourceName]().GetDefaultConfig)
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_ip_set" "example" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						name = cloudavenue_edgegateway_ip_set.example.name
					}`,
					Checks: GetResourceConfig()[EdgeGatewayIPSetResourceName]().GetDefaultChecks(),
				},
				Destroy: true,
			}
		},
		"example_with_vdc_group": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayIPSetResourceName]().GetSpecificConfig("example_with_vdc_group"))
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_ip_set" "example_with_vdc_group" {
						edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
						name = cloudavenue_edgegateway_ip_set.example_with_vdc_group.name
					}`,

					Checks: GetResourceConfig()[EdgeGatewayIPSetResourceName]().GetSpecificChecks("example_with_vdc_group"),
				},
			}
		},
	}
}

func TestAccEdgeGatewayIPSetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayIPSetDataSource{}),
	})
}
