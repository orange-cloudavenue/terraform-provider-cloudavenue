// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &CatalogResource{}

const (
	CatalogResourceName = ResourceName("cloudavenue_catalog")
)

type CatalogResource struct{}

func NewCatalogResourceTest() testsacc.TestACC {
	return &CatalogResource{}
}

// GetResourceName returns the name of the resource.
func (r *CatalogResource) GetResourceName() string {
	return CatalogResourceName.String()
}

func (r *CatalogResource) DependenciesConfig() (configs testsacc.TFData) {
	// configs.Append(GetResourceConfig()[EdgeGatewatResourceName]().GetDefaultConfig())
	return
}

func (r *CatalogResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Catalog)),
					resource.TestCheckResourceAttr(resourceName, "name", "example"),
					resource.TestCheckResourceAttr(resourceName, "delete_recursive", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_force", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_catalog" "example" {
						name             = "example"
						description      = "catalog for files"
						delete_recursive = true
						delete_force     = true
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "description", "catalog for files"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_catalog" "example" {
							name             = "example"
							description      = "updated catalog for files"
							delete_recursive = true
							delete_force     = true
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", "updated catalog for files"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateID:           "example",
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"delete_force", "delete_recursive"},
					},
				},
			}
		},
		// * Test Two
		// "test_two": func(_ context.Context) testsacc.Test {
		// 	return testsacc.Test{
		// 		Create: testsacc.TFConfig{
		// 			TFConfig: ``,
		// 		},
		// 		Updates: []testsacc.TFConfig{
		// 			{
		// 				TFConfig: ``,
		// 			},
		// 			{
		// 				TFConfig: ``,
		// 			},
		// 		},
		// 	}
		// },
	}
}

func TestAccCatalogResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&CatalogResource{}),
	})
}
