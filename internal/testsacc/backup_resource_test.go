package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &BackupResource{}

const (
	BackupResourceName      = testsacc.ResourceName("cloudavenue_backup")
	TestAccVMResourceConfig = `
	data "cloudavenue_catalog_vapp_template" "example" {
	  catalog_name  = "Orange-Linux"
	  template_name = "debian_10_X64"
	}
	resource "cloudavenue_vm" "example_backup" {
	  name        = "example"
	  description = "This is a example vm"
	  vapp_name = cloudavenue_vapp.example.name
	  vdc = cloudavenue_vdc.example.name
	  deploy_os = {
	    vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
	  }
	  settings = {
	  	customization = {
	  	  auto_generate_password = true
	  	}
	  }
	  resource = {}
	  state = {}
	}`
)

type BackupResource struct{}

func NewBackupResourceTest() testsacc.TestACC {
	return &BackupResource{}
}

// GetResourceName returns the name of the resource.
func (r *BackupResource) GetResourceName() string {
	return BackupResourceName.String()
}

func (r *BackupResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VAppResourceName]().GetDefaultConfig)
	resp.Append(AddConstantConfig(TestAccVMResourceConfig))
	return
}

func (r *BackupResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First Test For a VDC Backup named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_name"),
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
						    },{
						      policy_name = "D30"
						    }
						  ]
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
						ImportStateIDBuilder:    []string{"type", "target_name"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"target_id"},
					},
				},
			}
		},
		// * Test For a VAPP Backup named "example"
		"example_vapp": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_backup" "example_vapp" {
						type = "vapp"
						target_id = cloudavenue_vapp.example.id
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
						resource "cloudavenue_backup" "example_vapp" {
							type = "vapp"
							target_id = cloudavenue_vapp.example.id
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
						ImportStateIDBuilder:    []string{"type", "target_name"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"target_id"},
					},
				},
			}
		},
		// * Test For a VM Backup named "example"
		"example_vm": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_name"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_backup" "example_vm" {
					  type = "vm"
					  target_name = cloudavenue_vm.example_backup.name
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
						resource "cloudavenue_backup" "example_vm" {
						  type = "vm"
						  target_name = cloudavenue_vm.example_backup.name
						  policies = [{
						      policy_name = "D6"
						    },{
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
						ImportStateIDBuilder:    []string{"type", "target_name"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"target_id"},
					},
				},
			}
		},
	}
}

func TestAccBackupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&BackupResource{}),
	})
}
