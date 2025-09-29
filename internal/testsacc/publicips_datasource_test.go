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

var _ testsacc.TestACC = &PublicIPsDataSource{}

const (
	PublicIPsDataSourceName = testsacc.ResourceName("data.cloudavenue_publicips")
)

type PublicIPsDataSource struct{}

func NewPublicIPsDataSourceTest() testsacc.TestACC {
	return &PublicIPsDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *PublicIPsDataSource) GetResourceName() string {
	return PublicIPsDataSourceName.String()
}

func (r *PublicIPsDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[PublicIPResourceName]().GetDefaultConfig)
	return resp
}

func (r *PublicIPsDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_publicips" "example" {
						depends_on = [
							cloudavenue_publicip.example
						]

					 	edge_gateway_id = cloudavenue_edgegateway.example.id
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						// check if public_ips are not empty
						resource.TestCheckResourceAttrSet(resourceName, "public_ips.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "public_ips.0.public_ip"),
					},
				},
			}
		},
	}
}

func TestAccPublicIPsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&PublicIPsDataSource{}),
	})
}
