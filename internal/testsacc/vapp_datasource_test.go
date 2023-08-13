// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccVappDataSourceConfig = `
resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example vApp"
  }

data "cloudavenue_vapp" "test" {
	name = cloudavenue_vapp.example.name
}
`

func TestAccVappDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vapp.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccVappDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "MyVapp"),
					resource.TestCheckResourceAttr(dataSourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
				),
			},
		},
	})
}
