// Package edgegw provides the acceptance tests for the provider.
package edgegw

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccEdgeGatewayDataSourceConfig = `
data "cloudavenue_edgegateway" "test" {
	name = "tn01e02ocb0006205spt101"
}
`

func TestAccEdgeGatewayDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEdgeGatewayDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute
					resource.TestMatchResourceAttr(dataSourceName, "id", regexp.MustCompile(`(urn:vcloud:gateway:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tier0_vrf_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "owner_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "owner_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "lb_enabled"),
				),
			},
		},
	})
}
