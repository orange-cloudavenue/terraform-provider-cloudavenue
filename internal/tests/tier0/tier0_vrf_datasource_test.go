// Package tier0 provides the acceptance tests for the provider.
package tier0

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccTier0VrfDataSourceConfig = `
data "cloudavenue_tier0_vrf" "test" {
	name = "prvrf01eocb0006205allsp01"
}
`

func TestAccTier0VrfDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccTier0VrfDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "name", "prvrf01eocb0006205allsp01"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "id", "ca606aba-4bd2-5e66-a975-1ebb3ae2eca9"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "class_service", "VRF_STANDARD"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "tier0_provider", "pr01e02t0sp16"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "services.#", "3"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "services.0.service", "OBJECT_STORAGE"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "services.0.vlan_id", ""),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "services.1.service", "INTERNET"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "services.1.vlan_id", ""),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "services.2.service", "ADMIN"),
					resource.TestCheckResourceAttr("data.cloudavenue_tier0_vrf.test", "services.2.vlan_id", ""),
				),
			},
		},
	})
}
