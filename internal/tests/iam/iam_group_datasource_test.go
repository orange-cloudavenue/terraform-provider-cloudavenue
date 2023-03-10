// Package iam provides the acceptance tests for the provider.
package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccOrgGroupDataSourceConfig = `
resource "cloudavenue_iam_group" "example" {
  name          = "OrgTest"
  role          = "Organization Administrator"
  description   = "org test from go test"
}

data "cloudavenue_iam_group" "example" {
	name = cloudavenue_iam_group.example.name
}
`

func TestAccOrgGroupDataSource(t *testing.T) {
	resourceName := "data.cloudavenue_iam_group.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
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
