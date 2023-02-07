package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccJobsDataSourceConfig = `
data "cloudavenue_jobs" "test" {
	id = "fb064495-457d-40d4-8e53-79fe3824ca96"
	}
`

func TestAccJobsDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_jobs.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccJobsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr(dataSourceName, "id", "fb064495-457d-40d4-8e53-79fe3824ca96"),
				),
			},
		},
	})
}
