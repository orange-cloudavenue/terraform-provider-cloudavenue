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

const testAccAlbPoolDataSourceConfig = `
resource "cloudavenue_alb_pool" "example" {
	edge_gateway_name = "tn01e02ocb0006205spt102"
	name              = "Example"
  
	persistence_profile = {
	  type = "CLIENT_IP"
	}
  
	members = [
	  {
	    ip_address = "192.168.1.1"
	    port       = "80"
	  },
	  {
		ip_address = "192.168.1.2"
		port       = "80"
	  },
	  {
		ip_address = "192.168.1.3"
		port       = "80"
	  }
	]
  
	health_monitors = ["UDP", "TCP"]
  }

data "cloudavenue_alb_pool" "example" {
	edge_gateway_name = cloudavenue_alb_pool.example.edge_gateway_name
	name              = cloudavenue_alb_pool.example.name
}
`

func TestAccAlbPoolDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_alb_pool.example"
	resourceName := "cloudavenue_alb_pool.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccAlbPoolDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerPool)),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "persistence_profile.#", resourceName, "persistence_profile.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "members.#", resourceName, "members.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "health_monitors.#", resourceName, "health_monitors.#"),
				),
			},
		},
	})
}
