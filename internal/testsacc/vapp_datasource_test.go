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

var _ testsacc.TestACC = &VAppDatasource{}

const (
	VAppDatasourceName = testsacc.ResourceName("data.cloudavenue_vapp")
)

type VAppDatasource struct{}

func NewVAppDatasourceTest() testsacc.TestACC {
	return &VAppDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *VAppDatasource) GetResourceName() string {
	return VAppDatasourceName.String()
}

func (r *VAppDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VAppResourceName]().GetDefaultConfig)
	return
}

func (r *VAppDatasource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vapp" "example" {
						name = cloudavenue_vapp.example.name
						vdc = cloudavenue_vapp.example.vdc
					}`,
					// Here use resource config test to test the data source
					Checks: GetResourceConfig()[VAppResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVAppDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VAppDatasource{}),
	})
}
