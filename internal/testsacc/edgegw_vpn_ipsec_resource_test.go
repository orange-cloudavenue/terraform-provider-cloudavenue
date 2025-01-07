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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

const MytestAccVDCResourceConfig = `
resource "cloudavenue_vdc" "example" {
  name                  = "MyVDC1"
  vdc_group             = "MyGroup"
  description           = "Example vDC created by Terraform"
  cpu_allocated         = 22000
  memory_allocated      = 30
  cpu_speed_in_mhz      = 2200
  billing_model         = "PAYG"
  disponibility_class   = "ONE-ROOM"
  service_class         = "STD"
  storage_billing_model = "PAYG"
  storage_profiles = [{
	class   = "gold"
    default = true
    limit   = 500
  }]
  }
`

const MytestAccEdgeGatewayResourceConfig = `
data "cloudavenue_tier0_vrfs" "example" {}

resource "cloudavenue_edgegateway" "example" {
  depends_on = [ cloudavenue_vdc.example ]
  owner_name     = cloudavenue_vdc.example.name
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc"
  lb_enabled     = false
}
`

const MytestAccEdgeGatewayGroupResourceConfig = `
data "cloudavenue_tier0_vrfs" "example" {}

resource "cloudavenue_edgegateway" "example" {
  owner_name     = "MyVDCGroup"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc-group"
}
`

const testAccVPNIPSecResourceConfigDefault = `
resource "cloudavenue_edgegateway_vpn_ipsec" "example" {
  edge_gateway_id = cloudavenue_edgegateway.example.id

  name        = "example-default"
  description = "example VPN IPSec"
  enabled   = false

  pre_shared_key = "my-preshared-key"

  # Primary IP address of Edge Gateway pulled from data source
  # local_ip_address = cloudavenue_publicip.example.public_ip
  local_ip_address = "100.103.223.136"
  local_networks   = ["10.10.10.0/24", "30.30.30.0/28", "40.40.40.1/32"]

  # That is a fake remote IP address
  remote_ip_address = "1.2.3.4"
  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24", "192.168.20.0/28"]
}
`

const testAccVPNIPSecResourceConfigCustomize = `
resource "cloudavenue_edgegateway_vpn_ipsec" "example" {
  depends_on = [ cloudavenue_edgegateway.example ]
  edge_gateway_id = cloudavenue_edgegateway.example.id

  name        = "example-customize"
  description = "example VPN IPSec"
  enabled   = false

  pre_shared_key = "my-preshared-key"
  # Primary IP address of Edge Gateway pulled from data source
  local_ip_address = "100.103.223.136"
  local_networks   = ["10.10.10.0/24", "30.30.30.0/28", "40.40.40.1/32"]

  # That is a fake remote IP address
  remote_ip_address = "1.2.3.4"
  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24", "192.168.20.0/28"]

  security_profile = {
	ike_version = "IKE_V2"
	ike_dh_groups = "GROUP15"
	ike_encryption_algorithm = "AES_256"
	ike_digest_algorithm = "SHA2_384"
	ike_sa_lifetime = 86400
	tunnel_pfs = true
	tunnel_df_policy = "CLEAR"
	tunnel_encryption_algorithms = "AES_256"
	tunnel_digest_algorithms = "SHA2_256"
	tunnel_dh_groups = "GROUP2"
	tunnel_sa_lifetime = 3600
	tunnel_dpd = 45
  }
}
`

const testAccVPNIPSecResourceConfigDefaultUpdate = `
resource "cloudavenue_edgegateway_vpn_ipsec" "example" {
  edge_gateway_id = cloudavenue_edgegateway.example.id

  name        = "example-default"
  description = "example VPN IPSec Updated"
  enabled   = false

  pre_shared_key = "my-preshared-key-updated"

  # Primary IP address of Edge Gateway pulled from data source
  # local_ip_address = cloudavenue_publicip.example.public_ip
  local_ip_address = "100.103.223.136"
  local_networks   = ["10.10.10.0/24", "30.30.30.0/28", "40.40.40.1/32"]

  # That is a fake remote IP address
  remote_ip_address = "4.3.2.1"
  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24"]
}
`

