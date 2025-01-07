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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

const testAccProfilesDataSourceConfig = `
data "cloudavenue_storage_profiles" "example" {
}
`

func TestAccProfilesDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_storage_profiles.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccProfilesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(dataSourceName, "id", urn.TestIsType(urn.VDCStorageProfile)),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.vdc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.limit"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.used_storage"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.iops_allocated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.units"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.iops_limiting_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.maximum_disk_iops"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.default_disk_iops"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.disk_iops_per_gb_max"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.iops_limit"),
				),
			},
		},
	})
}
