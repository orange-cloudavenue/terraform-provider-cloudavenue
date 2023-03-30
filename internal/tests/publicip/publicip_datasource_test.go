// Package publicip provides the acceptance tests for the provider.
package publicip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccPublicIPsDataSourceConfig = `
data "cloudavenue_publicips" "test" {}
`

func TestAccPublicIPsDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_publicips.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccPublicIPsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					// Verify placeholder public_ips attribute
					resource.TestCheckResourceAttrSet(dataSourceName, "public_ips.#"),
				),
			},
		},
	})
}