const testAccVPNIPSecResourceConfigCustomizeUpdate = `
resource "cloudavenue_edgegateway_vpn_ipsec" "example" {
  edge_gateway_id = cloudavenue_edgegateway.example.id

  name        = "example-customize"
  description = "example VPN IPSec Updated"
  enabled   = false

  pre_shared_key = "my-preshared-key-updated"

  # Primary IP address of Edge Gateway pulled from data source
  local_ip_address = "100.103.223.136"
  local_networks   = ["10.10.10.0/24", "30.30.30.0/28"]

  # That is a fake remote IP address
  remote_ip_address = "4.3.2.1"
  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24"]

  security_profile = {
	ike_version = "IKE_V2"
	ike_dh_groups = "GROUP20"
	ike_encryption_algorithm = "AES_256"
	ike_digest_algorithm = "SHA2_384"
	ike_sa_lifetime = 86400
	tunnel_pfs = true
	tunnel_df_policy = "CLEAR"
	tunnel_encryption_algorithms = "AES_128"
	tunnel_digest_algorithms = "SHA2_512"
	tunnel_dh_groups = "GROUP2"
	tunnel_sa_lifetime = 7200
	tunnel_dpd = 60
  }
}
`

func vpnIPSecTestCheckDefault(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-default"),
		resource.TestCheckResourceAttr(resourceName, "description", "example VPN IPSec"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "my-preshared-key"),
		resource.TestCheckResourceAttrSet(resourceName, "local_ip_address"),
		resource.TestCheckResourceAttr(resourceName, "local_networks.#", "3"),
		resource.TestCheckResourceAttr(resourceName, "remote_ip_address", "1.2.3.4"),
		resource.TestCheckResourceAttr(resourceName, "remote_networks.#", "3"),
		resource.TestCheckResourceAttr(resourceName, "security_type", "DEFAULT"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_version"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_dh_groups"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_encryption_algorithm"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_digest_algorithm"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_sa_lifetime"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_pfs"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_df_policy"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_encryption_algorithms"),
		resource.TestCheckNoResourceAttr(resourceName, "security_profile.tunnel_digest_algorithms"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_dh_groups"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_sa_lifetime"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_dpd"),
	)
}

func vpnIPSecTestCheckDefaultUpdate(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-default"),
		resource.TestCheckResourceAttr(resourceName, "description", "example VPN IPSec Updated"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "my-preshared-key-updated"),
		resource.TestCheckResourceAttrSet(resourceName, "local_ip_address"),
		resource.TestCheckResourceAttr(resourceName, "local_networks.#", "3"),
		resource.TestCheckResourceAttr(resourceName, "remote_ip_address", "4.3.2.1"),
		resource.TestCheckResourceAttr(resourceName, "remote_networks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "security_type", "DEFAULT"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_version"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_dh_groups"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_encryption_algorithm"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_digest_algorithm"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.ike_sa_lifetime"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_pfs"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_df_policy"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_encryption_algorithms"),
		resource.TestCheckNoResourceAttr(resourceName, "security_profile.tunnel_digest_algorithms"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_dh_groups"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_sa_lifetime"),
		resource.TestCheckResourceAttrSet(resourceName, "security_profile.tunnel_dpd"),
	)
}

func vpnIPSecTestCheckCustomize(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-customize"),
		resource.TestCheckResourceAttr(resourceName, "description", "example VPN IPSec"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "my-preshared-key"),
		resource.TestCheckResourceAttrSet(resourceName, "local_ip_address"),
		resource.TestCheckResourceAttr(resourceName, "local_networks.#", "3"),
		resource.TestCheckResourceAttr(resourceName, "remote_ip_address", "1.2.3.4"),
		resource.TestCheckResourceAttr(resourceName, "remote_networks.#", "3"),
		resource.TestCheckResourceAttr(resourceName, "security_type", "CUSTOM"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_version", "IKE_V2"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_dh_groups", "GROUP15"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_encryption_algorithm", "AES_256"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_digest_algorithm", "SHA2_384"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_sa_lifetime", "86400"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_pfs", "true"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_df_policy", "CLEAR"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_encryption_algorithms", "AES_256"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_digest_algorithms", "SHA2_256"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_dh_groups", "GROUP2"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_sa_lifetime", "3600"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_dpd", "45"),
	)
}

