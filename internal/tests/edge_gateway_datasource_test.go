package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccEdgeGatewayDataSourceConfig = `
data "cloudavenue_edge_gateway" "test" {
	edge_id = "frangipane"
	}
`

func TestAccEdgeGatewayDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edge_gateway.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEdgeGatewayDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr(dataSourceName, "id", "frangipane"),
				),
			},
		},
	})
}
