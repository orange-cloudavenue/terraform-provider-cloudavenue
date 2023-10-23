package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketVersioningConfigurationDatasource{}

const (
	S3BucketVersioningConfigurationDatasourceName = testsacc.ResourceName("data.cloudavenue_s3_bucket_versioning_configuration")
)

type S3BucketVersioningConfigurationDatasource struct{}

func NewS3BucketVersioningConfigurationDatasourceTest() testsacc.TestACC {
	return &S3BucketVersioningConfigurationDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketVersioningConfigurationDatasource) GetResourceName() string {
	return S3BucketVersioningConfigurationDatasourceName.String()
}

func (r *S3BucketVersioningConfigurationDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[S3BucketVersioningConfigurationResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketVersioningConfigurationDatasource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_bucket_versioning_configuration" "example" {
						bucket = cloudavenue_s3_bucket.example.name
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[S3BucketVersioningConfigurationResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccS3BucketVersioningConfigurationDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketVersioningConfigurationDatasource{}),
	})
}
