// Package tier0 provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &Tier0VRFDataSource{}

const (
	Tier0VRFDataSourceName = testsacc.ResourceName("data.cloudavenue_tier0_vrf")
)

type Tier0VRFDataSource struct{}

func NewTier0VRFDataSourceTest() testsacc.TestACC {
	return &Tier0VRFDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *Tier0VRFDataSource) GetResourceName() string {
	return Tier0VRFDataSourceName.String()
}

func (r *Tier0VRFDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *Tier0VRFDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_tier0_vrf" "example" {
						name = "prvrf01eocb0006205allsp01"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", "prvrf01eocb0006205allsp01"),
						resource.TestCheckResourceAttr(resourceName, "class_service", "VRF_STANDARD"),
						resource.TestCheckResourceAttr(resourceName, "tier0_provider", "pr01e02t0sp16"),
						resource.TestCheckResourceAttr(resourceName, "services.#", "3"),
						resource.TestCheckResourceAttr(resourceName, "services.0.service", "OBJECT_STORAGE"),
						resource.TestCheckResourceAttr(resourceName, "services.0.vlan_id", ""),
						resource.TestCheckResourceAttr(resourceName, "services.1.service", "INTERNET"),
						resource.TestCheckResourceAttr(resourceName, "services.1.vlan_id", ""),
						resource.TestCheckResourceAttr(resourceName, "services.2.service", "ADMIN"),
						resource.TestCheckResourceAttr(resourceName, "services.2.vlan_id", ""),
					},
				},
			}
		},
	}
}

func TestAccTier0VrfDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&Tier0VRFDataSource{}),
	})
}
