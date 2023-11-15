package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &CatalogDataSource{}

const (
	CatalogDataSourceName = testsacc.ResourceName("data.cloudavenue_catalog")
)

type CatalogDataSource struct{}

func NewCatalogDataSourceTest() testsacc.TestACC {
	return &CatalogDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *CatalogDataSource) GetResourceName() string {
	return CatalogDataSourceName.String()
}

func (r *CatalogDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[CatalogResourceName]().GetDefaultConfig)
	return
}

func (r *CatalogDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_catalog" "example" {
						name = cloudavenue_catalog.example.name
					}`,
					Checks: GetResourceConfig()[CatalogResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccCatalogDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&CatalogDataSource{}),
	})
}
