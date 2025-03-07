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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

const testAccVMDataSourceConfig = `
data "cloudavenue_catalog_vapp_template" "example" {
	catalog_name  	= "Orange-Linux"
	template_name 	= "UBUNTU_20.04"
}

resource "cloudavenue_vm" "example" {
   name        = "example-vm"
   description = "This is a example vm"
 
   vapp_name = cloudavenue_vapp.example.name
 
   deploy_os = {
     vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
   }

   settings = {
	customization = {
		auto_generate_password = true
	}
   }
 
   state = {}
   resource = {}
}

resource "cloudavenue_vapp" "example" {
	name        = "example-vapp"
	description = "This is an example vApp"
}
  
data "cloudavenue_vm" "example" {
	name = cloudavenue_vm.example.name
	vapp_name = cloudavenue_vapp.example.name
}
`

func TestAccVMDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vm.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccVMDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// ! basic
					resource.TestCheckResourceAttrWith(dataSourceName, "id", urn.TestIsType(urn.VM)),
					resource.TestCheckResourceAttr(dataSourceName, "name", "example-vm"),
					resource.TestCheckResourceAttr(dataSourceName, "vapp_name", "example-vapp"),
					resource.TestCheckResourceAttrWith(dataSourceName, "vapp_id", urn.TestIsType(urn.VAPP)),
					resource.TestCheckResourceAttr(dataSourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					// ! resource
					resource.TestCheckResourceAttr(dataSourceName, "resource.cpus", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.cpus_cores", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.cpu_hot_add_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.memory", "1024"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.memory_hot_add_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.networks.#", "0"),
					// ! settings
					resource.TestCheckResourceAttrWith(dataSourceName, "settings.affinity_rule_id", urn.TestIsType(urn.VDCComputePolicy)),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.allow_local_admin_password"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.auto_generate_password"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.change_sid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.force"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.hostname"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.join_domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.join_org_domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.must_change_password_on_first_login"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.customization.number_of_auto_logons"),
					resource.TestCheckResourceAttrSet(dataSourceName, "settings.expose_hardware_virtualization"),
					resource.TestCheckResourceAttr(dataSourceName, "settings.os_type", "ubuntu64Guest"),
					resource.TestCheckResourceAttr(dataSourceName, "settings.storage_profile", "gold"),
					// ! state
					resource.TestCheckResourceAttr(dataSourceName, "state.power_on", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "state.status", "POWERED_ON"),
				),
			},
		},
	})
}
