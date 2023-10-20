package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketVersioningConfigurationResource{}

const (
	S3BucketVersioningConfigurationResourceName = testsacc.ResourceName("cloudavenue_s3_bucket_versioning_configuration")
)

type S3BucketVersioningConfigurationResource struct{}

func NewS3BucketVersioningConfigurationResourceTest() testsacc.TestACC {
	return &S3BucketVersioningConfigurationResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketVersioningConfigurationResource) GetResourceName() string {
	return S3BucketVersioningConfigurationResourceName.String()
}

func (r *S3BucketVersioningConfigurationResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketVersioningConfigurationResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_s3_bucket_versioning_configuration" "example" {
						bucket = cloudavenue_s3_bucket.example.name
						status = "Enabled"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "status", "Enabled"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_versioning_configuration" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							status = "Suspended"
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "status", "Suspended"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_versioning_configuration" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							status = "Enabled"
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "status", "Enabled"),
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

func TestAccS3BucketVersioningConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketVersioningConfigurationResource{}),
	})
}
