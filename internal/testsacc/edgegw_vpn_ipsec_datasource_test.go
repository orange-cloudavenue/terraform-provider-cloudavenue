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

var _ testsacc.TestACC = &VPNIPSecDataSource{}

const EdgeGatewayVPNIPSecDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_vpn_ipsec")

type VPNIPSecDataSource struct{}

func NewEdgeGatewayVPNIPSecDataSourceTest() testsacc.TestACC {
	return &VPNIPSecDataSource{}
}

func (r *VPNIPSecDataSource) GetResourceName() string {
	return EdgeGatewayVPNIPSecDataSourceName.String()
}

func (r *VPNIPSecDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayVPNIPSecResourceName]().GetDefaultConfig)
	return resp
}

func (r *VPNIPSecDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_vpn_ipsec" "example" {
					  edge_gateway_id = cloudavenue_edgegateway.example.id
					  name = cloudavenue_edgegateway_vpn_ipsec.example.name
					}`,
					Checks: GetResourceConfig()[EdgeGatewayVPNIPSecResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVPNIPSecDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VPNIPSecDataSource{}),
	})
}
