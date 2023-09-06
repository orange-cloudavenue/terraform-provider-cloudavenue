package testsacc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccSecurityGroupResourceConfig = `
resource "cloudavenue_edgegateway_security_group" "example" {

  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
  name            = "example"
  description     = "This is an example security group"
  member_org_network_ids = [
    cloudavenue_network_routed.example.id
  ]
}

data "cloudavenue_edgegateways" "example" {}


resource "cloudavenue_network_routed" "example" {
	name        = "MyOrgNet"
	description = "This is an example Net"
  
	edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
  
	gateway       = "192.168.1.254"
	prefix_length = 24
  
	dns1 = "1.1.1.1"
	dns2 = "8.8.8.8"
  
	dns_suffix = "example"
  
	static_ip_pool = [
	  {
		start_address = "192.168.1.10"
		end_address   = "192.168.1.20"
	  }
	]
  }
`

const testAccSecurityGroupResourceConfigUpdate = `
resource "cloudavenue_edgegateway_security_group" "example" {

	edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
	name            = "example-updated"
	description     = "This is an example security group updated"
	member_org_network_ids = [
	  cloudavenue_network_routed.example.id
	]
  }
  
  data "cloudavenue_edgegateways" "example" {}
  
  
  resource "cloudavenue_network_routed" "example" {
	  name        = "MyOrgNet"
	  description = "This is an example Net"
	
	  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
	
	  gateway       = "192.168.1.254"
	  prefix_length = 24
	
	  dns1 = "1.1.1.1"
	  dns2 = "8.8.8.8"
	
	  dns_suffix = "example"
	
	  static_ip_pool = [
		{
		  start_address = "192.168.1.10"
		  end_address   = "192.168.1.20"
		}
	  ]
	}
`

func securityGroupTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
		resource.TestCheckResourceAttr(resourceName, "description", "This is an example security group"),
		resource.TestCheckResourceAttr(resourceName, "name", "example"),
		resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "1"),
		resource.TestCheckResourceAttrSet(resourceName, "member_org_network_ids.0"),
	)
}

func TestAccSecurityGroupResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_security_group.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccSecurityGroupResourceConfig,
				Check:  securityGroupTestCheck(resourceName),
			},
			{
				// Update test
				Config: testAccSecurityGroupResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is an example security group updated"),
					resource.TestCheckResourceAttr(resourceName, "name", "example-updated"),
					resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "member_org_network_ids.0"),
				),
			},
			// ImportruetState testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSecurityGroupResourceImportStateIDFuncWithID(resourceName),
			},
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSecurityGroupResourceImportStateIDFuncWithName(resourceName),
			},
		},
	})
}

func testAccSecurityGroupResourceImportStateIDFuncWithID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccSecurityGroupResourceImportStateIDFuncWithName(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["edge_gateway_name"], rs.Primary.Attributes["name"]), nil
	}
}
