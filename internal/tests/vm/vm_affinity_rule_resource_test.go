package vm

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccVMAffinityRuleResourceConfig = `
resource "cloudavenue_vm_affinity_rule" "example" {
  name     = "test"
  polarity = "Affinity"

  vm_ids = [
    "urn:vcloud:vm:70b78935-cb64-4418-9607-4e3aeabbd168",
    "urn:vcloud:vm:c3912ae5-bbd1-45ae-8b1e-694d0a405a95"
  ]
}
`

func TestAccVmAffinityRuleResource(t *testing.T) {
	const resourceName = "cloudavenue_vm_affinity_rule.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVMAffinityRuleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "polarity", "Affinity"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "vm_ids.0", "urn:vcloud:vm:70b78935-cb64-4418-9607-4e3aeabbd168"),
					resource.TestCheckResourceAttr(resourceName, "vm_ids.1", "urn:vcloud:vm:c3912ae5-bbd1-45ae-8b1e-694d0a405a95"),
				),
			},
			// Uncomment if you want to test update or delete this block
			{
				// Update test
				Config: strings.Replace(testAccVMAffinityRuleResourceConfig, "test", "test2", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "test2"),
					resource.TestCheckResourceAttr(resourceName, "polarity", "Affinity"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "vm_ids.0", "urn:vcloud:vm:70b78935-cb64-4418-9607-4e3aeabbd168"),
					resource.TestCheckResourceAttr(resourceName, "vm_ids.1", "urn:vcloud:vm:c3912ae5-bbd1-45ae-8b1e-694d0a405a95"),
				),
			},
			// ImportruetState testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
