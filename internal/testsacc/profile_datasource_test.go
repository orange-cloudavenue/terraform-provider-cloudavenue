package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccProfileDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(dataSourceName, "id", uuid.TestIsType(uuid.VDCStorageProfile)),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdc"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "gold"),
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
