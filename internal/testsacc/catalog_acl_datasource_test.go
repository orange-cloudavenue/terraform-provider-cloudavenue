package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccACLDataSourceConfig = `
data "cloudavenue_catalog_acl" "example" {
	catalog_id = cloudavenue_catalog_acl.example.catalog_id
}
`

func TestAccACLDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalog_acl.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: ConcatTests(testAccCatalogResourceConfig, testAccACLResourceConfig, testAccACLDataSourceConfig),
				Check:  aclTestCheck(dataSourceName),
			},
			{
				Config: ConcatTests(testAccCatalogResourceConfig, testAccACLResourceConfigUpdateShareWithUsers, testAccACLDataSourceConfig),
				Check:  aclTestCheckShareWithUsers(dataSourceName),
			},
		},
	})
}
