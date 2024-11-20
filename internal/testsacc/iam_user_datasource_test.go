package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
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
	return
}

func (r *IAMUserDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[IAMUserResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_iam_user" "example" {
						name = cloudavenue_iam_user.example.name
					}`,
					Checks: GetResourceConfig()[IAMUserResourceName]().GetDefaultChecks(),
				},
			}
		},
		"example_saml_user": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[IAMUserSAMLResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_iam_user" "example_saml_user" {
						name = cloudavenue_iam_user_saml.example.user_name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.User)),
						resource.TestCheckResourceAttr(resourceName, "name", "mickael.stanislas.ext"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "0"),
						resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "0"),
						resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
					},
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
