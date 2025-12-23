/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGAppPortProfileDatasource{}

const (
	VDCGAppPortProfileDatasourceName = testsacc.ResourceName("data.cloudavenue_vdcg_app_port_profile")
)

type VDCGAppPortProfileDatasource struct{}

func NewVDCGAppPortProfileDatasourceTest() testsacc.TestACC {
	return &VDCGAppPortProfileDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGAppPortProfileDatasource) GetResourceName() string {
	return VDCGAppPortProfileDatasourceName.String()
}

func (r *VDCGAppPortProfileDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *VDCGAppPortProfileDatasource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGAppPortProfileResourceName]().GetDefaultConfig)
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_app_port_profile" "example" {
						vdc_group_name = cloudavenue_vdcg_app_port_profile.example.vdc_group_id
						name = cloudavenue_vdcg_app_port_profile.example.name
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[VDCGAppPortProfileResourceName]().GetDefaultChecks(),
				},
				Destroy: true,
			}
		},
		"example_by_id": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGAppPortProfileResourceName]().GetDefaultConfig)
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_app_port_profile" "example_by_id" {
						vdc_group_id = cloudavenue_vdcg_app_port_profile.example.vdc_group_id
						id = cloudavenue_vdcg_app_port_profile.example.id
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[VDCGAppPortProfileResourceName]().GetDefaultChecks(),
				},
				Destroy: true,
			}
		},
		"example_provider_scope": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_app_port_profile" "example_provider_scope" {
						vdc_group_id = cloudavenue_vdcg.example.id
						name = "BKP_TCP_bpcd"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.AppPortProfile)),
						resource.TestCheckResourceAttr(resourceName, "name", "BKP_TCP_bpcd"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.protocol", "TCP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.ports.0", "13782"),
						resource.TestCheckResourceAttr(resourceName, "scope", "PROVIDER"),
					},
				},
				Destroy: true,
			}
		},
		"example_system_scope": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_app_port_profile" "example_system_scope" {
						vdc_group_id = cloudavenue_vdcg.example.id
						name = "HTTP"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.AppPortProfile)),
						resource.TestCheckResourceAttr(resourceName, "name", "HTTP"),
						resource.TestCheckResourceAttr(resourceName, "description", "HTTP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.protocol", "TCP"),
						resource.TestCheckResourceAttr(resourceName, "app_ports.0.ports.0", "80"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
					},
				},
				Destroy: true,
			}
		},
		"example_two_app_ports_with_same_name": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGAppPortProfileResourceName]().GetSpecificConfig("example_http_scope_tenant"))
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_app_port_profile" "example_two_app_ports_with_same_name" {
						vdc_group_id = cloudavenue_vdcg_app_port_profile.example_http_scope_tenant.vdc_group_id
						name = cloudavenue_vdcg_app_port_profile.example_http_scope_tenant.name
					}`,
					TFAdvanced: testsacc.TFAdvanced{
						ExpectError: regexp.MustCompile(`Multiple App Port Profiles found`),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						data "cloudavenue_vdcg_app_port_profile" "example_two_app_ports_with_same_name" {
							vdc_group_id = cloudavenue_vdcg_app_port_profile.example_http_scope_tenant.vdc_group_id
							name = cloudavenue_vdcg_app_port_profile.example_http_scope_tenant.name
							scope = "TENANT"
						}`,
						Checks: GetResourceConfig()[VDCGAppPortProfileResourceName]().GetSpecificChecks("example_http_scope_tenant"),
					},
					{
						TFConfig: `
							data "cloudavenue_vdcg_app_port_profile" "example_two_app_ports_with_same_name" {
								vdc_group_id = cloudavenue_vdcg_app_port_profile.example_http_scope_tenant.vdc_group_id
								name = cloudavenue_vdcg_app_port_profile.example_http_scope_tenant.name
								scope = "SYSTEM"
							}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.AppPortProfile)),
							resource.TestCheckResourceAttr(resourceName, "name", "HTTP"),
							resource.TestCheckResourceAttr(resourceName, "description", "HTTP"),
							resource.TestCheckResourceAttr(resourceName, "app_ports.0.protocol", "TCP"),
							resource.TestCheckResourceAttr(resourceName, "app_ports.0.ports.0", "80"),
							resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
						},
					},
				},
				Destroy: true,
			}
		},
	}
}

func TestAccVDCGAppPortProfileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGAppPortProfileDatasource{}),
	})
}
