package catalog

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccACLDataSourceConfig = `
data "cloudavenue_catalog_acl" "example" {
	catalog_id = cloudavenue_catalog_acl.example.catalog_id
}
`

func TestAccACLDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalog_acl.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccCatalogResourceConfig, testAccACLResourceConfig, testAccACLDataSourceConfig),
				Check:  aclTestCheck(dataSourceName),
			},
			{
				Config: tests.ConcatTests(testAccCatalogResourceConfig, testAccACLResourceConfigUpdateShareWithUsers, testAccACLDataSourceConfig),
				Check:  aclTestCheckShareWithUsers(dataSourceName),
			},
		},
	})
}
