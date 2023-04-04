package network

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccNetworkRoutedDataSourceConfig = `
resource "cloudavenue_network_routed" "example" {
	name = "ExampleNetworkRouted"
	gateway       = "192.168.10.254"
	prefix_length = 24
	edge_gateway_id = "urn:vcloud:gateway:dde5d31a-2f32-43ef-b3b3-127245958298"
}

data "cloudavenue_network_routed" "example" {
	name = "ExampleNetworkRouted"
  	edge_gateway_id = "urn:vcloud:gateway:dde5d31a-2f32-43ef-b3b3-127245958298"
}
`

func TestAccNetworkRoutedDataSource(t *testing.T) {
	const dataSourceName = "data.cloudavenue_network_routed.example"
	const resourceName = "cloudavenue_network_routed.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccNetworkRoutedDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "interface_type", resourceName, "interface_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "gateway", resourceName, "gateway"),
				),
			},
		},
	})
}
