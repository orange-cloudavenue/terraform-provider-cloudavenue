package testsacc

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &VDCGroupFirewallResource{}

const (
	VDCGroupFirewallResourceName = testsacc.ResourceName("cloudavenue_vdc_group_firewall")
)

type VDCGroupFirewallResource struct{}

func NewVDCGroupFirewallResourceTest() testsacc.TestACC {
	return &VDCGroupFirewallResource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGroupFirewallResource) GetResourceName() string {
	return VDCGroupFirewallResourceName.String()
}

func (r *VDCGroupFirewallResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGroupResourceName]().GetDefaultConfig)
	return
}

func (r *VDCGroupFirewallResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VDCGroup)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc_group_firewall" "example" {
						vdc_group = cloudavenue_vdc_group.example.name
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
						resource "cloudavenue_vdc_group_firewall" "example" {
							vdc_group = cloudavenue_vdc_group.example.name
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
						resource "cloudavenue_vdc_group_firewall" "example" {
							vdc_group = cloudavenue_vdc_group.example.name
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
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "DROUP",
								"name":        "drop in IPv4 traffic",
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
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							PlanOnly:           true,
							ExpectError:        regexp.MustCompile(`Invalid Attribute Value Match`),
						},
					},
					// * Test invalid ip_protocol
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc_group_firewall" "example" {
							vdc_group = cloudavenue_vdc_group.example.name
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
								"ip_protocol": "IPV5",
								"enabled":     "true",
							}),
						},
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							PlanOnly:           true,
							ExpectError:        regexp.MustCompile(`Invalid Attribute Value Match`),
						},
					},
					// * Test to disable the firewall
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc_group_firewall" "example" {
							vdc_group = cloudavenue_vdc_group.example.name
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
						resource "cloudavenue_vdc_group_firewall" "example" {
							vdc_group = cloudavenue_vdc_group.example.name
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
						resource "cloudavenue_vdc_group_firewall" "example" {
							vdc_group = cloudavenue_vdc_group.example.name
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
						ImportStateIDBuilder: []string{"vdc_group"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},

		"example_with_app_port_profile": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VDCGroup)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetDataSourceConfig()[EdgeGatewayAppPortProfileDatasourceName]().GetSpecificConfig("example_provider_scope"))
					resp.Append(GetDataSourceConfig()[EdgeGatewayAppPortProfileDatasourceName]().GetSpecificConfig("example_system_scope"))
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc_group_firewall" "example_with_app_port_profile" {
						vdc_group = cloudavenue_vdc_group.example.name
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
						resource "cloudavenue_vdc_group_firewall" "example_with_app_port_profile" {
							vdc_group = cloudavenue_vdc_group.example.name
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
						ImportStateIDBuilder: []string{"vdc_group"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},

		"example_with_source_ids": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VDCGroup)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayIPSetResourceName]().GetSpecificConfig("example_with_vdc_group"))
					resp.Append(GetResourceConfig()[EdgeGatewaySecurityGroupResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc_group_firewall" "example_with_source_ids" {
						vdc_group = cloudavenue_vdc_group.example.name
						  rules = [
							{
								action      = "ALLOW"
								name        = "allow all IPv4 traffic"
								direction   = "IN_OUT"
								source_ids = [
									cloudavenue_edgegateway_ip_set.example.id,
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
						resource "cloudavenue_vdc_group_firewall" "example_with_source_ids" {
							vdc_group = cloudavenue_vdc_group.example.name
							enabled = true
							  rules = [
								{
									action      = "ALLOW"
									name        = "allow in IPv4 traffic"
									direction   = "IN"
									source_ids = [
										cloudavenue_edgegateway_ip_set.example.id,
										cloudavenue_edgegateway_security_group.example.id
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
						ImportStateIDBuilder: []string{"vdc_group"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},

		"example_with_destination_ids": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VDCGroup)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[EdgeGatewayIPSetResourceName]().GetDefaultConfig)
					resp.Append(GetResourceConfig()[EdgeGatewaySecurityGroupResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc_group_firewall" "example_with_destination_ids" {
						vdc_group = cloudavenue_vdc_group.example.name
						  rules = [
							{
								action      = "ALLOW"
								name        = "allow all IPv4 traffic"
								direction   = "IN_OUT"
								destination_ids = [
									cloudavenue_edgegateway_ip_set.example.id,
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
						resource "cloudavenue_vdc_group_firewall" "example_with_destination_ids" {
							vdc_group = cloudavenue_vdc_group.example.name
							enabled = true
							  rules = [
								{
									action      = "ALLOW"
									name        = "allow out IPv4 traffic"
									direction   = "OUT"
									destination_ids = [
										cloudavenue_edgegateway_ip_set.example.id,
										cloudavenue_edgegateway_security_group.example.id
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
						ImportStateIDBuilder: []string{"vdc_group"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},

		"example_all": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.VDCGroup)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetDataSourceConfig()[EdgeGatewayIPSetResourceName]().GetDefaultConfig)
					resp.Append(GetDataSourceConfig()[EdgeGatewaySecurityGroupResourceName]().GetDefaultConfig)
					resp.Append(GetDataSourceConfig()[EdgeGatewayAppPortProfileDatasourceName]().GetSpecificConfig("example_provider_scope"))
					resp.Append(GetDataSourceConfig()[EdgeGatewayAppPortProfileDatasourceName]().GetSpecificConfig("example_system_scope"))
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc_group_firewall" "example_with_app_port_profile" {
						vdc_group = cloudavenue_vdc_group.example.name
						  rules = [
							{
								action      = "ALLOW"
								name        = {{ generate . "name"}}
								direction   = "IN"
								source_ids = [
									data.cloudavenue_edgegateway_ip_set.example.id,
								],
								destination_ids = [
									data.cloudavenue_edgegateway_security_group.example.id,
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
						resource "cloudavenue_vdc_group_firewall" "example_with_app_port_profile" {
							vdc_group = cloudavenue_vdc_group.example.name
							enabled = true
							  rules = [
								{
									action      = "DROP"
									name        = {{ generate . "name"}}
									direction   = "IN"
									source_ids = [
										data.cloudavenue_edgegateway_ip_set.example.id,
									],
									destination_ids = [
										data.cloudavenue_edgegateway_security_group.example.id,
									],
									app_port_profile_ids = [
										data.cloudavenue_edgegateway_app_port_profile.example_provider_scope.id,
										data.cloudavenue_edgegateway_app_port_profile.example_system_scope.id
									]
									source_groups_excluded = false
									destination_groups_excluded = true
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":                      "DROP",
								"name":                        testsacc.GetValueFromTemplate(resourceName, "name"),
								"direction":                   "IN",
								"source_ids.#":                "1",
								"destination_ids.#":           "1",
								"app_port_profile_ids.#":      "2",
								"source_groups_excluded":      "false",
								"destination_groups_excluded": "true",
							}),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccVDCGroupFirewallResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGroupFirewallResource{}),
	})
}
