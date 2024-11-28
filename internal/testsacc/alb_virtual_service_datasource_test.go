package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ALBVirtualServiceDataSource{}

const (
	ALBVirtualServiceDataSourceName = testsacc.ResourceName("data.cloudavenue_alb_virtual_service")
)

type ALBVirtualServiceDataSource struct{}

func NewALBVirtualServiceDataSourceTest() testsacc.TestACC {
	return &ALBVirtualServiceDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *ALBVirtualServiceDataSource) GetResourceName() string {
	return ALBVirtualServiceDataSourceName.String()
}

func (r *ALBVirtualServiceDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[ALBVirtualServiceResourceName]().GetDefaultConfig)
	resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_with_id"))
	return
}

func (r *ALBVirtualServiceDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_alb_virtual_service" "example" {
						edge_gateway_name = data.cloudavenue_edgegateway.example_with_id.name
						name              = cloudavenue_alb_virtual_service.example.name
					}`,
					Checks: GetResourceConfig()[ALBVirtualServiceResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccALBVirtualServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ALBVirtualServiceDataSource{}),
	})
}
