package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPublicIPDataSourceConfig = `
data "cloudavenue_public_ip" "test" {}
`

func TestAccPublicIPDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_public_ip.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccPublicIPDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr(dataSourceName, "id", "cf9a73f3-4eb5-546b-bf2f-bdc63a439128"),
					resource.TestCheckResourceAttr(dataSourceName, "internal_ip", "10.10.10.0/30"),
					resource.TestCheckResourceAttr(dataSourceName, "network_config.0.edge_gateway_name", "edgeGatewayName"),
					resource.TestCheckResourceAttr(dataSourceName, "network_config.0.uplink_ip", "196.26.50.90"),
					resource.TestCheckResourceAttr(dataSourceName, "network_config.0.translated_ip", "10.10.10.0/30"),
					resource.TestCheckResourceAttr(dataSourceName, "network_config.1.edge_gateway_name", "edgeGatewayName2"),
					resource.TestCheckResourceAttr(dataSourceName, "network_config.1.uplink_ip", "196.26.50.91"),
					resource.TestCheckResourceAttr(dataSourceName, "network_config.1.translated_ip", "10.10.10.1/30"),
				),
			},
		},
	})
}
