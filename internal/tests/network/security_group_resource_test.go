package network

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccSecurityGroupResourceConfig = `
resource "cloudavenue_network_security_group" "example" {
}
`


func TestAccSecurityGroupResource(t *testing.T) {
	resourceName := "cloudavenue_network_security_group.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccSecurityGroupResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Uncomment if you want to test update or delete this block
			// {
			// 	// Update test
			// 	Config: strings.Replace(testAccSecurityGroupResourceConfig, "old", "new", 1),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttrSet(resourceName, "id"),
			// 	),
			// },
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
