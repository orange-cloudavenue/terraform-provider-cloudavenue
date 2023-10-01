// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCDataSource{}

const (
	VDCDataSourceName = testsacc.ResourceName("data.cloudavenue_vdc")
)

type VDCDataSource struct{}

func NewVDCDataSourceTest() testsacc.TestACC {
	return &VDCDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCDataSource) GetResourceName() string {
	return VDCDataSourceName.String()
}

func (r *VDCDataSource) DependenciesConfig() (configs testsacc.TFData) {
	configs.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig())
	return
}

func (r *VDCDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdc" "example" {
						name = cloudavenue_vdc.example.name
					}`,
					Checks: GetResourceConfig()[VDCResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVDCDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCDataSource{}),
	})
}
