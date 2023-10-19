package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketOwnershipControlsResource{}

const (
	S3BucketOwnershipControlsResourceName = testsacc.ResourceName("cloudavenue_s3_bucket_ownership_controls")
)

type S3BucketOwnershipControlsResource struct{}

func NewS3BucketOwnershipControlsResourceTest() testsacc.TestACC {
	return &S3BucketOwnershipControlsResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketOwnershipControlsResource) GetResourceName() string {
	return S3BucketOwnershipControlsResourceName.String()
}

func (r *S3BucketOwnershipControlsResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketOwnershipControlsResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// TODO : Complete tests
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_s3_bucket_ownership_controls" "example" {
						bucket = cloudavenue_s3_bucket.example.name
						rule = {
							object_ownership = "BucketOwnerPreferred"
						}
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "rule.object_ownership", "BucketOwnerPreferred"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_ownership_controls" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							rule = [{
								object_ownership = "BucketOwnerPreferred"
							}]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rule.object_ownership", "BucketOwnerPreferred"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_ownership_controls" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							rule = [{
								object_ownership = "ObjectWriter"
							}]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rule.object_ownership", "ObjectWriter"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"bucket"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccS3BucketOwnershipControlsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketOwnershipControlsResource{}),
	})
}
