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

var _ testsacc.TestACC = &StorageProfilesDataSource{}

const (
	StorageProfilesDataSourceName = testsacc.ResourceName("data.cloudavenue_storage_profile")
)

type StorageProfilesDataSource struct{}

func NewStorageProfilesDataSourceTest() testsacc.TestACC {
	return &StorageProfilesDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *StorageProfilesDataSource) GetResourceName() string {
	return StorageProfilesDataSourceName.String()
}

func (r *StorageProfilesDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
	return
}

func (r *StorageProfilesDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_storage_profiles" "example" {
						vdc = cloudavenue_vdc.example.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDCStorageProfile)),
						resource.TestCheckResourceAttrSet(resourceName, "vdc"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.vdc"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.name"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.limit"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.used_storage"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.default"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.enabled"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.iops_allocated"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.units"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.iops_limiting_enabled"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.maximum_disk_iops"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.default_disk_iops"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.disk_iops_per_gb_max"),
						resource.TestCheckResourceAttrSet(resourceName, "storage_profiles.0.iops_limit"),
					},
				},
			}
		},
	}
}

func TestAccStorageProfilesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3UserDataSource{}),
	})
}
