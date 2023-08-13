// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccUserDataSourceConfig = `
data "cloudavenue_iam_user" "example" {
	name = cloudavenue_iam_user.example.name
}
`

func TestAccUserDataSource(t *testing.T) {
	datasourceName := "data.cloudavenue_iam_user.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: ConcatTests(TestAccOrgUserResourceConfig, testAccUserDataSourceConfig),
				Check:  testsOrgUserResourceConfig(datasourceName),
			},
			{
				Config: ConcatTests(testAccOrgUserResourceConfigFull, testAccUserDataSourceConfig),
				Check:  testsOrgUserResourceConfigFull(datasourceName, true),
			},
		},
	})
}
