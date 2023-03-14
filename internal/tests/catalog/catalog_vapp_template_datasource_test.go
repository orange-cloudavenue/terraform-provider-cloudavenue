// Package catalog provides the acceptance tests for the provider.
package catalog

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccCatalogVappTemplateDataSourceConfig = `
data "cloudavenue_catalog_vapp_template" "example" {
	catalog_name  	= "Orange-Linux"
	template_name 	= "UBUNTU_20.04"
}
`

func TestAccCatalogVappTemplateDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalog_vapp_template.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogVappTemplateDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "catalog_name", "Orange-Linux"),
					resource.TestCheckResourceAttr(dataSourceName, "template_name", "UBUNTU_20.04"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vm_names.#"),
				),
			},
		},
	})
}
