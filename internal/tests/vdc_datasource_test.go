// Package tests provides the acceptance tests for the provider.
package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccVDCDataSourceConfig = `
data "cloudavenue_vdc" "test" {
	name = "VDC_Frangipane"
}
`

func TestAccVDCDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vdc.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccVDCDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "VDC_Frangipane"),
				),
			},
		},
	})
}
