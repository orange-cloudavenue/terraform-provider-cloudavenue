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

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/helpers"
)

var _ testsacc.TestACC = &VDCGResource{}

const (
	VDCGResourceName = testsacc.ResourceName("cloudavenue_vdcg")
)

type VDCGResource struct{}

func NewVDCGResourceTest() testsacc.TestACC {
	return &VDCGResource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGResource) GetResourceName() string {
	return VDCGResourceName.String()
}

func (r *VDCGResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetSpecificConfig("example_vdc_group_1"))
	resp.Append(GetResourceConfig()[VDCResourceName]().GetSpecificConfig("example_vdc_group_2"))
	return resp
}

func (r *VDCGResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", helpers.TestIsType(urn.VDCGroup)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg" "example" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						vdc_ids = [
							cloudavenue_vdc.example_vdc_group_1.id,
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "1"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// Update description and add a new vdc_id
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg" "example" {
							name = {{ get . "name" }}
							description = {{ generate . "description" }}
							vdc_ids = [
								cloudavenue_vdc.example_vdc_group_1.id,
								cloudavenue_vdc.example_vdc_group_2.id,
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "2"),
						},
					},
					// update name
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg" "example" {
							name = {{ generate . "name" }}
							description = {{ get . "description" }}
							vdc_ids = [
								cloudavenue_vdc.example_vdc_group_1.id,
								cloudavenue_vdc.example_vdc_group_2.id,
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "2"),
						},
					},
					// remove vdc ids
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg" "example" {
							name = {{ get . "name" }}
							description = {{ get . "description" }}
							vdc_ids = [
								cloudavenue_vdc.example_vdc_group_2.id,
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "1"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateID:     testsacc.GetValueFromTemplate(resourceName, "name"),
						ImportState:       true,
						ImportStateVerify: true,
					},
					{
						ImportStateID:     testsacc.GetValueFromTemplate(resourceName, "id"),
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			}
		},
	}
}

func TestAccVDCGResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGResource{}),
	})
}
