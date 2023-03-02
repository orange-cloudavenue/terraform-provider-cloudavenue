// Package edgegw provides the acceptance tests for the provider.
package edgegw

import (
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccEdgeGatewayResourceConfig = `
resource "cloudavenue_edgegateway" "test" {
	owner_type = "vdc"
	owner_name = "MyVDC"
	tier0_vrf_name = "prvrf01eocb0006205allsp01"
}
`

// const testAccEdgeGatewayResourceWithBadOwnerTypeConfig = `
// resource "cloudavenue_edge_gateway" "test" {
// 	owner_type = "vdc-bad"
//   owner_name = "myVDC01"
// 	tier0_vrf_name = "prvrf01iocb0000001allsp01"
// }
// `

func TestAccEdgeGatewayResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway.test"

	// resourceNameVDCGroup := "cloudavenue_edge_gateway.test-group"
	edgegw.ConfigEdgeGateway = func() edgegw.EdgeGatewayConfig {
		return edgegw.EdgeGatewayConfig{
			CheckJobDelay: 10 * time.Millisecond,
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Destroy: false,
				Config:  testAccEdgeGatewayResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`(urn:vcloud:gateway:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttr(resourceName, "owner_type", "vdc"),
					resource.TestCheckResourceAttr(resourceName, "owner_name", "MyVDC"),
					resource.TestCheckResourceAttr(resourceName, "tier0_vrf_name", "prvrf01eocb0006205allsp01"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`tn01e02ocb0006205spt[0-9]{3}`)),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateId:     "edgeName",
			// 	ImportStateVerify: true,
			// },
			// check bad owner_type
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/609
			// {
			// 	Config:      testAccEdgeGatewayResourceWithBadOwnerTypeConfig,
			// 	ExpectError: regexp.MustCompile(`.*`),
			// 	Destroy:     true,
			// },
			// {
			// 	Config: testAccEdgeGatewayGroupResourceConfig,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceNameVDCGroup, "id", "cc1f35c2-90a2-48d1-9359-62794faf44ae"),
			// 		resource.TestCheckResourceAttr(resourceNameVDCGroup, "edge_id", "cc1f35c2-90a2-48d1-9359-62794faf44ae"),
			// 		resource.TestCheckResourceAttr(resourceNameVDCGroup, "owner_type", "vdc-group"),
			// 		resource.TestCheckResourceAttr(resourceNameVDCGroup, "owner_name", "myVDC02"),
			// 		resource.TestCheckResourceAttr(resourceNameVDCGroup, "tier0_vrf_id", "prvrf01iocb0000001allsp01"),
			// 		resource.TestCheckResourceAttr(resourceNameVDCGroup, "edge_name", "edgeName2"),
			// 		resource.TestCheckResourceAttr(resourceNameVDCGroup, "description", "description"),
			// 	),
			// },
		},
	})
}
