package edgegw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccPortProfilesResourceConfig = `
data "cloudavenue_vdc" "example" {
	name = "VDC_Test"
}

resource "cloudavenue_edgegateway_app_port_profile" "example" {
	name        = "example-rule"
	description = "Application port profile for example"
	vdc  		= data.cloudavenue_vdc.example.id
  
	app_ports = [
	  {
		protocol = "ICMPv4"
	  },
	  {
		protocol = "TCP"
		ports = [
			"80",
			"443",
		]
	  },
	]
  }
  
`

const testAccPortProfilesResourceConfigUpdate = `
data "cloudavenue_vdc" "example" {
	name = "VDC_Test"
}

resource "cloudavenue_edgegateway_app_port_profile" "example" {
	name        = "example-rule"
	description = "Application port profile for example"
	vdc  		= data.cloudavenue_vdc.example.id
  
	app_ports = [
	  {
		protocol = "ICMPv4"
	  },
	  {
		protocol = "TCP"
		ports = [
			"80",
			"443",
			"8080",
		]
	  },
	  {
		protocol = "UDP"
		ports = [
			"53",
		]
	  }
	]
  }
  
`

func TestAccPortProfilesResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_app_port_profile.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccPortProfilesResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "example-rule"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
					resource.TestCheckResourceAttr(resourceName, "description", "Application port profile for example"),
					resource.TestCheckResourceAttr(resourceName, "app_ports.#", "2"),
				),
			},
			{
				// Apply test
				Config: testAccPortProfilesResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "example-rule"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
					resource.TestCheckResourceAttr(resourceName, "description", "Application port profile for example"),
					resource.TestCheckResourceAttr(resourceName, "app_ports.#", "3"),
				),
			},
			// ImportruetState testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPortProfilesResourceImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccPortProfilesResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["vdc"], rs.Primary.Attributes["name"]), nil
	}
}
