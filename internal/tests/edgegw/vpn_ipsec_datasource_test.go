package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccVpnIpsecDataSourceConfig = `
data "cloudavenue_edgegw_vpn_ipsec" "example" {
}
`

func TestAccVpnIpsecDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_edgegw_vpn_ipsec.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccVpnIpsecResourceConfig, testAccVpnIpsecDataSourceConfig),
				Check: vpnIpsecTestCheck(dataSourceName),
			},
		},
	})
}
