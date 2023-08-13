// package testsacc provides the acceptance tests for the provider.
// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccRoleResourceConfig = `
resource "cloudavenue_iam_role" "example" {
    name        = "roletest"
    description = "A test role"
	rights = [
		"Catalog: Add vApp from My Cloud",
		"Catalog: Edit Properties",
		"Catalog: View Private and Shared Catalogs",
		"Organization vDC Compute Policy: View",
		"vApp Template / Media: Edit",
		"vApp Template / Media: View",
	]
}
`

func TestAccRoleResource(t *testing.T) {
	resourceName := "cloudavenue_iam_role.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccRoleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "roletest"),
					resource.TestCheckResourceAttr(resourceName, "description", "A test role"),
					resource.TestCheckTypeSetElemAttr(
						resourceName,
						"rights.*",
						"Catalog: Add vApp from My Cloud",
					),
					resource.TestCheckResourceAttr(resourceName, "rights.#", "6"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "roletest",
				ImportStateVerify: true,
			},
		},
	})
}
