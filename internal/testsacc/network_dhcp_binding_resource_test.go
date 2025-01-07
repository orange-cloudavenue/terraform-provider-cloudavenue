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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccDhcpBindingResourceConfig = `
resource "cloudavenue_network_dhcp_binding" "example" {
	name           = "example"
	org_network_id = cloudavenue_network_dhcp.example.id
	mac_address    = "00:50:56:01:01:01"
	ip_address     = "192.168.1.231"
  }
`

const testAccDhcpBindingResourceConfigUpdate = `
resource "cloudavenue_network_dhcp_binding" "example" {
	name           = "example2"
	org_network_id = cloudavenue_network_dhcp.example.id
	mac_address    = "00:50:56:01:01:01"
	ip_address     = "192.168.1.232"
	lease_time     = 86000
}
`

func dhcpBindingTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "org_network_id"),
		resource.TestCheckResourceAttr(resourceName, "name", "example"),
		resource.TestCheckResourceAttr(resourceName, "mac_address", "00:50:56:01:01:01"),
		resource.TestCheckResourceAttr(resourceName, "ip_address", "192.168.1.231"),
		resource.TestCheckResourceAttr(resourceName, "lease_time", "86400"),
	)
}

func TestAccDhcpBindingResource(t *testing.T) {
	resourceName := "cloudavenue_network_dhcp_binding.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: fmt.Sprintf("%s\n%s", testAccDhcpResourceConfig, testAccDhcpBindingResourceConfig),
				Check:  dhcpBindingTestCheck(resourceName),
			},
			{
				// Update test
				Config: fmt.Sprintf("%s\n%s", testAccDhcpResourceConfig, testAccDhcpBindingResourceConfigUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "org_network_id"),
					resource.TestCheckResourceAttr(resourceName, "name", "example2"),
					resource.TestCheckResourceAttr(resourceName, "mac_address", "00:50:56:01:01:01"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "192.168.1.232"),
					resource.TestCheckResourceAttr(resourceName, "lease_time", "86000"),
				),
			},
			// Import State testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDHCPBindingResourceImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccDHCPBindingResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["org_network_id"], rs.Primary.Attributes["name"]), nil
	}
}
