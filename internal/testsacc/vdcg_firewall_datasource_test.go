package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGFirewallDataSource{}

const (
	VDCGFirewallDataSourceName = testsacc.ResourceName("data.cloudavenue_vdcg_firewall")
)

type VDCGFirewallDataSource struct{}

func NewVDCGFirewallDataSourceTest() testsacc.TestACC {
	return &VDCGFirewallDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGFirewallDataSource) GetResourceName() string {
	return VDCGFirewallDataSourceName.String()
}

func (r *VDCGFirewallDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGFirewallResourceName]().GetDefaultConfig)
	return
}

func (r *VDCGFirewallDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_firewall" "example" {
						vdc_group_name = cloudavenue_vdcg_firewall.example.vdc_group_name
					}`,
					Checks: GetResourceConfig()[VDCGFirewallResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVDCGFirewallDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGFirewallDataSource{}),
	})
}
