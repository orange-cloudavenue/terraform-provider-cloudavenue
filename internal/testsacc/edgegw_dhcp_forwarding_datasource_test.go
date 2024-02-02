package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &EdgeGatewayDhcpForwardingDataSource{}

const (
	EdgeGatewayDhcpForwardingDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_dhcp_forwarding")
)

type EdgeGatewayDhcpForwardingDataSource struct{}

func NewEdgeGatewayDhcpForwardingDataSourceTest() testsacc.TestACC {
	return &EdgeGatewayDhcpForwardingDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayDhcpForwardingDataSource) GetResourceName() string {
	return EdgeGatewayDhcpForwardingDataSourceName.String()
}

func (r *EdgeGatewayDhcpForwardingDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayDhcpForwardingResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayDhcpForwardingDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_dhcp_forwarding" "example" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
					}`,
					Checks: GetResourceConfig()[EdgeGatewayDhcpForwardingResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccEdgeGatewayDhcpForwardingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayDhcpForwardingDataSource{}),
	})
}
