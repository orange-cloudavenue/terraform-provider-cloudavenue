package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &VAppOrgNetworkResource{}

const (
	VAppOrgNetworkResourceName = testsacc.ResourceName("cloudavenue_vapp_org_network")
)

type VAppOrgNetworkResource struct{}

func NewVAppOrgNetworkResourceTest() testsacc.TestACC {
	return &VAppOrgNetworkResource{}
}

// GetResourceName returns the name of the resource.
func (r *VAppOrgNetworkResource) GetResourceName() string {
	return VAppOrgNetworkResourceName.String()
}

func (r *VAppOrgNetworkResource) DependenciesConfig() (configs testsacc.TFData) {
	configs.Append(GetResourceConfig()[VAppResourceName]().GetDefaultConfig())
	configs.Append(GetResourceConfig()[NetworkRoutedResourceName]().GetDefaultConfig())
	return
}

func (r *VAppOrgNetworkResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Network)),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
					// TODO : https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/527
					// resource.TestCheckResourceAttrSet(resourceName, "vapp_id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vapp_org_network" "example" {
						vapp_name    = cloudavenue_vapp.example.name
						network_name = cloudavenue_network_routed.example.name
						vdc          = cloudavenue_vdc.example.name
					  }`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "is_fenced", "false"),
						resource.TestCheckResourceAttr(resourceName, "retain_ip_mac_enabled", "false"),
					},
				},
				// TODO : Impossible to update due to https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/531 and https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/529
				// // ! Updates testing
				// Updates: []testsacc.TFConfig{
				// 	{
				// 		TFConfig: testsacc.GenerateFromTemplate(resourceName, `
				// 		resource "cloudavenue_vapp_org_network" "example" {
				// 			vapp_name    = cloudavenue_vapp.example.name
				// 			network_name = cloudavenue_network_routed.example.name
				// 			vdc          = cloudavenue_vdc.example.name

				// 			retain_ip_mac_enabled = true
				// 		  }`),
				// 		Checks: []resource.TestCheckFunc{
				// 			resource.TestCheckResourceAttr(resourceName, "retain_ip_mac_enabled", "true"),
				// 		},
				// 	},
				// },
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc", "vapp_name", "network_name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
	// TODO: ADD Test with VDC Group
}

func TestAccOrgNetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VAppOrgNetworkResource{}),
	})
}
