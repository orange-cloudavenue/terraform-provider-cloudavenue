package testsacc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

const testAccStaticRouteResourceConfig = `
resource "cloudavenue_edgegateway_static_route" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	name = "example"
	network_cidr = "192.168.1.0/24"
	next_hops = [
		{
			ip_address = "192.168.1.254"
		}
	]
}
`

const testAccStaticRouteResourceConfigUpdate = `
resource "cloudavenue_edgegateway_static_route" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	name = "example"
	description = "example description"
	network_cidr = "192.168.2.0/24"
	next_hops = [
		{
			ip_address = "192.168.2.254"
		},
		{
			ip_address = "192.168.2.253"
			admin_distance = 2
		}
	]
}
`

const testAccStaticRouteResourceConfigWithVDCGroup = `
resource "cloudavenue_edgegateway_static_route" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_group.id
	name = "example"
	network_cidr = "192.168.1.0/24"
	next_hops = [
		{
			ip_address = "192.168.1.254"
		}
	]
}
`

const testAccStaticRouteResourceConfigUpdateWithVDCGroup = `
resource "cloudavenue_edgegateway_static_route" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_group.id
	name = "example"
	description = "example description"
	network_cidr = "192.168.2.0/24"
	next_hops = [
		{
			ip_address = "192.168.2.254"
		},
		{
			ip_address = "192.168.2.253"
			admin_distance = 2
		}
	]
}
`

func staticRouteTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
		resource.TestCheckResourceAttr(resourceName, "name", "example"),
		resource.TestCheckNoResourceAttr(resourceName, "description"),
		resource.TestCheckResourceAttr(resourceName, "network_cidr", "192.168.1.0/24"),
		resource.TestCheckResourceAttr(resourceName, "next_hops.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "next_hops.0.ip_address", "192.168.1.254"),
		resource.TestCheckResourceAttr(resourceName, "next_hops.0.admin_distance", "1"),
	)
}

func staticRouteTestCheckUpdated(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
		resource.TestCheckResourceAttr(resourceName, "name", "example"),
		resource.TestCheckResourceAttr(resourceName, "description", "example description"),
		resource.TestCheckResourceAttr(resourceName, "network_cidr", "192.168.2.0/24"),
		resource.TestCheckResourceAttr(resourceName, "next_hops.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "next_hops.0.ip_address", "192.168.2.254"),
		resource.TestCheckResourceAttr(resourceName, "next_hops.0.admin_distance", "1"),
		resource.TestCheckResourceAttr(resourceName, "next_hops.1.ip_address", "192.168.2.253"),
		resource.TestCheckResourceAttr(resourceName, "next_hops.1.admin_distance", "2"),
	)
}

func TestAccStaticRouteResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_static_route.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// * Test with VDC
			// Read testing
			{
				// Apply test
				Config: ConcatTests(testAccEdgeGatewayResourceConfig, testAccStaticRouteResourceConfig),
				Check:  staticRouteTestCheck(resourceName),
			},
			// Update testing
			{
				// Update test
				Config: ConcatTests(testAccEdgeGatewayResourceConfig, testAccStaticRouteResourceConfigUpdate),
				Check:  staticRouteTestCheckUpdated(resourceName),
			},
			// Import State testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStaticRouteResourceImportStateIDFuncWithID(resourceName),
			},
			{
				// Delete test
				Destroy: true,
				Config:  ConcatTests(testAccEdgeGatewayResourceConfig, testAccStaticRouteResourceConfigUpdate),
			},

			// * Test with VDCGroup
			{
				// Apply test
				Config: ConcatTests(testAccEdgeGatewayGroupResourceConfig, testAccStaticRouteResourceConfigWithVDCGroup),
				Check:  staticRouteTestCheck(resourceName),
			},
			// Update testing
			{
				// Update test
				Config: ConcatTests(testAccEdgeGatewayGroupResourceConfig, testAccStaticRouteResourceConfigUpdateWithVDCGroup),
				Check:  staticRouteTestCheckUpdated(resourceName),
			},
			// Import State testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStaticRouteResourceImportStateIDFuncWithName(resourceName),
			},
		},
	})
}

func testAccStaticRouteResourceImportStateIDFuncWithID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccStaticRouteResourceImportStateIDFuncWithName(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_name"], rs.Primary.Attributes["name"]), nil
	}
}
