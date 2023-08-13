package testsacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccACLResourceConfig = `
resource "cloudavenue_catalog_acl" "example" {
	catalog_id = cloudavenue_catalog.example.id
	shared_with_everyone = true
	everyone_access_level = "ReadOnly"
}
`

const testAccACLResourceConfigUpdate = `
resource "cloudavenue_catalog_acl" "example" {
	catalog_id = cloudavenue_catalog.example.id
	shared_with_everyone = true
	everyone_access_level = "FullControl"
}
`

const testAccACLResourceConfigUpdateShareWithUsers = `
resource "cloudavenue_iam_user" "example" {
	name        = "example"
	role_name   = "Organization Administrator"
	password    = "Th!s1sSecur3P@ssword"
}

resource "cloudavenue_iam_user" "example2" {
	name        = "example2"
	role_name   = "Organization Administrator"
	password    = "Th!s1sSecur3P@ssword"
}

resource "cloudavenue_catalog_acl" "example" {
	catalog_id = cloudavenue_catalog.example.id
	shared_with_everyone = false
	shared_with_users = [
		{
			user_id = cloudavenue_iam_user.example.id
			access_level = "ReadOnly"
		},
		{
			user_id = cloudavenue_iam_user.example2.id
			access_level = "FullControl"
		}
	]
}
`

var (
	aclTestCheck = func(resourceName string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Catalog)),
			resource.TestCheckResourceAttr(resourceName, "shared_with_everyone", "true"),
			resource.TestCheckResourceAttr(resourceName, "everyone_access_level", "ReadOnly"),
			resource.TestCheckNoResourceAttr(resourceName, "shared_with_users.#"),
			resource.TestCheckResourceAttrWith(resourceName, "catalog_id", uuid.TestIsType(uuid.Catalog)),
			resource.TestCheckResourceAttrSet(resourceName, "catalog_name"),
		)
	}

	aclTestCheckShareWithUsers = func(resourceName string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Catalog)),
			resource.TestCheckResourceAttr(resourceName, "shared_with_everyone", "false"),
			resource.TestCheckNoResourceAttr(resourceName, "everyone_access_level"),
			resource.TestCheckResourceAttrWith(resourceName, "catalog_id", uuid.TestIsType(uuid.Catalog)),
			resource.TestCheckResourceAttrSet(resourceName, "catalog_name"),
			resource.TestCheckResourceAttr(resourceName, "shared_with_users.#", "2"),
			// shared_with_users it's a SetNestedAttribute, so we can't be sure of the order of the elements in the list is not possible to test each attribute
		)
	}
)

func TestCatalogAccACLResource(t *testing.T) {
	resourceName := "cloudavenue_catalog_acl.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: ConcatTests(testAccCatalogResourceConfig, testAccACLResourceConfig),
				Check:  aclTestCheck(resourceName),
			},
			// Update testing
			{
				// Update test
				Config: ConcatTests(testAccCatalogResourceConfig, testAccACLResourceConfigUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Catalog)),
					resource.TestCheckResourceAttr(resourceName, "shared_with_everyone", "true"),
					resource.TestCheckResourceAttr(resourceName, "everyone_access_level", "FullControl"),
					resource.TestCheckNoResourceAttr(resourceName, "shared_with_users.#"),
					resource.TestCheckResourceAttrWith(resourceName, "catalog_id", uuid.TestIsType(uuid.Catalog)),
					resource.TestCheckResourceAttrSet(resourceName, "catalog_name"),
				),
			},
			// Update testing
			{
				Config: ConcatTests(testAccCatalogResourceConfig, testAccACLResourceConfigUpdateShareWithUsers),
				Check:  aclTestCheckShareWithUsers(resourceName),
			},
			// Import State testing
			{
				// Import test
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
