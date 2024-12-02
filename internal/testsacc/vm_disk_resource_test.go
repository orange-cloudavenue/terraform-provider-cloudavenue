// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// * EXTERNAL DISK
			{
				Config: testAccVMDiskResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceNameDetachable, "id", urn.TestIsType(urn.Disk)),
					resource.TestCheckResourceAttr(resourceNameDetachable, "name", "disk-example-detachable"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "size_in_mb", "2048"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "is_detachable", "true"),
					resource.TestCheckNoResourceAttr(resourceNameDetachable, "vm_name"),
					resource.TestCheckNoResourceAttr(resourceNameDetachable, "vm_id"),
					resource.TestCheckNoResourceAttr(resourceNameDetachable, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceNameDetachable, "vapp_id"),
				),
			},
			{
				Config: strings.Replace(testAccVMDiskResourceConfig, "2048", "4096", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceNameDetachable, "id", urn.TestIsType(urn.Disk)),
					resource.TestCheckResourceAttr(resourceNameDetachable, "name", "disk-example-detachable"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "size_in_mb", "4096"),
					resource.TestCheckResourceAttr(resourceNameDetachable, "is_detachable", "true"),
					resource.TestCheckNoResourceAttr(resourceNameDetachable, "vm_name"),
					resource.TestCheckNoResourceAttr(resourceNameDetachable, "vm_id"),
					resource.TestCheckNoResourceAttr(resourceNameDetachable, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceNameDetachable, "vapp_id"),
				),
			},
			{
				// Import test
				ResourceName:      resourceNameDetachable,
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccVMDiskResourceConfig,
				ImportStateIdFunc: testAccVMDiskResourceImportStateIDFunc(resourceNameDetachable, 1),
			},
			{
				// Import test
				ResourceName:      resourceNameDetachable,
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccVMDiskResourceConfig,
				ImportStateIdFunc: testAccVMDiskResourceImportStateIDFunc(resourceNameDetachable, 2),
			},
			{
				// Import test
				Config:       testAccVMDiskResourceConfig,
				ResourceName: resourceNameDetachable,
				Destroy:      true,
			},

			// * EXTERNAL DISK WITH VM
			{
				Config: testAccVMDiskWithVMResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceNameDetachableWithVM, "id", urn.TestIsType(urn.Disk)),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "name", "disk-example-detachable-with-vm"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "size_in_mb", "2048"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "is_detachable", "true"),
					resource.TestCheckNoResourceAttr(resourceNameDetachableWithVM, "vm_name"),
					resource.TestCheckResourceAttrSet(resourceNameDetachableWithVM, "vm_id"),
					resource.TestCheckNoResourceAttr(resourceNameDetachableWithVM, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceNameDetachableWithVM, "vapp_id"),
				),
			},
			{
				Config: strings.Replace(testAccVMDiskWithVMResourceConfig, "2048", "4096", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceNameDetachableWithVM, "id", urn.TestIsType(urn.Disk)),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "name", "disk-example-detachable-with-vm"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "size_in_mb", "4096"),
					resource.TestCheckResourceAttr(resourceNameDetachableWithVM, "is_detachable", "true"),
					resource.TestCheckNoResourceAttr(resourceNameDetachableWithVM, "vm_name"),
					resource.TestCheckResourceAttrSet(resourceNameDetachableWithVM, "vm_id"),
					resource.TestCheckNoResourceAttr(resourceNameDetachableWithVM, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceNameDetachableWithVM, "vapp_id"),
				),
			},
			{
				// Import test
				ResourceName:      resourceNameDetachableWithVM,
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccVMDiskWithVMResourceConfig,
				ImportStateIdFunc: testAccVMDiskResourceImportStateIDFunc(resourceNameDetachableWithVM, 3),
			},
			{
				// Import test
				Config:       testAccVMDiskWithVMResourceConfig,
				ResourceName: resourceNameDetachableWithVM,
				Destroy:      true,
			},

			// * INTERNAL DISK
			{
				Config: testAccVMDiskInternalResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameInternal, "id"), // Internal Disk has ID 123456
					resource.TestCheckResourceAttr(resourceNameInternal, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameInternal, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameInternal, "size_in_mb", "2048"),
					resource.TestCheckResourceAttr(resourceNameInternal, "is_detachable", "false"),
					resource.TestCheckNoResourceAttr(resourceNameInternal, "vm_name"),
					resource.TestCheckResourceAttrSet(resourceNameInternal, "vm_id"),
					resource.TestCheckNoResourceAttr(resourceNameInternal, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceNameInternal, "vapp_id"),
				),
			},
			{
				Config: strings.Replace(testAccVMDiskInternalResourceConfig, "2048", "4096", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameInternal, "id"), // Internal Disk has ID 123456
					resource.TestCheckResourceAttr(resourceNameInternal, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameInternal, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceNameInternal, "size_in_mb", "4096"),
					resource.TestCheckResourceAttr(resourceNameInternal, "is_detachable", "false"),
					resource.TestCheckNoResourceAttr(resourceNameInternal, "vm_name"),
					resource.TestCheckResourceAttrSet(resourceNameInternal, "vm_id"),
					resource.TestCheckNoResourceAttr(resourceNameInternal, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceNameInternal, "vapp_id"),
				),
			},
			{
				// Import test
				ResourceName:      resourceNameInternal,
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccVMDiskInternalResourceConfig,
				ImportStateIdFunc: testAccVMDiskResourceImportStateIDFunc(resourceNameInternal, 4),
			},
			{
				// Import test
				Config:       testAccVMDiskInternalResourceConfig,
				ResourceName: resourceNameInternal,
				Destroy:      true,
			},
		},
	})
}

// testAccVMDiskResourceConfig is a helper function that returns id of import
//
//	`resourceName` is the name of the resource
//	`typeOfImportID` is the type of import ID that we want to test:
//	- Option 1: `vapp_id` and `disk_id` -> Detachable disk
//	- Option 2: `vdc`, `vapp_id` and `disk_id` -> Detachable disk with VDC Parameter
//	- Option 3: `vapp_id`, `vm_id` and `disk_id` -> Internal disk or Detachable disk with VM Parameter
//	- Option 4: `vdc`, `vapp_id`, `vm_id` and `disk_id` -> Internal disk with VDC Parameter or Detachable disk with VDC Parameter and VM Parameter
func testAccVMDiskResourceImportStateIDFunc(resourceName string, typeOfImportID int) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		switch typeOfImportID {
		case 1:
			return fmt.Sprintf("%s.%s", rs.Primary.Attributes["vapp_id"], rs.Primary.Attributes["id"]), nil
		case 2:
			return fmt.Sprintf("%s.%s.%s", rs.Primary.Attributes["vdc"], rs.Primary.Attributes["vapp_id"], rs.Primary.Attributes["id"]), nil
		case 3:
			return fmt.Sprintf("%s.%s.%s", rs.Primary.Attributes["vapp_id"], rs.Primary.Attributes["vm_id"], rs.Primary.Attributes["id"]), nil
		case 4:
			return fmt.Sprintf("%s.%s.%s.%s", rs.Primary.Attributes["vdc"], rs.Primary.Attributes["vapp_id"], rs.Primary.Attributes["vm_id"], rs.Primary.Attributes["id"]), nil
		default:
			return "", fmt.Errorf("Invalid type of import ID")
		}
	}
}
