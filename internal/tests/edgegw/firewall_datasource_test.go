package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccFirewallDataSourceConfig = `
data "cloudavenue_edgegateway_firewall" "example" {
	  edge_gateway_id = cloudavenue_edgegateway_firewall.example.edge_gateway_id
}
`

func TestAccFirewallDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway_firewall.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccFirewallResourceConfig, testAccFirewallDataSourceConfig),
				Check:  firewallTestCheck(dataSourceName),
			},
		},
	})
}
