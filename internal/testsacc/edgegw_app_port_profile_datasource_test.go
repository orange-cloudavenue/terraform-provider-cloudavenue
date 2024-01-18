package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &EdgeGatewayAppPortProfileDatasource{}

const (
	EdgeGatewayAppPortProfileDatasourceName = testsacc.ResourceName("data.cloudavenue_edgegateway_app_port_profile")
)

type EdgeGatewayAppPortProfileDatasource struct{}

func NewEdgeGatewayAppPortProfileDatasourceTest() testsacc.TestACC {
	return &EdgeGatewayAppPortProfileDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayAppPortProfileDatasource) GetResourceName() string {
	return EdgeGatewayAppPortProfileDatasourceName.String()
}

func (r *EdgeGatewayAppPortProfileDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayAppPortProfileDatasource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example" {
						name = cloudavenue_edgegateway_app_port_profile.example.name
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultChecks(),
				},
			}
		},
		"example_by_id": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_by_id" {
						id = cloudavenue_edgegateway_app_port_profile.example.id
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccEdgeGatewayAppPortProfileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayAppPortProfileDatasource{}),
	})
}
