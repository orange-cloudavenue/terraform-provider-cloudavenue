package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCNetworkIsolatedDataSource{}

const (
	VDCNetworkIsolatedDataSourceName = testsacc.ResourceName("data.cloudavenue_vdc_network_isolated")
)

type VDCNetworkIsolatedDataSource struct{}

func NewVDCNetworkIsolatedDataSourceTest() testsacc.TestACC {
	return &VDCNetworkIsolatedDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCNetworkIsolatedDataSource) GetResourceName() string {
	return VDCNetworkIsolatedDataSourceName.String()
}

func (r *VDCNetworkIsolatedDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCNetworkIsolatedResourceName]().GetDefaultConfig)
	return
}

func (r *VDCNetworkIsolatedDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdc_network_isolated" "example" {
						vdc = cloudavenue_vdc.example.name
						name = cloudavenue_vdc_network_isolated.example.name
					}`,
					Checks: GetResourceConfig()[VDCNetworkIsolatedResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVDCNetworkIsolatedDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCNetworkIsolatedDataSource{}),
	})
}
