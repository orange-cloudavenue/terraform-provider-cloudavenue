package network

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccDhcpDataSourceConfig = `
data "cloudavenue_network_dhcp" "example" {
	  org_network_id = cloudavenue_network_dhcp.example.id
}
`

func TestAccDhcpDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_network_dhcp.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccDhcpResourceConfig, testAccDhcpDataSourceConfig),
				Check:  dhcpTestCheck(dataSourceName),
			},
		},
	})
}
