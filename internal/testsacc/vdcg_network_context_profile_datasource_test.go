/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at HTTP://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGNetworkContextProfileDatasource{}

const VDCGNetworkContextProfileDatasourceName = testsacc.ResourceName("data.cloudavenue_vdcg_network_context_profile")

type VDCGNetworkContextProfileDatasource struct{}

func NewVDCGNetworkContextProfileDatasourceTest() testsacc.TestACC {
	return &VDCGNetworkContextProfileDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGNetworkContextProfileDatasource) GetResourceName() string {
	return VDCGNetworkContextProfileDatasourceName.String()
}

func (r *VDCGNetworkContextProfileDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
	return resp
}

func (r *VDCGNetworkContextProfileDatasource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_vdcg_network_context_profile" "example" {
						vdc_group_name = cloudavenue_vdcg.example.name
						name           = "SSL"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", "SSL"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
						resource.TestCheckResourceAttrSet(resourceName, "description"),
						resource.TestCheckResourceAttrPair(resourceName, "vdc_group_name", VDCGResourceName.String()+".example", "name"),
					},
				},
				Destroy: true,
			}
		},
		"example_by_vdc_group_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_vdcg_network_context_profile" "example_by_vdc_group_id" {
						vdc_group_id = cloudavenue_vdcg.example.id
						name         = "CIFS"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", "CIFS"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
						resource.TestCheckResourceAttrSet(resourceName, "description"),
						resource.TestCheckResourceAttrPair(resourceName, "vdc_group_id", VDCGResourceName.String()+".example", "id"),
					},
				},
				Destroy: true,
			}
		},
		"web_tier_http": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_vdcg_network_context_profile" "web_tier_http" {
						vdc_group_name = cloudavenue_vdcg.example.name
						name           = "HTTP"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", "HTTP"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
						resource.TestCheckResourceAttrSet(resourceName, "description"),
						resource.TestCheckResourceAttrPair(resourceName, "vdc_group_name", VDCGResourceName.String()+".example", "name"),
					},
				},
				Destroy: true,
			}
		},
		"database_tier_mysql": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_vdcg_network_context_profile" "database_tier_mysql" {
						vdc_group_id = cloudavenue_vdcg.example.id
						name         = "MYSQL"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", "MYSQL"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
						resource.TestCheckResourceAttrSet(resourceName, "description"),
						resource.TestCheckResourceAttrPair(resourceName, "vdc_group_id", VDCGResourceName.String()+".example", "id"),
					},
				},
				Destroy: true,
			}
		},
		"centralized_dns": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_vdcg_network_context_profile" "centralized_dns" {
						vdc_group_name = cloudavenue_vdcg.example.name
						name           = "DNS"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", "DNS"),
						resource.TestCheckResourceAttr(resourceName, "scope", "SYSTEM"),
						resource.TestCheckResourceAttrSet(resourceName, "description"),
						resource.TestCheckResourceAttrPair(resourceName, "vdc_group_name", VDCGResourceName.String()+".example", "name"),
					},
				},
				Destroy: true,
			}
		},
		"custom_erp_application": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGNetworkContextProfileResourceName]().GetSpecificConfig("custom_erp"))
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_vdcg_network_context_profile" "custom_erp_application" {
						vdc_group_name = cloudavenue_vdcg.example.name
						id             = cloudavenue_vdcg_network_context_profile.custom_erp.id
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttr(resourceName, "scope", "TENANT"),
						resource.TestCheckResourceAttr(resourceName, "description", "Internal ERP system (HTTP front end, MSSQL back end)"),
					},
				},
				Destroy: true,
			}
		},
		"lookup_consistency": func(_ context.Context, _ string) testsacc.Test {
			byName := VDCGNetworkContextProfileDatasourceName.String() + ".lookup_consistency_by_name"
			byGroupID := VDCGNetworkContextProfileDatasourceName.String() + ".lookup_consistency_by_group_id"

			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_vdcg_network_context_profile" "lookup_consistency_by_name" {
						vdc_group_name = cloudavenue_vdcg.example.name
						name           = "SSL"
					}

					data "cloudavenue_vdcg_network_context_profile" "lookup_consistency_by_group_id" {
						vdc_group_id = cloudavenue_vdcg.example.id
						name         = "SSL"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrPair(byName, "id", byGroupID, "id"),
						resource.TestCheckResourceAttrPair(byName, "description", byGroupID, "description"),
						resource.TestCheckResourceAttrPair(byName, "scope", byGroupID, "scope"),
					},
				},
				Destroy: true,
			}
		},
		"not_found_typo": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: testsacc.TFData(`
					data "cloudavenue_vdcg_network_context_profile" "not_found_typo" {
						vdc_group_name = cloudavenue_vdcg.example.name
						name           = "HTTPP"
					}`),
					TFAdvanced: testsacc.TFAdvanced{
						ExpectError: regexp.MustCompile(`No Network Context Profile found with name or ID`),
					},
				},
				Destroy: false,
			}
		},
	}
}

func TestAccVDCGNetworkContextProfileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGNetworkContextProfileDatasource{}),
	})
}
