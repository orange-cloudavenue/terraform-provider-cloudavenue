// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCatalogDataSourceConfig = `
resource "cloudavenue_catalog" "test" {
	name             = "test-catalog"
	description      = "catalog for files"
	delete_recursive = true
	delete_force     = true
}

data "cloudavenue_catalog" "test" {
	name = cloudavenue_catalog.test.name
}
`

func TestAccCatalogDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalog.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", "test-catalog"),
					resource.TestCheckResourceAttr(dataSourceName, "description", "catalog for files"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "owner_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
				),
			},
		},
	})
}
