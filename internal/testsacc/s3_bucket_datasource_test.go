package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketDatasource{}

const (
	S3BucketDatasourceName = testsacc.ResourceName("data.cloudavenue_s3_bucket")
)

type S3BucketDatasource struct{}

func NewS3BucketDatasourceTest() testsacc.TestACC {
	return &S3BucketDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketDatasource) GetResourceName() string {
	return S3BucketDatasourceName.String()
}

func (r *S3BucketDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[S3BucketResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketDatasource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_bucket" "example" {
						name = cloudavenue_s3_bucket.example.name
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[S3BucketResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccS3BucketDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketDatasource{}),
	})
}
