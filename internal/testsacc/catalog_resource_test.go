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

var _ testsacc.TestACC = &CatalogResource{}

const (
	CatalogResourceName = testsacc.ResourceName("cloudavenue_catalog")
)

type CatalogResource struct{}

func NewCatalogResourceTest() testsacc.TestACC {
	return &CatalogResource{}
}

// GetResourceName returns the name of the resource.
func (r *CatalogResource) GetResourceName() string {
	return CatalogResourceName.String()
}

func (r *CatalogResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *CatalogResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Catalog)),
					resource.TestCheckResourceAttr(resourceName, "delete_recursive", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_force", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_catalog" "example" {
						name             = {{ generate . "name" }}
						description      = {{ generate . "description" "longString" }}
						delete_recursive = true
						delete_force     = true
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_catalog" "example" {
							name             = {{ get . "name" }}
							description      = {{ generate . "description" "longString"}}
							delete_recursive = true
							delete_force     = true
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateID:           testsacc.GetValueFromTemplate(resourceName, "name"),
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"delete_force", "delete_recursive"},
					},
				},
			}
		},
	}
}

func TestAccCatalogResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&CatalogResource{}),
	})
}
