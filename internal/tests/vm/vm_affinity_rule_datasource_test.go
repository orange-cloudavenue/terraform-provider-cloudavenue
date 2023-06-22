package vm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccVMAffinityRuleDataSourceConfig = `
data "cloudavenue_vm_affinity_rule" "example" {
	id = cloudavenue_vm_affinity_rule.example.id
}

resource "cloudavenue_vm_affinity_rule" "example" {
	name     = "example-affinity-rule"
	polarity = "Affinity"
  
	vm_ids = [
	  cloudavenue_vm.example.id,
	  cloudavenue_vm.example2.id,
	]
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
  
  resource "cloudavenue_vm" "example2" {
	  name      = "example-vm2"
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
  
  data "cloudavenue_catalog_vapp_template" "example" {
	  catalog_name = "Orange-Linux"
	  template_name    = "debian_10_X64"
  }
  
  resource "cloudavenue_vapp" "example" {
	  name = "vapp_example"
	  description = "This is a example vapp"
  }
`

func TestAccVmAffinityRuleDatasourceDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vm_affinity_rule.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVMAffinityRuleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "example-affinity-rule"),
					resource.TestCheckResourceAttr(dataSourceName, "polarity", "Affinity"),
					resource.TestCheckResourceAttr(dataSourceName, "required", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vm_ids.0"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vm_ids.1"),
				),
			},
		},
	})
}
