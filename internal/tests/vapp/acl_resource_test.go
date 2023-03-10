package vapp

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccACLResourceConfig = `
resource "cloudavenue_vapp_acl" "example" {
	vdc                  = "MyVDC" # Optional
	vapp_name            = "MyVapp"
	shared_with = [{
	  access_level = "ReadOnly"
	  user_id      = "urn:vcloud:user:53665519-7036-43ea-ba97-63fc5a2aabe7"
	  },
	  {
		access_level = "FullControl"
		group_id     = "urn:vcloud:group:cd04ff68-688a-4ccb-87c1-905bbe4dba7e"
	}]
  }
`

const testAccACLResourceUpdateConfig = `
resource "cloudavenue_vapp_acl" "example" {
	vdc                   = "MyVDC" # Optional
	vapp_name             = "MyVapp"
	everyone_access_level = "Change"
  }
`

func TestAccACLResource(t *testing.T) {
	const resourceName = "cloudavenue_vapp_acl.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccACLResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`(urn:vcloud:vapp:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttr(resourceName, "vdc", "MyVDC"),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "MyVapp"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_with.0.subject_name"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_with.1.subject_name"),
				),
			},
			// Uncomment if you want to test update or delete this block
			{
				// Update test
				Config: testAccACLResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`(urn:vcloud:vapp:[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)),
					resource.TestCheckResourceAttr(resourceName, "vdc", "MyVDC"),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "MyVapp"),
					resource.TestCheckResourceAttr(resourceName, "everyone_access_level", "Change"),
				),
			},
			// ImportruetState testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "MyVDC.MyVapp",
			},
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "MyVapp",
			},
		},
	})
}
