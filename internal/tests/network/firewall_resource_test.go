package network

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccFirewallResourceConfig = `
data "cloudavenue_edgegateways" "example" {}

resource "cloudavenue_network_firewall" "example" {

  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
  rules = [
    {
      action      = "ALLOW"
      name        = "allow all IPv4 traffic"
      direction   = "IN_OUT"
      ip_protocol = "IPV4"
    }
  ]
}
`

const testAccFirewallResourceConfigUpdate = `
data "cloudavenue_edgegateways" "example" {}

resource "cloudavenue_network_firewall" "example" {

  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
  rules = [
    {
      action      = "ALLOW"
      name        = "allow IN IPv4 traffic"
      direction   = "IN"
      ip_protocol = "IPV4"
    },
	{
		action      = "ALLOW"
		name        = "allow OUT IPv4 traffic"
		direction   = "OUT"
		ip_protocol = "IPV4"
	}
  ]
}
`

func TestAccFirewallResource(t *testing.T) {
	resourceName := "cloudavenue_network_firewall.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccFirewallResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.name", "allow all IPv4 traffic"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.direction", "IN_OUT"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.ip_protocol", "IPV4"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.logging", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.sources_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.destinations_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.app_port_profile_ids.#", "0"),
				),
			},
			{
				// Update test
				Config: testAccFirewallResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					// * Rule 1
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.name", "allow IN IPv4 traffic"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.direction", "IN"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.ip_protocol", "IPV4"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.logging", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.sources_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.destinations_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.app_port_profile_ids.#", "0"),
					// * Rule 2
					resource.TestCheckResourceAttr(resourceName, "rules.1.action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.name", "allow OUT IPv4 traffic"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.direction", "OUT"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.ip_protocol", "IPV4"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.logging", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.sources_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.destinations_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.app_port_profile_ids.#", "0"),
				),
			},
			// ImportruetState testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFirewallResourceImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccFirewallResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.Attributes["edge_gateway_id"], nil
	}
}
