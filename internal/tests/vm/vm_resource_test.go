// Package vm provides the acceptance tests for the provider.
package vm

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const resourceNameVM = "cloudavenue_vm.example"

const testAccVMResourceConfigFromVappTemplate = `
data "cloudavenue_catalog_vapp_template" "example" {
	catalog_name = "Orange-Linux"
	vapp_name    = "debian_10_X64"
}

resource "cloudavenue_vapp" "example" {
	vapp_name = "vapp_example"
	description = "This is a example vapp"
}

resource "cloudavenue_vm" "example" {
	vm_name         	= "example-vm"
	description 		= "This is a example vm"
	accept_all_eulas 	= true
	vapp_name 			= cloudavenue_vapp.example.vapp_name
	vapp_template_id 	= data.cloudavenue_catalog_vapp_template.example.id
}
`

func TestAccVMResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVMResourceConfigFromVappTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVM, "id"),
					resource.TestCheckResourceAttr(resourceNameVM, "vapp_name", "vapp_example"),
					resource.TestCheckResourceAttr(resourceNameVM, "vm_name", "example-vm"),
					resource.TestCheckResourceAttr(resourceNameVM, "accept_all_eulas", "true"),
					resource.TestCheckResourceAttr(resourceNameVM, "expose_hardware_virtualization", "false"),
					resource.TestCheckResourceAttr(resourceNameVM, "os_type", "debian10_64Guest"),
					resource.TestCheckResourceAttr(resourceNameVM, "description", "This is a example vm"),

					resource.TestCheckResourceAttr(resourceNameVM, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					resource.TestCheckResourceAttr(resourceNameVM, "storage_profile", "gold"),

					resource.TestCheckResourceAttr(resourceNameVM, "status_code", "4"),
					resource.TestCheckResourceAttr(resourceNameVM, "status_text", "POWERED_ON"),

					// Resource
					resource.TestCheckResourceAttr(resourceNameVM, "resource.cpu_cores", "1"),
					resource.TestCheckResourceAttr(resourceNameVM, "resource.cpu_hot_add_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameVM, "resource.cpus", "1"),
					resource.TestCheckResourceAttr(resourceNameVM, "resource.memory", "1024"),
					resource.TestCheckResourceAttr(resourceNameVM, "resource.memory_hot_add_enabled", "false"),

					resource.TestCheckResourceAttrSet(resourceNameVM, "vapp_template_id"),
					resource.TestCheckResourceAttrSet(resourceNameVM, "id"),
				),
			},
		},
	})
}
