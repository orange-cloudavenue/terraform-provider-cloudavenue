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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &EdgeGatewayNetworkContextProfileResource{}

const (
	EdgeGatewayNetworkContextProfileResourceName = testsacc.ResourceName("cloudavenue_edgegateway_network_context_profile")
)

type EdgeGatewayNetworkContextProfileResource struct{}

func NewEdgeGatewayNetworkContextProfileResourceTest() testsacc.TestACC {
	return &EdgeGatewayNetworkContextProfileResource{}
}

func (r *EdgeGatewayNetworkContextProfileResource) GetResourceName() string {
	return EdgeGatewayNetworkContextProfileResourceName.String()
}

func (r *EdgeGatewayNetworkContextProfileResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayNetworkContextProfileResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// Basic profile with multiple App IDs.
		testNameExample: func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.NetworkContextProfile)),
					resource.TestCheckResourceAttr(resourceName, "scope", "TENANT"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
				},
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_network_context_profile" "example" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
						name              = {{ generate . "name" }}
						description       = {{ generate . "description" }}
						app_id = {
							values = ["SSH", "DNS"]
						}
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "app_id.values.#", "2"),
						resource.TestCheckTypeSetElemAttr(resourceName, "app_id.values.*", "SSH"),
						resource.TestCheckTypeSetElemAttr(resourceName, "app_id.values.*", "DNS"),
					},
				},
				// Update: switch to SSL with TLS_VERSION sub-attribute.
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_network_context_profile" "example" {
							edge_gateway_name = cloudavenue_edgegateway.example.name
							name              = {{ get . "name" }}
							description       = {{ generate . "description" }}
							app_id = {
								values = ["SSL"]
								sub_attribute = [
									{
										type   = "TLS_VERSION"
										values = ["TLS_V12", "TLS_V13"]
									}
								]
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "app_id.values.#", "1"),
							resource.TestCheckTypeSetElemAttr(resourceName, "app_id.values.*", "SSL"),
							resource.TestCheckResourceAttr(resourceName, "app_id.sub_attribute.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "app_id.sub_attribute.0.type", "TLS_VERSION"),
						},
					},
				},
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{testAttrEdgeGatewayName, "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccEdgeGatewayNetworkContextProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayNetworkContextProfileResource{}),
	})
}
