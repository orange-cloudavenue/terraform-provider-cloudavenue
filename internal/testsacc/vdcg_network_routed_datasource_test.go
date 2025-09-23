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

var _ testsacc.TestACC = &VDCGNetworkRoutedDataSource{}

const (
	VDCGNetworkRoutedDataSourceName = testsacc.ResourceName("data.cloudavenue_vdcg_network_routed")
)

type VDCGNetworkRoutedDataSource struct{}

func NewVDCGNetworkRoutedDataSourceTest() testsacc.TestACC {
	return &VDCGNetworkRoutedDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGNetworkRoutedDataSource) GetResourceName() string {
	return VDCGNetworkRoutedDataSourceName.String()
}

func (r *VDCGNetworkRoutedDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[VDCGNetworkRoutedResourceName]().GetDefaultConfig)
	return resp
}

func (r *VDCGNetworkRoutedDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_vdcg_network_routed" "example" {
						name = cloudavenue_vdcg_network_routed.example.name
						vdc_group_id = cloudavenue_vdcg.example.id
						edge_gateway_id = cloudavenue_edgegateway.example_with_vdc_group.id
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[VDCGNetworkRoutedResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccVDCGNetworkRoutedDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGNetworkRoutedDataSource{}),
	})
}
