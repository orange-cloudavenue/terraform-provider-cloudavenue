// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccEdgeGatewayResourceConfig = `
data "cloudavenue_tier0_vrfs" "example_with_vdc" {}

resource "cloudavenue_edgegateway" "example_with_vdc" {
  owner_name     = "MyEdgeGateway"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example_with_vdc.names.0
  owner_type     = "vdc"
  lb_enabled     = false
}
`

const testAccEdgeGatewayGroupResourceConfig = `
data "cloudavenue_tier0_vrfs" "example_with_group" {}

resource "cloudavenue_edgegateway" "example_with_group" {
  owner_name     = "MyEdgeGatewayGroup"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example_with_group.names.0
  owner_type     = "vdc-group"
}
`

var _ testsacc.TestACC = &EdgeGatewayResource{}

const (
	EdgeGatewayResourceName = testsacc.ResourceName("cloudavenue_edgegateway")
)

type EdgeGatewayResource struct{}

func NewEdgeGatewayResourceTest() testsacc.TestACC {
	return &EdgeGatewayResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayResource) GetResourceName() string {
	return EdgeGatewayResourceName.String()
}

func (r *EdgeGatewayResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
	resp.Append(GetDataSourceConfig()[Tier0VRFDataSourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "owner_name"),
					resource.TestCheckResourceAttr(resourceName, "owner_type", "vdc"),

					// Read-Only attributes
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway" "example" {
						owner_name     = cloudavenue_vdc.example.name
						tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
						owner_type     = "vdc"
					  }`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "bandwidth"),
						resource.TestCheckResourceAttr(resourceName, "lb_enabled", "false"), // Deprecated attribute
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// Test one of range value allowed in bandwidth attribute
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway" "example" {
							owner_name     = cloudavenue_vdc.example.name
							tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
							owner_type     = "vdc"
							bandwidth      = 20
						  }`),
						TFAdvanced: testsacc.TFAdvanced{
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
							ExpectError:        regexp.MustCompile(`Invalid Bandwidth value`),
						},
					},
					// Test overcommit bandwidth
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway" "example" {
							owner_name     = cloudavenue_vdc.example.name
							tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
							owner_type     = "vdc"
							bandwidth      = 300
						  }`),
						TFAdvanced: testsacc.TFAdvanced{
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
							ExpectError:        regexp.MustCompile(`Overcommitting bandwidth`),
						},
					},
					// Update bandwidth
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
											resource "cloudavenue_edgegateway" "example" {
												owner_name     = cloudavenue_vdc.example.name
												tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
												owner_type     = "vdc"
												bandwidth      = 25
											  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "bandwidth", "25"),
							resource.TestCheckResourceAttr(resourceName, "lb_enabled", "false"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"name"},
						ImportState:          true,
					},
				},
			}
		},
		// TODO After the implementation of the VDC_GROUP resource we can use new resourceConfig
		// "example_with_group": func(_ context.Context, resourceName string) testsacc.Test {
		// 	return testsacc.Test{
		// 		CommonChecks: []resource.TestCheckFunc{
		// 			resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
		// 			resource.TestCheckResourceAttrSet(resourceName, "owner_name"),
		// 			resource.TestCheckResourceAttr(resourceName, "owner_type", "vdc"),

		// 			// Read-Only attributes
		// 			resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`tn01e02ocb0006205spt[0-9]{3}`)),
		// 			resource.TestCheckResourceAttrSet(resourceName, "description"),
		// 		},
		// 		// ! Create testing
		// 		Create: testsacc.TFConfig{
		// 			TFConfig: testsacc.GenerateFromTemplate(resourceName, `
		// 			resource "cloudavenue_edgegateway" "example_with_vdc" {
		// 				owner_name     = cloudavenue_vdc.example.name
		// 				tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
		// 				owner_type     = "vdc-group"
		// 				lb_enabled     = false
		// 			  }`),
		// 			Checks: []resource.TestCheckFunc{
		// 				resource.TestCheckResourceAttr(resourceName, "lb_enabled", "false"),
		// 			},
		// 		},
		// 		// ! Updates testing
		// 		Updates: []testsacc.TFConfig{
		// 			// Update lb_enabled
		// 			{
		// 				TFConfig: testsacc.GenerateFromTemplate(resourceName, `
		// 				resource "cloudavenue_edgegateway" "example_with_vdc" {
		// 					owner_name     = cloudavenue_vdc.example.name
		// 					tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
		// 					owner_type     = "vdc"
		// 					lb_enabled     = true
		// 				  }`),
		// 				Checks: []resource.TestCheckFunc{
		// 					resource.TestCheckResourceAttr(resourceName, "lb_enabled", "false"),
		// 				},
		// 			},
		// 		},
		// 		// ! Imports testing
		// 		Imports: []testsacc.TFImport{
		// 			{
		// 				ImportStateIDBuilder: []string{"name"},
		// 				ImportState:          true,
		// 				ImportStateVerify:    true,
		// 			},
		// 		},
		// 	}
		// },
	}
}

func TestAccEdgeGatewayResource(t *testing.T) {
	edgegw.ConfigEdgeGateway = func() edgegw.EdgeGatewayConfig {
		return edgegw.EdgeGatewayConfig{
			CheckJobDelay: 10 * time.Second,
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayResource{}),
	})
}
