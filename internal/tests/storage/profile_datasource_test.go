package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccProfileDataSourceConfig = `
data "cloudavenue_storage_profile" "example" {
	name = "gold"
}
`

func TestAccProfileDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_storage_profile.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccProfileDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "limit"),
					resource.TestCheckResourceAttrSet(dataSourceName, "used_storage"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "iops_allocated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "units"),
					resource.TestCheckResourceAttrSet(dataSourceName, "iops_limiting_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "maximum_disk_iops"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default_disk_iops"),
					resource.TestCheckResourceAttrSet(dataSourceName, "disk_iops_per_gb_max"),
					resource.TestCheckResourceAttrSet(dataSourceName, "iops_limit"),
				),
			},
		},
	})
}
