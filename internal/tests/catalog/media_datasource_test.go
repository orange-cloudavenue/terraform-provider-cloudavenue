package catalog

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccCatalogMediaDataSourceConfig = `
data "cloudavenue_catalogs" "test" {}

data "cloudavenue_catalog_media" "example" {
	catalog_name = data.cloudavenue_catalogs.test.catalogs["catalog-example"].catalog_name
	name = "debian-9.9.0-amd64-netinst.iso"
}
`

func TestAccCatalogMediaDatasource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalog_media.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogMediaDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "catalog_id"),
					resource.TestCheckResourceAttr(dataSourceName, "catalog_name", "catalog-example"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "debian-9.9.0-amd64-netinst.iso"),
					resource.TestCheckResourceAttr(dataSourceName, "description", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_iso"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_published"),
					resource.TestCheckResourceAttrSet(dataSourceName, "size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profile"),
					resource.TestCheckResourceAttrSet(dataSourceName, "owner_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
				),
			},
		},
	})
}
