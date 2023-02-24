// Package catalog provides the acceptance tests for the provider.
package catalog

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccCatalogsDataSourceConfig = `
data "cloudavenue_catalogs" "test" {}
`

func TestAccCatalogsDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_catalogs.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "catalogs_name.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}
