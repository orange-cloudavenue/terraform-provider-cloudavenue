// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package testsacc

import (
	"fmt"
	"os"
	"testing"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// testAccOrgNetworkDataSourceConfig is formatted with the name of an existing VDC (CLOUDAVENUE_VDC).
const testAccOrgNetworkDataSourceConfig = `
data "cloudavenue_tier0_vrfs" "example" {}

resource "cloudavenue_edgegateway" "example" {
  owner_name     = "%[1]s"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
}

resource "cloudavenue_network_routed" "example" {
  name        = "MyOrgNet"
  description = "This is an example Net"

  edge_gateway_id = cloudavenue_edgegateway.example.id

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

resource "cloudavenue_vapp" "example" {
  name        = "MyVapp"
  vdc         = "%[1]s"
  description = "This is an example vApp"
}

resource "cloudavenue_vapp_org_network" "example" {
  vapp_name    = cloudavenue_vapp.example.name
  vdc          = cloudavenue_vapp.example.vdc
  network_name = cloudavenue_network_routed.example.name
}

data "cloudavenue_vapp_org_network" "example" {
	vapp_name    = cloudavenue_vapp.example.name
	network_name = cloudavenue_network_routed.example.name
	vdc          = cloudavenue_vapp.example.vdc

	// The network must be attached to the vApp before it is read.
	depends_on = [cloudavenue_vapp_org_network.example]
}
`

func TestAccOrgNetworkDataSource(t *testing.T) {
	const dataSourceName = "data.cloudavenue_vapp_org_network.example"
	const resourceName = "cloudavenue_vapp_org_network.example"

	vdc := os.Getenv("CLOUDAVENUE_VDC")
	if vdc == "" {
		t.Fatal("CLOUDAVENUE_VDC must be set for this acceptance test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: fmt.Sprintf(testAccOrgNetworkDataSourceConfig, vdc),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testIsURNType(urn.Network)),
					resource.TestCheckResourceAttrPair(dataSourceName, "network_name", resourceName, "network_name"),
				),
			},
		},
	})
}
