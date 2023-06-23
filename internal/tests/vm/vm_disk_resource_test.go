// Package vm provides the acceptance tests for the provider.
package vm

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccVMDiskResourceConfig = `
resource "cloudavenue_vapp" "example" {
	name = "vapp_example"
	description = "This is a example vapp"
}

resource "cloudavenue_vm_disk" "example-detachable" {
	vapp_id = cloudavenue_vapp.example.id
	name = "disk-example-detachable"
	bus_type = "SATA"
	size_in_mb = 2048
	is_detachable = true
}
`

const testAccVMDiskWithVMResourceConfig = `
resource "cloudavenue_vapp" "example" {
	name = "vapp_example"
	description = "This is a example vapp"
}

resource "cloudavenue_vm_disk" "example-detachable-with-vm" {
	vapp_id = cloudavenue_vapp.example.id
	name = "disk-example-detachable-with-vm"
	bus_type = "SATA"
	size_in_mb = 2048
	is_detachable = true
	vm_id = cloudavenue_vm.example.id
}

data "cloudavenue_catalog_vapp_template" "example" {
	catalog_name = "Orange-Linux"
	template_name    = "debian_10_X64"
}

resource "cloudavenue_vm" "example" {
	name      = "example-vm"
	vapp_name = cloudavenue_vapp.example.name
	deploy_os = {
	  vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
	}
	settings = {
	  customization = {}
	}

	resource = {}
	state = {}
}
`

const testAccVMDiskInternalResourceConfig = `
resource "cloudavenue_vapp" "example" {
	name = "vapp_example"
	description = "This is a example vapp"
}

resource "cloudavenue_vm_disk" "example-internal" {
	vapp_id = cloudavenue_vapp.example.id
	bus_type = "SATA"
	size_in_mb = 2048
	is_detachable = false
	vm_id = cloudavenue_vm.example.id
}

data "cloudavenue_catalog_vapp_template" "example" {
	catalog_name = "Orange-Linux"
	template_name    = "debian_10_X64"
}

resource "cloudavenue_vm" "example" {
	name      = "example-vm"
	vapp_name = cloudavenue_vapp.example.name
	deploy_os = {
	  vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
	}
	settings = {
	  customization = {}
	}

	resource = {}
	state = {}
}
`

func TestAccVMDiskResource(t *testing.T) {
	const (
		resourceNameDetachable       = "cloudavenue_vm_disk.example-detachable"
		resourceNameDetachableWithVM = "cloudavenue_vm_disk.example-detachable-with-vm"
		resourceNameInternal         = "cloudavenue_vm_disk.example-internal"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// * EXTERNAL DISK
			{
				Config: testAccVMDiskResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDetachable, "id"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "name", "disk-example-detachable"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "size_in_mb", "2048"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "is_detachable", "true"),
					resource.TestCheckNoResourceAttr(resourceNameDetachable, "vm_id"),
				),
			},
			{
				Config: strings.Replace(testAccVMDiskResourceConfig, "2048", "4096", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDetachable, "id"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "name", "disk-example-detachable"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "size_in_mb", "4096"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "is_detachable", "true"),
					resource.TestCheckNoResourceAttr(resourceNameDetachable, "vm_id"),
				),
			},

			// * EXTERNAL DISK WITH VM
			{
				Config: testAccVMDiskWithVMResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDetachableWithVM, "id"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "name", "disk-example-detachable-with-vm"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "size_in_mb", "2048"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "is_detachable", "true"),
					resource.TestCheckResourceAttrSet(resourceNameDetachableWithVM, "vm_id"),
				),
			},
			{
				Config: strings.Replace(testAccVMDiskWithVMResourceConfig, "2048", "4096", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDetachableWithVM, "id"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "name", "disk-example-detachable-with-vm"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "size_in_mb", "4096"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "is_detachable", "true"),
					resource.TestCheckResourceAttrSet(resourceNameDetachableWithVM, "vm_id"),
				),
			},

			// * INTERNAL DISK
			{
				Config: testAccVMDiskInternalResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameInternal, "id"),
					resource.TestCheckResourceAttr(resourceNameInternal, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameInternal, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameInternal, "size_in_mb", "2048"),
					resource.TestCheckResourceAttr(resourceNameInternal, "is_detachable", "false"),
					resource.TestCheckResourceAttrSet(resourceNameInternal, "vm_id"),
				),
			},
			{
				Config: strings.Replace(testAccVMDiskInternalResourceConfig, "2048", "4096", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameInternal, "id"),
					resource.TestCheckResourceAttr(resourceNameInternal, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameInternal, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameInternal, "size_in_mb", "4096"),
					resource.TestCheckResourceAttr(resourceNameInternal, "is_detachable", "false"),
					resource.TestCheckResourceAttrSet(resourceNameInternal, "vm_id"),
				),
			},
		},
	})
}

// func getInternalVMDiskID(args ...string) resource.ImportStateIdFunc {
// 	return func(s *terraform.State) (string, error) {
// 		disk, ok := s.RootModule().Resources[resourceName]
// 		if !ok {
// 			return "", fmt.Errorf("Disk not found: %s", disk)
// 		}

// 		if disk.Primary.ID == "" {
// 			return "", fmt.Errorf("no ID is set for %s", disk)
// 		}
// 		if len(args) == 1 {
// 			return fmt.Sprintf("%s.vapp_test3.TestRomain.%s", args[0], disk.Primary.ID), nil
// 		}
// 		return fmt.Sprintf("vapp_test3.TestRomain.%s", disk.Primary.ID), nil
// 	}
// }
