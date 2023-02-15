// Package tests provides the acceptance tests for the provider.
package tests

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
)

const testAccEdgeGatewayResourceConfig = `
resource "cloudavenue_edge_gateway" "test" {
	owner_type = "vdc"
  owner_name = "myVDC01"
	tier0_vrf_name = "prvrf01iocb0000001allsp01"
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
	resourceName := "cloudavenue_edge_gateway.test"

	// resourceNameVDCGroup := "cloudavenue_edge_gateway.test-group"
	edgegw.ConfigEdgeGateway = func() edgegw.EdgeGatewayConfig {
		return edgegw.EdgeGatewayConfig{
			CheckJobDelay: 10 * time.Millisecond,
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Destroy: false,
				Config:  testAccEdgeGatewayResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "cc1f35c2-90a2-48d1-9359-62794faf44ad"),
					resource.TestCheckResourceAttr(resourceName, "edge_id", "cc1f35c2-90a2-48d1-9359-62794faf44ad"),
					resource.TestCheckResourceAttr(resourceName, "owner_type", "vdc"),
					resource.TestCheckResourceAttr(resourceName, "owner_name", "myVDC01"),
					resource.TestCheckResourceAttr(resourceName, "tier0_vrf_name", "prvrf01iocb0000001allsp01"),
					resource.TestCheckResourceAttr(resourceName, "edge_name", "edgeName"),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "edgeName",
				ImportStateVerify: true,
			},
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
