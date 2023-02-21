package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../examples -test
const testAccVMInternalDiskResourceConfig = `
resource "cloudavenue_vapp_vm_internal_disk" "example" {
  vapp_name       = "vapp_test3"
  vm_name         = "TestRomain"
	allow_vm_reboot = true
	internal_disk {
		bus_type      = "sata"
		size_in_mb    = "500"
		bus_number    = 0
		unit_number   = 1
	}
}
`

const resourceName = "cloudavenue_vapp_vm_internal_disk.example"

func TestAccVMInternalDiskResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVMInternalDiskResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "vapp_test3"),
					resource.TestCheckResourceAttr(resourceName, "vm_name", "TestRomain"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.bus_type", "sata"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.size_in_mb", "500"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.bus_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.unit_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "allow_vm_reboot", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.storage_profile", "gold"),
				),
			},
			{
				// Update test
				Config: strings.Replace(testAccVMInternalDiskResourceConfig, "500", "600", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "vapp_test3"),
					resource.TestCheckResourceAttr(resourceName, "vm_name", "TestRomain"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.bus_type", "sata"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.size_in_mb", "600"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.bus_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.unit_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "allow_vm_reboot", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal_disk.storage_profile", "gold"),
				),
			},
			// ImportruetState testing
			{
				// Import test
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_vm_reboot"},
				ImportStateIdFunc:       getInternalVMDiskID(),
			},
			{
				// Import test with vdc
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_vm_reboot"},
				ImportStateIdFunc:       getInternalVMDiskID("VDC_Frangipane"),
			},
		},
	})
}

func getInternalVMDiskID(args ...string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		disk, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Disk not found: %s", disk)
		}

		if disk.Primary.ID == "" {
			return "", fmt.Errorf("no ID is set for %s", disk)
		}
		if len(args) == 1 {
			return fmt.Sprintf("%s.vapp_test3.TestRomain.%s", args[0], disk.Primary.ID), nil
		}
		return fmt.Sprintf("vapp_test3.TestRomain.%s", disk.Primary.ID), nil
	}
}
