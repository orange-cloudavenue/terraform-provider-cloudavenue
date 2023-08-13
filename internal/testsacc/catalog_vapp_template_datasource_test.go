// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccCatalogVappTemplateDataSourceConfig = `
data "cloudavenue_catalog_vapp_template" "example" {
	catalog_name  	= "Orange-Linux"
	template_name 	= "UBUNTU_20.04"
}
`

func TestAccCatalogVappTemplateDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalog_vapp_template.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogVappTemplateDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(dataSourceName, "id", uuid.TestIsType(uuid.VAPPTemplate)),
					// Catalog
					resource.TestCheckResourceAttr(dataSourceName, "catalog_name", "Orange-Linux"),
					resource.TestCheckResourceAttrWith(dataSourceName, "catalog_id", uuid.TestIsType(uuid.Catalog)),

					resource.TestCheckResourceAttr(dataSourceName, "template_name", "UBUNTU_20.04"),
					resource.TestCheckResourceAttrWith(dataSourceName, "template_id", uuid.TestIsType(uuid.VAPPTemplate)),
					// Other
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vm_names.#"),
				),
			},
		},
	})
}
