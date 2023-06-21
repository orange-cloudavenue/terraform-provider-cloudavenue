package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccVmDataSourceConfig = `
data "cloudavenue_vm_vm" "example" {
}
`


func TestAccVmDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vm_vm.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVmResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
