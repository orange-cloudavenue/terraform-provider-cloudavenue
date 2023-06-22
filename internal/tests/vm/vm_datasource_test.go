package vm

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
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
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccVMDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// ! basic
					resource.TestMatchResourceAttr(dataSourceName, "id", regexp.MustCompile(`(urn:vcloud:vm:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttr(dataSourceName, "name", "example-vm"),
					resource.TestCheckResourceAttr(dataSourceName, "vapp_name", "example-vapp"),
					resource.TestMatchResourceAttr(dataSourceName, "vapp_id", regexp.MustCompile(`(urn:vcloud:vapp:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttr(dataSourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					// ! resource
					resource.TestCheckResourceAttr(dataSourceName, "resource.cpus", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.cpus_cores", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.cpu_hot_add_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.memory", "1024"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.memory_hot_add_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.networks.#", "0"),
					// ! settings
					resource.TestMatchResourceAttr(dataSourceName, "settings.affinity_rule_id", regexp.MustCompile(`(urn:vcloud:vdcComputePolicy:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
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
