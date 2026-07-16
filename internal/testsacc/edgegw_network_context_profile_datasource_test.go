/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &EdgeGatewayNetworkContextProfileDatasource{}

const (
	EdgeGatewayNetworkContextProfileDatasourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_network_context_profile")
)

type EdgeGatewayNetworkContextProfileDatasource struct{}

func NewEdgeGatewayNetworkContextProfileDatasourceTest() testsacc.TestACC {
	return &EdgeGatewayNetworkContextProfileDatasource{}
}

func (r *EdgeGatewayNetworkContextProfileDatasource) GetResourceName() string {
	return EdgeGatewayNetworkContextProfileDatasourceName.String()
}

func (r *EdgeGatewayNetworkContextProfileDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayNetworkContextProfileDatasource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"system_profile_by_edge_gateway_name": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_edgegateway_network_context_profile" "system_profile_by_edge_gateway_name" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
						name              = "SSL"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", "SSL"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
						resource.TestCheckResourceAttrSet(resourceName, "description"),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					},
				},
				Destroy: true,
			}
		},
		"system_profile_by_edge_gateway_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_edgegateway_network_context_profile" "system_profile_by_edge_gateway_id" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						name            = "CIFS"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", "CIFS"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
						resource.TestCheckResourceAttrSet(resourceName, "description"),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
					},
				},
				Destroy: true,
			}
		},
		"custom_erp_application": func(_ context.Context, resourceName string) testsacc.Test {
			const (
				depResourceLabel = "existing_erp"
				description      = "Test profile for lookup by ID"
			)

			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
					resp.Append(func() map[string]testsacc.TFData {
						return map[string]testsacc.TFData{
							"cloudavenue_edgegateway_network_context_profile." + depResourceLabel: testsacc.TFData(`
							resource "cloudavenue_edgegateway_network_context_profile" "existing_erp" {
							  edge_gateway_name = cloudavenue_edgegateway.example.name
							  name              = {{ generate . "name" }}
							  description       = "Test profile for lookup by ID"

							  app_id = {
								values = ["SSL"]
							  }
							}
							`),
						}
					})
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(fmt.Sprintf(`
					data "cloudavenue_edgegateway_network_context_profile" "custom_erp_application" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
						id                = cloudavenue_edgegateway_network_context_profile.%s.id
					}`,
						depResourceLabel,
					)),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttr(resourceName, "scope", "TENANT"),
						resource.TestCheckResourceAttr(resourceName, "description", description),
					},
				},
				Destroy: true,
			}
		},
		"not_found_typo": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_edgegateway_network_context_profile" "not_found_typo" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
						name              = "THIS_PROFILE_DOES_NOT_EXIST"
					}`),
					TFAdvanced: testsacc.TFAdvanced{
						ExpectError: regexp.MustCompile(`Network Context Profile not found`),
					},
				},
				Destroy: false,
			}
		},
	}
}

func TestAccEdgeGatewayNetworkContextProfileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayNetworkContextProfileDatasource{}),
	})
}
