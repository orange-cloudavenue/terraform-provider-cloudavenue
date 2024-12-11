package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ALBPoolDataSource{}

const (
	ALBPoolDataSourceName = testsacc.ResourceName("data.cloudavenue_alb_pool")
)

type ALBPoolDataSource struct{}

func NewALBPoolDataSourceTest() testsacc.TestACC {
	return &ALBPoolDataSource{}
}

func (r *ALBPoolDataSource) GetResourceName() string {
	return ALBPoolDataSourceName.String()
}

func (r *ALBPoolDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[ALBPoolResourceName]().GetDefaultConfig)
	resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_with_id"))
	return
}

func (r *ALBPoolDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_alb_pool" "example" {
						edge_gateway_name = data.cloudavenue_edgegateway.example_with_id.name
						name              = cloudavenue_alb_pool.example.name
					}`,
					Checks: GetResourceConfig()[ALBPoolResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccALBPoolDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ALBPoolDataSource{}),
	})
}
