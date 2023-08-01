// Package vm provides the acceptance tests for the provider.
package vm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const TestAccVMResourceConfigFromVappTemplate = `
data "cloudavenue_catalog_vapp_template" "example" {
	catalog_name = "Orange-Linux"
	template_name    = "debian_10_X64"
}

resource "cloudavenue_vapp" "example" {
	name = "vapp_example"
	description = "This is a example vapp"
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

const testAccVMResourceConfigFromVappTemplateUpdate = `
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
		guest_properties = {
		  "guestinfo.hostname" = "example-vm"
		}
		customization = {
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
			ip                 = "192.168.0.111"
			ip_allocation_mode = "MANUAL"
			is_primary         = true
		  },
		  {
			type               = "org"
			name               = cloudavenue_vapp_org_network.example.network_name
			ip_allocation_mode = "DHCP"
		  }
		]
	  }
  
	state = {
	}
  }
`

func TestAccVMResource(t *testing.T) {
	const resourceNameVM = "cloudavenue_vm.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: TestAccVMResourceConfigFromVappTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					// ! base
					// ? id
					resource.TestCheckResourceAttrSet(resourceNameVM, "id"),
					// ? vapp_name
					resource.TestCheckResourceAttr(resourceNameVM, "vapp_name", "vapp_example"),
					// ? vapp_id
					resource.TestCheckResourceAttrSet(resourceNameVM, "vapp_id"),
					// ? vdc
					resource.TestCheckResourceAttr(resourceNameVM, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					// ? name
					resource.TestCheckResourceAttr(resourceNameVM, "name", "example-vm"),
					// ? description
					resource.TestCheckNoResourceAttr(resourceNameVM, "description"),

					// ! deploy_os
					// ? accept_all_eulas
					resource.TestCheckNoResourceAttr(resourceNameVM, "deploy_os.accept_all_eulas"),
					// ? boot_image_id
					resource.TestCheckNoResourceAttr(resourceNameVM, "deploy_os.boot_image_id"),
					// ? vapp_template_id (No check value becaus it's provided by the catalog)
					resource.TestCheckResourceAttrSet(resourceNameVM, "deploy_os.vapp_template_id"),
					// ? vm_name_in_template
					resource.TestCheckNoResourceAttr(resourceNameVM, "deploy_os.vm_name_in_template"),

					// ! settings
					// ? affinity_rule_id
					resource.TestCheckResourceAttrSet(resourceNameVM, "settings.affinity_rule_id"),
					// ? expose_hardware_virtualization
					resource.TestCheckResourceAttr(resourceNameVM, "settings.expose_hardware_virtualization", "false"),
					// ? os_type
					resource.TestCheckResourceAttr(resourceNameVM, "settings.os_type", "debian10_64Guest"),
					// ? storage_profile
					resource.TestCheckResourceAttr(resourceNameVM, "settings.storage_profile", "gold"),
					// * customization
					// ? enabled
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.enabled", "false"),
					// ? admin_password
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.admin_password"),
					// ? allow_local_admin_password
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.allow_local_admin_password", "false"),
					// ? change_sid
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.change_sid", "false"),
					// ? force
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.force", "false"),
					// ? hostname
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.hostname", "example-vm"),
					// ? init_script
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.init_script"),
					// ? join_domain
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.join_domain", "false"),
					// ? join_domain_account_ou
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.join_domain_account_ou"),
					// ? join_domain_name
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.join_domain_name"),
					// ? join_domain_password
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.join_domain_password"),
					// ? join_domain_user
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.join_domain_user"),
					// ? join_org_domain
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.join_org_domain", "false"),
					// ? must_change_password_on_first_login
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.must_change_password_on_first_login", "false"),
					// ? number_of_auto_logons
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.number_of_auto_logons", "0"),

					// ! state
					// ? power_on
					resource.TestCheckResourceAttr(resourceNameVM, "state.power_on", "true"),
					// ? status
					resource.TestCheckResourceAttr(resourceNameVM, "state.status", "POWERED_ON"),

					// ! resource
					// ? cpus
					resource.TestCheckResourceAttr(resourceNameVM, "resource.cpus", "1"),
					// ? cpu_hot_add_enabled
					resource.TestCheckResourceAttr(resourceNameVM, "resource.cpu_hot_add_enabled", "true"),
					// ? cpus_cores
					resource.TestCheckResourceAttr(resourceNameVM, "resource.cpus_cores", "1"),
					// ? memory
					resource.TestCheckResourceAttr(resourceNameVM, "resource.memory", "1024"),
					// ? memory_hot_add_enabled
					resource.TestCheckResourceAttr(resourceNameVM, "resource.memory_hot_add_enabled", "true"),
					// * networks

				),
			},
			{
				// Apply test
				Config: testAccVMResourceConfigFromVappTemplateUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					// ! base
					// ? id
					resource.TestCheckResourceAttrSet(resourceNameVM, "id"),
					// ? vapp_name
					resource.TestCheckResourceAttr(resourceNameVM, "vapp_name", "vapp_example"),
					// ? vapp_id
					resource.TestCheckResourceAttrSet(resourceNameVM, "vapp_id"),
					// ? vdc
					resource.TestCheckResourceAttr(resourceNameVM, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					// ? name
					resource.TestCheckResourceAttr(resourceNameVM, "name", "example-vm"),
					// ? description
					resource.TestCheckNoResourceAttr(resourceNameVM, "description"),

					// ! deploy_os
					// ? accept_all_eulas
					resource.TestCheckNoResourceAttr(resourceNameVM, "deploy_os.accept_all_eulas"),
					// ? boot_image_id
					resource.TestCheckNoResourceAttr(resourceNameVM, "deploy_os.boot_image_id"),
					// ? vapp_template_id (No check value becaus it's provided by the catalog)
					resource.TestCheckResourceAttrSet(resourceNameVM, "deploy_os.vapp_template_id"),
					// ? vm_name_in_template
					resource.TestCheckNoResourceAttr(resourceNameVM, "deploy_os.vm_name_in_template"),

					// ! settings
					// ? affinity_rule_id
					resource.TestCheckResourceAttrSet(resourceNameVM, "settings.affinity_rule_id"),
					// ? expose_hardware_virtualization
					resource.TestCheckResourceAttr(resourceNameVM, "settings.expose_hardware_virtualization", "false"),
					// ? os_type
					resource.TestCheckResourceAttr(resourceNameVM, "settings.os_type", "debian10_64Guest"),
					// ? storage_profile
					resource.TestCheckResourceAttr(resourceNameVM, "settings.storage_profile", "gold"),
					// ? guest_properties
					resource.TestCheckResourceAttrSet(resourceNameVM, "settings.guest_properties.%"),
					// * customization
					// ? enabled
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.enabled", "false"),
					// ? admin_password
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.admin_password"),
					// ? allow_local_admin_password
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.allow_local_admin_password", "false"),
					// ? change_sid
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.change_sid", "false"),
					// ? force
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.force", "false"),
					// ? hostname
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.hostname", "example-vm"),
					// ? init_script
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.init_script"),
					// ? join_domain
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.join_domain", "false"),
					// ? join_domain_account_ou
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.join_domain_account_ou"),
					// ? join_domain_name
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.join_domain_name"),
					// ? join_domain_password
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.join_domain_password"),
					// ? join_domain_user
					resource.TestCheckNoResourceAttr(resourceNameVM, "settings.customization.join_domain_user"),
					// ? join_org_domain
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.join_org_domain", "false"),
					// ? must_change_password_on_first_login
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.must_change_password_on_first_login", "false"),
					// ? number_of_auto_logons
					resource.TestCheckResourceAttr(resourceNameVM, "settings.customization.number_of_auto_logons", "0"),

					// ! state
					// ? power_on
					resource.TestCheckResourceAttr(resourceNameVM, "state.power_on", "true"),
					// ? status
					resource.TestCheckResourceAttr(resourceNameVM, "state.status", "POWERED_ON"),

					// ! resource
					// ? cpus
					resource.TestCheckResourceAttr(resourceNameVM, "resource.cpus", "2"),
					// ? cpu_hot_add_enabled
					resource.TestCheckResourceAttr(resourceNameVM, "resource.cpu_hot_add_enabled", "true"),
					// ? cpus_cores
					resource.TestCheckResourceAttr(resourceNameVM, "resource.cpus_cores", "1"),
					// ? memory
					resource.TestCheckResourceAttr(resourceNameVM, "resource.memory", "2048"),
					// ? memory_hot_add_enabled
					resource.TestCheckResourceAttr(resourceNameVM, "resource.memory_hot_add_enabled", "true"),
					// * networks
					// # 0
					// ? ip
					resource.TestCheckResourceAttr(resourceNameVM, "resource.networks.0.ip", "192.168.0.111"),
					// ? ip_allocation_mode
					resource.TestCheckResourceAttr(resourceNameVM, "resource.networks.0.ip_allocation_mode", "MANUAL"),
					// ? is_primary
					resource.TestCheckResourceAttr(resourceNameVM, "resource.networks.0.is_primary", "true"),
					// ? type
					resource.TestCheckResourceAttr(resourceNameVM, "resource.networks.0.type", "org"),
					// ? name
					resource.TestCheckResourceAttr(resourceNameVM, "resource.networks.0.name", "INET"),
					// # 1
					// ? ip (DHCP)
					// ? ip_allocation_mode
					resource.TestCheckResourceAttr(resourceNameVM, "resource.networks.1.ip_allocation_mode", "DHCP"),
					// ? is_primary
					resource.TestCheckResourceAttr(resourceNameVM, "resource.networks.1.is_primary", "false"),
					// ? type
					resource.TestCheckResourceAttr(resourceNameVM, "resource.networks.1.type", "org"),
					// ? name
					resource.TestCheckResourceAttr(resourceNameVM, "resource.networks.1.name", "INET"),
				),
			},
			// ImportruetState testing
			{
				// Import test
				ResourceName:      resourceNameVM,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testAccVMResourceImportStateIDFunc(resourceNameVM),
			},
		},
	})
}

func testAccVMResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["vapp_name"], rs.Primary.ID), nil
	}
}
