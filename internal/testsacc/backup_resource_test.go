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
	// configs.Append(GetResourceConfig()[VAPPResourceName]().GetDefaultConfig())
	// configs.Append(GetResourceConfig()[VMResourceName]().GetDefaultConfig())
	return
}

func (r *BackupResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// TODO : Complete tests
		// * First Test For a VDC Backup named "example"
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
		// * Second Test For a VAPP Backup named "example"
		"example2": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "target_name", "example"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_backup" "example" {
						type = "vapp"
						target_name = cloudavenue_vapp.example.name
						policies = [{
								policy_name = "D6"
							}]
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "type", "vapp"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.policy_name", "D6"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_backup" "example" {
							type = "vapp"
							target_name = cloudavenue_vapp.example.name
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
		// * Second Test For a VM Backup named "example"
		"example3": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "target_name", "example"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_backup" "example" {
						type = "vm"
						target_name = cloudavenue_vm.example.name
						policies = [{
								policy_name = "D6"
							}]
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "type", "vm"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.policy_name", "D6"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_backup" "example" {
							type = "vm"
							target_name = cloudavenue_vm.example.name
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
