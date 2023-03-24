// Package iam provides the acceptance tests for the provider.
package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccUserDataSourceConfig = `
resource "cloudavenue_iam_user" "test" {
	name   = "testuser"
	role_name   = "Organization Administrator"
	password    = "Th!s1sSecur3P@ssword"
}

data "cloudavenue_iam_user" "test" {
	name = cloudavenue_iam_user.test.name
}
`

const testAccUserDataSourceConfigFull = `
resource "cloudavenue_iam_user" "test" {
	name              = "testuserfull"
	role_name         = "Organization Administrator"
	password          = "Th!s1sSecur3P@ssword"
	enabled           = true
	email             = "foo@bar.com"
	telephone         = "1234567890"
	full_name         = "Test User"
	take_ownership    = true
	deployed_vm_quota = 10
	stored_vm_quota   = 5
}

data "cloudavenue_iam_user" "test" {
	name = cloudavenue_iam_user.test.name
}
`

func TestAccUserDataSource(t *testing.T) {
	resourceName := "data.cloudavenue_iam_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccUserDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testuser"),
					resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccUserDataSourceConfigFull,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testuserfull"),
					resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email", "foo@bar.com"),
					resource.TestCheckResourceAttr(resourceName, "telephone", "1234567890"),
					resource.TestCheckResourceAttr(resourceName, "full_name", "Test User"),
					resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "10"),
					resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "5"),
				),
			},
		},
	})
}
