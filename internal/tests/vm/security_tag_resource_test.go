package vm

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../examples -test
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

func TestAccSecurityTagResource(t *testing.T) {
	const resourceName = "cloudavenue_vm_security_tag.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccSecurityTagResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "tag-example"),
					resource.TestMatchResourceAttr(resourceName, "vm_ids.0", regexp.MustCompile(`(urn:vcloud:vm:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
				),
			},
			// Uncomment if you want to test update or delete this block
			// {
			// 	// Update test
			// 	Config: strings.Replace(testAccSecurityTagResourceConfig, "tag-example", "example-tag", 1),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceName, "id", "example-tag"),
			// 		resource.TestMatchResourceAttr(resourceName, "vm_ids.0", regexp.MustCompile(`(urn:vcloud:network:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
			// 	),
			// },
			// ImportruetState testing
			// {
			// 	// Import test
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}
