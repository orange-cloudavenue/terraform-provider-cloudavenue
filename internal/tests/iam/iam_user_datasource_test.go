// Package iam provides the acceptance tests for the provider.
package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

const testAccOrgUserDataSourceConfig = `
resource "cloudavenue_iam_user" "test" {
	user_name   = "testuser"
	description = "A test user"
	role        = "Organization Administrator"
	password    = "Th!s1sSecur3P@ssword"
}

data "cloudavenue_iam_user" "test" {
	user_name = cloudavenue_iam_user.test.user_name
}
`

const testAccOrgUserDataSourceConfigFull = `
resource "cloudavenue_iam_user" "test" {
	user_name         = "testuserfull"
	description       = "A test user"
	role              = "Organization Administrator"
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
	user_name = cloudavenue_iam_user.test.user_name
}
`

func TestAccOrgUserDataSource(t *testing.T) {
	resourceName := "cloudavenue_iam_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOrgUserDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_name", "testuser"),
					resource.TestCheckResourceAttr(resourceName, "description", "A test user"),
					resource.TestCheckResourceAttr(resourceName, "role", "Organization Administrator"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccOrgUserDataSourceConfigFull,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_name", "testuserfull"),
					resource.TestCheckResourceAttr(resourceName, "description", "A test user"),
					resource.TestCheckResourceAttr(resourceName, "role", "Organization Administrator"),
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
