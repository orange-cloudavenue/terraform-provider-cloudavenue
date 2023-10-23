package testsacc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccVAPPACLResourceConfig = `
resource "cloudavenue_iam_user" "example" {
	name   		= "example"
	role_name   = "Organization Administrator"
	password    = "Th!s1sSecur3P@ssword"
}

resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example vApp"
}

resource "cloudavenue_vapp_acl" "example" {
	vapp_name      = cloudavenue_vapp.example.name
	shared_with = [{
	  access_level = "ReadOnly"
	  user_id      = cloudavenue_iam_user.example.id
	  }]
}
`

const testAccVAPPACLResourceConfigUpdate = `
resource "cloudavenue_iam_user" "example" {
	name   		= "example"
	role_name   = "Organization Administrator"
	password    = "Th!s1sSecur3P@ssword"
}

resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example vApp"
}

resource "cloudavenue_vapp_acl" "example" {
	vapp_name     		  = cloudavenue_vapp.example.name
	everyone_access_level = "Change"
  }
`

func TestAccVAPPACLResource(t *testing.T) {
	const resourceName = "cloudavenue_vapp_acl.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVAPPACLResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VAPP)),
					resource.TestCheckResourceAttr(resourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "MyVapp"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_with.0.subject_name"),
				),
			},
			// Uncomment if you want to test update or delete this block
			{
				// Update test
				Config: testAccVAPPACLResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VAPP)),
					resource.TestCheckResourceAttr(resourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "MyVapp"),
					resource.TestCheckResourceAttr(resourceName, "everyone_access_level", "Change"),
				),
			},
			// Import testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccACLResourceImportStateIDFunc(resourceName),
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

func testAccACLResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s.%s", rs.Primary.Attributes["vdc"], rs.Primary.Attributes["vapp_name"]), nil
	}
}
