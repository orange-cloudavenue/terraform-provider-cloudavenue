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

var _ testsacc.TestACC = &NATRuleResource{}

const (
	EdgeGatewayNATRuleResourceName = testsacc.ResourceName("cloudavenue_edgegateway_nat_rule")
)

type NATRuleResource struct{}

func NewEdgeGatewayNATRuleResourceTest() testsacc.TestACC {
	return &NATRuleResource{}
}

func (r *NATRuleResource) GetResourceName() string {
	return EdgeGatewayNATRuleResourceName.String()
}

func (r *NATRuleResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return
}

func (r *NATRuleResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "SNAT"),
					resource.TestCheckNoResourceAttr(resourceName, "dnat_external_port"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_nat_rule" "example" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						
						name        = {{ generate . "name" }}
						rule_type   = "SNAT"
						description = {{ generate . "description" }}
						
						# Using primary_ip from edge gateway
						external_address         = "89.32.25.10"
						internal_address         = "11.11.11.0/24"
						snat_destination_address = "8.8.8.8"
						
						priority = 10
						}`),
					Checks: []resource.TestCheckFunc{
						// ! base checks
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "external_address", "89.32.25.10"),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.0/24"),
						resource.TestCheckResourceAttr(resourceName, "snat_destination_address", "8.8.8.8"),
						resource.TestCheckResourceAttr(resourceName, "priority", "10"),
					},
				},
				// ! Update testing
				Updates: []testsacc.TFConfig{
					// * Update name / description
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_nat_rule" "example" {
							edge_gateway_id = cloudavenue_edgegateway.example.id

							name        = {{ generate . "name" }}
							rule_type   = "SNAT"
							description = {{ generate . "description" }}

							external_address         = "89.32.25.10"
							internal_address         = "11.11.11.0/24"
							snat_destination_address = "9.9.9.9"

							priority = 0
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "external_address", "89.32.25.10"),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.0/24"),
							resource.TestCheckResourceAttr(resourceName, "snat_destination_address", "9.9.9.9"),
							resource.TestCheckResourceAttr(resourceName, "priority", "0"),
						},
					},
					// Update external_address / priority
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_nat_rule" "example" {
							edge_gateway_id = cloudavenue_edgegateway.example.id

							name        = {{ get . "name" }}
							rule_type   = "SNAT"
							description = {{ get . "description" }}

							external_address         = "89.32.25.11"
							internal_address         = "11.11.11.0/24"
							snat_destination_address = "9.9.9.9"

							priority = 0
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "external_address", "89.32.25.11"),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.0/24"),
							resource.TestCheckResourceAttr(resourceName, "snat_destination_address", "9.9.9.9"),
							resource.TestCheckResourceAttr(resourceName, "priority", "0"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_no_snat": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "NO_SNAT"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_nat_rule" "example_no_snat" {
						edge_gateway_id = cloudavenue_edgegateway.example.id

						name        = {{ generate . "name" }}
						rule_type   = "NO_SNAT"
						description = {{ generate . "description" }}

						# Using primary_ip from edge gateway
						internal_address         = "11.11.11.0/24"

						priority = 10
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.0/24"),
						resource.TestCheckResourceAttr(resourceName, "priority", "10"),
					},
				},
			}
		},
		"example_dnat": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_nat_rule" "example_dnat" {
						edge_gateway_id = cloudavenue_edgegateway.example.id

						name        = {{ generate . "name" }}
						rule_type   = "DNAT"
						description = {{ generate . "description" }}

						# Using primary_ip from edge gateway
						external_address         = "89.32.25.10"
						internal_address         = "11.11.11.4"

						dnat_external_port = "8080"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "external_address"),
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.4"),
						resource.TestCheckResourceAttr(resourceName, "dnat_external_port", "8080"),
						resource.TestCheckResourceAttr(resourceName, "priority", "0"),
						resource.TestCheckNoResourceAttr(resourceName, "snat_destination_address"),
					},
				},
				// ! Update testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_nat_rule" "example_dnat" {
							edge_gateway_id = cloudavenue_edgegateway.example.id

							name        = {{ get . "name" }}
							rule_type   = "DNAT"
							description = {{ generate . "description" }}

							# Using primary_ip from edge gateway
							external_address         = "89.32.25.10"
							internal_address         = "4.11.11.11"

							priority = 25
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttrSet(resourceName, "external_address"),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "internal_address", "4.11.11.11"),
							resource.TestCheckNoResourceAttr(resourceName, "dnat_external_port"),
							resource.TestCheckResourceAttr(resourceName, "priority", "25"),
							resource.TestCheckNoResourceAttr(resourceName, "snat_destination_address"),
						},
					},
				},
			}
		},
		"example_dnat_with_app_port_profile": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_nat_rule" "example_dnat_with_app_port_profile" {
						edge_gateway_id = cloudavenue_edgegateway.example.id

						name        = {{ generate . "name" }}
						rule_type   = "DNAT"
						description = {{ generate . "description" }}

						# Using primary_ip from edge gateway
						external_address         = "89.32.25.10"
						internal_address         = "11.11.11.4"

						app_port_profile_id = cloudavenue_edgegateway_app_port_profile.example.id
				}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttrSet(resourceName, "external_address"),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.4"),
						resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
					},
				},
				// ! Update testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_nat_rule" "example_dnat_with_app_port_profile" {
							edge_gateway_id = cloudavenue_edgegateway.example.id
					
							name        = {{ get . "name" }}
							rule_type   = "DNAT"
							description = {{ generate . "description" }}

							# Using primary_ip from edge gateway
							external_address         = "89.32.25.10"
							internal_address         = "11.11.11.2"

							app_port_profile_id = cloudavenue_edgegateway_app_port_profile.example.id
				}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttrSet(resourceName, "external_address"),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.2"),
							resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
						},
					},
				},
			}
		},
		"example_reflexive": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "external_address"),
					resource.TestCheckResourceAttr(resourceName, "priority", "0"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_nat_rule" "example_reflexive" {
						edge_gateway_id = cloudavenue_edgegateway.example.id

						name        = {{ generate . "name" }}
						rule_type   = "REFLEXIVE"
						description = {{ generate . "description" }}

						# Using primary_ip from edge gateway
						external_address         = "89.32.25.10"
						internal_address         = "192.168.0.1"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "internal_address", "192.168.0.1"),
						resource.TestCheckResourceAttr(resourceName, "external_address", "89.32.25.10"),
						resource.TestCheckResourceAttr(resourceName, "rule_type", "REFLEXIVE"),
					},
				},
			}
		},
		"example_with_vdc_group": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetSpecificConfig("example_with_vdc_group"))
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_nat_rule" "example_with_vdc_group" {
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id

							name        = {{ generate . "name" }}
							rule_type   = "DNAT"
							description = {{ generate . "description" }}

							# Using primary_ip from edge gateway
							external_address         = "89.32.25.10"
							internal_address         = "11.11.11.4"

							dnat_external_port = "8080"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttrSet(resourceName, "external_address"),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.4"),
						resource.TestCheckResourceAttr(resourceName, "dnat_external_port", "8080"),
						resource.TestCheckResourceAttr(resourceName, "priority", "0"),
					},
				},
				// ! Update testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_nat_rule" "example_with_vdc_group" {
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id

							name        = {{ get . "name" }}
							rule_type   = "DNAT"
							description = {{ generate . "description" }}

							# Using primary_ip from edge gateway
							external_address         = "89.32.25.10"
							internal_address         = "4.11.11.11"

							priority = 25
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttrSet(resourceName, "external_address"),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "internal_address", "4.11.11.11"),
							resource.TestCheckNoResourceAttr(resourceName, "dnat_external_port"),
							resource.TestCheckResourceAttr(resourceName, "priority", "25"),
						},
					},
				},
			}
		},
		"example_two_rules_with_same_name": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_nat_rule" "example_two_rules_with_same_name" {
						edge_gateway_id = cloudavenue_edgegateway.example.id

						name        = "SAMENAME"
						rule_type   = "DNAT"
						description = {{ generate . "description" }}

						external_address         = "89.32.25.10"
						internal_address         = "4.11.11.11"
						priority = 25
						}
							
					resource "cloudavenue_edgegateway_nat_rule" "example_two_rules_with_same_name_2" {
						edge_gateway_id = cloudavenue_edgegateway.example.id

						name        = "SAMENAME"
						rule_type   = "DNAT"
						description = {{ generate . "description2" }}

						# Using primary_ip from edge gateway
						external_address         = "189.32.25.10"
						internal_address         = "5.11.11.11"
						priority = 25
						}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", "SAMENAME"),
						resource.TestCheckResourceAttrSet(resourceName, "external_address"),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "external_address", "89.32.25.10"),
						resource.TestCheckResourceAttr(resourceName, "internal_address", "4.11.11.11"),
						resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
						resource.TestCheckResourceAttr(resourceName, "priority", "25"),
					},
				},
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccEdgeGatewayNATRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&NATRuleResource{}),
	})
}
