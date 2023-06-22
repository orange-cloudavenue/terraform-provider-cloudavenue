package vapp

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccACLResourceConfig = `
resource "cloudavenue_iam_user" "example" {
	user_name   = "exampleuser"
	description = "An example user"
	role        = "Organization Administrator"
	password    = "Th!s1sSecur3P@ssword"
  }

resource "cloudavenue_iam_group" "example" {
	name        = "examplegroup"
	role        = "Organization Administrator"
	description = "An example group"
  }

resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example vApp"
  }

resource "cloudavenue_vapp_acl" "example" {
	vdc                  = "MyVDC" # Optional
	vapp_name            = cloudavenue_vapp.example.name
	shared_with = [{
	  access_level = "ReadOnly"
	  user_id      = cloudavenue_iam_user.example.id
	  },
	  {
		access_level = "FullControl"
		group_id     = cloudavenue_iam_group.example.id
	}]
  }
`

const testACLResourceUpdateConfig = `
resource "cloudavenue_iam_user" "example" {
	user_name   = "exampleuser"
	description = "An example user"
	role        = "Organization Administrator"
	password    = "Th!s1sSecur3P@ssword"
  }

resource "cloudavenue_iam_group" "example" {
	name        = "examplegroup"
	role        = "Organization Administrator"
	description = "An example group"
  }

resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example vApp"
  }

resource "cloudavenue_vapp_acl" "example" {
	vdc                   = "MyVDC" # Optional
	vapp_name             = cloudavenue_vapp.example.name
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
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.VAPP.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					resource.TestCheckResourceAttr(resourceName, "vapp_name", "MyVapp"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_with.0.subject_name"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_with.1.subject_name"),
				),
			},
			// Uncomment if you want to test update or delete this block
			{
				// Update test
				Config: testACLResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.VAPP.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
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
