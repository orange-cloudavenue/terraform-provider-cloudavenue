// Package iam provides the acceptance tests for the provider.
package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccUserDataSourceConfig = `
data "cloudavenue_iam_user" "example" {
	name = cloudavenue_iam_user.example.name
}
`

func TestAccUserDataSource(t *testing.T) {
	datasourceName := "data.cloudavenue_iam_user.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: tests.ConcatTests(TestAccOrgUserResourceConfig, testAccUserDataSourceConfig),
				Check:  testsOrgUserResourceConfig(datasourceName),
			},
			{
				Config: tests.ConcatTests(testAccOrgUserResourceConfigFull, testAccUserDataSourceConfig),
				Check:  testsOrgUserResourceConfigFull(datasourceName, true),
			},
		},
	})
}
