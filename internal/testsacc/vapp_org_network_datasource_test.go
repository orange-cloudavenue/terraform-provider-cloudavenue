package testsacc

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccOrgNetworkDataSourceConfig = `
data "cloudavenue_tier0_vrfs" "example" {}

resource "cloudavenue_edgegateway" "example" {
  owner_name     = "MyVDC"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc"
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
  vdc         = "MyVDC"
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
}
`

func TestAccOrgNetworkDataSource(t *testing.T) {
	const dataSourceName = "data.cloudavenue_vapp_org_network.example"
	const resourceName = "cloudavenue_vapp_org_network.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccOrgNetworkDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.Network.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttrPair(dataSourceName, "network_name", resourceName, "network_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_fenced", resourceName, "is_fenced"),
					resource.TestCheckResourceAttrPair(dataSourceName, "retain_ip_mac_enabled", resourceName, "retain_ip_mac_enabled"),
				),
			},
		},
	})
}
