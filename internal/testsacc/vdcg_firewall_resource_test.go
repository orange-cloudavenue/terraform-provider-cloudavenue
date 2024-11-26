package testsacc

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGFirewallResource{}

const (
	VDCGFirewallResourceName = testsacc.ResourceName("cloudavenue_vdcg_firewall")
)

type VDCGFirewallResource struct{}

func NewVDCGFirewallResourceTest() testsacc.TestACC {
	return &VDCGFirewallResource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGFirewallResource) GetResourceName() string {
	return VDCGFirewallResourceName.String()
}

func (r *VDCGFirewallResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
	return
}

func (r *VDCGFirewallResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDCGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_firewall" "example" {
						vdc_group_name = cloudavenue_vdcg.example.name
						  rules = [
							{
							action      = "ALLOW"
							name        = "allow all IPv4 traffic"
							direction   = "IN_OUT"
							ip_protocol = "IPV4"
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
							"action":      "ALLOW",
							"name":        "allow all IPv4 traffic",
							"direction":   "IN_OUT",
							"ip_protocol": "IPV4",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Test invalid direction
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example" {
							vdc_group_name = cloudavenue_vdcg.example.name
							enabled = true
							rules = [
								{
									action      = "ALLOW"
									name        = "allow out IPv4 traffic"
									direction   = "OUT"
									ip_protocol = "IPV4"
								},
								{
									action      = "REJECT"
									name        = "reject in IPv4 traffic"
									direction   = "INN"
									ip_protocol = "IPV4"
									enabled     = false
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "REJECT",
								"name":        "reject in IPv4 traffic",
								"direction":   "INN",
								"ip_protocol": "IPV4",
								"enabled":     "false",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "ALLOW",
								"name":        "allow out IPv4 traffic",
								"direction":   "OUT",
								"ip_protocol": "IPV4",
								"enabled":     "true",
							}),
						},
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							PlanOnly:           true,
							ExpectError:        regexp.MustCompile(`Invalid Attribute Value Match`),
						},
					},
					// * Test invalid action
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example" {
							vdc_group_name = cloudavenue_vdcg.example.name
							enabled = true
							  rules = [
								{
									action      = "DROUP"
									name        = "allow out IPv4 traffic"
									direction   = "OUT"
									ip_protocol = "IPV4"
								},
								{
									action      = "DROUP"
									name        = "drop in IPv4 traffic"
									direction   = "IN"
									ip_protocol = "IPV4"
									enabled     = false
								}
							]
						}`),
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							PlanOnly:           true,
							ExpectError:        regexp.MustCompile(`Invalid Attribute Value Match`),
						},
					},
					// * Test invalid ip_protocol
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example" {
							vdc_group_name = cloudavenue_vdcg.example.name
							enabled = true
							  rules = [
								{
									action      = "ALLOW"
									name        = "allow out IPv4 traffic"
									direction   = "OUT"
									ip_protocol = "IPV5"
								},
								{
									action      = "REJECT"
									name        = "reject in IPv4 traffic"
									direction   = "IN"
									ip_protocol = "IPV4"
									enabled     = false
								}
							]
						}`),
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							PlanOnly:           true,
							ExpectError:        regexp.MustCompile(`Invalid Attribute Value Match`),
						},
					},
					// * Test to disable the firewall
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example" {
							vdc_group_name = cloudavenue_vdcg.example.name
							enabled = false
							  rules = [
								{
								action      = "ALLOW"
								name        = "allow all IPv4 traffic"
								direction   = "IN_OUT"
								ip_protocol = "IPV4"
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "ALLOW",
								"name":        "allow all IPv4 traffic",
								"direction":   "IN_OUT",
								"ip_protocol": "IPV4",
							}),
						},
					},
					// * Test to add a new rules
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example" {
							vdc_group_name = cloudavenue_vdcg.example.name
							enabled = true
							  rules = [
								{
									action      = "ALLOW"
									name        = "allow out IPv4 traffic"
									direction   = "OUT"
									ip_protocol = "IPV4"
								},
								{
									action      = "REJECT"
									name        = "reject in IPv4 traffic"
									direction   = "IN"
									ip_protocol = "IPV4"
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "REJECT",
								"name":        "reject in IPv4 traffic",
								"direction":   "IN",
								"ip_protocol": "IPV4",
								"enabled":     "true",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "ALLOW",
								"name":        "allow out IPv4 traffic",
								"direction":   "OUT",
								"ip_protocol": "IPV4",
								"enabled":     "true",
							}),
						},
					},
					// * Test to disable one rule
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example" {
							vdc_group_name = cloudavenue_vdcg.example.name
							enabled = true
							  rules = [
								{
									action      = "ALLOW"
									name        = "allow out IPv4 traffic"
									direction   = "OUT"
									ip_protocol = "IPV4"
								},
								{
									action      = "REJECT"
									name        = "reject in IPv4 traffic"
									direction   = "IN"
									ip_protocol = "IPV4"
									enabled     = false
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "REJECT",
								"name":        "reject in IPv4 traffic",
								"direction":   "IN",
								"ip_protocol": "IPV4",
								"enabled":     "false",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "ALLOW",
								"name":        "allow out IPv4 traffic",
								"direction":   "OUT",
								"ip_protocol": "IPV4",
								"enabled":     "true",
							}),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group_name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},

		// TODO migrate to cloudavenue_vdcg_app_port_profile
		"example_with_app_port_profile": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDCGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetDataSourceConfig()[EdgeGatewayAppPortProfileDatasourceName]().GetSpecificConfig("example_provider_scope"))
					resp.Append(GetDataSourceConfig()[EdgeGatewayAppPortProfileDatasourceName]().GetSpecificConfig("example_system_scope"))
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_firewall" "example_with_app_port_profile" {
						vdc_group_id = cloudavenue_vdcg.example.id
						  rules = [
							{
								action      = "ALLOW"
								name        = "allow all IPv4 traffic"
								direction   = "IN_OUT"
								app_port_profile_ids = [
									data.cloudavenue_edgegateway_app_port_profile.example_provider_scope.id,
								]
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
							"action":                 "ALLOW",
							"name":                   "allow all IPv4 traffic",
							"direction":              "IN_OUT",
							"app_port_profile_ids.#": "1",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example_with_app_port_profile" {
							vdc_group_id = cloudavenue_vdcg.example.id
							enabled = true
							  rules = [
								{
									action      = "ALLOW"
									name        = "allow out IPv4 traffic"
									direction   = "OUT"
									app_port_profile_ids = [
										data.cloudavenue_edgegateway_app_port_profile.example_provider_scope.id,
										data.cloudavenue_edgegateway_app_port_profile.example_system_scope.id
									]
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":                 "ALLOW",
								"name":                   "allow out IPv4 traffic",
								"direction":              "OUT",
								"app_port_profile_ids.#": "2",
							}),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group_name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},

		"example_with_source_ids": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDCGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGIPSetResourceName]().GetDefaultConfig)
					resp.Append(GetResourceConfig()[VDCGSecurityGroupResourceName]().GetDefaultConfig)
					// TODO Add Dynamic security Group
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_firewall" "example_with_source_ids" {
						vdc_group_name = cloudavenue_vdcg.example.name
						  rules = [
							{
								action      = "ALLOW"
								name        = "allow all IPv4 traffic"
								direction   = "IN_OUT"
								source_ids = [
									cloudavenue_vdcg_ip_set.example.id,
								],
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
							"action":       "ALLOW",
							"name":         "allow all IPv4 traffic",
							"direction":    "IN_OUT",
							"source_ids.#": "1",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example_with_source_ids" {
							vdc_group_name = cloudavenue_vdcg.example.name
							enabled = true
							  rules = [
								{
									action      = "ALLOW"
									name        = "allow in IPv4 traffic"
									direction   = "IN"
									source_ids = [
										cloudavenue_vdcg_ip_set.example.id,
										cloudavenue_vdcg_security_group.example.id
									]
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":       "ALLOW",
								"name":         "allow in IPv4 traffic",
								"direction":    "IN",
								"source_ids.#": "2",
							}),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group_name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},

		"example_with_destination_ids": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDCGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGIPSetResourceName]().GetDefaultConfig)
					resp.Append(GetResourceConfig()[VDCGSecurityGroupResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_firewall" "example_with_destination_ids" {
						vdc_group_name = cloudavenue_vdcg.example.name
						  rules = [
							{
								action      = "ALLOW"
								name        = "allow all IPv4 traffic"
								direction   = "IN_OUT"
								destination_ids = [
									cloudavenue_vdcg_ip_set.example.id,
								],
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
							"action":            "ALLOW",
							"name":              "allow all IPv4 traffic",
							"direction":         "IN_OUT",
							"destination_ids.#": "1",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example_with_destination_ids" {
							vdc_group_name = cloudavenue_vdcg.example.name
							enabled = true
							  rules = [
								{
									action      = "ALLOW"
									name        = "allow out IPv4 traffic"
									direction   = "OUT"
									destination_ids = [
										cloudavenue_vdcg_ip_set.example.id,
										cloudavenue_vdcg_security_group.example.id
									]
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":            "ALLOW",
								"name":              "allow out IPv4 traffic",
								"direction":         "OUT",
								"destination_ids.#": "2",
							}),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group_name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},

		"example_all": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDCGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGIPSetResourceName]().GetDefaultConfig)
					resp.Append(GetResourceConfig()[VDCGSecurityGroupResourceName]().GetDefaultConfig)
					resp.Append(GetDataSourceConfig()[EdgeGatewayAppPortProfileDatasourceName]().GetSpecificConfig("example_provider_scope"))
					resp.Append(GetDataSourceConfig()[EdgeGatewayAppPortProfileDatasourceName]().GetSpecificConfig("example_system_scope"))
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_firewall" "example_all" {
						vdc_group_name = cloudavenue_vdcg.example.name
						  rules = [
							{
								action      = "ALLOW"
								name        = {{ generate . "name"}}
								direction   = "IN"
								source_ids = [
									cloudavenue_vdcg_ip_set.example.id,
								],
								destination_ids = [
									cloudavenue_vdcg_security_group.example.id,
								],
								app_port_profile_ids = [
									data.cloudavenue_edgegateway_app_port_profile.example_provider_scope.id,
								]
								source_groups_excluded = true
								destination_groups_excluded = true
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
							"action":                      "ALLOW",
							"name":                        testsacc.GetValueFromTemplate(resourceName, "name"),
							"direction":                   "IN",
							"source_ids.#":                "1",
							"destination_ids.#":           "1",
							"app_port_profile_ids.#":      "1",
							"source_groups_excluded":      "true",
							"destination_groups_excluded": "true",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_firewall" "example_all" {
							vdc_group_name = cloudavenue_vdcg.example.name
							enabled = true
							  rules = [
								{
									action      = "ALLOW"
									name        = {{ get . "name"}}
									direction   = "IN"
									source_ids = [
										cloudavenue_vdcg_ip_set.example.id,
									],
									destination_ids = [
										cloudavenue_vdcg_security_group.example.id,
									],
									app_port_profile_ids = [
										data.cloudavenue_edgegateway_app_port_profile.example_provider_scope.id,
									]
									source_groups_excluded = true
									destination_groups_excluded = true
								},
							    {
									action      = "DROP"
									name        = {{ generate . "name2"}}
									direction   = "IN"
									source_ids = [
										cloudavenue_vdcg_ip_set.example.id,
									],
									destination_ids = [
										cloudavenue_vdcg_security_group.example.id,
									],
									app_port_profile_ids = [
										data.cloudavenue_edgegateway_app_port_profile.example_provider_scope.id,
										data.cloudavenue_edgegateway_app_port_profile.example_system_scope.id
									]
									source_groups_excluded = false
									destination_groups_excluded = false
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":                      "ALLOW",
								"name":                        testsacc.GetValueFromTemplate(resourceName, "name"),
								"direction":                   "IN",
								"source_ids.#":                "1",
								"destination_ids.#":           "1",
								"app_port_profile_ids.#":      "1",
								"source_groups_excluded":      "true",
								"destination_groups_excluded": "true",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":                      "DROP",
								"name":                        testsacc.GetValueFromTemplate(resourceName, "name2"),
								"direction":                   "IN",
								"source_ids.#":                "1",
								"destination_ids.#":           "1",
								"app_port_profile_ids.#":      "2",
								"source_groups_excluded":      "false",
								"destination_groups_excluded": "false",
							}),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group_name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccVDCGFirewallResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGFirewallResource{}),
	})
}
