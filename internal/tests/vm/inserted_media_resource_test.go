// Package vm provides the acceptance tests for the provider.
package vm

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccVMInsertedMediaResourceConfig = `
resource "cloudavenue_vm_inserted_media" "example" {
	catalog = "catalog-example"
	name    = "debian-9.9.0-amd64-netinst.iso"
	vapp_name = "vapp-example"
	vm_name   = "vm-example"
  }
`

func TestAccVMInsertedMediaResource(t *testing.T) {
	resourceName := "cloudavenue_vm_inserted_media.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVMInsertedMediaResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					resource.TestCheckResourceAttr(resourceName, "catalog", "catalog-example"),
					resource.TestCheckResourceAttr(resourceName, "name", "debian-9.9.0-amd64-netinst.iso"),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "vapp-example"),
					resource.TestCheckResourceAttr(resourceName, "vm_name", "vm-example"),
				),
			},
		},
	})
}
