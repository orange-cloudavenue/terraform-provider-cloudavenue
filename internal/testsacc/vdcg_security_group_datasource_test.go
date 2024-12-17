package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGSecurityGroupDataSource{}

const (
	VDCGSecurityGroupDataSourceName = testsacc.ResourceName("data.cloudavenue_vdcg_security_group")
)

type VDCGSecurityGroupDataSource struct{}

func NewVDCGSecurityGroupDataSourceTest() testsacc.TestACC {
	return &VDCGSecurityGroupDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGSecurityGroupDataSource) GetResourceName() string {
	return VDCGSecurityGroupDataSourceName.String()
}

func (r *VDCGSecurityGroupDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGSecurityGroupResourceName]().GetDefaultConfig)
	return
}

func (r *VDCGSecurityGroupDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_security_group" "example" {
						vdc_group_name = cloudavenue_vdcg.example.name
						name = cloudavenue_vdcg_security_group.example.name
					}`,
					Checks: GetResourceConfig()[VDCGSecurityGroupResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVDCGSecurityGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGSecurityGroupDataSource{}),
	})
}
