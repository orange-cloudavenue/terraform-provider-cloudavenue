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
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VPNIPSecResource{}

const EdgeGatewayVPNIPSecResourceName = testsacc.ResourceName("cloudavenue_edgegateway_vpn_ipsec")

type VPNIPSecResource struct{}

func NewEdgeGatewayVPNIPSecResourceTest() testsacc.TestACC {
	return &VPNIPSecResource{}
}

func (r *VPNIPSecResource) GetResourceName() string {
	return EdgeGatewayVPNIPSecResourceName.String()
}

func (r *VPNIPSecResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	resp.Append(GetResourceConfig()[PublicIPResourceName]().GetDefaultConfig)
	return resp
}

func (r *VPNIPSecResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{ // vpnIPSecTestCheckCommon(resourceName)},
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					// add a test to check id an format string ID
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				},

				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_vpn_ipsec" "example" {
					  edge_gateway_id = cloudavenue_edgegateway.example.id
					  name        = {{ generate . "name" }}
					  description = {{ generate . "description" }}
					  enabled   = false
					  pre_shared_key = "my-preshared-key"
					  local_ip_address = cloudavenue_publicip.example.public_ip
					  local_networks   = ["10.10.10.0/24", "30.30.30.0/28", "40.40.40.1/32"]
					  remote_ip_address = "1.2.3.4"
					  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24", "192.168.20.0/28"]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "my-preshared-key"),
						resource.TestCheckResourceAttrSet(resourceName, "local_ip_address"),
						resource.TestCheckResourceAttr(resourceName, "local_networks.#", "3"),
						resource.TestCheckResourceAttr(resourceName, "remote_ip_address", "1.2.3.4"),
						resource.TestCheckResourceAttr(resourceName, "remote_id", "1.2.3.4"),
						resource.TestCheckResourceAttr(resourceName, "remote_networks.#", "3"),
						resource.TestCheckResourceAttr(resourceName, "security_type", "DEFAULT"),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						// Update the resource with new values
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
							resource "cloudavenue_edgegateway_vpn_ipsec" "example" {
							  edge_gateway_id = cloudavenue_edgegateway.example.id
							  name        = {{ get . "name" }}
							  description = {{ generate . "description" }}
							  enabled   = true
							  pre_shared_key = "my-preshared-key-updated"
							  local_ip_address = cloudavenue_publicip.example.public_ip
							  local_networks   = ["10.10.10.0/24", "30.30.30.0/28", "40.40.40.1/32"]
							  remote_id        = "my-remote-id"
							  remote_ip_address = "4.3.2.1"
							  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24"]
							}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "my-preshared-key-updated"),
							resource.TestCheckResourceAttrSet(resourceName, "local_ip_address"),
							resource.TestCheckResourceAttr(resourceName, "local_networks.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "remote_id", "my-remote-id"),
							resource.TestCheckResourceAttr(resourceName, "remote_ip_address", "4.3.2.1"),
							resource.TestCheckResourceAttr(resourceName, "remote_networks.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "security_type", "DEFAULT"),
						},
					},
				},
				Imports: []testsacc.TFImport{
					{ImportStateIDBuilder: []string{"edge_gateway_id", "name"}, ImportState: true, ImportStateVerify: true},
					{ImportStateIDBuilder: []string{"edge_gateway_name", "id"}, ImportState: true, ImportStateVerify: true},
				},
			}
		},
		"customize": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{ // vpnIPSecTestCheckCommon(resourceName)},
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				},
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_vpn_ipsec" "customize" {
						  edge_gateway_id = cloudavenue_edgegateway.example.id
						  name        = {{ generate . "name" }}
						  description = {{ generate . "description" }}
						  enabled   = false
						  pre_shared_key = "my-preshared-key"
						  local_ip_address = cloudavenue_publicip.example.public_ip
						  local_networks   = ["10.10.10.0/24", "30.30.30.0/28", "40.40.40.1/32"]
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
						}`),
					Checks: []resource.TestCheckFunc{ // vpnIPSecTestCheckExample(resourceName)},
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "my-preshared-key"),
						resource.TestCheckResourceAttrSet(resourceName, "local_ip_address"),
						resource.TestCheckResourceAttr(resourceName, "local_networks.#", "3"),
						resource.TestCheckResourceAttr(resourceName, "remote_ip_address", "1.2.3.4"),
						resource.TestCheckResourceAttr(resourceName, "remote_networks.#", "3"),

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
						resource.TestCheckResourceAttr(resourceName, "security_type", "CUSTOM"),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_vpn_ipsec" "customize" {
					  edge_gateway_id = cloudavenue_edgegateway.example.id
					  name        = {{ get . "name" }}
					  description = {{ generate . "description" }}
					  enabled   = true
					  pre_shared_key = "my-preshared-key-updated"
					  local_ip_address = cloudavenue_publicip.example.public_ip
					  local_networks   = ["10.10.10.0/24", "30.30.30.0/28"]
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
					}`),
						Checks: []resource.TestCheckFunc{ // vpnIPSecTestCheckCustomizeUpdate(resourceName)}},
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "my-preshared-key-updated"),
							resource.TestCheckResourceAttrSet(resourceName, "local_ip_address"),
							resource.TestCheckResourceAttr(resourceName, "local_networks.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "remote_ip_address", "4.3.2.1"),
							resource.TestCheckResourceAttr(resourceName, "remote_networks.#", "2"),
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
							resource.TestCheckResourceAttr(resourceName, "security_type", "CUSTOM"),
						},
					},
				},
			}
		},
	}
}

func TestAccVPNIPSecResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VPNIPSecResource{}),
	})
}
