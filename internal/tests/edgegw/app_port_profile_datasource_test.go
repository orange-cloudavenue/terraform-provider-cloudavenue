package edgegw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -test
const testAccPortProfileDataSourceConfig = `
data "cloudavenue_app_port_profile" "example" {
}
`

func TestAccPortProfileDataSource(t *testing.T) {
	dataSourceName := "data.cloudavenue_app_port_profile.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: tests.ConcatTests(testAccPortProfileResourceConfig, testAccPortProfileDataSourceConfig),
				Check: portProfileTestCheck(dataSourceName),
			},
		},
	})
}
