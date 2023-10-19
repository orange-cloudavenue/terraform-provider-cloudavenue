package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketOwnershipControlsDataSource{}

const (
	S3BucketOwnershipControlsDataSourceName = testsacc.ResourceName("data.cloudavenue_s3_bucket_ownership_controls")
)

type S3BucketOwnershipControlsDataSource struct{}

func NewS3BucketOwnershipControlsDataSourceTest() testsacc.TestACC {
	return &S3BucketOwnershipControlsDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketOwnershipControlsDataSource) GetResourceName() string {
	return S3BucketOwnershipControlsDataSourceName.String()
}

func (r *S3BucketOwnershipControlsDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[S3BucketOwnershipControlsResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketOwnershipControlsDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_bucket_ownership_controls" "example" {
						bucket = cloudavenue_s3_bucket_ownership_controls.example.id
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[S3BucketOwnershipControlsResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccS3BucketOwnershipControlsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketOwnershipControlsDataSource{}),
	})
}
