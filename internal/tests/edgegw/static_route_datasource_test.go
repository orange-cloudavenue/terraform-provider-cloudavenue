package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccStaticRouteDataSourceConfig = `
data "cloudavenue_edgegateway_static_route" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	name = cloudavenue_edgegateway_static_route.example.name
}
`

func TestAccStaticRouteDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway_static_route.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccStaticRouteResourceConfig, testAccStaticRouteDataSourceConfig),
				Check:  staticRouteTestCheck(dataSourceName),
			},
		},
	})
}
