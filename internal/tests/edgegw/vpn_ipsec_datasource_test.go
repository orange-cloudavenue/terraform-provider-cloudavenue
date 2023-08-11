package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccVPNIPSecDataSourceConfig = `
data "cloudavenue_edgegateway_vpn_ipsec" "example" {
  depends_on = [ cloudavenue_edgegateway_vpn_ipsec.example ]
  edge_gateway_id = cloudavenue_edgegateway_vpn_ipsec.example.edge_gateway_id
  name = cloudavenue_edgegateway_vpn_ipsec.example.name
}
`

func TestAccVPNIPSecDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegateway_vpn_ipsec.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccVPNIPSecDataSourceConfig, testAccVPNIPSecResourceConfigCustomize, MytestAccEdgeGatewayGroupResourceConfig, MytestAccVDCResourceConfig),
				Check:  vpnIPSecTestCheckCustomize(dataSourceName),
			},
		},
	})
}
