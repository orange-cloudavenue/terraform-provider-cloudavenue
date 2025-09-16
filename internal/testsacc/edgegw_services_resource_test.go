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

var _ testsacc.TestACC = &EdgeGatewayServicesResource{}

const (
	EdgeGatewayServicesResourceName = testsacc.ResourceName("cloudavenue_edgegateway_services")
)

type EdgeGatewayServicesResource struct{}

func NewEdgeGatewayServicesResourceTest() testsacc.TestACC {
	return &EdgeGatewayServicesResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayServicesResource) GetResourceName() string {
	return EdgeGatewayServicesResourceName.String()
}

func (r *EdgeGatewayServicesResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayServicesResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_services" "example" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
						resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
						resource.TestCheckResourceAttrSet(resourceName, "network"),
					},
				},
				// ! Update is not supported
				// ! Import is not supported. Create a new resource instead.
			}
		},
	}
}

func TestAccEdgeGatewayServicesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayServicesResource{}),
	})
}
