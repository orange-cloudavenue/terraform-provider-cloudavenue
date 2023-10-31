package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketLifecycleConfigurationDataSource{}

const (
	S3BucketLifecycleConfigurationDataSourceName = testsacc.ResourceName("data.cloudavenue_s3_bucket_lifecycle_configuration")
)

type S3BucketLifecycleConfigurationDataSource struct{}

func NewS3BucketLifecycleConfigurationDataSourceTest() testsacc.TestACC {
	return &S3BucketLifecycleConfigurationDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketLifecycleConfigurationDataSource) GetResourceName() string {
	return S3BucketLifecycleConfigurationDataSourceName.String()
}

func (r *S3BucketLifecycleConfigurationDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[S3BucketLifecycleConfigurationResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketLifecycleConfigurationDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
						bucket = cloudavenue_s3_bucket_lifecycle_configuration.example.bucket
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[S3BucketLifecycleConfigurationResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccS3BucketLifecycleConfigurationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketLifecycleConfigurationDataSource{}),
	})
}
