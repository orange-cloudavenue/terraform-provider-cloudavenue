package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDhcpBindingDataSourceConfig = `
data "cloudavenue_network_dhcp_binding" "example" {
	org_network_id = cloudavenue_network_dhcp.example.id
	name = cloudavenue_network_dhcp_binding.example.name
}
`

func TestAccDhcpBindingDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_network_dhcp_binding.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: ConcatTests(testAccDhcpResourceConfig, testAccDhcpBindingResourceConfig, testAccDhcpBindingDataSourceConfig),
				Check:  dhcpBindingTestCheck(dataSourceName),
			},
		},
	})
}
