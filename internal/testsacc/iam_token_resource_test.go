package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccTokenResourceConfig = `
resource "cloudavenue_iam_token" "example" {
	name            = "example"

	save_in_tfstate = true
	save_in_file    = true
	print_token     = true
}
`

const testAccTokenResourceConfigUpdate = `
resource "cloudavenue_iam_token" "example" {
	name            = "exampleUpdated"

	save_in_tfstate = true
	save_in_file    = true
	print_token     = true
}
`

func TestAccTokenResource(t *testing.T) {
	resourceName := "cloudavenue_iam_token.example"

	t.Cleanup(deleteFile("token.json", t))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccTokenResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Token)),
					resource.TestCheckResourceAttr(resourceName, "name", "example"),
					resource.TestCheckResourceAttr(resourceName, "save_in_tfstate", "true"),
					resource.TestCheckResourceAttr(resourceName, "save_in_file", "true"),
					resource.TestCheckResourceAttr(resourceName, "print_token", "true"),
					resource.TestCheckResourceAttr(resourceName, "file_name", "token.json"),
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					testCheckFileExists("token.json"),
				),
			},
			{
				// Update test
				// Any change generates replacement
				Config:             testAccTokenResourceConfigUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
