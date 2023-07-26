package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccRuleResourceConfig = `
resource "cloudavenue_nat_rule" "example" {
}
`

const testAccRuleResourceConfigUpdate = `
resource "cloudavenue_nat_rule" "example" {
}
`

func ruleTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	)
}

func TestAccRuleResource(t *testing.T) {
	resourceName := "cloudavenue_nat_rule.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccRuleResourceConfig,
				Check: ruleTestCheck(resourceName),
			},
			// Update testing
			{
				// Update test
				Config: testAccRuleResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Import State testing
			{
				// Import test
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}
