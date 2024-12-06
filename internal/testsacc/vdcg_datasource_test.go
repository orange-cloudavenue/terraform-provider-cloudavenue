package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGDataSource{}

const (
	VDCGDataSourceName = testsacc.ResourceName("data.cloudavenue_vdcg")
)

type VDCGDataSource struct{}

func NewVDCGDataSourceTest() testsacc.TestACC {
	return &VDCGDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGDataSource) GetResourceName() string {
	return VDCGDataSourceName.String()
}

func (r *VDCGDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
	return
}

func (r *VDCGDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg" "example" {
						name = cloudavenue_vdcg.example.name
					}`,
					Checks: GetResourceConfig()[VDCGResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVDCGDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGDataSource{}),
	})
}
