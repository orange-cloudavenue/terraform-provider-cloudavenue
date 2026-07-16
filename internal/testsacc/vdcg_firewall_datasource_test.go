/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGFirewallDataSource{}

const (
	VDCGFirewallDataSourceName = testsacc.ResourceName("data.cloudavenue_vdcg_firewall")
)

type VDCGFirewallDataSource struct{}

func NewVDCGFirewallDataSourceTest() testsacc.TestACC {
	return &VDCGFirewallDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGFirewallDataSource) GetResourceName() string {
	return VDCGFirewallDataSourceName.String()
}

func (r *VDCGFirewallDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGFirewallResourceName]().GetDefaultConfig)
	return resp
}

func (r *VDCGFirewallDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		testNameExample: func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_firewall" "example" {
						vdc_group_name = cloudavenue_vdcg_firewall.example.vdc_group_name
					}`,
					Checks: GetResourceConfig()[VDCGFirewallResourceName]().GetDefaultChecks(),
				},
			}
		},
		testNameExampleWithContextProfile: func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGFirewallResourceName]().GetSpecificConfig(testNameExampleWithContextProfile))
					return resp
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_firewall" "example_with_context_profile" {
					  vdc_group_name = cloudavenue_vdcg_firewall.example_with_context_profile.vdc_group_name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.name", "allow outbound SSL traffic"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.action", "ALLOW"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.direction", "OUT"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.network_context_profile_ids.#", "1"),
						resource.TestCheckResourceAttrSet(resourceName, "rules.0.network_context_profile_ids.0"),
						resource.TestCheckResourceAttr(resourceName, "rules.1.name", "block all inbound"),
						resource.TestCheckResourceAttr(resourceName, "rules.1.network_context_profile_ids.#", "0"),
					},
				},
				Destroy: true,
			}
		},
	}
}

func TestAccVDCGFirewallDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGFirewallDataSource{}),
	})
}
