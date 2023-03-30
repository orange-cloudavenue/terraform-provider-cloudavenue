package catalog

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccCatalogMediasDataSourceConfig = `
data "cloudavenue_catalogs" "test" {}

data "cloudavenue_catalog_medias" "example" {
	catalog_name = data.cloudavenue_catalogs.test.catalogs["catalog-example"].name
}
`

func TestAccCatalogMediasDatasource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalog_medias.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogMediasDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "catalog_name", "catalog-example"),
					resource.TestCheckResourceAttr(dataSourceName, "medias_name.0", "debian-9.9.0-amd64-netinst.iso"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.is_iso", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.is_published", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.storage_profile", "gold"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.status", "RESOLVED"),
					resource.TestCheckResourceAttr(dataSourceName, "medias.debian-9.9.0-amd64-netinst.iso.created_at", "2023-02-28T15:10:57.136Z"),
				),
			},
		},
	})
}
