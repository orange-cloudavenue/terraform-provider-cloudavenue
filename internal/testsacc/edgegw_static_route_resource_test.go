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

var _ testsacc.TestACC = &EdgeGatewayStaticRouteResource{}

const EdgeGatewayStaticRouteResourceName = testsacc.ResourceName("cloudavenue_edgegateway_static_route")

type EdgeGatewayStaticRouteResource struct{}

func NewEdgeGatewayStaticRouteResourceTest() testsacc.TestACC {
	return &EdgeGatewayStaticRouteResource{}
}

func (r *EdgeGatewayStaticRouteResource) GetResourceName() string {
	return EdgeGatewayStaticRouteResourceName.String()
}

func (r *EdgeGatewayStaticRouteResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayStaticRouteResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					// Check that the resource has been created and has an ID formatted as uuid v4
					resource.TestCheckResourceAttrWith(resourceName, "id", ToValidate("uuid4")),
					// Check that the edge_gateway_id is a valid urn of type Gateway
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
				},
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
		resource "cloudavenue_edgegateway_static_route" "example" {
		  edge_gateway_id = cloudavenue_edgegateway.example.id
		  name = {{ generate . "name" }}
		  network_cidr = "192.168.1.0/24"
		  next_hops = [
		{ ip_address = "192.168.1.254" }
		  ]
		}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckNoResourceAttr(resourceName, "description"),
						resource.TestCheckResourceAttr(resourceName, "network_cidr", "192.168.1.0/24"),
						resource.TestCheckResourceAttr(resourceName, "next_hops.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "next_hops.0.ip_address", "192.168.1.254"),
						resource.TestCheckResourceAttr(resourceName, "next_hops.0.admin_distance", "1"),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
		resource "cloudavenue_edgegateway_static_route" "example" {
		  edge_gateway_id = cloudavenue_edgegateway.example.id
		  name = {{ generate . "name" }}
		  description = {{ generate . "description" }}
		  network_cidr = "192.168.2.0/24"
		  next_hops = [
		{ ip_address = "192.168.2.254" },
		{ ip_address = "192.168.2.253", admin_distance = 2 }
		  ]
		}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "network_cidr", "192.168.2.0/24"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.0.ip_address", "192.168.2.254"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.0.admin_distance", "1"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.1.ip_address", "192.168.2.253"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.1.admin_distance", "2"),
						},
					},
				},
				Imports: []testsacc.TFImport{
					{ImportStateIDBuilder: []string{"edge_gateway_id", "id"}, ImportState: true, ImportStateVerify: true},
				},
				Destroy: true,
			}
		},
		"example_with_vdc_group": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetSpecificConfig("example_with_vdc_group"))
					return resp
				},
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", ToValidate("uuid4")),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
				},
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_static_route" "example_with_vdc_group" {
					  edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
					  name = {{ generate . "name" }}
					  network_cidr = "192.168.1.0/24"
					  next_hops = [
					    { ip_address = "192.168.1.254" }
					  ]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckNoResourceAttr(resourceName, "description"),
						resource.TestCheckResourceAttr(resourceName, "network_cidr", "192.168.1.0/24"),
						resource.TestCheckResourceAttr(resourceName, "next_hops.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "next_hops.0.ip_address", "192.168.1.254"),
						resource.TestCheckResourceAttr(resourceName, "next_hops.0.admin_distance", "1"),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
							resource "cloudavenue_edgegateway_static_route" "example_with_vdc_group" {
							  edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
							  name = {{ generate . "name" }}
							  description = {{ generate . "description" }}
							  network_cidr = "192.168.2.0/24"
							  next_hops = [
							    { ip_address = "192.168.2.254" },
							    { ip_address = "192.168.2.253", admin_distance = 2 }
							  ]
							}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "network_cidr", "192.168.2.0/24"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.0.ip_address", "192.168.2.254"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.0.admin_distance", "1"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.1.ip_address", "192.168.2.253"),
							resource.TestCheckResourceAttr(resourceName, "next_hops.1.admin_distance", "2"),
						},
					},
				},
				Imports: []testsacc.TFImport{
					{ImportStateIDBuilder: []string{"edge_gateway_name", "name"}, ImportState: true, ImportStateVerify: true},
				},
			}
		},
	}
}

func TestAccEdgeGatewayStaticRouteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayStaticRouteResource{}),
	})
}
