package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccIPSetDataSourceConfig = `
data "cloudavenue_edgegateway_ip_set" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	name = cloudavenue_edgegateway_ip_set.example.name
}
`

func TestAccIPSetDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway_ip_set.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccIPSetResourceConfig, testAccIPSetDataSourceConfig),
				Check:  ipSetTestCheck(dataSourceName),
			},
		},
	})
}
