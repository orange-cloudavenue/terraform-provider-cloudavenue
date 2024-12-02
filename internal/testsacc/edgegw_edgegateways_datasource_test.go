// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &EdgeGatewaysDataSource{}

const (
	EdgeGatewaysDataSourceName = testsacc.ResourceName("data.cloudavenue_edgegateways")
)

type EdgeGatewaysDataSource struct{}

func NewEdgeGatewaysDataSourceTest() testsacc.TestACC {
	return &EdgeGatewaysDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewaysDataSource) GetResourceName() string {
	return EdgeGatewaysDataSourceName.String()
}

func (r *EdgeGatewaysDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewaysDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_edgegateways" "example" {}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateways.0.id", urn.TestIsType(urn.Gateway)),
					},
				},
			}
		},
	}
}

func TestAccEdgeGatewaysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewaysDataSource{}),
	})
}
