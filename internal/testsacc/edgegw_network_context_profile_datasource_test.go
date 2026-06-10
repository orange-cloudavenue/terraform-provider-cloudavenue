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

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayNetworkContextProfileDatasource) GetResourceName() string {
	return EdgeGatewayNetworkContextProfileDatasourceName.String()
}

func (r *EdgeGatewayNetworkContextProfileDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayNetworkContextProfileDatasource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// Lookup a well-known SYSTEM profile by name.
		testNameExample: func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_network_context_profile" "example" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
						name              = "SSL"
					}`,
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
		// Lookup by edge gateway ID instead of name.
		"example_by_edge_gateway_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_network_context_profile" "example_by_edge_gateway_id" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						name            = "CIFS"
					}`,
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
		// Lookup by the profile ID directly.
		"example_by_profile_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_network_context_profile" "example_by_profile_id" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
						id                = "urn:vcloud:networkContextProfile:45b67f48-0e35-3e97-98c7-ace4276a17dc"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "id", "urn:vcloud:networkContextProfile:45b67f48-0e35-3e97-98c7-ace4276a17dc"),
						resource.TestCheckResourceAttr(resourceName, "name", "AMQP"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
					},
				},
				Destroy: true,
			}
		},
		// Verify that a non-existent profile returns a proper error.
		"example_not_found": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_network_context_profile" "example_not_found" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
						name              = "THIS_PROFILE_DOES_NOT_EXIST"
					}`,
					TFAdvanced: testsacc.TFAdvanced{
						ExpectError: regexp.MustCompile(`Network Context Profile not found`),
					},
				},
				Destroy: true,
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
