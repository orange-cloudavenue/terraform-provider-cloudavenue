package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccSecurityTagResourceConfig = `
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

resource "cloudavenue_vm_security_tag" "example" {
	id = "tag-example"
	vm_ids = [
    cloudavenue_vm.example.id,
  ]
}
`

const testAccSecurityTagResourceConfigUpdate = `
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

resource "cloudavenue_vm_security_tag" "example" {
	id = "tag-example"
	vm_ids = [
    cloudavenue_vm.example.id,
    cloudavenue_vm.example2.id,
  ]
}
`

func TestAccSecurityTagResource(t *testing.T) {
	const resourceName = "cloudavenue_vm_security_tag.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccSecurityTagResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "tag-example"),
					resource.TestCheckResourceAttrWith(resourceName, "vm_ids.0", uuid.TestIsType(uuid.VM)),
				),
			},
			{
				// Apply test
				Config: testAccSecurityTagResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "tag-example"),
					resource.TestCheckResourceAttrWith(resourceName, "vm_ids.0", uuid.TestIsType(uuid.VM)),
					resource.TestCheckResourceAttrWith(resourceName, "vm_ids.1", uuid.TestIsType(uuid.VM)),
				),
			},
			// Import testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
