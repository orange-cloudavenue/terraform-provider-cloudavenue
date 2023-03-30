package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccAlbPoolResourceConfig = `
resource "cloudavenue_alb_pool" "example" {
}
`

func TestAccAlbPoolResource(t *testing.T) {
	const resourceName = "cloudavenue_alb_pool.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccAlbPoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Uncomment if you want to test update or delete this block
			// {
			// 	// Update test
			// 	Config: strings.Replace(testAccAlbPoolResourceConfig, "old", "new", 1),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttrSet(resourceName, "id"),
			// 	),
			// },
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
