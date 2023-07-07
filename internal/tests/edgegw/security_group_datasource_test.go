package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccSecurityGroupDataSourceConfig = `
data "cloudavenue_edgegateway_security_group" "example" {
	name            = cloudavenue_edgegateway_security_group.example.name
	edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
}
`

func TestAccSecurityGroupDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway_security_group.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccSecurityGroupResourceConfig, testAccSecurityGroupDataSourceConfig),
				Check:  securityGroupTestCheck(dataSourceName),
			},
		},
	})
}
