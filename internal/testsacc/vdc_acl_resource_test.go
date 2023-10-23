package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccVDCACLResourceConfig = `
resource "cloudavenue_vdc_acl" "example" {
  vdc                   = "VDC_Test" # Optional
  everyone_access_level = "ReadOnly"
}
`

const testAccVDCACLResourceSharedWithConfig = `
resource "cloudavenue_vdc_acl" "example" {
  vdc                   = "VDC_Test" # Optional
	shared_with = [
	{
	  access_level = "ReadOnly"
	  user_id      = "urn:vcloud:user:53665519-7036-43ea-ba97-63fc5a2aabe7"
	}
	]
}
`

func TestAccVDCACLResource(t *testing.T) {
	const resourceName = "cloudavenue_vdc_acl.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVDCACLResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VDC)),
					resource.TestCheckResourceAttr(resourceName, "vdc", "VDC_Test"),
					resource.TestCheckResourceAttr(resourceName, "everyone_access_level", "ReadOnly"),
				),
			},
			{
				// Apply test
				Config: testAccVDCACLResourceSharedWithConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VDC)),
					resource.TestCheckResourceAttr(resourceName, "vdc", "VDC_Test"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_with.0.subject_name"),
				),
			},
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "VDC_Test",
			},
		},
	})
}
