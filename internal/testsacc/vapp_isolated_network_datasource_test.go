package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VAppIsolatedNetworkDataSource{}

const (
	VAppIsolatedNetworkDataSourceName = testsacc.ResourceName("data.cloudavenue_vapp_isolated_network")
)

type VAppIsolatedNetworkDataSource struct{}

func NewVAppIsolatedNetworkDataSourceTest() testsacc.TestACC {
	return &VAppIsolatedNetworkDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VAppIsolatedNetworkDataSource) GetResourceName() string {
	return VAppIsolatedNetworkDataSourceName.String()
}

func (r *VAppIsolatedNetworkDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[VAppIsolatedNetworkResourceName]().GetDefaultConfig)
	return
}

func (r *VAppIsolatedNetworkDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (with vapp_name)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vapp_isolated_network" "example" {
						vdc = cloudavenue_vdc.example.name
						vapp_name = cloudavenue_vapp.example.name
						name = cloudavenue_vapp_isolated_network.example.name
					}`,
					Checks: GetResourceConfig()[VAppIsolatedNetworkResourceName]().GetDefaultChecks(),
				},
			}
		},
		// * Test Two (with vapp_id)
		"example_2": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vapp_isolated_network" "example" {
						vdc = cloudavenue_vdc.example.name
						vapp_id = cloudavenue_vapp.example.id
						name = cloudavenue_vapp_isolated_network.example.name
					}`,
					Checks: GetResourceConfig()[VAppIsolatedNetworkResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVAppIsolatedNetworkDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VAppIsolatedNetworkDataSource{}),
	})
}
