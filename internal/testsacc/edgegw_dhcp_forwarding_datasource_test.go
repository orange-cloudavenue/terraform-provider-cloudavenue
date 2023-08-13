package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDhcpForwardingDataSourceConfig = `
data "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
}
`

func TestAccDhcpForwardingDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway_dhcp_forwarding.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: ConcatTests(testAccEdgeGatewayResourceConfig, testAccDhcpForwardingResourceConfig, testAccDhcpForwardingDataSourceConfig),
				Check:  dhcpForwardingTestCheck(dataSourceName),
			},
		},
	})
}
