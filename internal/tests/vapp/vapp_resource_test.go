// Package vapp provides the acceptance tests for the provider.
package vapp

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccEVappResourceConfig = `
resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example vApp"

	lease = {
		runtime_lease_in_sec = 3600
		storage_lease_in_sec = 3600
	}

	guest_properties = {
		"key" = "Value"
	}
  }
`

const testAccEVappResourceUpdatedConfig = `
resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example modified vApp"
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
					resource.TestCheckResourceAttr(resourceName, "lease.runtime_lease_in_sec", "3600"),
					resource.TestCheckResourceAttr(resourceName, "lease.storage_lease_in_sec", "3600"),
					resource.TestCheckResourceAttr(resourceName, "guest_properties.key", "Value"),
				),
			},
			// Update
			{
				Destroy: false,
				Config:  testAccEVappResourceUpdatedConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "MyVapp"),
					resource.TestCheckResourceAttr(resourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					resource.TestCheckResourceAttr(resourceName, "description", "This is an example modified vApp"),
					resource.TestCheckResourceAttr(resourceName, "lease.runtime_lease_in_sec", "0"),
					resource.TestCheckResourceAttr(resourceName, "lease.storage_lease_in_sec", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "guest_properties.#"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "MyVapp",
				ImportStateVerify: true,
			},
		},
	})
}
