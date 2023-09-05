package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &CatalogACLResource{}

const (
	CatalogACLResourceName = ResourceName("cloudavenue_catalog_acl")
)

type CatalogACLResource struct{}

func NewCatalogACLResourceTest() testsacc.TestACC {
	return &CatalogACLResource{}
}

// GetResourceName returns the name of the resource.
func (r *CatalogACLResource) GetResourceName() string {
	return CatalogACLResourceName.String()
}

func (r *CatalogACLResource) DependenciesConfig() (configs testsacc.TFData) {
	configs.Append(GetResourceConfig()[CatalogResourceName]().GetDefaultConfig())

	// TODO After the implementation of the user resource we can use new resourceConfig
	configs.Append(AddConstantConfig(testAccOrgUserResourceConfigForCatalogACL))
	return
}

func (r *CatalogACLResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "catalog_name"),

					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Catalog)),
					resource.TestCheckResourceAttrWith(resourceName, "catalog_id", uuid.TestIsType(uuid.Catalog)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_catalog_acl" "example" {
						catalog_id = cloudavenue_catalog.example.id
						shared_with_everyone = true
						everyone_access_level = "ReadOnly"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckNoResourceAttr(resourceName, "shared_with_users.#"),

						resource.TestCheckResourceAttr(resourceName, "everyone_access_level", "ReadOnly"),
						resource.TestCheckResourceAttr(resourceName, "shared_with_everyone", "true"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_catalog_acl" "example" {
							catalog_id = cloudavenue_catalog.example.id
							shared_with_everyone = true
							everyone_access_level = "FullControl"
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckNoResourceAttr(resourceName, "shared_with_users.#"),

							resource.TestCheckResourceAttr(resourceName, "everyone_access_level", "FullControl"),
							resource.TestCheckResourceAttr(resourceName, "shared_with_everyone", "true"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_catalog_acl" "example" {
							catalog_id = cloudavenue_catalog.example.id
							shared_with_everyone = false
							shared_with_users = [
								{
									user_id = cloudavenue_iam_user.example.id
									access_level = "ReadOnly"
								},
								{
									user_id = cloudavenue_iam_user.example2.id
									access_level = "FullControl"
								}
							]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckNoResourceAttr(resourceName, "everyone_access_level"),

							resource.TestCheckResourceAttr(resourceName, "shared_with_everyone", "false"),
							resource.TestCheckResourceAttr(resourceName, "shared_with_users.#", "2"),

							resource.TestCheckResourceAttrWith(resourceName, "shared_with_users.0.user_id", uuid.TestIsType(uuid.User)),
							resource.TestCheckResourceAttrWith(resourceName, "shared_with_users.1.user_id", uuid.TestIsType(uuid.User)),
							// shared_with_users it's a SetNestedAttribute, so we can't be sure of the order of the elements in the list is not possible to test each attribute
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccCatalogACLResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&CatalogACLResource{}),
	})
}
