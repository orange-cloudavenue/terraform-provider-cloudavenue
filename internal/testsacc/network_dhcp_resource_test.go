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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDhcpResourceConfig = `
resource "cloudavenue_network_dhcp" "example" {
	org_network_id = cloudavenue_network_routed.example.id
	mode           = "EDGE"
	pools = [
	  {
		start_address = "192.168.1.30"
		end_address   = "192.168.1.100"
	  }
	]
	dns_servers = [
	  "1.1.1.1",
	  "1.0.0.1"
	]
}

data "cloudavenue_edgegateways" "example" {}

resource "cloudavenue_network_routed" "example" {
	name        = "MyOrgNet"
	description = "This is an example Net"
  
	edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
  
	gateway       = "192.168.1.254"
	prefix_length = 24
  
	dns1 = "1.1.1.1"
	dns2 = "8.8.8.8"
  
	dns_suffix = "example"
  
	static_ip_pool = [
	  {
		start_address = "192.168.1.10"
		end_address   = "192.168.1.20"
	  }
	]
}
`

const testAccDhcpResourceConfigUpdate = `
resource "cloudavenue_network_dhcp" "example" {
	org_network_id = cloudavenue_network_routed.example.id
	mode           = "EDGE"
	pools = [
	  {
		start_address = "192.168.1.30"
		end_address   = "192.168.1.100"
	  },
	  {
		start_address = "192.168.1.200"
		end_address   = "192.168.1.230"
	  }
	]
	dns_servers = [
	  "8.8.8.8",
	  "9.9.9.9"
	]
}

data "cloudavenue_edgegateways" "example" {}

resource "cloudavenue_network_routed" "example" {
	name        = "MyOrgNet"
	description = "This is an example Net"
  
	edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
  
	gateway       = "192.168.1.254"
	prefix_length = 24
  
	dns1 = "1.1.1.1"
	dns2 = "8.8.8.8"
  
	dns_suffix = "example"
  
	static_ip_pool = [
	  {
		start_address = "192.168.1.10"
		end_address   = "192.168.1.20"
	  }
	]
  }
`

func dhcpTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "org_network_id"),
		resource.TestCheckResourceAttr(resourceName, "mode", "EDGE"),
		resource.TestCheckResourceAttr(resourceName, "lease_time", "86400"),
		resource.TestCheckResourceAttr(resourceName, "pools.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "pools.0.start_address", "192.168.1.30"),
		resource.TestCheckResourceAttr(resourceName, "pools.0.end_address", "192.168.1.100"),
		resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "1.1.1.1"),
		resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "1.0.0.1"),
	)
}

func TestAccDhcpResource(t *testing.T) {
	resourceName := "cloudavenue_network_dhcp.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccDhcpResourceConfig,
				Check:  dhcpTestCheck(resourceName),
			},
			// Uncomment if you want to test update or delete this block
			{
				// Update test
				Config: testAccDhcpResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "org_network_id"),
					resource.TestCheckResourceAttr(resourceName, "mode", "EDGE"),
					resource.TestCheckResourceAttr(resourceName, "lease_time", "86400"),

					resource.TestCheckResourceAttr(resourceName, "pools.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "pools.0.start_address", "192.168.1.30"),
					resource.TestCheckResourceAttr(resourceName, "pools.0.end_address", "192.168.1.100"),
					resource.TestCheckResourceAttr(resourceName, "pools.1.start_address", "192.168.1.200"),
					resource.TestCheckResourceAttr(resourceName, "pools.1.end_address", "192.168.1.230"),

					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "9.9.9.9"),
				),
			},
			// Import State testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
