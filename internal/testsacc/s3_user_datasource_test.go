package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &S3UserDataSource{}

const (
	S3UserDataSourceName = testsacc.ResourceName("data.cloudavenue_s3_user")
)

type S3UserDataSource struct{}

func NewS3UserDataSourceTest() testsacc.TestACC {
	return &S3UserDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *S3UserDataSource) GetResourceName() string {
	return S3UserDataSourceName.String()
}

func (r *S3UserDataSource) DependenciesConfig() (deps testsacc.DependenciesConfigResponse) {
	deps.Append(AddConstantConfig(testAccOrgUserResourceConfigFull))
	return
}

func (r *S3UserDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_s3_user" "example" {
						user_name = cloudavenue_iam_user.example.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.User)),
						resource.TestCheckResourceAttrSet(resourceName, "user_id"),
						resource.TestCheckResourceAttrSet(resourceName, "user_name"),
						resource.TestCheckResourceAttrSet(resourceName, "full_name"),
						resource.TestCheckResourceAttrSet(resourceName, "canonical_id"),
					},
				},
			}
		},
	}
}

func TestAccS3UserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3UserDataSource{}),
	})
}
