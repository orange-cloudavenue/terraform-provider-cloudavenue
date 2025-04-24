/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

// Used into old tests.
const TestAccVMResourceConfigFromVappTemplate = `
data "cloudavenue_catalog_vapp_template" "example" {
	catalog_name = "Orange-Linux"
	template_name    = "debian_10_X64"
}

resource "cloudavenue_vapp" "example" {
	name = "vapp_example"
	description = "This is a example vapp"
}

resource "cloudavenue_vapp_org_network" "example" {
	vapp_name    = cloudavenue_vapp.example.name
	network_name = "INET"
}

resource "cloudavenue_vm" "example" {
	name      = "example-vm"
	vapp_name = cloudavenue_vapp.example.name
	deploy_os = {
	  vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
	}
	settings = {
	  customization = {
		auto_generate_password = true
	  }
	}
	resource = {
	}

	state = {
	}
}
`

var _ testsacc.TestACC = &VMResource{}

const (
	VMResourceName = testsacc.ResourceName("cloudavenue_vm")
)

type VMResource struct{}

func NewVMResourceTest() testsacc.TestACC {
	return &VMResource{}
}

// GetResourceName returns the name of the resource.
func (r *VMResource) GetResourceName() string {
	return VMResourceName.String()
}

func (r *VMResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VAppOrgNetworkResourceName]().GetDefaultConfig)
	resp.Append(GetResourceConfig()[CatalogVAppTemplateDataSourceName]().GetDefaultConfig)
	return
}

