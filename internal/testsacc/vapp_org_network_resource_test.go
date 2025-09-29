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
