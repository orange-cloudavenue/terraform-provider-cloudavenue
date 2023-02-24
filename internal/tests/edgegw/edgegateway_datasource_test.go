// Package edgegw provides the acceptance tests for the provider.
package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccEdgeGatewayDataSourceConfig = `
data "cloudavenue_edgegateway" "test" {
	edge_id = "frangipane"
	}
`

func TestAccEdgeGatewayDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
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
