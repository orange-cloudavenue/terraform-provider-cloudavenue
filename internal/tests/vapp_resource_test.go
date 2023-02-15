// Package tests provides the acceptance tests for the provider.
package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccEVappResourceConfig = `
resource "cloudavenue_vapp" "test" {
	vapp_name = "vapp_name"
	description = "This is a test vapp"
  }
`

func TestAccVappResource(t *testing.T) {
	resourceName := "cloudavenue_vapp.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Destroy: false,
				Config:  testAccEVappResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "vapp_name"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test vapp"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "vapp_name",
				ImportStateVerify: true,
				Destroy:           true,
			},
		},
	})
}
