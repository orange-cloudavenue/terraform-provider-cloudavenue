// Package tests provides the acceptance tests for the provider.
package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCatalogResourceConfig = `
resource "cloudavenue_catalog" "test" {
	catalog_name     = "test-catalog"
	description      = "catalog for files"
	delete_recursive = true
	delete_force     = true
}
`

func TestAccCatalogResource(t *testing.T) {
	resourceName := "cloudavenue_catalog.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "catalog_name", "test-catalog"),
					resource.TestCheckResourceAttr(resourceName, "description", "catalog for files"),
					resource.TestCheckResourceAttr(resourceName, "delete_recursive", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_force", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "href"),
					resource.TestCheckResourceAttrSet(resourceName, "catalog_version"),
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
