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

const testAccVPNIPSecDataSourceConfig = `
data "cloudavenue_edgegateway_vpn_ipsec" "example" {
  depends_on = [ cloudavenue_edgegateway_vpn_ipsec.example ]
  edge_gateway_id = cloudavenue_edgegateway_vpn_ipsec.example.edge_gateway_id
  name = cloudavenue_edgegateway_vpn_ipsec.example.name
}
`

func TestAccVPNIPSecDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway_vpn_ipsec.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: ConcatTests(testAccVPNIPSecDataSourceConfig, testAccVPNIPSecResourceConfigCustomize, MytestAccEdgeGatewayGroupResourceConfig, MytestAccVDCResourceConfig),
				Check:  vpnIPSecTestCheckCustomize(dataSourceName),
			},
		},
	})
}
