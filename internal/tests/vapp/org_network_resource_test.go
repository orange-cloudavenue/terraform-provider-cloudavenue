package vapp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccOrgNetworkResourceConfig = `
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
  description = "This is an example vApp"
  vdc         = "MyVDC"
}

resource "cloudavenue_vapp_org_network" "example" {
  vapp_name    = cloudavenue_vapp.example.name
  network_name = cloudavenue_network_routed.example.name
  vdc          = "MyVDC"
}
`

func TestAccOrgNetworkResource(t *testing.T) {
	resourceName := "cloudavenue_vapp_org_network.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccOrgNetworkResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "vdc", "MyVDC"),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "MyVapp"),
					resource.TestCheckResourceAttr(resourceName, "network_name", "MyOrgNet"),
				),
			},
			// Import
			{
				// Import test with vdc
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "MyVDC.MyVapp.MyOrgNet",
			},
		},
	})
}
