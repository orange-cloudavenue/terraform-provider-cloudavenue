package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &vdcg_SecurityGroupResource{}

const (
	vdcg_SecurityGroupResourceName = testsacc.ResourceName("cloudavenue_vdcg_security_group")
)

type vdcg_SecurityGroupResource struct{}

func Newvdcg_SecurityGroupResourceTest() testsacc.TestACC {
	return &vdcg_SecurityGroupResource{}
}

// GetResourceName returns the name of the resource.
func (r *vdcg_SecurityGroupResource) GetResourceName() string {
	return vdcg_SecurityGroupResourceName.String()
}

func (r *vdcg_SecurityGroupResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[NetworkRoutedResourceName]().GetDefaultConfig)
	return
}

func (r *vdcg_SecurityGroupResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.SecurityGroup)),
					resource.TestCheckResourceAttr(resourceName, "member_org_network_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_security_group" "example" {
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
						resource "cloudavenue_vdcg_security_group" "example" {
							foo = "barUpdated"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "foo", "barUpdated"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_security_group" "example" {
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

func TestAccvdcg_SecurityGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&vdcg_SecurityGroupResource{}),
	})
}
