// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &CatalogsDataSource{}

const (
	CatalogsDataSourceName = testsacc.ResourceName("data.cloudavenue_catalogs")
)

type CatalogsDataSource struct{}

func NewCatalogsDataSourceTest() testsacc.TestACC {
	return &CatalogsDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *CatalogsDataSource) GetResourceName() string {
	return CatalogsDataSourceName.String()
}

func (r *CatalogsDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *CatalogsDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `data "cloudavenue_catalogs" "example" {}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs_name.#"),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs.%"),

						resource.TestCheckResourceAttrWith(resourceName, "catalogs.Orange-Linux.id", uuid.TestIsType(uuid.Catalog)),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs.Orange-Linux.name"),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs.Orange-Linux.created_at"),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs.Orange-Linux.preserve_identity_information"),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs.Orange-Linux.number_of_media"),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs.Orange-Linux.media_item_list.#"),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs.Orange-Linux.is_shared"),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs.Orange-Linux.is_published"),
						resource.TestCheckResourceAttrSet(resourceName, "catalogs.Orange-Linux.is_local"),

						resource.TestCheckNoResourceAttr(resourceName, "catalogs.Orange-Linux.owner_name"),  // In Orange-Linux catalog, owner_name is empty
						resource.TestCheckNoResourceAttr(resourceName, "catalogs.Orange-Linux.description"), // In Orange-Linux catalog, description is empty
						resource.TestCheckNoResourceAttr(resourceName, "catalogs.Orange-Linux.is_cached"),   // In Orange-Linux catalog, is_cached is false
					},
				},
			}
		},
	}
}

func TestACCCatalogsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&CatalogsDataSource{}),
	})
}
