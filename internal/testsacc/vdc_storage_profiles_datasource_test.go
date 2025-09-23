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

var _ testsacc.TestACC = &VDCStorageProfilesDataSource{}

const (
	VDCStorageProfilesDataSourceName = testsacc.ResourceName("data.cloudavenue_vdc_storage_profiles")
)

type VDCStorageProfilesDataSource struct{}

func NewVDCStorageProfilesDataSourceTest() testsacc.TestACC {
	return &VDCStorageProfilesDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCStorageProfilesDataSource) GetResourceName() string {
	return VDCStorageProfilesDataSourceName.String()
}

func (r *VDCStorageProfilesDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetSpecificConfig("example_storage_profiles"))
	return resp
}

func (r *VDCStorageProfilesDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "vdc_id", helpers.TestIsType(urn.VDC)),
					resource.TestMatchResourceAttr(resourceName, "vdc_name", regex.VDCNameRegex()),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					data "cloudavenue_vdc_storage_profiles" "example" {
						vdc_id = cloudavenue_vdc.example_storage_profiles.id
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.#"),
						resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.class", "gold"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.limit"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.used"),
						resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.default", "true"),
						resource.TestCheckResourceAttr(resourceName, "storage_profiles.1.class", "silver"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.1.id"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.1.limit"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.1.used"),
						resource.TestCheckResourceAttr(resourceName, "storage_profiles.1.default", "false"),
					},
				},
			}
		},
	}
}

func TestAccVDCStorageProfilesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCStorageProfilesDataSource{}),
	})
}
