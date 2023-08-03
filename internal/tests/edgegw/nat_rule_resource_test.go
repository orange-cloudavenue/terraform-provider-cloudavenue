package edgegw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccNATRuleResourceConfigSnat = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
  
	name        = "example-snat"
	rule_type   = "SNAT"
	description = "description SNAT example"
  
	# Using primary_ip from edge gateway
	external_address         = "89.32.25.10"
	internal_address         = "11.11.11.0/24"
	snat_destination_address = "8.8.8.8"
	
	priority = 10
}
`

const testAccNATRuleResourceConfigDnat = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
  
	name        = "example-dnat"
	rule_type   = "DNAT"
	description = "description DNAT example"
  
	# Using primary_ip from edge gateway
	external_address         = "89.32.25.10"
	internal_address         = "11.11.11.4"
  
	dnat_external_port = "8080"
}
`

const testAccNATRuleResourceConfigReflexive = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
  
	name        = "example-reflexive"
	rule_type   = "REFLEXIVE"
	description = "description REFLEXIVE example"
  
	# Using primary_ip from edge gateway
	external_address         = "89.32.25.10"
	internal_address         = "192.168.0.1"
}
`

const testAccNATRuleResourceConfigUpdateSnat = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
  
	name        = "example-snat"
	rule_type   = "SNAT"
	description = "description SNAT example Updated!!"
  
	# Using primary_ip from edge gateway
	external_address         = "89.32.25.10"
	internal_address         = "11.11.11.0/24"
	snat_destination_address = "9.9.9.9"
	
	priority = 0
}
`

const testAccNATRuleResourceConfigUpdateDnat = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
  
	name        = "example-dnat"
	rule_type   = "DNAT"
	description = "description DNAT example Updated!!"
  
	# Using primary_ip from edge gateway
	external_address         = "89.32.25.10"
	internal_address         = "4.11.11.11"

	priority = 25
}
`

const testAccNATRuleResourceConfigDnatWithVDCGroup = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_group.id
  
	name        = "example-dnat"
	rule_type   = "DNAT"
	description = "description DNAT example"
  
	# Using primary_ip from edge gateway
	external_address         = "89.32.25.10"
	internal_address         = "11.11.11.4"

	dnat_external_port = "8080"
}
`

const testAccNATRuleResourceConfigUpdateDnatWithVDCGroup = `
resource "cloudavenue_edgegateway_nat_rule" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_group.id
  
	name        = "example-dnat"
	rule_type   = "DNAT"
	description = "description DNAT example Updated!!"
  
	# Using primary_ip from edge gateway
	external_address         = "89.32.25.10"
	internal_address         = "4.11.11.11"

	priority = 25
}
`

func natRuleSnatTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "external_address"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-snat"),
		resource.TestCheckResourceAttr(resourceName, "description", "description SNAT example"),
		resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.0/24"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "SNAT"),
		resource.TestCheckResourceAttr(resourceName, "snat_destination_address", "8.8.8.8"),
		resource.TestCheckResourceAttr(resourceName, "priority", "10"),
		resource.TestCheckNoResourceAttr(resourceName, "dnat_external_port"),
	)
}

func natRuleSnatUpdateTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "external_address"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-snat"),
		resource.TestCheckResourceAttr(resourceName, "description", "description SNAT example Updated!!"),
		resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.0/24"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "SNAT"),
		resource.TestCheckResourceAttr(resourceName, "snat_destination_address", "9.9.9.9"),
		resource.TestCheckResourceAttr(resourceName, "priority", "0"),
		resource.TestCheckNoResourceAttr(resourceName, "dnat_external_port"),
	)
}

func natRuleDnatTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "external_address"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-dnat"),
		resource.TestCheckResourceAttr(resourceName, "description", "description DNAT example"),
		resource.TestCheckResourceAttr(resourceName, "internal_address", "11.11.11.4"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
		resource.TestCheckResourceAttr(resourceName, "dnat_external_port", "8080"),
		resource.TestCheckResourceAttr(resourceName, "priority", "0"),
		resource.TestCheckNoResourceAttr(resourceName, "snat_destination_address"),
	)
}

func natRuleDnatUpdateTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "external_address"),
		resource.TestCheckResourceAttr(resourceName, "name", "example-dnat"),
		resource.TestCheckResourceAttr(resourceName, "description", "description DNAT example Updated!!"),
		resource.TestCheckResourceAttr(resourceName, "internal_address", "4.11.11.11"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "DNAT"),
		resource.TestCheckNoResourceAttr(resourceName, "dnat_external_port"),
		resource.TestCheckResourceAttr(resourceName, "priority", "25"),
		resource.TestCheckNoResourceAttr(resourceName, "snat_destination_address"),
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
		resource.TestCheckNoResourceAttr(resourceName, "snat_destination_address"),
		resource.TestCheckNoResourceAttr(resourceName, "dnat_external_port"),
	)
}

func TestAccNATRuleResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_nat_rule.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// * Test with VDC
			{
				// Apply test Snat
				Config: tests.ConcatTests(testAccNATRuleResourceConfigSnat, testAccEdgeGatewayResourceConfig),
				Check:  natRuleSnatTestCheck(resourceName),
			},
			{
				// Update test Snat
				Config: tests.ConcatTests(testAccNATRuleResourceConfigUpdateSnat, testAccEdgeGatewayResourceConfig),
				Check:  natRuleSnatUpdateTestCheck(resourceName),
			},
			{
				// Delete test Snat
				Destroy: true,
				Config:  tests.ConcatTests(testAccNATRuleResourceConfigUpdateSnat, testAccEdgeGatewayResourceConfig),
			},
			{
				// Apply test Dnat
				Config: tests.ConcatTests(testAccNATRuleResourceConfigDnat, testAccEdgeGatewayResourceConfig),
				Check:  natRuleDnatTestCheck(resourceName),
			},
			{
				// Update test Dnat
				Config: tests.ConcatTests(testAccNATRuleResourceConfigUpdateDnat, testAccEdgeGatewayResourceConfig),
				Check:  natRuleDnatUpdateTestCheck(resourceName),
			},
			{
				// Delete test Dnat
				Destroy: true,
				Config:  tests.ConcatTests(testAccNATRuleResourceConfigUpdateDnat, testAccEdgeGatewayResourceConfig),
			},
			{
				// Apply test Reflexive
				Config: tests.ConcatTests(testAccNATRuleResourceConfigReflexive, testAccEdgeGatewayResourceConfig),
				Check:  natRuleReflexiveTestCheck(resourceName),
			},
			{
				// Import test ID and Name
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNATRuleResourceImportStateIDFuncWithIDAndName(resourceName),
			},
			{
				// Import test Name and ID
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNATRuleResourceImportStateIDFuncWithNameAndID(resourceName),
			},
			// * Test with VDCGroup
			{
				// Apply test Dnat
				Config: tests.ConcatTests(testAccNATRuleResourceConfigDnatWithVDCGroup, testAccEdgeGatewayGroupResourceConfig),
				Check:  natRuleDnatTestCheck(resourceName),
			},
			{
				// Update test Dnat
				Config: tests.ConcatTests(testAccNATRuleResourceConfigUpdateDnatWithVDCGroup, testAccEdgeGatewayGroupResourceConfig),
				Check:  natRuleDnatUpdateTestCheck(resourceName),
			},
			{
				// Import test ID and Name
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNATRuleResourceImportStateIDFuncWithIDAndName(resourceName),
			},
			{
				// Import test Name and ID
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNATRuleResourceImportStateIDFuncWithNameAndID(resourceName),
			},
		},
	})
}

func testAccNATRuleResourceImportStateIDFuncWithIDAndName(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		// edgeGatewayIDOrName.ipSetName
		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccNATRuleResourceImportStateIDFuncWithNameAndID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		// edgeGatewayIDOrName.ipSetName
		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_name"], rs.Primary.Attributes["id"]), nil
	}
}
