package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ALBServiceEngineGroupDataSource{}

const (
	ALBServiceEngineGroupDataSourceName = testsacc.ResourceName("data.cloudavenue_alb_service_engine_group")
)

type ALBServiceEngineGroupDataSource struct{}

func NewALBServiceEngineGroupDataSourceTest() testsacc.TestACC {
	return &ALBServiceEngineGroupDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *ALBServiceEngineGroupDataSource) GetResourceName() string {
	return ALBServiceEngineGroupDataSourceName.String()
}

func (r *ALBServiceEngineGroupDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *ALBServiceEngineGroupDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_alb_service_engine_group" "example" {
						name = "v010w02eprnxcdshrdsegp04"
						edge_gateway_name = "tn01e02ocb0006205spt101"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.ServiceEngineGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "deployed_virtual_services"),
					},
				},
			}
		},
		"example_with_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_alb_service_engine_group" "example_with_id" {
						id = "urn:vcloud:serviceEngineGroup:737b9768-95a0-4955-bbbe-d5eab846e8dc"
						edge_gateway_name = "tn01e02ocb0006205spt101"
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.ServiceEngineGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "deployed_virtual_services"),
					},
				},
			}
		},
		"example_with_edge_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_alb_service_engine_group" "example_with_edge_id" {
						id = "urn:vcloud:serviceEngineGroup:737b9768-95a0-4955-bbbe-d5eab846e8dc"
						edge_gateway_id = "urn:vcloud:gateway:d3c42a20-96b9-4452-91dd-f71b71dfe314"
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.ServiceEngineGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "deployed_virtual_services"),
					},
				},
			}
		},
		"example_with_name_and_edge_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_alb_service_engine_group" "example_with_name_and_edge_id" {
						name = "v010w02eprnxcdshrdsegp04"
						edge_gateway_id = "urn:vcloud:gateway:d3c42a20-96b9-4452-91dd-f71b71dfe314"
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.ServiceEngineGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "deployed_virtual_services"),
					},
				},
			}
		},
	}
}

func TestAccALBServiceEngineGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ALBServiceEngineGroupDataSource{}),
	})
}
