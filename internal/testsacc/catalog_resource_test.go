// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCatalogResourceConfig = `
resource "cloudavenue_catalog" "example" {
	name             = "test-catalog"
	description      = "catalog for files"
	delete_recursive = true
	delete_force     = true
}
`

func TestAccCatalogResource(t *testing.T) {
	resourceName := "cloudavenue_catalog.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-catalog"),
					resource.TestCheckResourceAttr(resourceName, "description", "catalog for files"),
					resource.TestCheckResourceAttr(resourceName, "delete_recursive", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_force", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccCatalogResourceUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-catalog"),
					resource.TestCheckResourceAttr(resourceName, "description", "catalog for ISO"),
					resource.TestCheckResourceAttr(resourceName, "delete_recursive", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_force", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "test-catalog",
				ImportStateVerify: true,
				// These fields can't be retrieved from catalog data
				ImportStateVerifyIgnore: []string{"delete_force", "delete_recursive"},
			},
		},
	})
}

func testAccCatalogResourceUpdate() string {
	return strings.Replace(testAccCatalogResourceConfig, "catalog for files", "catalog for ISO", 1)
}
