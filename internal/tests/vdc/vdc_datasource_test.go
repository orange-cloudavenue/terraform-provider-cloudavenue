// Package vdc provides the acceptance tests for the provider.
package vdc

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccVDCDataSourceConfig = `
data "cloudavenue_vdc" "test" {
	name = "MyVDC"
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
					resource.TestMatchResourceAttr(dataSourceName, "id", regexp.MustCompile(uuid.VDC.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "service_class"),
					resource.TestCheckResourceAttrSet(dataSourceName, "disponibility_class"),
					resource.TestCheckResourceAttrSet(dataSourceName, "billing_model"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cpu_speed_in_mhz"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cpu_allocated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "memory_allocated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_billing_model"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_profiles.0.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdc_group"),
				),
			},
		},
	})
}
