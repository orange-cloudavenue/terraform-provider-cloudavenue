package network

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccNetworkRoutedResourceConfig = `
data "cloudavenue_edgegateway" "example" {
  name = "tn01e02ocb0006205spt101"
}

resource "cloudavenue_network_routed" "example" {
  name        = "OrgNetExampleOnVDCGroup"
  description = "Org Net Example"

  edge_gateway_id = data.cloudavenue_edgegateway.example.id

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
`

const testAccNetworkRoutedResourceOnVDCConfig = `
	data "cloudavenue_tier0_vrfs" "example_with_vdc" {}

	resource "cloudavenue_edgegateway" "example_with_vdc" {
	  owner_name     = "MyVDC"
	  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example_with_vdc.names.0
	  owner_type     = "vdc"
	}

	resource "cloudavenue_network_routed" "example" {
	  name        = "OrgNetExampleOnVDC"
	  description = "Org Net Example"

	  edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id

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
`

func TestAccNetworkRoutedResource(t *testing.T) {
	const resourceName = "cloudavenue_network_routed.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test on VDC Group
			{
				// Apply test
				Config: testAccNetworkRoutedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.Network.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "name", "OrgNetExampleOnVDCGroup"),
					resource.TestCheckResourceAttr(resourceName, "description", "Org Net Example"),
					resource.TestMatchResourceAttr(resourceName, "edge_gateway_id", regexp.MustCompile(uuid.Gateway.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.1.254"),
					resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "dns2", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example"),
					resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.start_address", "192.168.1.10"),
					resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.end_address", "192.168.1.20"),
				),
			},
			{
				// Update test
				Config: newUpdatedConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.Network.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "description", "Example"),
					resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.2"),
					resource.TestCheckResourceAttr(resourceName, "dns2", "8.8.8.9"),
					resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.end_address", "192.168.1.30"),
				),
			},
			// Import test
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "HackathonShared.OrgNetExampleOnVDCGroup",
			},
			// Test on VDC
			{
				// Apply test
				Config: testAccNetworkRoutedResourceOnVDCConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.Network.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "name", "OrgNetExampleOnVDC"),
					resource.TestCheckResourceAttr(resourceName, "description", "Org Net Example"),
					// resource.TestCheckResourceAttr(resourceName, "vdc", "MyVDC"),
					resource.TestMatchResourceAttr(resourceName, "edge_gateway_id", regexp.MustCompile(uuid.Gateway.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.1.254"),
					resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "dns2", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example"),
					resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.start_address", "192.168.1.10"),
					resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.end_address", "192.168.1.20"),
				),
			},
			// Import test
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "MyVDC.OrgNetExampleOnVDC",
			},
		},
	})
}

func newUpdatedConfig() string {
	s := strings.Replace(testAccNetworkRoutedResourceConfig, "Org Net Example", "Example", 1)
	s = strings.Replace(s, "1.1.1.1", "1.1.1.2", 1)
	s = strings.Replace(s, "8.8.8.8", "8.8.8.9", 1)
	s = strings.Replace(s, "192.168.1.20", "192.168.1.30", 1)

	return s
}
