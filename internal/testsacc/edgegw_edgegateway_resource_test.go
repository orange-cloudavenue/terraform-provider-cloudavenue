// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
)

const testAccEdgeGatewayResourceConfig = `
data "cloudavenue_tier0_vrfs" "example_with_vdc" {}

resource "cloudavenue_edgegateway" "example_with_vdc" {
  owner_name     = "MyEdgeGateway"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example_with_vdc.names.0
  owner_type     = "vdc"
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
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
					resource.TestMatchResourceAttr(resourceName, "tier0_vrf_name", regexp.MustCompile(regexpTier0VRFName)),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway" "example" {
						owner_name     = cloudavenue_vdc.example.name
						tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
						owner_type     = "vdc"
						bandwidth      = 25
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "bandwidth"),
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
							bandwidth      = 5
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "bandwidth", "5"),
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
				Destroy: true,
			}
		},
		"example_with_vdc_group": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "owner_name"),
					resource.TestCheckResourceAttr(resourceName, "owner_type", "vdc-group"),

					// Read-Only attributes
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
					resource.TestMatchResourceAttr(resourceName, "tier0_vrf_name", regexp.MustCompile(regexpTier0VRFName)),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGroupResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway" "example_with_vdc_group" {
						owner_name     = cloudavenue_vdc_group.example.name
						tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
						owner_type     = "vdc-group"
						bandwidth      = 25
					  }`),
					Checks: []resource.TestCheckFunc{},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway" "example_with_vdc_group" {
							owner_name     = cloudavenue_vdc_group.example.name
							tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
							owner_type     = "vdc-group"
							bandwidth      = 5
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "bandwidth", "5"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
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
