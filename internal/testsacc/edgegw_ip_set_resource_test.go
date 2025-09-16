/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &EdgeGatewayIPSetResource{}

const (
	EdgeGatewayIPSetResourceName = testsacc.ResourceName("cloudavenue_edgegateway_ip_set")
)

type EdgeGatewayIPSetResource struct{}

func NewEdgeGatewayIPSetResourceTest() testsacc.TestACC {
	return &EdgeGatewayIPSetResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayIPSetResource) GetResourceName() string {
	return EdgeGatewayIPSetResourceName.String()
}

func (r *EdgeGatewayIPSetResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *EdgeGatewayIPSetResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_ip_set" "example" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						ip_addresses = [
							"192.168.1.1",
							"192.168.1.2",
						]
						edge_gateway_name = cloudavenue_edgegateway.example.name
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "2"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_ip_set" "example" {
							name = {{ get . "name" }}
							description = {{ generate . "description" }}
							ip_addresses = [
								"192.168.1.1",
								"192.168.1.2",
								"192.168.1.3",
							]
							edge_gateway_name = cloudavenue_edgegateway.example.name
						}`),
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
						ImportState:          true,
					},
				},
				Destroy: true,
			}
		},
		"example_for_elb": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_for_elb"))
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_ip_set" "example_for_elb" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						ip_addresses = [
							"192.168.1.1",
							"192.168.1.2",
							"192.168.1.3",
						]
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
					}`),
				},
			}
		},
	}
}

func TestAccEdgeGatewayIPSetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayIPSetResource{}),
	})
}
