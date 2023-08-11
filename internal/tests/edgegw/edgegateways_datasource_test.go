// Package edgegw provides the acceptance tests for the provider.
package edgegw

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccEdgeGatewaysDataSourceConfig = `
data "cloudavenue_edgegateways" "example" {}
`

func TestAccEdgeGatewaysDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateways.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEdgeGatewaysDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute
					resource.TestMatchResourceAttr(dataSourceName, "id", regexp.MustCompile(`([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttrWith(dataSourceName, "edge_gateways.0.id", uuid.TestIsType(uuid.Gateway)),
				),
			},
		},
	})
}
