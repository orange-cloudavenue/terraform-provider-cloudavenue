package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../examples -test
const testAccOrgGroupDataSourceConfig = `
resource "cloudavenue_org_group" "example" {
  name          = "OrgTest"
  role          = "Organization Administrator"
  description   = "org test from go test"
}

data "cloudavenue_org_group" "example" {
	name = cloudavenue_org_group.example.name
}
`

func TestAccOrgGroupDataSource(t *testing.T) {
	resourceName := "data.cloudavenue_org_group.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOrgGroupDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "OrgTest"),
					resource.TestCheckResourceAttr(resourceName, "role", "Organization Administrator"),
					resource.TestCheckResourceAttr(resourceName, "description", "org test from go test"),
					resource.TestCheckResourceAttr(resourceName, "user_names.#", "0"),
				),
			},
		},
	})
}
