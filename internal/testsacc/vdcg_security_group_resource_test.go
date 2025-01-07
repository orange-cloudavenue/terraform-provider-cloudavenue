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

var _ testsacc.TestACC = &VDCGSecurityGroupResource{}

const (
	VDCGSecurityGroupResourceName = testsacc.ResourceName("cloudavenue_vdcg_security_group")
)

type VDCGSecurityGroupResource struct{}

func NewVDCGSecurityGroupResourceTest() testsacc.TestACC {
	return &VDCGSecurityGroupResource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGSecurityGroupResource) GetResourceName() string {
	return VDCGSecurityGroupResourceName.String()
}

func (r *VDCGSecurityGroupResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *VDCGSecurityGroupResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
					return
				},
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.SecurityGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_security_group" "example" {
						vdc_group_name  = cloudavenue_vdcg.example.name
						name            = {{ generate . "name" }}
						description     = {{ generate . "description" }}
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "0"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Update name
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_security_group" "example" {
							vdc_group_name  = cloudavenue_vdcg.example.name
							name            = {{ generate . "name" }}
							description     = {{ get . "description" }}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "0"),
						},
					},
					// * Update description
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_security_group" "example" {
							vdc_group_name  = cloudavenue_vdcg.example.name
							name            = {{ get . "name" }}
							description     = {{ generate . "description" }}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "0"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_with_networks": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGNetworkIsolatedResourceName]().GetDefaultConfig)
					resp.Append(GetResourceConfig()[VDCGNetworkIsolatedResourceName]().GetSpecificConfig("example_without_static_ip_pool"))
					return
				},
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.SecurityGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_security_group" "example_with_networks" {
						vdc_group_id  = cloudavenue_vdcg_network_isolated.example.vdc_group_id
						name            = {{ generate . "name" }}
						description     = {{ generate . "description" }}
						member_org_network_ids = [
							cloudavenue_vdcg_network_isolated.example.id
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "1"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Add routed network
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_security_group" "example_with_networks" {
							vdc_group_name  = cloudavenue_vdcg_network_isolated.example.vdc_group_name
							name            = {{ get . "name" }}
							description     = {{ get . "description" }}
							member_org_network_ids = [
								cloudavenue_vdcg_network_isolated.example.id,
								cloudavenue_vdcg_network_isolated.example_without_static_ip_pool.id
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "2"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccVDCGSecurityGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGSecurityGroupResource{}),
	})
}
