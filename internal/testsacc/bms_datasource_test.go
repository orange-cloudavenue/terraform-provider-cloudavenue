package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &BMSDataSource{}

const (
	BMSDataSourceName = testsacc.ResourceName("data.cloudavenue_bms")
)

type BMSDataSource struct{}

func NewBMSDataSourceTest() testsacc.TestACC {
	return &BMSDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *BMSDataSource) GetResourceName() string {
	return BMSDataSourceName.String()
}

func (r *BMSDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *BMSDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_bms" "example" {}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
					},
				},
			}
		},
	}
}

func TestAccBMSDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&BMSDataSource{}),
	})
}
