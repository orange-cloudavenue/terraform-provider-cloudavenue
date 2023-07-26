package edgegw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccNATRuleResourceConfigSnat = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[1].id
  
	name        = "example-snat"
	rule_type   = "SNAT"
	description = "description SNAT example"
  
	# Using primary_ip from edge gateway
	external_address         = data.cloudavenue_publicips.example.public_ips[2].public_ip
	internal_address         = "11.11.11.0/24"
	snat_destination_address = "8.8.8.8"
	
	priority = 10
}
`

const testAccNATRuleResourceConfigDnat = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[1].id
  
	name        = "example-dnat"
	rule_type   = "DNAT"
	description = "description DNAT example"
  
	# Using primary_ip from edge gateway
	external_address         = data.cloudavenue_publicips.example.public_ips[2].public_ip
	internal_address         = "11.11.11.4"
  
	dnat_external_port = "8080"
}
`

const testAccNATRuleResourceConfigReflexive = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[1].id
  
	name        = "example-reflexive"
	rule_type   = "REFLEXIVE"
	description = "description REFLEXIVE example"
  
	# Using primary_ip from edge gateway
	external_address         = data.cloudavenue_publicips.example.public_ips[2].public_ip
	internal_address         = "192.168.0.1"
}
`

const testAccNATRuleResourceConfigDataSource = `
data "cloudavenue_edgegateways" "example" {}

data "cloudavenue_publicips" "example" {}
`

const testAccNATRuleResourceConfigUpdateSnat = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[1].id
  
	name        = "example-snat"
	rule_type   = "SNAT"
	description = "description SNAT example Updated!!"
  
	# Using primary_ip from edge gateway
	external_address         = data.cloudavenue_publicips.example.public_ips[2].public_ip
	internal_address         = "11.11.11.0/24"
	snat_destination_address = "9.9.9.9"
	
	priority = 0
}
`

func natRuleSnatTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
		resource.TestCheckResourceAttrSet(resourceName, "external_address"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-snat"),
		resource.TestCheckResourceAttr(resourceName, "description", "description SNAT example"),
		resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.0/24"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "SNAT"),
		resource.TestCheckResourceAttr(resourceName, "snat_destination_address", "8.8.8.8"),
		resource.TestCheckResourceAttr(resourceName, "priority", "10"),
	)
}

func natRuleSnatUpdateTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
		resource.TestCheckResourceAttrSet(resourceName, "external_address"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-snat"),
		resource.TestCheckResourceAttr(resourceName, "description", "description SNAT example Updated!!"),
		resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.0/24"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "SNAT"),
		resource.TestCheckResourceAttr(resourceName, "snat_destination_address", "9.9.9.9"),
		resource.TestCheckResourceAttr(resourceName, "priority", "0"),
	)
}

func natRuleDnatTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
		resource.TestCheckResourceAttrSet(resourceName, "external_address"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-dnat"),
		resource.TestCheckResourceAttr(resourceName, "description", "description DNAT example"),
		resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.4"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
		resource.TestCheckResourceAttr(resourceName, "dnat_external_port", "8080"),
		resource.TestCheckResourceAttr(resourceName, "priority", "0"),
	)
}

func natRuleReflexiveTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
		resource.TestCheckResourceAttrSet(resourceName, "external_address"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-reflexive"),
		resource.TestCheckResourceAttr(resourceName, "description", "description REFLEXIVE example"),
		resource.TestCheckResourceAttr(resourceName, "internal_address", "192.168.0.1"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "REFLEXIVE"),
		resource.TestCheckResourceAttr(resourceName, "priority", "0"),
	)
}

func TestAccNATRuleResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_nat_rule.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test Snat
				Config: tests.ConcatTests(testAccNATRuleResourceConfigSnat, testAccNATRuleResourceConfigDataSource),
				Check:  natRuleSnatTestCheck(resourceName),
			},
			// Update testing Snat
			{
				// Update test
				Config: tests.ConcatTests(testAccNATRuleResourceConfigUpdateSnat, testAccNATRuleResourceConfigDataSource),
				Check:  natRuleSnatUpdateTestCheck(resourceName),
			},
			{
				// Apply test Dnat
				Config: tests.ConcatTests(testAccNATRuleResourceConfigDnat, testAccNATRuleResourceConfigDataSource),
				Check:  natRuleDnatTestCheck(resourceName),
			},
			{
				// Delete test
				Destroy: true,
				Config:  tests.ConcatTests(testAccNATRuleResourceConfigSnat, testAccNATRuleResourceConfigDataSource),
			},
			{
				// Apply test Reflexive
				Config: tests.ConcatTests(testAccNATRuleResourceConfigReflexive, testAccNATRuleResourceConfigDataSource),
				Check:  natRuleReflexiveTestCheck(resourceName),
			},
			// Import State testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNATRuleResourceImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccNATRuleResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		// edgeGatewayIDOrName.ipSetName
		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_id"], rs.Primary.Attributes["name"]), nil
	}
}
