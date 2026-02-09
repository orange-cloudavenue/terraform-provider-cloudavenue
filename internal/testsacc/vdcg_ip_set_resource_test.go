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

var _ testsacc.TestACC = &VDCGIPSetResource{}

const (
	VDCGIPSetResourceName = testsacc.ResourceName("cloudavenue_vdcg_ip_set")
)

type VDCGIPSetResource struct{}

func NewVDCGIPSetResourceTest() testsacc.TestACC {
	return &VDCGIPSetResource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGIPSetResource) GetResourceName() string {
	return VDCGIPSetResourceName.String()
}

func (r *VDCGIPSetResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *VDCGIPSetResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.SecurityGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_ip_set" "example" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						ip_addresses = [
							"192.168.1.1",
							"192.168.1.2",
						]
						vdc_group_name = cloudavenue_vdcg.example.name
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "2"),
						resource.TestCheckTypeSetElemAttr(resourceName, "ip_addresses.*", "192.168.1.1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "ip_addresses.*", "192.168.1.2"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_ip_set" "example" {
							name = {{ get . "name" }}
							description = {{ generate . "description" }}
							ip_addresses = [
								"192.168.1.1",
								"192.168.1.2",
								"192.168.1.3",
							]
							vdc_group_name = cloudavenue_vdcg.example.name
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "3"),
							resource.TestCheckTypeSetElemAttr(resourceName, "ip_addresses.*", "192.168.1.1"),
							resource.TestCheckTypeSetElemAttr(resourceName, "ip_addresses.*", "192.168.1.2"),
							resource.TestCheckTypeSetElemAttr(resourceName, "ip_addresses.*", "192.168.1.3"),
						},
					},
					// * Empty ip_addresses and description
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_ip_set" "example" {
							name = {{ generate . "name" }}
							vdc_group_name = cloudavenue_vdcg.example.name
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckNoResourceAttr(resourceName, "description"),
							resource.TestCheckNoResourceAttr(resourceName, "ip_addresses"),
						},
					},
					// * Update name
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_ip_set" "example" {
							name = {{ generate . "name" }}
							description = {{ get . "description" }}
							ip_addresses = [
								"192.168.1.1",
								"192.168.1.2",
								"192.168.1.3",
							]
							vdc_group_name = cloudavenue_vdcg.example.name
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "3"),
							resource.TestCheckTypeSetElemAttr(resourceName, "ip_addresses.*", "192.168.1.1"),
							resource.TestCheckTypeSetElemAttr(resourceName, "ip_addresses.*", "192.168.1.2"),
							resource.TestCheckTypeSetElemAttr(resourceName, "ip_addresses.*", "192.168.1.3"),
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

func TestAccVDCGIPSetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGIPSetResource{}),
	})
}
