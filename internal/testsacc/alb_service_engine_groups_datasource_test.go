package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ALBServiceEngineGroupsDataSource{}

const (
	ALBServiceEngineGroupsDataSourceName = testsacc.ResourceName("data.cloudavenue_alb_service_engine_groups")
)

type ALBServiceEngineGroupsDataSource struct{}

func NewALBServiceEngineGroupsDataSourceTest() testsacc.TestACC {
	return &ALBServiceEngineGroupsDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *ALBServiceEngineGroupsDataSource) GetResourceName() string {
	return ALBServiceEngineGroupsDataSourceName.String()
}

func (r *ALBServiceEngineGroupsDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_with_id"))
	return
}

func (r *ALBServiceEngineGroupsDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_alb_service_engine_groups" "example" {
						edge_gateway_name = data.cloudavenue_edgegateway.example_with_id.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.#"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.edge_gateway_id"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.deployed_virtual_services"),
					},
				},
			}
		},
		"example_with_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_alb_service_engine_groups" "example_with_id" {
						edge_gateway_id = data.cloudavenue_edgegateway.example_with_id.id
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.#"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.edge_gateway_id"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.deployed_virtual_services"),
					},
				},
			}
		},
	}
}

func TestAccALBServiceEngineGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ALBServiceEngineGroupsDataSource{}),
	})
}
