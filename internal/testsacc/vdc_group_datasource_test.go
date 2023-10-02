package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGroupDataSource{}

const (
	VDCGroupDataSourceName = testsacc.ResourceName("data.cloudavenue_vdc_group")
)

type VDCGroupDataSource struct{}

func NewVDCGroupDataSourceTest() testsacc.TestACC {
	return &VDCGroupDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGroupDataSource) GetResourceName() string {
	return VDCGroupDataSourceName.String()
}

func (r *VDCGroupDataSource) DependenciesConfig() (configs testsacc.TFData) {
	configs.Append(GetResourceConfig()[VDCGroupResourceName]().GetDefaultConfig())
	return
}

func (r *VDCGroupDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdc_group" "example" {
						name = cloudavenue_vdc_group.example.name
					}`,
					Checks: GetResourceConfig()[VDCResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVDCGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGroupDataSource{}),
	})
}
