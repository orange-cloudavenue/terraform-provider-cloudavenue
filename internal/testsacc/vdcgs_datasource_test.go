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

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/helpers"
)

var _ testsacc.TestACC = &VDCGsDataSource{}

const (
	VDCGsDataSourceName = testsacc.ResourceName("data.cloudavenue_vdcgs")
)

type VDCGsDataSource struct{}

func NewVDCGsDataSourceTest() testsacc.TestACC {
	return &VDCGsDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGsDataSource) GetResourceName() string {
	return VDCGsDataSourceName.String()
}

func (r *VDCGsDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
	return resp
}

func (r *VDCGsDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					// No filter applied
					data "cloudavenue_vdcgs" "example" {}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckNoResourceAttr(resourceName, "filter_by_id"),
						resource.TestCheckNoResourceAttr(resourceName, "filter_by_name"),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_groups.#"),
						resource.TestCheckResourceAttrWith(resourceName, "vdc_groups.0.id", helpers.TestIsType(urn.VDCGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_groups.0.name"),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_groups.0.vdcs.#"),
						resource.TestCheckResourceAttrWith(resourceName, "vdc_groups.0.vdcs.0.id", helpers.TestIsType(urn.VDC)),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_groups.0.vdcs.0.name"),
					},
				},
			}
		},
		"example_with_filter": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					// Filter applied
					data "cloudavenue_vdcgs" "example_with_filter" {
						filter_by_name = "tftest*"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckNoResourceAttr(resourceName, "filter_by_id"),
						resource.TestCheckResourceAttr(resourceName, "filter_by_name", "tftest*"),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_groups.#"),
						resource.TestCheckResourceAttrWith(resourceName, "vdc_groups.0.id", helpers.TestIsType(urn.VDCGroup)),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_groups.0.name"),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_groups.0.vdcs.#"),
						resource.TestCheckResourceAttrWith(resourceName, "vdc_groups.0.vdcs.0.id", helpers.TestIsType(urn.VDC)),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_groups.0.vdcs.0.name"),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						// Filter applied
						data "cloudavenue_vdcgs" "example_with_filter" {
							filter_by_name = "nonexistent*"
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "id"),
							resource.TestCheckNoResourceAttr(resourceName, "filter_by_id"),
							resource.TestCheckResourceAttr(resourceName, "filter_by_name", "nonexistent*"),
							resource.TestCheckResourceAttr(resourceName, "vdc_groups.#", "0"),
						},
					},
				},
			}
		},
	}
}

func TestAccVDCGSDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGsDataSource{}),
	})
}
