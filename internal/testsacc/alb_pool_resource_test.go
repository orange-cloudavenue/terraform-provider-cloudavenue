package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccAlbPoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.LoadBalancerPool)),
					resource.TestCheckResourceAttr(resourceName, "name", "Example"),
					resource.TestCheckResourceAttr(resourceName, "persistence_profile.type", "CLIENT_IP"),
				),
			},
			{
				// Update test
				Config: testAccAlbPoolResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.LoadBalancerPool)),
					resource.TestCheckResourceAttr(resourceName, "name", "Example"),
					resource.TestCheckNoResourceAttr(resourceName, "persistence_profile"),
					resource.TestCheckNoResourceAttr(resourceName, "members"),
					resource.TestCheckNoResourceAttr(resourceName, "health_monitors"),
				),
			},
		},
	})
}
