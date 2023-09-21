// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const TestAccEVAppResourceConfig = `
resource "cloudavenue_vapp" "example" {
	name        = "MyVapp"
	description = "This is an example vApp"

	lease = {
		runtime_lease_in_sec = 3600
		storage_lease_in_sec = 3600
	}

	guest_properties = {
		"key" = "Value"
	}
  }
`

var _ testsacc.TestACC = &VAppResource{}

const (
	VAppResourceName = testsacc.ResourceName("cloudavenue_vapp")
)

type VAppResource struct{}

func NewVAppResourceTest() testsacc.TestACC {
	return &VAppResource{}
}

// GetResourceName returns the name of the resource.
func (r *VAppResource) GetResourceName() string {
	return VAppResourceName.String()
}

func (r *VAppResource) DependenciesConfig() (configs testsacc.TFData) {
	configs.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig())
	return
}

func (r *VAppResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VAPP)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vapp" "example" {
						name        = {{ generate . "name" }}
						description = {{ generate . "description" }}
						vdc 		= cloudavenue_vdc.example.name
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "lease.runtime_lease_in_sec", "0"),
						resource.TestCheckResourceAttr(resourceName, "lease.storage_lease_in_sec", "0"),
						resource.TestCheckNoResourceAttr(resourceName, "guest_properties.#"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vapp" "example" {
							name        = {{ get . "name" }}
							description = {{ generate . "description" }}
							vdc 		= cloudavenue_vdc.example.name

							lease = {
								runtime_lease_in_sec = 3600
								storage_lease_in_sec = 3600
							}
						
							guest_properties = {
								"key" = "Value"
							}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "lease.runtime_lease_in_sec", "3600"),
							resource.TestCheckResourceAttr(resourceName, "lease.storage_lease_in_sec", "3600"),
							resource.TestCheckResourceAttr(resourceName, "guest_properties.key", "Value"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccVAppResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VAppResource{}),
	})
}
