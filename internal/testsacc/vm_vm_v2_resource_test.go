package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &vm_V2Resource{}

const (
	vm_V2ResourceName = testsacc.ResourceName("cloudavenue_vm_v_2")
)

type vm_V2Resource struct{}

func Newvm_V2ResourceTest() testsacc.TestACC {
	return &vm_V2Resource{}
}

// GetResourceName returns the name of the resource.
func (r *vm_V2Resource) GetResourceName() string {
	return vm_V2ResourceName.String()
}

func (r *vm_V2Resource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// TODO : Add dependencies config
	// resp.Append(GetResourceConfig()[CatalogResourceName]().GetDefaultConfig)

	// This is method for add dependencies legacy config
	// resp.Append(AddConstantConfig(constantName))
	return
}

func (r *vm_V2Resource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// TODO : Complete tests
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Catalog)), // TODO : Change type
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vm_v_2" "example" {
						foo = "bar"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "foo", "bar"),

					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vm_v_2" "example" {
							foo = "barUpdated"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "foo", "barUpdated"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vm_v_2" "example" {
							foo = "barUpdated"
							bar = "foo"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "foo", "barUpdated"),
							resource.TestCheckResourceAttr(resourceName, "bar", "foo"),
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
		// It's possible to add multiple tests

		// Complete and functional example :
		/*
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
		*/
	}
}

func TestAccvm_V2Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&vm_V2Resource{}),
	})
}
