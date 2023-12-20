package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCatalogMediaDataSourceConfig = `
data "cloudavenue_catalog_media" "example" {
	catalog_name = "catalog-example"
	name = "debian-9.9.0-amd64-netinst.iso"
}
`

func TestAccCatalogMediaDatasource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalog_media.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogMediaDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
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
