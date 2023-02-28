package vdc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccVDCGroupDataSourceConfig = `
data "cloudavenue_vdc_group" "example" {
	name = "MyVDCGroup"
}
`

const dataSourceName = "data.cloudavenue_vdc_group.example"

func TestAccCloudavenueVdcGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccVDCGroupDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "urn:vcloud:vdcGroup:805a3699-7dba-419b-b731-daae4154617e"),
				),
			},
		},
	})
}
