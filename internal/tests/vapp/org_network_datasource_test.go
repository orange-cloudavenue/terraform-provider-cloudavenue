package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccOrgNetworkDataSourceConfig = `
data "cloudavenue_vapp_org_network" "example" {
}
`


func TestAccOrgNetworkDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_vapp_org_network.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccOrgNetworkResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
