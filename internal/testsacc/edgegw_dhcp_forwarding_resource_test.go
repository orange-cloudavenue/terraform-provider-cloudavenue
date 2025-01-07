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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &EdgeGatewayDhcpForwardingResource{}

const (
	EdgeGatewayDhcpForwardingResourceName = testsacc.ResourceName("cloudavenue_edgegateway_dhcp_forwarding")
)

type EdgeGatewayDhcpForwardingResource struct{}

func NewEdgeGatewayDhcpForwardingResourceTest() testsacc.TestACC {
	return &EdgeGatewayDhcpForwardingResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayDhcpForwardingResource) GetResourceName() string {
	return EdgeGatewayDhcpForwardingResourceName.String()
}

func (r *EdgeGatewayDhcpForwardingResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *EdgeGatewayDhcpForwardingResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						enabled = true
						dhcp_servers = [
							"192.168.10.10"
						]
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.10"),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
							edge_gateway_id = cloudavenue_edgegateway.example.id
							enabled = true
							dhcp_servers = [
								"192.168.10.10",
								"192.168.10.11"
							]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "2"),
							resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.10"),
							resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.11"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
							edge_gateway_id = cloudavenue_edgegateway.example.id
							enabled = false
							dhcp_servers = [
								"192.168.10.10",
								"192.168.10.11",
								"192.168.10.12"
							]
						}`,
						TFAdvanced: testsacc.TFAdvanced{
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
							ExpectError:        regexp.MustCompile("DHCP Servers cannot be edited"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			}
		},
		// This test is commented out because the previous test generate a error.
		// | ErrorCode:req-0006 - ErrorReason:Schema validation
		// | error - ErrorMessage:Error : 1 validation error for EdgeGatewayPut
		// | rateLimit
		// |   Value error, 0 is not a valid QosProfileId [type=value_error, input_value=0, input_type=int]
		// |     For further information visit https://errors.pydantic.dev/2.4/v/value_error
		// Actually no workaround is possible, the error is generated by the API.
		//
		// "example_with_vdc_group": func(_ context.Context, resourceName string) testsacc.Test {
		// 	return testsacc.Test{
		// 		// ! Create testing
		// 		CommonChecks: []resource.TestCheckFunc{
		// 			resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
		// 			resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
		// 			resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
		// 		},
		// 		CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
		// 			resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetSpecificConfig("example_with_vdc_group"))
		// 			return
		// 		},
		// 		Create: testsacc.TFConfig{
		// 			TFConfig: `
		// 			resource "cloudavenue_edgegateway_dhcp_forwarding" "example_with_vdc_group" {
		// 				edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
		// 				enabled = true
		// 				dhcp_servers = [
		// 					"192.168.10.10"
		// 				]
		// 			}`,
		// 			Checks: []resource.TestCheckFunc{
		// 				resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		// 				resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "1"),
		// 				resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.10"),
		// 			},
		// 		},
		// 		Updates: []testsacc.TFConfig{
		// 			{
		// 				TFConfig: `
		// 				resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
		// 					edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
		// 					enabled = true
		// 					dhcp_servers = [
		// 						"192.168.10.10",
		// 						"192.168.10.11"
		// 					]
		// 				}`,
		// 				Checks: []resource.TestCheckFunc{
		// 					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		// 					resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "2"),
		// 					resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.10"),
		// 					resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.11"),
		// 				},
		// 			},
		// 			{
		// 				TFConfig: `
		// 				resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
		// 					edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
		// 					enabled = false
		// 					dhcp_servers = [
		// 						"192.168.10.10",
		// 						"192.168.10.11",
		// 						"192.168.10.12"
		// 					]
		// 				}`,
		// 				TFAdvanced: testsacc.TFAdvanced{
		// 					PlanOnly:           true,
		// 					ExpectNonEmptyPlan: true,
		// 					ExpectError:        regexp.MustCompile("DHCP Servers cannot be edited"),
		// 				},
		// 			},
		// 		},
		// 		// ! Imports testing
		// 		Imports: []testsacc.TFImport{
		// 			{
		// 				ImportState:       true,
		// 				ImportStateVerify: true,
		// 			},
		// 		},
		// 	}
		// },
	}
}

func TestAccEdgeGatewayDhcpForwardingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayDhcpForwardingResource{}),
	})
}
