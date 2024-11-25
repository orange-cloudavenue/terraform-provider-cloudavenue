package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
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
	return
}

func (r *EdgeGatewayAppPortProfileDatasource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultConfig)
					return
				},
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
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultConfig)
					return
				},
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
		"example_provider_scope": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_provider_scope" {
						name = "BKP_TCP_bpcd"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.AppPortProfile)),
						resource.TestCheckResourceAttr(resourceName, "name", "BKP_TCP_bpcd"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.protocol", "TCP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.ports.0", "13782"),
					},
				},
			}
		},
		"example_system_scope": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_system_scope" {
						name = "HTTP"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.AppPortProfile)),
						resource.TestCheckResourceAttr(resourceName, "name", "HTTP"),
						resource.TestCheckResourceAttr(resourceName, "description", "HTTP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.protocol", "TCP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.ports.0", "80"),
					},
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
