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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ELBServiceEngineGroupsDataSource{}

const (
	ELBServiceEngineGroupsDataSourceName = testsacc.ResourceName("data.cloudavenue_elb_service_engine_groups")
)

type ELBServiceEngineGroupsDataSource struct{}

func NewELBServiceEngineGroupsDataSourceTest() testsacc.TestACC {
	return &ELBServiceEngineGroupsDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *ELBServiceEngineGroupsDataSource) GetResourceName() string {
	return ELBServiceEngineGroupsDataSourceName.String()
}

func (r *ELBServiceEngineGroupsDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_for_elb"))
	return resp
}

func (r *ELBServiceEngineGroupsDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_elb_service_engine_groups" "example" {
						edge_gateway_name = data.cloudavenue_edgegateway.example_for_elb.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.#"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.edge_gateway_id"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.deployed_virtual_services"),
					},
				},
			}
		},
		"example_with_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_elb_service_engine_groups" "example_with_id" {
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.#"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.edge_gateway_id"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "service_engine_groups.0.deployed_virtual_services"),
					},
				},
			}
		},
	}
}

func TestAccELBServiceEngineGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ELBServiceEngineGroupsDataSource{}),
	})
}
