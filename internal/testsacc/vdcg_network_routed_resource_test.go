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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGNetworkRoutedResource{}

const (
	VDCGNetworkRoutedResourceName = testsacc.ResourceName("cloudavenue_vdcg_network_routed")
)

type VDCGNetworkRoutedResource struct{}

func NewVDCGNetworkRoutedResourceTest() testsacc.TestACC {
	return &VDCGNetworkRoutedResource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGNetworkRoutedResource) GetResourceName() string {
	return VDCGNetworkRoutedResourceName.String()
}

func (r *VDCGNetworkRoutedResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetSpecificConfig(testNameExampleWithVDCGroup))
	return resp
}

func (r *VDCGNetworkRoutedResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		testNameExample: func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Network)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_network_routed" "example" {
						name        = {{ generate . "name" }}
						description = {{ generate . "description" }}
						vdc_group_id = cloudavenue_vdcg.example.id
						edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
					
						gateway       = "192.168.100.1"
						prefix_length = 24
					
						dns1 = "1.1.1.1"
						dns2 = "1.0.0.1"
						dns_suffix = "example.com"
					
						static_ip_pool = [
						{
							start_address = "192.168.100.10"
							end_address   = "192.168.100.20"
						},
						{
							start_address = "192.168.100.100"
							end_address   = "192.168.100.130"
						}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.100.1"),
						resource.TestCheckResourceAttr(resourceName, "prefix_length", "24"),
						resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.1"),
						resource.TestCheckResourceAttr(resourceName, "dns2", "1.0.0.1"),
						resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example.com"),
						resource.TestCheckResourceAttr(resourceName, "static_ip_pool.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "guest_vlan_allowed", "false"), // Default value
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
							testAttrStartAddress: testRoutedNetworkIP10,
							testAttrEndAddress:   testRoutedNetworkIP20,
						}),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
							testAttrStartAddress: testRoutedNetworkIP100,
							testAttrEndAddress:   testRoutedNetworkIP130,
						}),
					},
				},
				// ! Updates testing
				// * Update name
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_network_routed" "example" {
							name        = {{ generate . "name" }}
							description = {{ get . "description" }}
							vdc_group_id = cloudavenue_vdcg.example.id
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
						
							gateway       = "192.168.100.1"
							prefix_length = 24
						
							dns1 = "1.1.1.1"
							dns2 = "1.0.0.1"
							dns_suffix = "example.com"
						
							static_ip_pool = [
							{
								start_address = "192.168.100.10"
								end_address   = "192.168.100.20"
							},
							{
								start_address = "192.168.100.100"
								end_address   = "192.168.100.130"
							}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.100.1"),
							resource.TestCheckResourceAttr(resourceName, "prefix_length", "24"),
							resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.1"),
							resource.TestCheckResourceAttr(resourceName, "dns2", "1.0.0.1"),
							resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example.com"),
							resource.TestCheckResourceAttr(resourceName, "static_ip_pool.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "guest_vlan_allowed", "false"), // Default value
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: testRoutedNetworkIP10,
								testAttrEndAddress:   testRoutedNetworkIP20,
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: testRoutedNetworkIP100,
								testAttrEndAddress:   testRoutedNetworkIP130,
							}),
						},
					},
					// * Update description
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_network_routed" "example" {
							name        = {{ get . "name" }}
							description = {{ generate . "description" }}
							vdc_group_id = cloudavenue_vdcg.example.id
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
						
							gateway       = "192.168.100.1"
							prefix_length = 24
						
							dns1 = "1.1.1.1"
							dns2 = "1.0.0.1"
							dns_suffix = "example.com"
						
							static_ip_pool = [
							{
								start_address = "192.168.100.10"
								end_address   = "192.168.100.20"
							},
							{
								start_address = "192.168.100.100"
								end_address   = "192.168.100.130"
							}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.100.1"),
							resource.TestCheckResourceAttr(resourceName, "prefix_length", "24"),
							resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.1"),
							resource.TestCheckResourceAttr(resourceName, "dns2", "1.0.0.1"),
							resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example.com"),
							resource.TestCheckResourceAttr(resourceName, "static_ip_pool.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "guest_vlan_allowed", "false"), // Default value
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: testRoutedNetworkIP10,
								testAttrEndAddress:   testRoutedNetworkIP20,
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: testRoutedNetworkIP100,
								testAttrEndAddress:   testRoutedNetworkIP130,
							}),
						},
					},
					// * Update DNS
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_network_routed" "example" {
							name        = {{ get . "name" }}
							description = {{ generate . "description" }}
							vdc_group_id = cloudavenue_vdcg.example.id
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
						
							gateway       = "192.168.100.1"
							prefix_length = 24
						
							dns1 = "208.67.222.222"
							dns2 = "208.67.220.220"
							dns_suffix = "example.local"
						
							static_ip_pool = [
							{
								start_address = "192.168.100.10"
								end_address   = "192.168.100.20"
							},
							{
								start_address = "192.168.100.100"
								end_address   = "192.168.100.130"
							}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.100.1"),
							resource.TestCheckResourceAttr(resourceName, "prefix_length", "24"),
							resource.TestCheckResourceAttr(resourceName, "dns1", "208.67.222.222"),
							resource.TestCheckResourceAttr(resourceName, "dns2", "208.67.220.220"),
							resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example.local"),
							resource.TestCheckResourceAttr(resourceName, "static_ip_pool.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "guest_vlan_allowed", "false"), // Default value
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: testRoutedNetworkIP10,
								testAttrEndAddress:   testRoutedNetworkIP20,
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: testRoutedNetworkIP100,
								testAttrEndAddress:   testRoutedNetworkIP130,
							}),
						},
					},
					// * Update Static IP Pool
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_network_routed" "example" {
							name        = {{ get . "name" }}
							description = {{ generate . "description" }}
							vdc_group_id = cloudavenue_vdcg.example.id
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
						
							gateway       = "192.168.100.1"
							prefix_length = 24
						
							dns1 = "208.67.222.222"
							dns2 = "208.67.220.220"
							dns_suffix = "example.local"
						
							static_ip_pool = [
							{
								start_address = "192.168.100.40"
								end_address   = "192.168.100.60"
							},
							{
								start_address = "192.168.100.100"
								end_address   = "192.168.100.130"
							},
							{
								start_address = "192.168.100.200"
								end_address   = "192.168.100.220"
							}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.100.1"),
							resource.TestCheckResourceAttr(resourceName, "prefix_length", "24"),
							resource.TestCheckResourceAttr(resourceName, "dns1", "208.67.222.222"),
							resource.TestCheckResourceAttr(resourceName, "dns2", "208.67.220.220"),
							resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example.local"),
							resource.TestCheckResourceAttr(resourceName, "static_ip_pool.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "guest_vlan_allowed", "false"), // Default value
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: "192.168.100.40",
								testAttrEndAddress:   "192.168.100.60",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: testRoutedNetworkIP100,
								testAttrEndAddress:   testRoutedNetworkIP130,
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: "192.168.100.200",
								testAttrEndAddress:   "192.168.100.220",
							}),
						},
					},
					// * Update Guest VLAN Allowed
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_network_routed" "example" {
							name        = {{ get . "name" }}
							description = {{ get . "description" }}
							vdc_group_id = cloudavenue_vdcg.example.id
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id

							gateway       = "192.168.100.1"
							prefix_length = 24
						
							dns1 = "208.67.222.222"
							dns2 = "208.67.220.220"
							dns_suffix = "example.local"
						
							static_ip_pool = [
							{
								start_address = "192.168.100.40"
								end_address   = "192.168.100.60"
							},
							{
								start_address = "192.168.100.100"
								end_address   = "192.168.100.130"
							},
							{
								start_address = "192.168.100.200"
								end_address   = "192.168.100.220"
							}
							]
							guest_vlan_allowed = true
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.100.1"),
							resource.TestCheckResourceAttr(resourceName, "prefix_length", "24"),
							resource.TestCheckResourceAttr(resourceName, "dns1", "208.67.222.222"),
							resource.TestCheckResourceAttr(resourceName, "dns2", "208.67.220.220"),
							resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example.local"),
							resource.TestCheckResourceAttr(resourceName, "static_ip_pool.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "guest_vlan_allowed", "true"), // Default value
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: "192.168.100.40",
								testAttrEndAddress:   "192.168.100.60",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: testRoutedNetworkIP100,
								testAttrEndAddress:   testRoutedNetworkIP130,
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: "192.168.100.200",
								testAttrEndAddress:   "192.168.100.220",
							}),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{testAttrVDCGroupID, "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{testAttrVDCGroupID, testAttrName},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{testAttrVDCGroupName, "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{testAttrVDCGroupName, testAttrName},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		// ! MoveState testing
		testNameMoveState: func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{},
				// ! Create testing (old deprecated resource — must use hardcoded address since resourceName points to new type)
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_network_routed" "example_move_state" {
						name             = {{ generate . "name" }}
						description      = {{ generate . "description" }}
						edge_gateway_id  = cloudavenue_edgegateway.example_with_vdc_group.id

						gateway       = "192.168.200.1"
						prefix_length = 24

						dns1 = "1.1.1.1"
						dns2 = "1.0.0.1"
						dns_suffix = "example.com"

						static_ip_pool = [
						{
							start_address = "192.168.200.10"
							end_address   = "192.168.200.20"
						}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr("cloudavenue_network_routed.example_move_state", "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr("cloudavenue_network_routed.example_move_state", "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr("cloudavenue_network_routed.example_move_state", "gateway", "192.168.200.1"),
						resource.TestCheckResourceAttr("cloudavenue_network_routed.example_move_state", "prefix_length", "24"),
						resource.TestCheckResourceAttr("cloudavenue_network_routed.example_move_state", "dns1", "1.1.1.1"),
						resource.TestCheckResourceAttr("cloudavenue_network_routed.example_move_state", "dns2", "1.0.0.1"),
						resource.TestCheckResourceAttr("cloudavenue_network_routed.example_move_state", "dns_suffix", "example.com"),
						resource.TestCheckResourceAttr("cloudavenue_network_routed.example_move_state", "static_ip_pool.#", "1"),
						resource.TestCheckTypeSetElemNestedAttrs("cloudavenue_network_routed.example_move_state", "static_ip_pool.*", map[string]string{
							testAttrStartAddress: "192.168.200.10",
							testAttrEndAddress:   "192.168.200.20",
						}),
					},
				},
				// ! Updates testing (moved block to new resource type)
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_network_routed" "example_move_state" {
							name        = {{ get . "name" }}
							description = {{ get . "description" }}
							vdc_group_id    = cloudavenue_vdcg.example.id
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id

							gateway       = "192.168.200.1"
							prefix_length = 24

							dns1 = "1.1.1.1"
							dns2 = "1.0.0.1"
							dns_suffix = "example.com"

							static_ip_pool = [
							{
								start_address = "192.168.200.10"
								end_address   = "192.168.200.20"
							}
							]
						}

						moved {
							from = cloudavenue_network_routed.example_move_state
							to   = cloudavenue_vdcg_network_routed.example_move_state
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.200.1"),
							resource.TestCheckResourceAttr(resourceName, "prefix_length", "24"),
							resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.1"),
							resource.TestCheckResourceAttr(resourceName, "dns2", "1.0.0.1"),
							resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example.com"),
							resource.TestCheckResourceAttr(resourceName, "static_ip_pool.#", "1"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								testAttrStartAddress: "192.168.200.10",
								testAttrEndAddress:   "192.168.200.20",
							}),
							resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
							resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
						},
					},
				},
			}
		},
	}
}

func TestAccVDCGNetworkRoutedResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGNetworkRoutedResource{}),
	})
}
