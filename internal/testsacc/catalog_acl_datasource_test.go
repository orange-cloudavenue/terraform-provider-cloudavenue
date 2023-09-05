package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &CatalogACLDataSource{}

const (
	CatalogACLDataSourceName = ResourceName("data.cloudavenue_catalog_acl")
)

type CatalogACLDataSource struct{}

func NewCatalogACLDataSourceTest() testsacc.TestACC {
	return &CatalogACLDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *CatalogACLDataSource) GetResourceName() string {
	return CatalogACLDataSourceName.String()
}

func (r *CatalogACLDataSource) DependenciesConfig() (configs testsacc.TFData) {
	configs.Append(GetResourceConfig()[CatalogACLResourceName]().GetDefaultConfig())
	return
}

func (r *CatalogACLDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_catalog_acl" "example" {
						catalog_id = cloudavenue_catalog_acl.example.catalog_id
					}`,
					Checks: NewCatalogACLResourceTest().Tests(ctx)["example"](ctx, resourceName).GenerateCheckWithCommonChecks(),
				},
			}
		},
	}
}

func TestCatalogAccACLDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&CatalogACLDataSource{}),
	})
}
