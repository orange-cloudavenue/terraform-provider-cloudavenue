// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &VDCDataSource{}

const (
	VDCDataSourceName = testsacc.ResourceName("data.cloudavenue_vdc")
)

type VDCDataSource struct{}

func NewVDCDataSourceTest() testsacc.TestACC {
	return &VDCDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCDataSource) GetResourceName() string {
	return VDCDataSourceName.String()
}

func (r *VDCDataSource) DependenciesConfig() (configs testsacc.TFData) {
	configs.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig())
	return
}

func (r *VDCDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdc" "example" {
						name = cloudavenue_vdc.example.name
					}`,
					// TFConfig: testsacc.GenerateFromTemplate("cloudavenue_vdc",`
					// data "cloudavenue_vdc" "example" {
					// 	name = {{ get . "name" }}
					// }`,
					// Checks: NewVDCResourceTest().Tests(ctx)["example"](ctx, resourceName).GenerateCheckWithCommonChecks(),
					Checks: []resource.TestCheckFunc{
						resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.VDC.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrSet(resourceName, "description"),
						resource.TestCheckResourceAttrSet(resourceName, "service_class"),
						resource.TestCheckResourceAttrSet(resourceName, "disponibility_class"),
						resource.TestCheckResourceAttrSet(resourceName, "billing_model"),
						resource.TestCheckResourceAttrSet(resourceName, "cpu_speed_in_mhz"),
						resource.TestCheckResourceAttrSet(resourceName, "cpu_allocated"),
						resource.TestCheckResourceAttrSet(resourceName, "memory_allocated"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_billing_model"),
						resource.TestCheckResourceAttr(resourceName, "storage_profiles.#", "1"),
					},
				},
			}
		},
	}
}

func TestVDCDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCDataSource{}),
	})
}
