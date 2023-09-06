// Package tier0 provides the acceptance tests for the provider.
package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccTier0VrfsDataSourceConfig = `
data "cloudavenue_tier0_vrfs" "test" {}
`

func TestAccTier0VrfsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccTier0VrfsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrfs.test", "names.0", "prvrf01eocb0006205allsp01"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrfs.test", "id", "d767aafb-f919-5cc7-8d97-0287f2d672ab"),
				),
			},
		},
	})
}
