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

var _ testsacc.TestACC = &ELBPoolResource{}

const (
	ELBPoolResourceName = testsacc.ResourceName("cloudavenue_elb_pool")
)

type ELBPoolResource struct{}

func NewELBPoolResourceTest() testsacc.TestACC {
	return &ELBPoolResource{}
}

// GetResourceName returns the name of the resource.
func (r *ELBPoolResource) GetResourceName() string {
	return ELBPoolResourceName.String()
}

func (r *ELBPoolResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_for_elb"))
	return
}

func (r *ELBPoolResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerPool)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_pool" "example" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
						enabled = true
						default_port = 80
						members = {
							targets = [
								{
									ip_address = "192.168.0.1"
									port = 80
								}
							]
						}
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "default_port", "80"),
						resource.TestCheckResourceAttr(resourceName, "members.targets.#", "1"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.targets.*", map[string]string{
							"ip_address": "192.168.0.1",
							"port":       "80",
							// Default values
							"ratio":   "1",
							"enabled": "true",
						}),

						// Default values
						resource.TestCheckResourceAttr(resourceName, "members.graceful_timeout_period", "1"),
						resource.TestCheckResourceAttr(resourceName, "algorithm", "LEAST_CONNECTIONS"),
						resource.TestCheckResourceAttr(resourceName, "health.passive_monitoring_enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "tls.common_name_check_enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "tls.enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "persistence.type", "CLIENT_IP"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Update name and add a new disabled target
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_pool" "example" {
							name = {{ generate . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							default_port = 80
							members = {
								targets = [
									{
										ip_address = "192.168.0.1"
										port = 80
									},
									{
										ip_address = "192.168.0.2"
										port = 8080
										enabled = false
										ratio = 2
									}
								]
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "default_port", "80"),
							resource.TestCheckResourceAttr(resourceName, "members.targets.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.targets.*", map[string]string{
								"ip_address": "192.168.0.1",
								"port":       "80",
								// Default values
								"ratio":   "1",
								"enabled": "true",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.targets.*", map[string]string{
								"ip_address": "192.168.0.2",
								"port":       "8080",
								"ratio":      "2",
								"enabled":    "false",
							}),

							// Default values
							resource.TestCheckResourceAttr(resourceName, "members.graceful_timeout_period", "1"),
							resource.TestCheckResourceAttr(resourceName, "algorithm", "LEAST_CONNECTIONS"),
							resource.TestCheckResourceAttr(resourceName, "health.passive_monitoring_enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "tls.common_name_check_enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "tls.enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "persistence.type", "CLIENT_IP"),
						},
					},
					// * Update remove the disabled target, change algorithm and add health monitors
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_pool" "example" {
							name = {{ generate . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							default_port = 80
							algorithm = "ROUND_ROBIN"
							members = {
								targets = [
									{
										ip_address = "192.168.0.1"
										port = 80
									}
								]
							}
							health = {
								monitors = ["HTTP", "TCP"]
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "default_port", "80"),
							resource.TestCheckResourceAttr(resourceName, "algorithm", "ROUND_ROBIN"),
							resource.TestCheckResourceAttr(resourceName, "health.monitors.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "health.monitors.0", "HTTP"),
							resource.TestCheckResourceAttr(resourceName, "health.monitors.1", "TCP"),
							resource.TestCheckResourceAttr(resourceName, "members.targets.#", "1"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.targets.*", map[string]string{
								"ip_address": "192.168.0.1",
								"port":       "80",
								// Default values
								"ratio":   "1",
								"enabled": "true",
							}),

							// Default values
							resource.TestCheckResourceAttr(resourceName, "members.graceful_timeout_period", "1"),
							resource.TestCheckResourceAttr(resourceName, "health.passive_monitoring_enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "tls.common_name_check_enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "tls.enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "persistence.type", "CLIENT_IP"),
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
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
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
		"example_with_edge_name": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerPool)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_pool" "example_with_edge_name" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_name = data.cloudavenue_edgegateway.example_for_elb.name
						enabled = true
						default_port = 80
						members = {
							targets = [
								{
									ip_address = "192.168.0.1"
									port = 80
								}
							]
						}
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "default_port", "80"),
						resource.TestCheckResourceAttr(resourceName, "members.targets.#", "1"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.targets.*", map[string]string{
							"ip_address": "192.168.0.1",
							"port":       "80",
							// Default values
							"ratio":   "1",
							"enabled": "true",
						}),

						// Default values
						resource.TestCheckResourceAttr(resourceName, "members.graceful_timeout_period", "1"),
						resource.TestCheckResourceAttr(resourceName, "algorithm", "LEAST_CONNECTIONS"),
						resource.TestCheckResourceAttr(resourceName, "health.passive_monitoring_enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "tls.common_name_check_enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "tls.enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "persistence.type", "CLIENT_IP"),
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
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
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
		"example_with_ip_set": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerPool)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayIPSetResourceName]().GetSpecificConfig("example_for_elb"))
					resp.Append(GetResourceConfig()[ORGCertificateLibraryResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_pool" "example_with_ip_set" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
						enabled = true
						default_port = 80
						members = {
							target_group = cloudavenue_edgegateway_ip_set.example_for_elb.id
						}
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "default_port", "80"),
						resource.TestCheckNoResourceAttr(resourceName, "members.targets"),
						resource.TestCheckResourceAttrWith(resourceName, "members.target_group", urn.TestIsType(urn.SecurityGroup)),
						// Default values
						resource.TestCheckResourceAttr(resourceName, "members.graceful_timeout_period", "1"),
						resource.TestCheckResourceAttr(resourceName, "algorithm", "LEAST_CONNECTIONS"),
						resource.TestCheckResourceAttr(resourceName, "health.passive_monitoring_enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "tls.common_name_check_enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "tls.enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "persistence.type", "CLIENT_IP"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Update add TLS
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_pool" "example_with_ip_set" {
							name = {{ get . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							default_port = 80
							members = {
								target_group = cloudavenue_edgegateway_ip_set.example_for_elb.id
							}
							tls = {
								enabled = true
								ca_certificate_refs = [
									cloudavenue_org_certificate_library.example.id
								]
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "default_port", "80"),
							resource.TestCheckNoResourceAttr(resourceName, "members.targets"),
							resource.TestCheckResourceAttrWith(resourceName, "members.target_group", urn.TestIsType(urn.SecurityGroup)),
							resource.TestCheckResourceAttr(resourceName, "tls.enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "tls.common_name_check_enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "tls.ca_certificate_refs.#", "1"),

							// Default values
							resource.TestCheckResourceAttr(resourceName, "members.graceful_timeout_period", "1"),
							resource.TestCheckResourceAttr(resourceName, "algorithm", "LEAST_CONNECTIONS"),
							resource.TestCheckResourceAttr(resourceName, "health.passive_monitoring_enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "persistence.type", "CLIENT_IP"),
						},
					},
					// * Update add persistence
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_pool" "example_with_ip_set" {
							name = {{ get . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							default_port = 80
							members = {
								target_group = cloudavenue_edgegateway_ip_set.example_for_elb.id
							}
							tls = {
								enabled = true
								ca_certificate_refs = [
									cloudavenue_org_certificate_library.example.id
								]
							}
							persistence = {
								type = "TLS"
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "default_port", "80"),
							resource.TestCheckNoResourceAttr(resourceName, "members.targets"),
							resource.TestCheckResourceAttrWith(resourceName, "members.target_group", urn.TestIsType(urn.SecurityGroup)),
							resource.TestCheckResourceAttr(resourceName, "tls.enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "tls.common_name_check_enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "tls.ca_certificate_refs.#", "1"),

							resource.TestCheckResourceAttr(resourceName, "persistence.type", "TLS"),
							resource.TestCheckNoResourceAttr(resourceName, "persistence.value"),

							// Default values
							resource.TestCheckResourceAttr(resourceName, "members.graceful_timeout_period", "1"),
							resource.TestCheckResourceAttr(resourceName, "algorithm", "LEAST_CONNECTIONS"),
							resource.TestCheckResourceAttr(resourceName, "health.passive_monitoring_enabled", "true"),
						},
					},
					// * Update change persistence
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_pool" "example_with_ip_set" {
							name = {{ get . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							default_port = 80
							members = {
								target_group = cloudavenue_edgegateway_ip_set.example_for_elb.id
							}
							tls = {
								enabled = true
								ca_certificate_refs = [
									cloudavenue_org_certificate_library.example.id
								]
							}
							persistence = {
								type = "CUSTOM_HTTP_HEADER"
								value = "X-Custom"
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "default_port", "80"),
							resource.TestCheckNoResourceAttr(resourceName, "members.targets"),
							resource.TestCheckResourceAttrWith(resourceName, "members.target_group", urn.TestIsType(urn.SecurityGroup)),
							resource.TestCheckResourceAttr(resourceName, "tls.enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "tls.common_name_check_enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "tls.ca_certificate_refs.#", "1"),

							resource.TestCheckResourceAttr(resourceName, "persistence.type", "CUSTOM_HTTP_HEADER"),
							resource.TestCheckResourceAttr(resourceName, "persistence.value", "X-Custom"),

							// Default values
							resource.TestCheckResourceAttr(resourceName, "members.graceful_timeout_period", "1"),
							resource.TestCheckResourceAttr(resourceName, "algorithm", "LEAST_CONNECTIONS"),
							resource.TestCheckResourceAttr(resourceName, "health.passive_monitoring_enabled", "true"),
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
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
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
	}
}

func TestAccELBPoolResource(t *testing.T) {
	cleanup := orgCertificateLibraryResourcePreCheck()
	defer cleanup()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ELBPoolResource{}),
	})
}
