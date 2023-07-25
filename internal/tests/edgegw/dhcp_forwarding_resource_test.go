package edgegw

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccDhcpForwardingResourceConfig = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10"
	]
}
`

const testAccDhcpForwardingResourceConfigUpdate = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10",
		"192.168.10.11"
	]
}
`

const testAccDhcpForwardingResourceConfigUpdateError = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	enabled = false
	dhcp_servers = [
		"192.168.10.10"
	]
}
`

const testAccDhcpForwardingResourceConfigWithVDCGroup = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_group.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10",
	]
}
`

const testAccDhcpForwardingResourceConfigWithVDCGroupUpdate = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_group.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10",
		"192.168.10.11"
	]
}
`

func dhcpForwardingTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "1"),
	)
}

func dhcpForwardingTestCheckWithVDCGroup(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "2"),
	)
}

func TestAccDhcpForwardingResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_dhcp_forwarding.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// * Test with VDC
			{
				// Apply
				Config: tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccDhcpForwardingResourceConfig),
				Check:  dhcpForwardingTestCheck(resourceName),
			},
			{
				// Update
				Config: tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccDhcpForwardingResourceConfigUpdate),
				Check:  dhcpForwardingTestCheckWithVDCGroup(resourceName),
			},
			{
				// Import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// check error when updating dhcp_servers if enabled is false
				Config:             tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccDhcpForwardingResourceConfigUpdateError),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile("DHCP servers cannot be edited"),
			},
			// Destroy test with VDC
			{
				Destroy: true,
				Config:  tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccDhcpForwardingResourceConfigUpdate),
			},

			// * Test with VDC group
			{
				// Apply
				Config: tests.ConcatTests(testAccEdgeGatewayGroupResourceConfig, testAccDhcpForwardingResourceConfigWithVDCGroup),
				Check:  dhcpForwardingTestCheck(resourceName),
			},
			{
				// Update
				Config: tests.ConcatTests(testAccEdgeGatewayGroupResourceConfig, testAccDhcpForwardingResourceConfigWithVDCGroupUpdate),
				Check:  dhcpForwardingTestCheckWithVDCGroup(resourceName),
			},
			{
				// Import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
