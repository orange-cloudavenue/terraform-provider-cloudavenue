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
						edge_gateway_name = cloudavenue_edgegateway_app_port_profile.example.edge_gateway_id
						name = cloudavenue_edgegateway_app_port_profile.example.name
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultChecks(),
				},
				Destroy: true,
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
						edge_gateway_id = cloudavenue_edgegateway_app_port_profile.example.edge_gateway_id
						id = cloudavenue_edgegateway_app_port_profile.example.id
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[EdgeGatewayAppPortProfileResourceName]().GetDefaultChecks(),
				},
				Destroy: true,
			}
		},
		"example_provider_scope": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_provider_scope" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
						name = "BKP_TCP_bpcd"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.AppPortProfile)),
						resource.TestCheckResourceAttr(resourceName, "name", "BKP_TCP_bpcd"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.protocol", "TCP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.ports.0", "13782"),
					},
				},
				Destroy: true,
			}
		},
		"example_system_scope": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_system_scope" {
						edge_gateway_id = cloudavenue_edgegateway.example.id
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
				Destroy: true,
			}
		},
		"example_with_vdc_group": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetSpecificConfig("example_with_vdc_group"))
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateway_app_port_profile" "example_with_vdc_group" {
						edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
						name = "Heartbeat"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.AppPortProfile)),
						resource.TestCheckResourceAttr(resourceName, "name", "Heartbeat"),
						resource.TestCheckResourceAttr(resourceName, "description", "Heartbeat"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.#", "2"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "app_ports.*", map[string]string{
							"protocol": "TCP",
							"ports.0":  "57348",
						}),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "app_ports.*", map[string]string{
							"protocol": "TCP",
							"ports.0":  "52267",
						}),
					},
				},
				Destroy: true,
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
