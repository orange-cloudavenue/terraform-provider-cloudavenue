package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &IAMRolesDataSource{}

const (
	IAMRolesDataSourceName = testsacc.ResourceName("data.cloudavenue_iam_roles")
)

type IAMRolesDataSource struct{}

func NewIAMRolesDataSourceTest() testsacc.TestACC {
	return &IAMRolesDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *IAMRolesDataSource) GetResourceName() string {
	return IAMRolesDataSourceName.String()
}

func (r *IAMRolesDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *IAMRolesDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `data "cloudavenue_iam_roles" "example" {}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "roles.*", map[string]string{
							"name":        "Organization Administrator",
							"description": "Built-in rights for administering an organization",
							"read_only":   "true",
						}),
					},
				},
			}
		},
	}
}

func TestAccIAMRolesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&IAMRolesDataSource{}),
	})
}
