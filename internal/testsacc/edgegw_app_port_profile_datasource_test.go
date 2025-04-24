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

var _ testsacc.TestACC = &EdgeGatewayAppPortProfileDatasource{}

const (
	EdgeGatewayAppPortProfileDatasourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_app_port_profile")
)

type EdgeGatewayAppPortProfileDatasource struct{}

func NewEdgeGatewayAppPortProfileDatasourceTest() testsacc.TestACC {
	return &EdgeGatewayAppPortProfileDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayAppPortProfileDatasource) GetResourceName() string {
	return EdgeGatewayAppPortProfileDatasourceName.String()
}

func (r *EdgeGatewayAppPortProfileDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *EdgeGatewayAppPortProfileDatasource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
						name = cloudavenue_edgegateway_app_port_profile.example.name
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultChecks(),
				},
				Destroy: true,
			}
		},
		"example_by_id": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_by_id" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						id = cloudavenue_edgegateway_app_port_profile.example.id
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultChecks(),
				},
				Destroy: true,
			}
		},
		"example_provider_scope": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_provider_scope" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						name = "BKP_TCP_bpcd"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.AppPortProfile)),
						resource.TestCheckResourceAttr(resourceName, "name", "BKP_TCP_bpcd"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.protocol", "TCP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.ports.0", "13782"),
						resource.TestCheckResourceAttr(resourceName, "scope", "PROVIDER"),
					},
				},
				Destroy: true,
			}
		},
		"example_system_scope": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_system_scope" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						name = "HTTP"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.AppPortProfile)),
						resource.TestCheckResourceAttr(resourceName, "name", "HTTP"),
						resource.TestCheckResourceAttr(resourceName, "description", "HTTP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.protocol", "TCP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.ports.0", "80"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
					},
				},
				Destroy: true,
			}
		},
		"example_two_app_ports_with_same_name": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetSpecificConfig("example_http_scope_tenant"))
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_two_app_ports_with_same_name" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						name = cloudavenue_edgegateway_app_port_profile.example_http_scope_tenant.name
					}`,
					TFAdvanced: testsacc.TFAdvanced{
						ExpectError: regexp.MustCompile(`Multiple App Port Profiles found`),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						data "cloudavenue_edgegateway_app_port_profile" "example_two_app_ports_with_same_name" {
							edge_gateway_id = cloudavenue_edgegateway.example.id
							name = cloudavenue_edgegateway_app_port_profile.example_http_scope_tenant.name
							scope = "TENANT"
						}`,
						Checks: GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetSpecificChecks("example_http_scope_tenant"),
					},
					{
						TFConfig: `
							data "cloudavenue_edgegateway_app_port_profile" "example_two_app_ports_with_same_name" {
								edge_gateway_id = cloudavenue_edgegateway.example.id
								name = cloudavenue_edgegateway_app_port_profile.example_http_scope_tenant.name
								scope = "SYSTEM"
							}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.AppPortProfile)),
							resource.TestCheckResourceAttr(resourceName, "name", "HTTP"),
							resource.TestCheckResourceAttr(resourceName, "description", "HTTP"),
							resource.TestCheckResourceAttr(resourceName, "app_ports.0.protocol", "TCP"),
							resource.TestCheckResourceAttr(resourceName, "app_ports.0.ports.0", "80"),
							resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
						},
					},
				},
				Destroy: true,
			}
		},
	}
}

func TestAccEdgeGatewayAppPortProfileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayAppPortProfileDatasource{}),
	})
}
