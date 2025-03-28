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

const testAccNetworkIsolatedResourceConfig = `
resource "cloudavenue_network_isolated" "example" {
	vdc 	= "VDC_Test"
	name        = "rsx-example-isolated-network"
	description = "My isolated Org VDC network"
  
	gateway       = "1.1.1.1"
	prefix_length = 24
  
	dns1 = "8.8.8.8"
	dns2 = "8.8.4.4"
	dns_suffix = "example.com"
  
	static_ip_pool = [
	  {
		start_address = "1.1.1.10"
		end_address   = "1.1.1.20"
	  },
	  {
		start_address = "1.1.1.100"
		end_address   = "1.1.1.103"
	  }
	]
}
`

const updateAccNetworkIsolatedResourceConfig = `
resource "cloudavenue_network_isolated" "example" {
	vdc 	= "VDC_Test"
	name        = "rsx-example-isolated-network"
	description = "Example"
  
	gateway       = "1.1.1.1"
	prefix_length = 24
  
	dns1 = "1.1.1.2"
	dns2 = "8.8.8.9"
	dns_suffix = "example.com"
  
	static_ip_pool = [
	  {
		start_address = "1.1.1.10"
		end_address   = "1.1.1.20"
	  },
	  {
		start_address = "1.1.1.100"
		end_address   = "1.1.1.130"
	  }
	]
}
`

func TestAccNetworkIsolatedResource(t *testing.T) {
	const resourceName = "cloudavenue_network_isolated.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccNetworkIsolatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Network)),
					resource.TestCheckResourceAttr(resourceName, "vdc", "VDC_Test"),
					resource.TestCheckResourceAttr(resourceName, "name", "rsx-example-isolated-network"),
					resource.TestCheckResourceAttr(resourceName, "description", "My isolated Org VDC network"),
					resource.TestCheckResourceAttr(resourceName, "gateway", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceName, "dns1", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "dns2", "8.8.4.4"),
					resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "static_ip_pool.#", "2"),
				),
			},
			// Update testing
			{
				// Apply test
				Config: updateAccNetworkIsolatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Network)),
					resource.TestCheckResourceAttr(resourceName, "vdc", "VDC_Test"),
					resource.TestCheckResourceAttr(resourceName, "description", "Example"),
					resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.2"),
					resource.TestCheckResourceAttr(resourceName, "dns2", "8.8.8.9"),
					resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.end_address", "1.1.1.130"),
				),
			},
			// ImportruetState testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "VDC_Test.rsx-example-isolated-network",
			},
		},
	})
}
