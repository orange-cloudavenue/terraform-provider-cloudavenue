package tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccAlbPoolResourceConfig = `
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
`

const testAccAlbPoolResourceConfigUpdate = `
resource "cloudavenue_alb_pool" "example" {
	edge_gateway_name = "tn01e02ocb0006205spt102"
	name              = "Example"
  }
`

func TestAccAlbPoolResource(t *testing.T) {
	const resourceName = "cloudavenue_alb_pool.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccAlbPoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`(urn:vcloud:loadBalancerPool:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttr(resourceName, "name", "Example"),
				),
			},
			// Uncomment if you want to test update or delete this block
			{
				// Update test
				Config: testAccAlbPoolResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`(urn:vcloud:loadBalancerPool:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttr(resourceName, "name", "Example"),
					resource.TestCheckNoResourceAttr(resourceName, "persistence_profile"),
					resource.TestCheckNoResourceAttr(resourceName, "members"),
					resource.TestCheckNoResourceAttr(resourceName, "health_monitors"),
				),
			},
		},
	})
}
