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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDisksDataSourceConfig = `
data "cloudavenue_vm_disks" "example" {
	vm_id = cloudavenue_vm.example.id
	vapp_id = cloudavenue_vapp.example.id
	depends_on = [cloudavenue_vm_disk.example-detachable-with-vm, cloudavenue_vm_disk.example-detachable, cloudavenue_vm_disk.example-internal]
}

resource "cloudavenue_vapp" "example" {
	name = "vapp_example"
	description = "This is a example vapp"
}

resource "cloudavenue_vm_disk" "example-detachable" {
	vapp_id = cloudavenue_vapp.example.id
	name = "disk-example-detachable"
	bus_type = "SATA"
	size_in_mb = 2048
	is_detachable = true
}


resource "cloudavenue_vm_disk" "example-detachable-with-vm" {
	vapp_id = cloudavenue_vapp.example.id
	name = "disk-example-detachable-with-vm"
	bus_type = "SATA"
	size_in_mb = 2048
	is_detachable = true
	vm_id = cloudavenue_vm.example.id
}

resource "cloudavenue_vm_disk" "example-internal" {
	vapp_id = cloudavenue_vapp.example.id
	bus_type = "SATA"
	size_in_mb = 2048
	is_detachable = false
	vm_id = cloudavenue_vm.example.id
}

data "cloudavenue_catalog_vapp_template" "example" {
	catalog_name = "Orange-Linux"
	template_name    = "debian_10_X64"
}

resource "cloudavenue_vm" "example" {
	name      = "example-vm"
	vapp_name = cloudavenue_vapp.example.name
	deploy_os = {
	  vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
	}
	settings = {
	  customization = {}
	}

	resource = {}
	state = {}
}
`

func TestAccDisksDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vm_disks.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: ConcatTests(testAccDisksDataSourceConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vapp_id"),
					resource.TestCheckResourceAttr(dataSourceName, "vapp_name", "vapp_example"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vm_id"),
					resource.TestCheckResourceAttr(dataSourceName, "vm_name", "example-vm"),
					// 1 disk system attached to the vm
					// 1 internal disk attached to the vm
					// 1 detachable disk attached to the vm
					// 1 detachable disk not attached to the vm
					resource.TestCheckResourceAttr(dataSourceName, "disks.#", "4"),

					resource.TestCheckResourceAttrSet(dataSourceName, "disks.0.id"),
					// resource.TestCheckResourceAttrSet(dataSourceName, "disks.0.name"), // Name is not set for internal disks
					resource.TestCheckResourceAttrSet(dataSourceName, "disks.0.size_in_mb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "disks.0.storage_profile"),
					resource.TestCheckResourceAttrSet(dataSourceName, "disks.0.is_detachable"),
				),
			},
		},
	})
}
