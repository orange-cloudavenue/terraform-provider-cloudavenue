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

const testAccNetworkIsolatedDataSourceConfig = `
data "cloudavenue_network_isolated" "example" {
	  name = "net-isolated"
	  vdc = "VDC_Test"
}
`

func TestAccNetworkIsolatedDataSource(t *testing.T) {
	const dataSourceName = "data.cloudavenue_network_isolated.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccNetworkIsolatedDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
				),
			},
		},
	})
}
