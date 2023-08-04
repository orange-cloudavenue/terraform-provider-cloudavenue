// Package catalog provides the acceptance tests for the provider.
package catalog

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccCatalogsDataSourceConfig = `
data "cloudavenue_catalogs" "example" {}
`

func TestAccCatalogsDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalogs.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs_name.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.%"),

					resource.TestCheckResourceAttrWith(dataSourceName, "catalogs.Orange-Linux.id", uuid.TestIsType(uuid.Catalog)),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.Orange-Linux.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.Orange-Linux.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.Orange-Linux.preserve_identity_information"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.Orange-Linux.number_of_media"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.Orange-Linux.media_item_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.Orange-Linux.is_shared"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.Orange-Linux.is_published"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.Orange-Linux.is_local"),

					resource.TestCheckNoResourceAttr(dataSourceName, "catalogs.Orange-Linux.owner_name"),  // In Orange-Linux catalog, owner_name is empty
					resource.TestCheckNoResourceAttr(dataSourceName, "catalogs.Orange-Linux.description"), // In Orange-Linux catalog, description is empty
					resource.TestCheckNoResourceAttr(dataSourceName, "catalogs.Orange-Linux.is_cached"),   // In Orange-Linux catalog, is_cached is false
				),
			},
		},
	})
}
