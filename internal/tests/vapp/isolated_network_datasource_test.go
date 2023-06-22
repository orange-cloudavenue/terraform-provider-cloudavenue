package vapp

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccIsolatedNetworkDataSourceConfig = `
resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example vApp"
  }
  
  resource "cloudavenue_vapp_isolated_network" "example" {
	name                  = "MyVappNet"
	vapp_name             = cloudavenue_vapp.example.name
	gateway               = "192.168.10.1"
	netmask               = "255.255.255.0"
	dns1                  = "192.168.10.1"
	dns2                  = "192.168.10.3"
	dns_suffix            = "myvapp.biz"
	guest_vlan_allowed    = true
	retain_ip_mac_enabled = true
  
	static_ip_pool = [{
	  start_address = "192.168.10.51"
	  end_address   = "192.168.10.101"
	  },
	  {
		start_address = "192.168.10.10"
		end_address   = "192.168.10.30"
	}]
  }
  
  data "cloudavenue_vapp_isolated_network" "example" {
	  vapp_name = cloudavenue_vapp.example.name
	  name      = cloudavenue_vapp_isolated_network.example.name
  }
`

func TestAccIsolatedNetworkDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vapp_isolated_network.example"
	resourceName := "cloudavenue_vapp_isolated_network.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccIsolatedNetworkDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.Network.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttrPair(dataSourceName, "vapp_name", resourceName, "vapp_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "guest_vlan_allowed", resourceName, "guest_vlan_allowed"),
					resource.TestCheckResourceAttrPair(dataSourceName, "retain_ip_mac_enabled", resourceName, "retain_ip_mac_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "gateway", resourceName, "gateway"),
					resource.TestCheckResourceAttrPair(dataSourceName, "netmask", resourceName, "netmask"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dns1", resourceName, "dns1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dns2", resourceName, "dns2"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dns_suffix", resourceName, "dns_suffix"),
					resource.TestCheckResourceAttrPair(dataSourceName, "static_ip_pool.0.start_address", resourceName, "static_ip_pool.0.start_address"),
					resource.TestCheckResourceAttrPair(dataSourceName, "static_ip_pool.0.end_address", resourceName, "static_ip_pool.0.end_address"),
				),
			},
		},
	})
}
