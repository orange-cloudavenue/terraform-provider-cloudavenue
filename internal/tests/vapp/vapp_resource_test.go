// Package vapp provides the acceptance tests for the provider.
package vapp

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccEVappResourceConfig = `
resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example vApp"
  }
`

func TestAccVappResource(t *testing.T) {
	resourceName := "cloudavenue_vapp.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Destroy: false,
				Config:  testAccEVappResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "MyVapp"),
					resource.TestCheckResourceAttr(resourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					resource.TestCheckResourceAttr(resourceName, "description", "This is an example vApp"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "MyVapp",
				ImportStateVerify: true,
				Destroy:           true,
			},
		},
	})
}
