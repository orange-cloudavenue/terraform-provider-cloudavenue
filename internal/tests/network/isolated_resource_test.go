package network

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccNetworkIsolatedResourceConfig = `
resource "cloudavenue_network_isolated" "example" {
	vdc 	= "VDC_Test"
	name        = "rsx-example-isolated-network"
	description = "My isolated Org VDC network"
  
	gateway       = "1.1.1.1"
	prefix_length = 24
  
	dns1 = "8.8.8.8"
	dns2 = "8.8.4.4"
	dns_suffix = "example.com"
  
	static_ip_pool = [
	  {
		start_address = "1.1.1.10"
		end_address   = "1.1.1.20"
	  },
	  {
		start_address = "1.1.1.100"
		end_address   = "1.1.1.103"
	  }
	]
}
`

const updateAccNetworkIsolatedResourceConfig = `
resource "cloudavenue_network_isolated" "example" {
	vdc 	= "VDC_Test"
	name        = "rsx-example-isolated-network"
	description = "Example"
  
	gateway       = "1.1.1.1"
	prefix_length = 24
  
	dns1 = "1.1.1.2"
	dns2 = "8.8.8.9"
	dns_suffix = "example.com"
  
	static_ip_pool = [
	  {
		start_address = "1.1.1.10"
		end_address   = "1.1.1.20"
	  },
	  {
		start_address = "1.1.1.100"
		end_address   = "1.1.1.130"
	  }
	]
}
`

func TestAccNetworkIsolatedResource(t *testing.T) {
	const resourceName = "cloudavenue_network_isolated.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccNetworkIsolatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`(urn:vcloud:network:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttr(resourceName, "vdc", "VDC_Test"),
					resource.TestCheckResourceAttr(resourceName, "name", "rsx-example-isolated-network"),
					resource.TestCheckResourceAttr(resourceName, "description", "My isolated Org VDC network"),
					resource.TestCheckResourceAttr(resourceName, "gateway", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceName, "dns1", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "dns2", "8.8.4.4"),
					resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "static_ip_pool.#", "2"),
				),
			},
			// Update testing
			{
				// Apply test
				Config: updateAccNetworkIsolatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`(urn:vcloud:network:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttr(resourceName, "vdc", "VDC_Test"),
					resource.TestCheckResourceAttr(resourceName, "description", "Example"),
					resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.2"),
					resource.TestCheckResourceAttr(resourceName, "dns2", "8.8.8.9"),
					resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.end_address", "1.1.1.130"),
				),
			},
			// ImportruetState testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "VDC_Test.rsx-example-isolated-network",
			},
		},
	})
}
