package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &bms_DatasourceDataSource{}

const (
	bms_DatasourceDataSourceName = testsacc.ResourceName("data.cloudavenue_bms_datasource")
)

type bms_DatasourceDataSource struct{}

func Newbms_DatasourceDataSourceTest() testsacc.TestACC {
	return &bms_DatasourceDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *bms_DatasourceDataSource) GetResourceName() string {
	return bms_DatasourceDataSourceName.String()
}

func (r *bms_DatasourceDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[bms_DatasourceResourceName]().GetDefaultConfig),
	return
}

func (r *bms_DatasourceDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_bms_datasource" "example" {
						foo_id = cloudavenue_foo_bar.example.id
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[bms_DatasourceResourceName]().GetDefaultChecks()
				},
			}
		},
	}
}

func TestAccbms_DatasourceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&bms_DatasourceDataSource{}),
	})
}