func (r *VMResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VM)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vm" "example" {
						name      = {{ generate . "name" }}
						vdc 	  = cloudavenue_vdc.example.name
						vapp_name = cloudavenue_vapp.example.name
						deploy_os = {
						  vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
						}
						settings = {
						  customization = {
							auto_generate_password = true
						  }
						}
						resource = {
						}

						state = {
						}
					}`),
					Checks: []resource.TestCheckFunc{
						// ! base
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckNoResourceAttr(resourceName, "description"),

						resource.TestCheckNoResourceAttr(resourceName, "deploy_os.accept_all_eulas"),
						resource.TestCheckNoResourceAttr(resourceName, "deploy_os.boot_image_id"),
						// ? vapp_template_id (No check value becaus it's provided by the catalog)
						resource.TestCheckResourceAttrSet(resourceName, "deploy_os.vapp_template_id"),
						resource.TestCheckNoResourceAttr(resourceName, "deploy_os.vm_name_in_template"),

						resource.TestCheckResourceAttr(resourceName, "settings.expose_hardware_virtualization", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.storage_profile", "gold"),
						resource.TestCheckResourceAttrSet(resourceName, "settings.affinity_rule_id"),
						resource.TestCheckResourceAttrSet(resourceName, "settings.os_type"),

						resource.TestCheckResourceAttr(resourceName, "settings.customization.enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.allow_local_admin_password", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.change_sid", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.force", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.hostname", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.join_domain", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.join_org_domain", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.must_change_password_on_first_login", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.number_of_auto_logons", "0"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.admin_password"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.init_script"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_account_ou"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_name"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_password"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_user"),

						resource.TestCheckResourceAttr(resourceName, "state.power_on", "true"),

						resource.TestCheckResourceAttr(resourceName, "resource.cpus", "1"),
						resource.TestCheckResourceAttr(resourceName, "resource.cpu_hot_add_enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "resource.cpus_cores", "1"),
						resource.TestCheckResourceAttr(resourceName, "resource.memory", "1024"),
						resource.TestCheckResourceAttr(resourceName, "resource.memory_hot_add_enabled", "true"),
					},
				},
				// ! Update testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vm" "example" {
							name      = {{ get . "name" }}
							vdc 	  = cloudavenue_vdc.example.name
							vapp_name = cloudavenue_vapp.example.name
							deploy_os = {
							  vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
							}
							settings = {
								guest_properties = {
								  "guestinfo.hostname" = {{ get . "name" }}
								}
								customization = {
								  enabled = true
								  auto_generate_password = true
								}
							  }
							  resource = {
								cpus   = 2
								memory = 2048
								networks = [
								  {
									type               = "org"
									name               = cloudavenue_vapp_org_network.example.network_name
									ip                 = "192.168.1.111"
									ip_allocation_mode = "MANUAL"
									is_primary         = true
								  },
								  {
									type               = "org"
									name               = cloudavenue_vapp_org_network.example.network_name
									ip_allocation_mode = "DHCP"
								  },
								  {
									type               = "org"
									name               = cloudavenue_vapp_org_network.example.network_name
									ip_allocation_mode = "POOL"
								  },
								  {
									type               = "org"
									name               = cloudavenue_vapp_org_network.example.network_name
									ip_allocation_mode = "NONE"
								  },
								]
							  }

							state = {
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckNoResourceAttr(resourceName, "description"),

							resource.TestCheckNoResourceAttr(resourceName, "deploy_os.accept_all_eulas"),
							resource.TestCheckNoResourceAttr(resourceName, "deploy_os.boot_image_id"),
							// ? vapp_template_id (No check value becaus it's provided by the catalog)
							resource.TestCheckResourceAttrSet(resourceName, "deploy_os.vapp_template_id"),
							resource.TestCheckNoResourceAttr(resourceName, "deploy_os.vm_name_in_template"),

							resource.TestCheckResourceAttr(resourceName, "settings.expose_hardware_virtualization", "false"),
							resource.TestCheckResourceAttr(resourceName, "settings.storage_profile", "gold"),
							resource.TestCheckResourceAttrSet(resourceName, "settings.affinity_rule_id"),
							resource.TestCheckResourceAttrSet(resourceName, "settings.os_type"),

							resource.TestCheckResourceAttr(resourceName, "settings.customization.enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "settings.customization.allow_local_admin_password", "false"),
							resource.TestCheckResourceAttr(resourceName, "settings.customization.change_sid", "false"),
							resource.TestCheckResourceAttr(resourceName, "settings.customization.force", "false"),
							resource.TestCheckResourceAttr(resourceName, "settings.customization.hostname", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "settings.customization.join_domain", "false"),
							resource.TestCheckResourceAttr(resourceName, "settings.customization.join_org_domain", "false"),
							resource.TestCheckResourceAttr(resourceName, "settings.customization.must_change_password_on_first_login", "false"),
							resource.TestCheckResourceAttr(resourceName, "settings.customization.number_of_auto_logons", "0"),
							resource.TestCheckNoResourceAttr(resourceName, "settings.customization.admin_password"),
							resource.TestCheckNoResourceAttr(resourceName, "settings.customization.init_script"),
							resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_account_ou"),
							resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_name"),
							resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_password"),
							resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_user"),

							resource.TestCheckResourceAttrSet(resourceName, "settings.guest_properties.%"),

							resource.TestCheckResourceAttr(resourceName, "state.power_on", "true"),

							resource.TestCheckResourceAttr(resourceName, "resource.cpus", "2"),
							resource.TestCheckResourceAttr(resourceName, "resource.cpu_hot_add_enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "resource.cpus_cores", "1"),
							resource.TestCheckResourceAttr(resourceName, "resource.memory", "2048"),
							resource.TestCheckResourceAttr(resourceName, "resource.memory_hot_add_enabled", "true"),

							// * networks
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "resource.networks.*", map[string]string{
								"ip_allocation_mode": "MANUAL",
								"is_primary":         "true",
								"type":               "org",
								"ip":                 "192.168.1.111",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "resource.networks.*", map[string]string{
								"ip_allocation_mode": "DHCP",
								"is_primary":         "false",
								"type":               "org",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "resource.networks.*", map[string]string{
								"ip_allocation_mode": "POOL",
								"is_primary":         "false",
								"type":               "org",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "resource.networks.*", map[string]string{
								"ip_allocation_mode": "NONE",
								"is_primary":         "false",
								"type":               "org",
							}),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc", "vapp_name", "id"},
						ImportState:          true,
						ImportStateVerify:    false,
					},
				},
			}
		},
		"example_with_password": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VM)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vm" "example_with_password" {
						name      = {{ generate . "name" }}
						vdc 	  = cloudavenue_vdc.example.name
						vapp_name = cloudavenue_vapp.example.name
						deploy_os = {
						  vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
						}
						settings = {
						  customization = {
						  	enabled = true
							allow_local_admin_password = true
							admin_password = "password"
						  }
						}
						resource = {
						}

						state = {
						}
					}`),
					Checks: []resource.TestCheckFunc{
						// ! base
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckNoResourceAttr(resourceName, "description"),

						resource.TestCheckNoResourceAttr(resourceName, "deploy_os.accept_all_eulas"),
						resource.TestCheckNoResourceAttr(resourceName, "deploy_os.boot_image_id"),
						// ? vapp_template_id (No check value becaus it's provided by the catalog)
						resource.TestCheckResourceAttrSet(resourceName, "deploy_os.vapp_template_id"),
						resource.TestCheckNoResourceAttr(resourceName, "deploy_os.vm_name_in_template"),

						resource.TestCheckResourceAttr(resourceName, "settings.expose_hardware_virtualization", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.storage_profile", "gold"),
						resource.TestCheckResourceAttrSet(resourceName, "settings.affinity_rule_id"),
						resource.TestCheckResourceAttrSet(resourceName, "settings.os_type"),

						resource.TestCheckResourceAttr(resourceName, "settings.customization.enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.allow_local_admin_password", "true"),
						resource.TestCheckResourceAttrSet(resourceName, "settings.customization.admin_password"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.change_sid", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.force", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.hostname", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.join_domain", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.join_org_domain", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.must_change_password_on_first_login", "false"),
						resource.TestCheckResourceAttr(resourceName, "settings.customization.number_of_auto_logons", "0"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.init_script"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_account_ou"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_name"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_password"),
						resource.TestCheckNoResourceAttr(resourceName, "settings.customization.join_domain_user"),

						resource.TestCheckResourceAttr(resourceName, "state.power_on", "true"),

						resource.TestCheckResourceAttr(resourceName, "resource.cpus", "1"),
						resource.TestCheckResourceAttr(resourceName, "resource.cpu_hot_add_enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "resource.cpus_cores", "1"),
						resource.TestCheckResourceAttr(resourceName, "resource.memory", "1024"),
						resource.TestCheckResourceAttr(resourceName, "resource.memory_hot_add_enabled", "true"),
					},
				},
			}
		},
	}
}

func TestAccVMResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VMResource{}),
	})
}
