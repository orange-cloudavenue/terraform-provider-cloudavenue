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

// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &CatalogVAppTemplateDataSource{}

const (
	CatalogVAppTemplateDataSourceName = testsacc.ResourceName("data.cloudavenue_catalog_vapp_template")
)

type CatalogVAppTemplateDataSource struct{}

func NewCatalogVAppTemplateDataSourceTest() testsacc.TestACC {
	return &CatalogVAppTemplateDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *CatalogVAppTemplateDataSource) GetResourceName() string {
	return CatalogVAppTemplateDataSourceName.String()
}

func (r *CatalogVAppTemplateDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *CatalogVAppTemplateDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_catalog_vapp_template" "example" {
						catalog_name  	= "Orange-Linux"
						template_name 	= "UBUNTU_20.04"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VAPPTemplate)),
						// Catalog
						resource.TestCheckResourceAttr(resourceName, "catalog_name", "Orange-Linux"),
						resource.TestCheckResourceAttrWith(resourceName, "catalog_id", urn.TestIsType(urn.Catalog)),

						resource.TestCheckResourceAttr(resourceName, "template_name", "UBUNTU_20.04"),
						resource.TestCheckResourceAttrWith(resourceName, "template_id", urn.TestIsType(urn.VAPPTemplate)),
						// Other
						resource.TestCheckResourceAttrSet(resourceName, "created_at"),
						resource.TestCheckResourceAttrSet(resourceName, "vm_names.#"),
					},
				},
			}
		},
	}
}

func TestACCCatalogVAppTemplateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&CatalogVAppTemplateDataSource{}),
	})
}
