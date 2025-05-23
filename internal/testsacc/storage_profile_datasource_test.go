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

var _ testsacc.TestACC = &StorageProfileDataSource{}

const (
	StorageProfileDataSourceName = testsacc.ResourceName("data.cloudavenue_storage_profile")
)

type StorageProfileDataSource struct{}

func NewStorageProfileDataSourceTest() testsacc.TestACC {
	return &StorageProfileDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *StorageProfileDataSource) GetResourceName() string {
	return StorageProfileDataSourceName.String()
}

func (r *StorageProfileDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
	return
}

func (r *StorageProfileDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_storage_profile" "example" {
						name = "gold"
						vdc = cloudavenue_vdc.example.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDCStorageProfile)),
						resource.TestCheckResourceAttrSet(resourceName, "vdc"),
						resource.TestCheckResourceAttr(resourceName, "name", "gold"),
						resource.TestCheckResourceAttrSet(resourceName, "limit"),
						resource.TestCheckResourceAttrSet(resourceName, "used_storage"),
						resource.TestCheckResourceAttrSet(resourceName, "default"),
						resource.TestCheckResourceAttrSet(resourceName, "enabled"),
						resource.TestCheckResourceAttrSet(resourceName, "iops_allocated"),
						resource.TestCheckResourceAttrSet(resourceName, "units"),
						resource.TestCheckResourceAttrSet(resourceName, "iops_limiting_enabled"),
						resource.TestCheckResourceAttrSet(resourceName, "maximum_disk_iops"),
						resource.TestCheckResourceAttrSet(resourceName, "default_disk_iops"),
						resource.TestCheckResourceAttrSet(resourceName, "disk_iops_per_gb_max"),
						resource.TestCheckResourceAttrSet(resourceName, "iops_limit"),
					},
				},
			}
		},
	}
}

func TestAccStorageProfileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3UserDataSource{}),
	})
}
