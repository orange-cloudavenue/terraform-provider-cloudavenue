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

var _ testsacc.TestACC = &CatalogACLDataSource{}

const (
	CatalogACLDataSourceName = testsacc.ResourceName("data.cloudavenue_catalog_acl")
)

type CatalogACLDataSource struct{}

func NewCatalogACLDataSourceTest() testsacc.TestACC {
	return &CatalogACLDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *CatalogACLDataSource) GetResourceName() string {
	return CatalogACLDataSourceName.String()
}

func (r *CatalogACLDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[CatalogACLResourceName]().GetDefaultConfig)
	return
}

func (r *CatalogACLDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_catalog_acl" "example" {
						catalog_id = cloudavenue_catalog_acl.example.catalog_id
					}`,
					Checks: GetResourceConfig()[CatalogACLResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestCatalogAccACLDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&CatalogACLDataSource{}),
	})
}
