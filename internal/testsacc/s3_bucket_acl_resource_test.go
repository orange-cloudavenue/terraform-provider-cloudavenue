package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketACLResource{}

const (
	S3BucketACLResourceName = testsacc.ResourceName("cloudavenue_s3_bucket_acl")
)

type S3BucketACLResource struct{}

func NewS3BucketACLResourceTest() testsacc.TestACC {
	return &S3BucketACLResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketACLResource) GetResourceName() string {
	return S3BucketACLResourceName.String()
}

func (r *S3BucketACLResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketACLResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// TODO : Complete tests
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_s3_bucket_acl" "example" {
						bucket = cloudavenue_s3_bucket.example.name
						acl = "public-read"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policies.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policies.0.grants.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policies.0.grants.1.permission", "READ"),
						resource.TestCheckResourceAttr(resourceName, "access_control_policies.0.grants.1.grantee.uri", "http://acs.amazonaws.com/groups/global/AllUsers"),
						resource.TestCheckResourceAttrSet(resourceName, "access_control_policies.0.owner.id"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_acl" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							acl = "public-read"
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policies.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policies.0.grants.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policies.0.grants.1.permission", "READ"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policies.0.grants.1.grantee.uri", "http://acs.amazonaws.com/groups/global/AllUsers"),
							resource.TestCheckResourceAttrSet(resourceName, "access_control_policies.0.owner.id"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_acl" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							acl = "private"
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "acl", "private"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policies.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policies.0.grants.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policies.0.grants.0.permission", "FULL_CONTROL"),
							resource.TestCheckResourceAttr(resourceName, "access_control_policies.0.grants.0.grantee.type", "CanonicalUser"),
							resource.TestCheckResourceAttrSet(resourceName, "access_control_policies.0.owner.id"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder:    []string{"bucket"},
						ImportStateVerifyIgnore: []string{"acl"},
						ImportState:             true,
						ImportStateVerify:       true,
					},
				},
			}
		},
		// It's possible to add multiple tests

		// Complete and functional example :
		/*
			// * Test One (example)
			"example": func(_ context.Context, resourceName string) testsacc.Test {
				return testsacc.Test{
					CommonChecks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "catalog_name"),

						resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Catalog)),
						resource.TestCheckResourceAttrWith(resourceName, "catalog_id", uuid.TestIsType(uuid.Catalog)),
					},
					// ! Create testing
					Create: testsacc.TFConfig{
						TFConfig: `
						resource "cloudavenue_catalog_ACL" "example" {
							catalog_id = cloudavenue_catalog.example.id
							shared_with_everyone = true
							everyone_access_level = "ReadOnly"
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckNoResourceAttr(resourceName, "shared_with_users.#"),

							resource.TestCheckResourceAttr(resourceName, "everyone_access_level", "ReadOnly"),
							resource.TestCheckResourceAttr(resourceName, "shared_with_everyone", "true"),
						},
					},
					// ! Updates testing
					Updates: []testsacc.TFConfig{
						{
							TFConfig: `
							resource "cloudavenue_catalog_ACL" "example" {
								catalog_id = cloudavenue_catalog.example.id
								shared_with_everyone = true
								everyone_access_level = "FullControl"
							}`,
							Checks: []resource.TestCheckFunc{
								resource.TestCheckNoResourceAttr(resourceName, "shared_with_users.#"),

								resource.TestCheckResourceAttr(resourceName, "everyone_access_level", "FullControl"),
								resource.TestCheckResourceAttr(resourceName, "shared_with_everyone", "true"),
							},
						},
						{
							TFConfig: `
							resource "cloudavenue_catalog_ACL" "example" {
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
							}`,
							Checks: []resource.TestCheckFunc{
								resource.TestCheckNoResourceAttr(resourceName, "everyone_access_level"),

								resource.TestCheckResourceAttr(resourceName, "shared_with_everyone", "false"),
								resource.TestCheckResourceAttr(resourceName, "shared_with_users.#", "2"),

								resource.TestCheckResourceAttrWith(resourceName, "shared_with_users.0.user_id", uuid.TestIsType(uuid.User)),
								resource.TestCheckResourceAttrWith(resourceName, "shared_with_users.1.user_id", uuid.TestIsType(uuid.User)),
								// shared_with_users it's a SetNestedAttribute, so we can't be sure of the order of the elements in the list is not possible to test each attribute
							},
						},
					},
					// ! Imports testing
					Imports: []testsacc.TFImport{
						{
							ImportStateIDBuilder: []string{"id"},
							ImportState:          true,
							ImportStateVerify:    true,
						},
					},
				}
			},
		*/
	}
}

func TestAccS3BucketACLResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketACLResource{}),
	})
}
