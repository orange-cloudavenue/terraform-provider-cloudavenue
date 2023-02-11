package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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
					resource.TestCheckResourceAttr(dataSourceName, "id", "urn:vcloud:vapp:83ce9843-f84c-4a72-8c16-3ec12835c4b4"),
					resource.TestCheckResourceAttr(dataSourceName, "vapp_id", "urn:vcloud:vapp:83ce9843-f84c-4a72-8c16-3ec12835c4b4"),
					resource.TestCheckResourceAttr(dataSourceName, "vapp_name", "vapp_test"),
					resource.TestCheckResourceAttr(dataSourceName, "status_text", "RESOLVED"),
					resource.TestCheckResourceAttr(dataSourceName, "vdc", os.Getenv("CLOUDAVENUE_VDC")),

					resource.TestCheckResourceAttr(dataSourceName, "href", os.Getenv("CLOUDAVENUE_URL")+"/api/vApp/vapp-83ce9843-f84c-4a72-8c16-3ec12835c4b4"),
				),
			},
		},
	})
}

const testAccVappDataSourceConfig = `
data "cloudavenue_vapp" "test" {
	name = "vapp_test"
}
`
