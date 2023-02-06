package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTier0VrfsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccTier0VrfsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrfs.test", "tier0_vrfs.0.name", "prvrf01eocb0006205allsp01"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrfs.test", "id", "474baabd-09ea-5f2f-ad74-c4a3778fd16e"),
				),
			},
		},
	})
}

const testAccTier0VrfsDataSourceConfig = `
data "cloudavenue_tier0_vrfs" "test" {}
`
