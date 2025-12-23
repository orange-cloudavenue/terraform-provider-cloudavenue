/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCDataSource{}

const (
	VDCDataSourceName = testsacc.ResourceName("data.cloudavenue_vdc")
)

type VDCDataSource struct{}

func NewVDCDataSourceTest() testsacc.TestACC {
	return &VDCDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCDataSource) GetResourceName() string {
	return VDCDataSourceName.String()
}

func (r *VDCDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
	return resp
}

func (r *VDCDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: GetResourceConfig()[VDCResourceName]().GetDefaultChecks(),
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdc" "example" {
						name = cloudavenue_vdc.example.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.class", "gold"),
						resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.default", "true"),
						resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.limit", "500"),
					},
				},
			}
		},
	}
}

func TestAccVDCDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCDataSource{}),
	})
}
