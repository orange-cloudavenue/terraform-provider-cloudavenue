package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VAppIsolatedNetworkResource{}

const (
	VAppIsolatedNetworkResourceName = testsacc.ResourceName("cloudavenue_vapp_isolated_network")
)

type VAppIsolatedNetworkResource struct{}

func NewVAppIsolatedNetworkResourceTest() testsacc.TestACC {
	return &VAppIsolatedNetworkResource{}
}

// GetResourceName returns the name of the resource.
func (r *VAppIsolatedNetworkResource) GetResourceName() string {
	return VAppIsolatedNetworkResourceName.String()
}

func (r *VAppIsolatedNetworkResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *VAppIsolatedNetworkResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "vdc"),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vapp_id"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VAppResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vapp_isolated_network" "example" {
						name                  = {{ generate . "name" }}
						vdc                   = cloudavenue_vdc.example.name
						vapp_name             = cloudavenue_vapp.example.name
						gateway               = "192.168.10.1"
						netmask               = "255.255.255.0"
						dns1                  = "192.168.10.1"
						dns2                  = "192.168.10.3"
						dns_suffix            = "myvapp.biz"
						guest_vlan_allowed    = true
						retain_ip_mac_enabled = true

						static_ip_pool = [
							{
								start_address = "192.168.10.51"
								end_address   = "192.168.10.101"
							},
							{
								start_address = "192.168.10.10"
								end_address   = "192.168.10.30"
							}
						]
					}
					`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.10.1"),
						resource.TestCheckResourceAttr(resourceName, "netmask", "255.255.255.0"),
						resource.TestCheckResourceAttr(resourceName, "dns1", "192.168.10.1"),
						resource.TestCheckResourceAttr(resourceName, "dns2", "192.168.10.3"),
						resource.TestCheckResourceAttr(resourceName, "dns_suffix", "myvapp.biz"),
						resource.TestCheckResourceAttr(resourceName, "guest_vlan_allowed", "true"),
						resource.TestCheckResourceAttr(resourceName, "retain_ip_mac_enabled", "true"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
							"start_address": "192.168.10.51",
							"end_address":   "192.168.10.101",
						}),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
							"start_address": "192.168.10.10",
							"end_address":   "192.168.10.30",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vapp_isolated_network" "example" {
							name                  = {{ get . "name" }}
							description           = {{ generate . "description" }}
							vdc                   = cloudavenue_vdc.example.name
							vapp_name             = cloudavenue_vapp.example.name
							gateway               = "192.168.10.1"
							netmask               = "255.255.255.0"
							dns1                  = "192.168.10.1"
							dns2                  = "192.168.10.3"
							dns_suffix            = "myvapp.biz"
							guest_vlan_allowed    = true
							retain_ip_mac_enabled = true

							static_ip_pool = [
								{
									start_address = "192.168.10.51"
									end_address   = "192.168.10.101"
								},
								{
									start_address = "192.168.10.10"
									end_address   = "192.168.10.30"
								},
								{
									start_address = "192.168.10.200"
									end_address   = "192.168.10.210"
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrSet(resourceName, "vdc"),
							resource.TestCheckResourceAttrSet(resourceName, "vapp_name"),
							resource.TestCheckResourceAttrSet(resourceName, "vapp_id"),
							resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.10.1"),
							resource.TestCheckResourceAttr(resourceName, "netmask", "255.255.255.0"),
							resource.TestCheckResourceAttr(resourceName, "dns1", "192.168.10.1"),
							resource.TestCheckResourceAttr(resourceName, "dns2", "192.168.10.3"),
							resource.TestCheckResourceAttr(resourceName, "dns_suffix", "myvapp.biz"),
							resource.TestCheckResourceAttr(resourceName, "guest_vlan_allowed", "true"),
							resource.TestCheckResourceAttr(resourceName, "retain_ip_mac_enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								"start_address": "192.168.10.51",
								"end_address":   "192.168.10.101",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								"start_address": "192.168.10.10",
								"end_address":   "192.168.10.30",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "static_ip_pool.*", map[string]string{
								"start_address": "192.168.10.200",
								"end_address":   "192.168.10.210",
							}),
						},
					},
				},
				// ! Import testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc", "vapp_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccVAppIsolatedNetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VAppIsolatedNetworkResource{}),
	})
}
