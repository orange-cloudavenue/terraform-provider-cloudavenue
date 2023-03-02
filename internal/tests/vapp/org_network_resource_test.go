package vapp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccOrgNetworkResourceConfig = `
resource "cloudavenue_vapp_org_network" "example" {
  vapp_name = "vapp_test3"
  network_name = "test_remi"
}
`

const resourceName = "cloudavenue_vapp_org_network.example"

func TestAccOrgNetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccOrgNetworkResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "vapp_test3"),
					resource.TestCheckResourceAttr(resourceName, "network_name", "test_remi"),
				),
			},
			// Uncomment if you want to test update or delete this block
			// {
			// 	// Update test
			// 	Config: strings.Replace(testAccOrgNetworkResourceConfig, "old", "new", 1),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttrSet(resourceName, "id"),
			// 	),
			// },
			// ImportruetState testing
			{
				// Import test without vdc
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "vapp_test3.test_remi",
			},
			{
				// Import test with vdc
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "VDC_Frangipane.vapp_test3.test_remi",
			},
		},
	})
}
