// Package vm provides the acceptance tests for the provider.
package vm

import (
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
	name = "disk-example"
	bus_type = "SATA"
	size_in_mb = 2048
	is_detachable = true
}
`

func TestAccVMDiskResource(t *testing.T) {
	const resourceName = "cloudavenue_vm_disk.example-detachable"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVMDiskResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "disk-example"),
					resource.TestCheckResourceAttr(resourceName, "bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile", "gold"),
					resource.TestCheckResourceAttr(resourceName, "size_in_mb", "2048"),
					resource.TestCheckResourceAttr(resourceName, "is_detachable", "true"),
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
