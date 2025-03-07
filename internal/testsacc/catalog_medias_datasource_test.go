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
)

const testAccCatalogMediasDataSourceConfig = `
data "cloudavenue_catalogs" "test" {}

data "cloudavenue_catalog_medias" "example" {
	catalog_name = data.cloudavenue_catalogs.test.catalogs["catalog-example"].name
}
`

func TestAccCatalogMediasDatasource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalog_medias.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogMediasDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "catalog_name", "catalog-example"),
					resource.TestCheckResourceAttr(dataSourceName, "medias_name.0", "debian-9.9.0-amd64-netinst.iso"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.is_iso", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.is_published", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.storage_profile", "gold"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.status", "RESOLVED"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.created_at", "2023-02-28T15:10:57.136Z"),
				),
			},
		},
	})
}
