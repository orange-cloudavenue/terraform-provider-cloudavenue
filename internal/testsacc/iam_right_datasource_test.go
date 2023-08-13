package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccIAMRightDataSourceConfig = `
data "cloudavenue_iam_right" "example" {
	name = "Catalog: Add vApp from My Cloud"
}
`

func TestAccIamRightDatasource(t *testing.T) {
	dataSourceName := "data.cloudavenue_iam_right.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIAMRightDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bundle_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "category_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "Catalog: Add vApp from My Cloud"),
					resource.TestCheckResourceAttrSet(dataSourceName, "right_type"),
				),
			},
		},
	})
}
