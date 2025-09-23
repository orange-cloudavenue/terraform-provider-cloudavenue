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

var _ testsacc.TestACC = &ELBServiceEngineGroupDataSource{}

const (
	ELBServiceEngineGroupDataSourceName = testsacc.ResourceName("data.cloudavenue_elb_service_engine_group")
)

type ELBServiceEngineGroupDataSource struct{}

func NewELBServiceEngineGroupDataSourceTest() testsacc.TestACC {
	return &ELBServiceEngineGroupDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *ELBServiceEngineGroupDataSource) GetResourceName() string {
	return ELBServiceEngineGroupDataSourceName.String()
}

func (r *ELBServiceEngineGroupDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_for_elb"))
	resp.Append(GetDataSourceConfig()[ELBServiceEngineGroupsDataSourceName]().GetDefaultConfig)
	return resp
}

func (r *ELBServiceEngineGroupDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_elb_service_engine_group" "example" {
						name = data.cloudavenue_elb_service_engine_groups.example.service_engine_groups.0.name
						edge_gateway_name = data.cloudavenue_edgegateway.example_for_elb.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.ServiceEngineGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "deployed_virtual_services"),
					},
				},
			}
		},
		"example_with_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_elb_service_engine_group" "example_for_elb" {
						id = data.cloudavenue_elb_service_engine_groups.example.service_engine_groups.0.id
						edge_gateway_name = data.cloudavenue_edgegateway.example_for_elb.name
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.ServiceEngineGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "deployed_virtual_services"),
					},
				},
			}
		},
		"example_with_edge_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_elb_service_engine_group" "example_with_edge_id" {
						id = data.cloudavenue_elb_service_engine_groups.example.service_engine_groups.0.id
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.ServiceEngineGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "deployed_virtual_services"),
					},
				},
			}
		},
		"example_with_name_and_edge_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_elb_service_engine_group" "example_with_name_and_edge_id" {
						name = data.cloudavenue_elb_service_engine_groups.example.service_engine_groups.0.name
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.ServiceEngineGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "max_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "reserved_virtual_services"),
						resource.TestCheckResourceAttrSet(resourceName, "deployed_virtual_services"),
					},
				},
			}
		},
	}
}

func TestAccELBServiceEngineGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ELBServiceEngineGroupDataSource{}),
	})
}
