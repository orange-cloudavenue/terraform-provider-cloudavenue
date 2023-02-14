package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVdcDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vdc.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccVdcDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "VDC_ECO_1"),
					resource.TestCheckResourceAttr(dataSourceName, "vdc_group", "myvDCGroup"),
					resource.TestCheckResourceAttr(dataSourceName, "memory_allocated", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "billing_model", "PAYG"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_profile.0.default", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_profile.0.limit", "500"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_profile.0.class", "gold"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_profile.1.default", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_profile.1.limit", "500"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_profile.1.class", "gold"),
					resource.TestCheckResourceAttr(dataSourceName, "cpu_allocated", "11000"),
					resource.TestCheckResourceAttr(dataSourceName, "disponibility_class", "ONE-ROOM"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_billing_model", "PAYG"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "VDC_ECO_1"),
					resource.TestCheckResourceAttr(dataSourceName, "cpu_speed_in_mhz", "2200"),
					resource.TestCheckResourceAttr(dataSourceName, "description", "Additionnal VDC for benchmarks"),
					resource.TestCheckResourceAttr(dataSourceName, "service_class", "ECO"),
				),
			},
		},
	})
}

const testAccVdcDataSourceConfig = `
data "cloudavenue_vdc" "test" {
	name = "VDC_Frangipane"
}
`
