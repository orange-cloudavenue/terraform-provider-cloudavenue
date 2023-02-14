package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccVappDataSourceConfig = `
resource "cloudavenue_vapp" "example" {
	vapp_name = "vapp_name"
	description = "This is a test vapp"
}

data "cloudavenue_vapp" "test" {
	vapp_name = cloudavenue_vapp.example.vapp_name
}
`

func TestAccVappDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vapp.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccVappDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vapp_id"),
					resource.TestCheckResourceAttr(dataSourceName, "vapp_name", "vapp_name"),
					resource.TestCheckResourceAttr(dataSourceName, "status_text", "RESOLVED"),
					resource.TestCheckResourceAttr(dataSourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),
					resource.TestCheckResourceAttrSet(dataSourceName, "href"),
				),
			},
		},
	})
}
