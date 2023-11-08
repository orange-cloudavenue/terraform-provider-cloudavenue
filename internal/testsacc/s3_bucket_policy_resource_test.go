package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketPolicyResource{}

const (
	S3BucketPolicyResourceName = testsacc.ResourceName("cloudavenue_s3_bucket_policy")
	S3Bucket4Policy            = `
	resource "cloudavenue_s3_bucket" "example" {
		name = "example-bucket"
	}`
)

type S3BucketPolicyResource struct{}

func NewS3BucketPolicyResourceTest() testsacc.TestACC {
	return &S3BucketPolicyResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketPolicyResource) GetResourceName() string {
	return S3BucketPolicyResourceName.String()
}

func (r *S3BucketPolicyResource) DependenciesConfig() (deps testsacc.DependenciesConfigResponse) {
	// Add constant dependencies config to give the good path in resource json policy
	deps.Append(AddConstantConfig(S3Bucket4Policy))
	// deps.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketPolicyResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttr(resourceName, "id", "example-bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_s3_bucket_policy" "example" {
					bucket = cloudavenue_s3_bucket.example.name
					policy = jsonencode({
						Version = "2012-10-17"
						Statement = [
						  {
						    Effect = "Allow"
						    Principal = "*"
						    Action = [
						      "s3:*"
						    ]
						    Resource = [
						      "arn:aws:s3:::example-bucket/*"
						    ]
						  }
						]
					})
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "policy"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_policy" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							policy = jsonencode({
								Version = "2012-10-17"
								Statement = [
								  {
									Effect = "Allow"
									Principal = "*"
									Action = [
									  "s3:*"
									]
									Resource = [
									  "arn:aws:s3:::example-bucket/*"
									]
								  }
								]
							  })
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "policy"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_s3_bucket_policy" "example" {
							bucket = cloudavenue_s3_bucket.example.name
							policy = jsonencode({
								Version = "2012-10-17"
								Statement = [
								  {
									Effect = "Allow"
									Principal = "*"
									Action = [ 
									  "s3:DeleteBucket",
									  "s3:GetObject",
									  "s3:ListBucketVersions",
									]
									Resource = [
									  "arn:aws:s3:::example-bucket",
									  "arn:aws:s3:::example-bucket/*",
									]
								  }
								]
							  })
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "policy"),
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

func TestAccS3BucketPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketPolicyResource{}),
	})
}
