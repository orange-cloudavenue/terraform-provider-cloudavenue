// Package publicip provides the acceptance tests for the provider.
package publicip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccPublicIPDataSourceConfig = `
data "cloudavenue_publicip" "test" {}
`

func TestAccPublicIPDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_publicip.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
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
