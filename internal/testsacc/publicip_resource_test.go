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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &PublicIPResource{}

const (
	PublicIPResourceName = testsacc.ResourceName("cloudavenue_publicip")
)

// regexpIPv4 matches a valid IPv4 address.
var regexpIPv4 = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)

type PublicIPResource struct{}

func NewPublicIPResourceTest() testsacc.TestACC {
	return &PublicIPResource{}
}

// GetResourceName returns the name of the resource.
func (r *PublicIPResource) GetResourceName() string {
	return PublicIPResourceName.String()
}

func (r *PublicIPResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return resp
}

func (r *PublicIPResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		testNameExample: func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(resourceName, "public_ip", regexpIPv4),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_publicip" "example" {
					  edge_gateway_id = cloudavenue_edgegateway.example.id
					}`,
					Checks: []resource.TestCheckFunc{},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_publicip" "example" {
						  edge_gateway_name = cloudavenue_edgegateway.example.name
					}`,
						Checks: []resource.TestCheckFunc{},
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
				},
			}
		},

		// Regression test for https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1233
		// Ensure that creating a public IP using edge_gateway_id (URN) referencing a
		// cloudavenue_edgegateway resource does not fail with "Edge not found" (err-0009).
		// Root cause: the SDK expects a bare UUID, not the full URN — the provider was
		// previously passing the URN directly to PublicIP.New().
		"example_with_edge_gateway_id_from_resource": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(resourceName, "public_ip", regexpIPv4),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
				},
				// ! Create testing — reproduces the exact config from issue #1233
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_publicip" "example_with_edge_gateway_id_from_resource" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
					}`,
					Checks: []resource.TestCheckFunc{},
				},
			}
		},

		"example_with_edge_name": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_publicip" "example_with_edge_name" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestMatchResourceAttr(resourceName, "public_ip", regexpIPv4),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccPublicIPResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&PublicIPResource{}),
	})
}
