package vcda

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccVCDAResourceConfig = `
resource "cloudavenue_vcda_ip" "example" {
	ip_address = "10.0.0.1"
}
`

func TestAccVCDAResource(t *testing.T) {
	const resourceName = "cloudavenue_vcda_ip.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccVCDAResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip_address", "10.0.0.1"),
				),
			},
		},
	})
}