func vpnIPSecTestCheckCustomizeUpdate(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-customize"),
		resource.TestCheckResourceAttr(resourceName, "description", "example VPN IPSec Updated"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "my-preshared-key-updated"),
		resource.TestCheckResourceAttrSet(resourceName, "local_ip_address"),
		resource.TestCheckResourceAttr(resourceName, "local_networks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "remote_ip_address", "4.3.2.1"),
		resource.TestCheckResourceAttr(resourceName, "remote_networks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "security_type", "CUSTOM"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_version", "IKE_V2"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_dh_groups", "GROUP20"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_encryption_algorithm", "AES_256"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_digest_algorithm", "SHA2_384"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.ike_sa_lifetime", "86400"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_pfs", "true"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_df_policy", "CLEAR"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_encryption_algorithms", "AES_128"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_digest_algorithms", "SHA2_512"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_dh_groups", "GROUP2"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_sa_lifetime", "7200"),
		resource.TestCheckResourceAttr(resourceName, "security_profile.tunnel_dpd", "60"),
	)
}

func TestAccVPNIPSecResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_vpn_ipsec.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// * Test with VPN IPSec default Profile
			{
				// Apply test
				Config: ConcatTests(MytestAccVDCResourceConfig, MytestAccEdgeGatewayResourceConfig, testAccVPNIPSecResourceConfigDefault),
				// Config: testAccVPNIPSecResourceConfigDefault,
				Check: vpnIPSecTestCheckDefault(resourceName),
			},
			{
				// Update test
				Config: ConcatTests(MytestAccVDCResourceConfig, MytestAccEdgeGatewayResourceConfig, testAccVPNIPSecResourceConfigDefaultUpdate),
				Check:  vpnIPSecTestCheckDefaultUpdate(resourceName),
			},
			{
				// Delete test
				Config:  ConcatTests(MytestAccVDCResourceConfig, MytestAccEdgeGatewayResourceConfig, testAccVPNIPSecResourceConfigDefaultUpdate),
				Destroy: true,
			},
			// * Test with VPN IPSec customize Profile
			{
				// Apply test
				Config: ConcatTests(MytestAccVDCResourceConfig, MytestAccEdgeGatewayResourceConfig, testAccVPNIPSecResourceConfigCustomize),
				Check:  vpnIPSecTestCheckCustomize(resourceName),
			},
			{
				// Update test
				Config: ConcatTests(MytestAccVDCResourceConfig, MytestAccEdgeGatewayResourceConfig, testAccVPNIPSecResourceConfigCustomizeUpdate),
				Check:  vpnIPSecTestCheckCustomizeUpdate(resourceName),
			},
			{
				// Import test ID and Name
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVPNIPSecResourceImportStateIDFuncWithIDAndName(resourceName),
			},
			{
				// Import test ID and Name
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVPNIPSecResourceImportStateIDFuncWithNameAndID(resourceName),
			},
			// * Test with VDCGroup
			// * Test with VPN IPSec default Profile
			{
				// Apply test
				Config: ConcatTests(MytestAccVDCResourceConfig, MytestAccEdgeGatewayGroupResourceConfig, testAccVPNIPSecResourceConfigDefault),
				Check:  vpnIPSecTestCheckDefault(resourceName),
			},
			{
				// Update test
				Config: ConcatTests(MytestAccVDCResourceConfig, MytestAccEdgeGatewayGroupResourceConfig, testAccVPNIPSecResourceConfigDefaultUpdate),
				Check:  vpnIPSecTestCheckDefaultUpdate(resourceName),
			},
			{
				// Delete test
				Config:  ConcatTests(MytestAccVDCResourceConfig, MytestAccEdgeGatewayGroupResourceConfig, testAccVPNIPSecResourceConfigDefaultUpdate),
				Destroy: true,
			},
		},
	})
}

func testAccVPNIPSecResourceImportStateIDFuncWithIDAndName(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		// edgeGatewayIDOrName.ipSetName
		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccVPNIPSecResourceImportStateIDFuncWithNameAndID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		// edgeGatewayIDOrName.ipSetName
		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_name"], rs.Primary.Attributes["id"]), nil
	}
}
