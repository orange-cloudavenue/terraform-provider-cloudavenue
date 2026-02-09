/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGDynamicSecurityGroupDataSource{}

const (
	VDCGDynamicSecurityGroupDataSourceName = testsacc.ResourceName("data.cloudavenue_vdcg_dynamic_security_group")
)

type VDCGDynamicSecurityGroupDataSource struct{}

func NewVDCGDynamicSecurityGroupDataSourceTest() testsacc.TestACC {
	return &VDCGDynamicSecurityGroupDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGDynamicSecurityGroupDataSource) GetResourceName() string {
	return VDCGDynamicSecurityGroupDataSourceName.String()
}

func (r *VDCGDynamicSecurityGroupDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGDynamicSecurityGroupResourceName]().GetDefaultConfig)
	return resp
}

func (r *VDCGDynamicSecurityGroupDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_dynamic_security_group" "example" {
						vdc_group_name = cloudavenue_vdcg.example.name
						name = cloudavenue_vdcg_dynamic_security_group.example.name
					}`,
					Checks: GetResourceConfig()[VDCGDynamicSecurityGroupResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVDCGDynamicSecurityGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGDynamicSecurityGroupDataSource{}),
	})
}
