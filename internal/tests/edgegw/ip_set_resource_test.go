package edgegw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccIPSetResourceConfig = `
resource "cloudavenue_edgegateway_ip_set" "example" {
	name = "example-ip-set"
	description = "example of ip set"
	ip_addresses = [
		"192.168.1.1",
		"192.168.1.2",
	]
	edge_gateway_name = cloudavenue_edgegateway.example_with_vdc.name
}
`

const testAccIPSetResourceConfigUpdate = `
resource "cloudavenue_edgegateway_ip_set" "example" {
	name = "example-ip-set"
	description = "example of ip set"
	ip_addresses = [
		"192.168.1.1",
		"192.168.1.2",
		"192.168.1.3"
	]
	edge_gateway_name = cloudavenue_edgegateway.example_with_vdc.name
}
`

func ipSetTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-ip-set"),
		resource.TestCheckResourceAttr(resourceName, "description", "example of ip set"),
		resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "2"),
	)
}

func TestAccIPSetResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_ip_set.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccIPSetResourceConfig),
				Check:  ipSetTestCheck(resourceName),
			},
			// Update testing
			{
				// Update test
				Config: tests.ConcatTests(testAccEdgeGatewayResourceConfig, testAccIPSetResourceConfigUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "name", "example-ip-set"),
					resource.TestCheckResourceAttr(resourceName, "description", "example of ip set"),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "3"),
				),
			},
			// Import State testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIPSetResourceImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccIPSetResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		// edgeGatewayIDOrName.ipSetName
		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_id"], rs.Primary.Attributes["name"]), nil
	}
}
