package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &IAMUserDataSource{}

const (
	IAMUserDataSourceName = testsacc.ResourceName("data.cloudavenue_iam_user")
)

type IAMUserDataSource struct{}

func NewIAMUserDataSourceTest() testsacc.TestACC {
	return &IAMUserDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *IAMUserDataSource) GetResourceName() string {
	return IAMUserDataSourceName.String()
}

func (r *IAMUserDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[IAMUserResourceName]().GetDefaultConfig)
	return
}

func (r *IAMUserDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_iam_user" "example" {
							name = cloudavenue_iam_user.example.name
					}`,
					Checks: GetResourceConfig()[IAMUserResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccIAMUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&IAMUserDataSource{}),
	})
}
