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

	"github.com/orange-cloudavenue/common-go/regex"
	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/helpers"
)

var _ testsacc.TestACC = &VDCStorageProfileDataSource{}

const (
	VDCStorageProfileDataSourceName = testsacc.ResourceName("data.cloudavenue_vdc_storage_profile")
)

type VDCStorageProfileDataSource struct{}

func NewVDCStorageProfileDataSourceTest() testsacc.TestACC {
	return &VDCStorageProfileDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCStorageProfileDataSource) GetResourceName() string {
	return VDCStorageProfileDataSourceName.String()
}

func (r *VDCStorageProfileDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetSpecificConfig("example_storage_profiles"))
	return resp
}

func (r *VDCStorageProfileDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", helpers.TestIsType(urn.VDCStorageProfile)),
					resource.TestCheckResourceAttrWith(resourceName, "vdc_id", helpers.TestIsType(urn.VDC)),
					resource.TestMatchResourceAttr(resourceName, "vdc_name", regex.VDCNameRegex()),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					data "cloudavenue_vdc_storage_profile" "example" {
						vdc_id = cloudavenue_vdc.example_storage_profiles.id
						class  = "gold"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "class", "gold"),
						resource.TestCheckResourceAttr(resourceName, "limit", "500"),
						resource.TestCheckResourceAttrSet(resourceName, "used"),
						resource.TestCheckResourceAttr(resourceName, "default", "true"),
					},
				},
				// ! Update testing
				// Use vdc_name instead of vdc_id
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						data "cloudavenue_vdc_storage_profile" "example" {
							vdc_name = cloudavenue_vdc.example_storage_profiles.name
							class    = "silver"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "class", "silver"),
							resource.TestCheckResourceAttr(resourceName, "limit", "500"),
							resource.TestCheckResourceAttrSet(resourceName, "used"),
							resource.TestCheckResourceAttr(resourceName, "default", "false"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						data "cloudavenue_vdc_storage_profiles" "example" {
							vdc_id = cloudavenue_vdc.example_storage_profiles.id
						}
						data "cloudavenue_vdc_storage_profile" "example" {
							vdc_name = cloudavenue_vdc.example_storage_profiles.name
							id    = data.cloudavenue_vdc_storage_profiles.example.storage_profiles[1].id
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "class", "silver"),
							resource.TestCheckResourceAttr(resourceName, "limit", "500"),
							resource.TestCheckResourceAttrSet(resourceName, "used"),
							resource.TestCheckResourceAttr(resourceName, "default", "false"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						data "cloudavenue_vdc_storage_profiles" "example" {
							vdc_name = cloudavenue_vdc.example_storage_profiles.name
						}
						data "cloudavenue_vdc_storage_profile" "example" {
							vdc_id = cloudavenue_vdc.example_storage_profiles.id
							class    = data.cloudavenue_vdc_storage_profiles.example.storage_profiles[0].class
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "class", "gold"),
							resource.TestCheckResourceAttr(resourceName, "limit", "500"),
							resource.TestCheckResourceAttrSet(resourceName, "used"),
							resource.TestCheckResourceAttr(resourceName, "default", "true"),
						},
					},
				},
			}
		},
	}
}

func TestAccVDCStorageProfileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCStorageProfileDataSource{}),
	})
}
