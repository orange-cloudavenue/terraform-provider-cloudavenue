package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccFirewallDataSourceConfig = `
data "cloudavenue_edgegateway_firewall" "example" {
	  edge_gateway_id = cloudavenue_edgegateway_firewall.example.edge_gateway_id
}
`

func TestAccFirewallDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway_firewall.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: ConcatTests(testAccFirewallResourceConfig, testAccFirewallDataSourceConfig),
				Check:  firewallTestCheck(dataSourceName),
			},
		},
	})
}
