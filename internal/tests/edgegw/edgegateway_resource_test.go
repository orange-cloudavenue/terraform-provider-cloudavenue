// Package edgegw provides the acceptance tests for the provider.
package edgegw

import (
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccEdgeGatewayResourceConfig = `
data "cloudavenue_tier0_vrfs" "example_with_vdc" {}

resource "cloudavenue_edgegateway" "example_with_vdc" {
  owner_name     = "MyVDC"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example_with_vdc.names.0
  owner_type     = "vdc"
  lb_enabled     = false
}
`

const testAccEdgeGatewayGroupResourceConfig = `
data "cloudavenue_tier0_vrfs" "example_with_group" {}

resource "cloudavenue_edgegateway" "example_with_group" {
  owner_name     = "MyVDCGroup"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example_with_group.names.0
  owner_type     = "vdc-group"
}
`

func TestAccEdgeGatewayResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway.example_with_vdc"
	resourceNameVDCGroup := "cloudavenue_edgegateway.example_with_group"

	edgegw.ConfigEdgeGateway = func() edgegw.EdgeGatewayConfig {
		return edgegw.EdgeGatewayConfig{
			CheckJobDelay: 10 * time.Second,
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeGatewayResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.Gateway.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "owner_type", "vdc"),
					resource.TestCheckResourceAttr(resourceName, "owner_name", "MyVDC"),
					resource.TestMatchResourceAttr(resourceName, "tier0_vrf_name", regexp.MustCompile(`prvrf01eocb0006205allsp[0-9]{2}`)),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`tn01e02ocb0006205spt[0-9]{3}`)),
					resource.TestCheckResourceAttr(resourceName, "lb_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
			{
				Destroy: true,
				Config:  testAccEdgeGatewayResourceConfig,
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "tn01e02ocb0006205spt101",
				ImportStateVerify: true,
			},
			// check bad owner_type
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/609
			// {
			// 	Config:      testAccEdgeGatewayResourceWithBadOwnerTypeConfig,
			// 	ExpectError: regexp.MustCompile(`.*`),
			// 	Destroy:     true,
			// },
			{
				Config: testAccEdgeGatewayGroupResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceNameVDCGroup, "id", regexp.MustCompile(uuid.Gateway.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceNameVDCGroup, "owner_type", "vdc-group"),
					resource.TestCheckResourceAttr(resourceNameVDCGroup, "owner_name", "MyVDCGroup"),
					resource.TestCheckResourceAttr(resourceNameVDCGroup, "tier0_vrf_name", "prvrf01eocb0006205allsp01"),
					resource.TestMatchResourceAttr(resourceNameVDCGroup, "name", regexp.MustCompile(`tn01e02ocb0006205spt[0-9]{3}`)),
					resource.TestCheckResourceAttr(resourceNameVDCGroup, "lb_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceNameVDCGroup, "description"),
				),
			},
		},
	})
}
