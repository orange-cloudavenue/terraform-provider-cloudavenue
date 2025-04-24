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

var _ testsacc.TestACC = &EdgeGatewaySecurityGroupDataSource{}

const (
	EdgeGatewaySecurityGroupDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_security_group")
)

type EdgeGatewaySecurityGroupDataSource struct{}

func NewEdgeGatewaySecurityGroupDataSourceTest() testsacc.TestACC {
	return &EdgeGatewaySecurityGroupDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewaySecurityGroupDataSource) GetResourceName() string {
	return EdgeGatewaySecurityGroupDataSourceName.String()
}

func (r *EdgeGatewaySecurityGroupDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewaySecurityGroupResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewaySecurityGroupDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_security_group" "example" {
						name            = cloudavenue_edgegateway_security_group.example.name
						edge_gateway_id = cloudavenue_edgegateway.example.id
					}`,
					Checks: GetResourceConfig()[EdgeGatewaySecurityGroupResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccEdgeGatewaySecurityGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewaySecurityGroupDataSource{}),
	})
}
