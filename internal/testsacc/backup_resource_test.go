/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

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

func (r *BackupResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *BackupResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First Test For a VDC Backup named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
					return resp
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
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VAppResourceName]().GetDefaultConfig)
					return resp
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
									policy_name = "D30"
								},
								{
									policy_name = "M3"
								}]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.policy_name", "D30"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.policy_name", "M3"),
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
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VMResourceName]().GetDefaultConfig)
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_backup" "example_vm" {
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
						resource "cloudavenue_backup" "example_vm" {
						  type = "vm"
						  target_name = cloudavenue_vm.example.name
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
