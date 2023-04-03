package tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccAlbPoolDataSourceConfig = `
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

data "cloudavenue_alb_pool" "example" {
	edge_gateway_name = cloudavenue_alb_pool.example.edge_gateway_name
	name              = cloudavenue_alb_pool.example.name
}
`

func TestAccAlbPoolDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_alb_pool.example"
	resourceName := "cloudavenue_alb_pool.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccAlbPoolDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`urn:vcloud:loadBalancerPool:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "persistence_profile.#", resourceName, "persistence_profile.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "members.#", resourceName, "members.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "health_monitors.#", resourceName, "health_monitors.#"),
				),
			},
		},
	})
}
