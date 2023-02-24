// Package tests provides the acceptance tests for the provider.
package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPublicIPDataSourceConfig = `
data "cloudavenue_publicip" "test" {}
`

func TestAccPublicIPDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_publicip.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccPublicIPDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr(dataSourceName, "id", "frangipane"),
				),
			},
		},
	})
}
