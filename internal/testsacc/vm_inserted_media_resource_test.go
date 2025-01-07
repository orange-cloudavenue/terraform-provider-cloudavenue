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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const TestAccVMInsertedMediaResourceConfig = `
resource "cloudavenue_vm_inserted_media" "example" {
	catalog = "catalog-example"
	name    = "debian-9.9.0-amd64-netinst.iso"
	vapp_name = cloudavenue_vapp.example.name
	vm_name   = cloudavenue_vm.example.name
}
`

func TestAccVMInsertedMediaResource(t *testing.T) {
	resourceName := "cloudavenue_vm_inserted_media.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: ConcatTests(TestAccVMResourceConfigFromVappTemplate, TestAccVMInsertedMediaResourceConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
					resource.TestCheckResourceAttr(resourceName, "catalog", "catalog-example"),
					resource.TestCheckResourceAttr(resourceName, "name", "debian-9.9.0-amd64-netinst.iso"),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vm_name"),
				),
			},
		},
	})
}
