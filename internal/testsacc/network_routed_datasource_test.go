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

const testAccNetworkRoutedDataSourceConfig = `
data "cloudavenue_edgegateway" "example" {
	name = "tn01e02ocb0006205spt101"
}

resource "cloudavenue_network_routed" "example" {
	name = "ExampleNetworkRouted"
	gateway       = "192.168.10.254"
	prefix_length = 24
	edge_gateway_id = data.cloudavenue_edgegateway.example.id
	dns1 = "1.1.1.1"
	dns2 = "8.8.8.8"

	dns_suffix = "example"

	static_ip_pool = [
	  {
		start_address = "192.168.10.10"
		end_address   = "192.168.10.20"
	  }
	]
}

data "cloudavenue_network_routed" "example" {
	name = cloudavenue_network_routed.example.name
  	edge_gateway_id = cloudavenue_network_routed.example.edge_gateway_id
}
`

func TestAccNetworkRoutedDataSource(t *testing.T) {
	const dataSourceName = "data.cloudavenue_network_routed.example"
	const resourceName = "cloudavenue_network_routed.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccNetworkRoutedDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Network)),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "static_ip_pool.#", resourceName, "static_ip_pool.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "gateway", resourceName, "gateway"),
				),
			},
		},
	})
}
