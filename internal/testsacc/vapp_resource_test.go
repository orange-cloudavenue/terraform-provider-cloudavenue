/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VAppResource{}

const (
	VAppResourceName = testsacc.ResourceName("cloudavenue_vapp")
)

type VAppResource struct{}

func NewVAppResourceTest() testsacc.TestACC {
	return &VAppResource{}
}

// GetResourceName returns the name of the resource.
func (r *VAppResource) GetResourceName() string {
	return VAppResourceName.String()
}

func (r *VAppResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
	return resp
}

func (r *VAppResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VAPP)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vapp" "example" {
						name        = {{ generate . "name" }}
						description = {{ generate . "description" }}
						vdc 		= cloudavenue_vdc.example.name
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "lease.runtime_lease_in_sec", "0"),
						resource.TestCheckResourceAttr(resourceName, "lease.storage_lease_in_sec", "0"),
						resource.TestCheckNoResourceAttr(resourceName, "guest_properties.#"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vapp" "example" {
							name        = {{ get . "name" }}
							description = {{ generate . "description" }}
							vdc 		= cloudavenue_vdc.example.name

							lease = {
								runtime_lease_in_sec = 3600
								storage_lease_in_sec = 3600
							}
						
							guest_properties = {
								"key" = "Value"
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "lease.runtime_lease_in_sec", "3600"),
							resource.TestCheckResourceAttr(resourceName, "lease.storage_lease_in_sec", "3600"),
							resource.TestCheckResourceAttr(resourceName, "guest_properties.key", "Value"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vapp" "example" {
							name        = {{ get . "name" }}
							description = {{ get . "description" }}
							vdc 		= cloudavenue_vdc.example.name

							lease = {
								runtime_lease_in_sec = 36000
								storage_lease_in_sec = 360000
							}
						
							guest_properties = {
								"key" = "Value"
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "lease.runtime_lease_in_sec", "36000"),
							resource.TestCheckResourceAttr(resourceName, "lease.storage_lease_in_sec", "360000"),
							resource.TestCheckResourceAttr(resourceName, "guest_properties.key", "Value"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccVAppResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VAppResource{}),
	})
}
