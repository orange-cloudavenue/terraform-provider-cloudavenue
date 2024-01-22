// Package tier0 provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &Tier0VRFsDataSource{}

const (
	Tier0VRFsDataSourceName = testsacc.ResourceName("data.cloudavenue_tier0_vrfs")
)

type Tier0VRFsDataSource struct{}

func NewTier0VRFsDataSourceTest() testsacc.TestACC {
	return &Tier0VRFsDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *Tier0VRFsDataSource) GetResourceName() string {
	return Tier0VRFsDataSourceName.String()
}

func (r *Tier0VRFsDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *Tier0VRFsDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `data "cloudavenue_tier0_vrfs" "example" {}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "names.#"),
					},
				},
			}
		},
	}
}

func TestAccTier0VrfsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&Tier0VRFsDataSource{}),
	})
}
