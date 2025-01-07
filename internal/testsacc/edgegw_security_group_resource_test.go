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

var _ testsacc.TestACC = &EdgeGatewaySecurityGroupResource{}

const (
	EdgeGatewaySecurityGroupResourceName = testsacc.ResourceName("cloudavenue_edgegateway_security_group")
)

type EdgeGatewaySecurityGroupResource struct{}

func NewEdgeGatewaySecurityGroupResourceTest() testsacc.TestACC {
	return &EdgeGatewaySecurityGroupResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewaySecurityGroupResource) GetResourceName() string {
	return EdgeGatewaySecurityGroupResourceName.String()
}

func (r *EdgeGatewaySecurityGroupResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[NetworkRoutedResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewaySecurityGroupResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.SecurityGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_security_group" "example" {
						name            = {{ generate . "name" }}
						description     = {{ generate . "description" }}
						
						edge_gateway_id = cloudavenue_edgegateway.example.id
						member_org_network_ids = [
						  cloudavenue_network_routed.example.id
						]
					  }`),
					Checks: []resource.TestCheckFunc{
						// id
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "1"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_security_group" "example" {
							name            = {{ generate . "name" }}
							description     = {{ generate . "description" }}
							
							edge_gateway_id = cloudavenue_edgegateway.example.id
							member_org_network_ids = [
							  cloudavenue_network_routed.example.id
							]
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "1"),
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
	}
}

func TestAccEdgeGatewaySecurityGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewaySecurityGroupResource{}),
	})
}
