package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGIPSetDataSource{}

const (
	VDCGIPSetDataSourceName = testsacc.ResourceName("data.cloudavenue_vdcg_ip_set")
)

type VDCGIPSetDataSource struct{}

func NewVDCGIPSetDataSourceTest() testsacc.TestACC {
	return &VDCGIPSetDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGIPSetDataSource) GetResourceName() string {
	return VDCGIPSetDataSourceName.String()
}

func (r *VDCGIPSetDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGIPSetResourceName]().GetDefaultConfig)
	return
}

func (r *VDCGIPSetDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_ip_set" "example" {
						name = cloudavenue_vdcg_ip_set.example.name
						vdc_group_name = cloudavenue_vdcg.example.name
					}`,
					Checks: GetResourceConfig()[VDCGIPSetResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVDCGIPSetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGIPSetDataSource{}),
	})
}
