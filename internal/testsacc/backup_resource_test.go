package testsacc

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &BackupResource{}

const (
	BackupResourceName = testsacc.ResourceName("cloudavenue_backup")
)

type BackupResource struct{}

func NewBackupResourceTest() testsacc.TestACC {
	return &BackupResource{}
}

// GetResourceName returns the name of the resource.
func (r *BackupResource) GetResourceName() string {
	return BackupResourceName.String()
}

func (r *BackupResource) DependenciesConfig() (configs testsacc.TFData) {
	// TODO : Add dependencies config
	configs.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig())

	// This is method for add dependencies legacy config
	// configs.Append(AddConstantConfig(constantName))
	return
}

func (r *BackupResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// TODO : Complete tests
		// * First For a VDC Backup named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "target_name", "example"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_backup" "example" {
						type = "vdc"
						target_name = cloudavenue_vdc.example.name
						policies = [{
								policy_name = "D6"
							}]
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "type", "vdc"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.policy_name", "D6"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_backup" "example" {
							type = "vdc"
							target_name = cloudavenue_vdc.example.name
							policies = [{
									policy_name = "D6"
								},
								{
									policy_name = "D30"
								}]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.policy_name", "D6"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.policy_name", "D30"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDFunc: testAccBackupResourceImportStateIDFuncWithTypeAndTargetName(resourceName),
						ImportState:       true,
						ImportStateVerify: true,
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

func testAccBackupResourceImportStateIDFuncWithTypeAndTargetName(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		// Type.Target_name
		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["type"], rs.Primary.Attributes["target_name"]), nil
	}
}

func TestAccBackupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&BackupResource{}),
	})
}
