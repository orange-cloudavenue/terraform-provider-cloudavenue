// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccVDCsDataSourceConfig = `
data "cloudavenue_vdcs" "test" {
}
`

func TestAccVDCsDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vdcs.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccVDCsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(dataSourceName, "id", uuid.TestIsType(uuid.VDC)),
				),
			},
		},
	})
}
