package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccNATRuleDataSourceConfig = `
data "cloudavenue_edgegateway_nat_rule" "example" {
	depends_on = [cloudavenue_edgegateway.example_with_vdc, cloudavenue_edgegateway_nat_rule.example]
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	name = "example-snat"
}
`

func TestAccNatRuleDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway_nat_rule.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccNATRuleResourceConfigSnat, testAccNATRuleDataSourceConfig, testAccEdgeGatewayResourceConfig),
				Check:  natRuleSnatTestCheck(dataSourceName),
			},
		},
	})
}
