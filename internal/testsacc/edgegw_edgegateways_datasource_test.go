// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccEdgeGatewaysDataSourceConfig = `
data "cloudavenue_edgegateways" "example" {}
`

func TestAccEdgeGatewaysDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateways.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEdgeGatewaysDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrWith(dataSourceName, "id", uuid.TestIsType(uuid.Gateway)),
					resource.TestCheckResourceAttrWith(dataSourceName, "edge_gateways.0.id", uuid.TestIsType(uuid.Gateway)),
				),
			},
		},
	})
}
