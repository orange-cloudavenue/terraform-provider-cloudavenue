// Package vdc provides the acceptance tests for the provider.
package vdc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccVDCDataSourceConfig = `
data "cloudavenue_vdc" "test" {
	name = "VDC_Frangipane"
}
`

func TestAccVDCDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vdc.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
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
