// Package tests provides the acceptance tests for the provider.
// Package iam provides the acceptance tests for the provider.
package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../examples -test
const testAccOrgRoleResourceConfig = `
resource "cloudavenue_iam_role" "example" {
	name   		= "roletest"
	description = "A test role"
	rights = [
    	"Catalog: Add vApp from My Cloud",
    	"Catalog: Edit Properties",
    	"Catalog: View Private and Shared Catalogs",
    	"Organization vDC Compute Policy: View",
    	"vApp Template / Media: Edit",
    	"vApp Template / Media: View",
  	]
  }
`

func TestAccOrgRoleResource(t *testing.T) {
	resourceName := "cloudavenue_iam_role.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOrgRoleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "roletest"),
					resource.TestCheckResourceAttr(resourceName, "description", "A test role"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rights.*", "Catalog: Add vApp from My Cloud"),
					resource.TestCheckResourceAttr(resourceName, "rights.#", "6"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "roletest",
				ImportStateVerify: true,
			},
		},
	})
}
