package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccAlbPoolDataSourceConfig = `
data "cloudavenue_alb_pool" "example" {
}
`

func TestAccAlbPoolDataSource(t *testing.T) {
	const dataSourceName = "data.cloudavenue_alb_pool.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccAlbPoolDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}
