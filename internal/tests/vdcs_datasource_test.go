package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVdcsDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vdcs.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccVdcsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "frangipane"),
					resource.TestCheckResourceAttr(dataSourceName, "vdcs.0.vdc_name", "vdc01"),
					resource.TestCheckResourceAttr(dataSourceName, "vdcs.0.vdc_uuid", "1bb33b06-8a3b-4a2c-a077-d4771794cf3c"),
					resource.TestCheckResourceAttr(dataSourceName, "vdcs.1.vdc_name", "vdc02"),
					resource.TestCheckResourceAttr(dataSourceName, "vdcs.1.vdc_uuid", "9b1243ea-ca4f-4a20-9ff4-844f0584b3e7"),
				),
			},
		},
	})
}

const testAccVdcsDataSourceConfig = `
data "cloudavenue_vdcs" "test" {
}
`
