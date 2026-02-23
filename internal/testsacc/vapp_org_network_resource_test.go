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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VAppOrgNetworkResource{}

const (
	VAppOrgNetworkResourceName = testsacc.ResourceName("cloudavenue_vapp_org_network")
)

type VAppOrgNetworkResource struct{}

func NewVAppOrgNetworkResourceTest() testsacc.TestACC {
	return &VAppOrgNetworkResource{}
}

// GetResourceName returns the name of the resource.
func (r *VAppOrgNetworkResource) GetResourceName() string {
	return VAppOrgNetworkResourceName.String()
}

func (r *VAppOrgNetworkResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VAppResourceName]().GetDefaultConfig)
	resp.Append(GetResourceConfig()[NetworkRoutedResourceName]().GetDefaultConfig)
	return resp
}

func (r *VAppOrgNetworkResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Network)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
					resource.TestCheckResourceAttrSet(resourceName, "network_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_name"),
					resource.TestCheckNoResourceAttr(resourceName, "vapp_id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vapp_org_network" "example" {
						vapp_name    = cloudavenue_vapp.example.name
						network_name = cloudavenue_network_routed.example.name
						vdc          = cloudavenue_vdc.example.name
					  }`),
					Checks: []resource.TestCheckFunc{},
				},
				// ! Update testing
				// * No update for this resource
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc", "vapp_name", "network_name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},

		// * Test named "example_multiple_networks"
		// Reproduces issue #1202: when multiple org networks are attached to the same vApp
		// in parallel (no explicit depends_on), only the first one actually attaches and
		// the others silently fail, causing the next plan to error with
		// "Unable to find network in the VApp".
		"example_multiple_networks": func(_ context.Context, resourceName string) testsacc.Test {
			// resourceName = "cloudavenue_vapp_org_network.example_multiple_networks" (first resource)
			const secondResourceName = "cloudavenue_vapp_org_network.example_multiple_networks_second"
			const secondNetworkResourceName = "cloudavenue_network_routed.example_multiple_networks_second"

			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(func() map[string]testsacc.TFData {
						return map[string]testsacc.TFData{
							secondNetworkResourceName: testsacc.GenerateFromTemplate(secondNetworkResourceName, `
							resource "cloudavenue_network_routed" "example_multiple_networks_second" {
								name        = {{ generate . "name_second" }}
								edge_gateway_id = cloudavenue_edgegateway.example.id

								gateway       = "192.168.2.254"
								prefix_length = 24

								dns1 = "1.1.1.1"
								dns2 = "8.8.8.8"

								dns_suffix = "example"

								static_ip_pool = [
								  {
									start_address = "192.168.2.10"
									end_address   = "192.168.2.20"
								  }
								]
							}`),
						}
					})
					return resp
				},
				CommonChecks: []resource.TestCheckFunc{
					// Checks for the first org network
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Network)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
					resource.TestCheckResourceAttrSet(resourceName, "network_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_name"),
					// Checks for the second org network
					resource.TestCheckResourceAttrWith(secondResourceName, "id", urn.TestIsType(urn.Network)),
					resource.TestCheckResourceAttrSet(secondResourceName, "vdc"),
					resource.TestCheckResourceAttrSet(secondResourceName, "network_name"),
					resource.TestCheckResourceAttrSet(secondResourceName, "vapp_name"),
				},
				// ! Create testing — both resources are created in parallel (no depends_on)
				// to reproduce the race condition described in issue #1202.
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vapp_org_network" "example_multiple_networks" {
						vapp_name    = cloudavenue_vapp.example.name
						network_name = cloudavenue_network_routed.example.name
						vdc          = cloudavenue_vdc.example.name
					}

					resource "cloudavenue_vapp_org_network" "example_multiple_networks_second" {
						vapp_name    = cloudavenue_vapp.example.name
						network_name = cloudavenue_network_routed.example_multiple_networks_second.name
						vdc          = cloudavenue_vdc.example.name
					}`),
					// Verify that the two resources received distinct IDs — a silent overwrite
					// would leave both pointing to the same network URN.
					Checks: []resource.TestCheckFunc{
						func(s *terraform.State) error {
							first, ok := s.RootModule().Resources[resourceName]
							if !ok {
								return fmt.Errorf("resource %s not found in state", resourceName)
							}
							second, ok := s.RootModule().Resources[secondResourceName]
							if !ok {
								return fmt.Errorf("resource %s not found in state", secondResourceName)
							}
							idFirst := first.Primary.ID
							idSecond := second.Primary.ID
							if idFirst == idSecond {
								return fmt.Errorf("expected distinct IDs for both org networks, got the same: %s", idFirst)
							}
							return nil
						},
					},
				},
				// ! Update testing
				// Run a plan-only step to trigger Read/refresh on both resources.
				// Before the bug fix, Read calls findOrgNetwork and errors with
				// "Unable to find network in the VApp" for the silently dropped resource.
				// After the fix, both reads succeed and the plan is empty.
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vapp_org_network" "example_multiple_networks" {
							vapp_name    = cloudavenue_vapp.example.name
							network_name = cloudavenue_network_routed.example.name
							vdc          = cloudavenue_vdc.example.name
						}

						resource "cloudavenue_vapp_org_network" "example_multiple_networks_second" {
							vapp_name    = cloudavenue_vapp.example.name
							network_name = cloudavenue_network_routed.example_multiple_networks_second.name
							vdc          = cloudavenue_vdc.example.name
						}`),
						TFAdvanced: testsacc.TFAdvanced{
							PlanOnly:           true,
							ExpectNonEmptyPlan: false,
						},
						Checks: []resource.TestCheckFunc{},
					},
				},
				// ! Imports testing
				// * No imports for this test
			}
		},
	}
	// TODO: ADD Test with VDC Group
}

func TestAccOrgNetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VAppOrgNetworkResource{}),
	})
}
