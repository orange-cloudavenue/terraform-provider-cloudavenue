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

var _ testsacc.TestACC = &VDCStorageProfilesDataSource{}

const (
	VDCStorageProfilesDataSourceName = testsacc.ResourceName("data.cloudavenue_vdc_storageprofiles")
)

type VDCStorageProfilesDataSource struct{}

func NewVDCStorageProfilesDataSourceTest() testsacc.TestACC {
	return &VDCStorageProfilesDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCStorageProfilesDataSource) GetResourceName() string {
	return VDCStorageProfilesDataSourceName.String()
}

func (r *VDCStorageProfilesDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
	return
}

func (r *VDCStorageProfilesDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					data "cloudavenue_vdc_storage_profiles" "example" {
						depends_on = [
							cloudavenue_vdc.example
						]
						vdc_id = cloudavenue_vdc.example.id
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_id"),
						resource.TestCheckResourceAttrSet(resourceName, "vdc_name"),
						resource.TestCheckTypeSetElemAttr(resourceName, "storageprofiles.*.id", "*"),
						resource.TestCheckTypeSetElemAttr(resourceName, "storageprofiles.*.class", "*"),
						resource.TestCheckTypeSetElemAttr(resourceName, "storageprofiles.*.limit", "*"),
						resource.TestCheckTypeSetElemAttr(resourceName, "storageprofiles.*.used", "*"),
						resource.TestCheckTypeSetElemAttr(resourceName, "storageprofiles.*.default", "*"),
					},
				},
			}
		},
	}
}

func TestAccVDCStorageProfilesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCStorageProfilesDataSource{}),
	})
}
