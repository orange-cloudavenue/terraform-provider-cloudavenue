package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const TestAccGroupResourceConfig = `
resource "cloudavenue_vdc_group" "example" {
	name = "example"
	vdc_ids = [
		cloudavenue_vdc.example.id,
	]
}
`

const TestAccGroupResourceConfigUpdate = `
resource "cloudavenue_vdc_group" "example" {
	name = "example2"
	description = "Description of example2"
	vdc_ids = [
		cloudavenue_vdc.example.id,
		cloudavenue_vdc.example2.id,
	]
}
`

func groupTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", "example"),
		resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "1"),
		resource.TestCheckResourceAttrSet(resourceName, "status"),
		resource.TestCheckResourceAttrSet(resourceName, "type"),
	)
}

func groupTestCheckUpdated(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", "example2"),
		resource.TestCheckResourceAttr(resourceName, "description", "Description of example2"),
		resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "2"),
		resource.TestCheckResourceAttrSet(resourceName, "status"),
		resource.TestCheckResourceAttrSet(resourceName, "type"),
	)
}

func TestAccGroupResource(t *testing.T) {
	resourceName := "cloudavenue_vdc_group.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: ConcatTests(TestAccVDCResourceConfigWithoutVDCGroup, TestAccGroupResourceConfig),
				Check:  groupTestCheck(resourceName),
			},
			// Update testing
			{
				// Update test
				Config: ConcatTests(TestAccVDCResourceConfigWithoutVDCGroup, TestAccGroupResourceConfigUpdate),
				Check:  groupTestCheckUpdated(resourceName),
			},
			// Import State testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
