package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccForwardingResourceConfig = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10"
	]
}
`

const testAccForwardingResourceConfigUpdate = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10",
		"192.168.10.11"
	]
}
`

func forwardingTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "1"),
	)
}

func TestAccForwardingResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_dhcp_forwarding.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccForwardingResourceConfig),
				Check:  forwardingTestCheck(resourceName),
			},
			// Update testing
			{
				// Update test
				Config: tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccForwardingResourceConfigUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "2"),
				),
			},
			// Import State testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
