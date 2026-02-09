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

var _ testsacc.TestACC = &IAMUserResource{}

const (
	IAMUserResourceName = testsacc.ResourceName("cloudavenue_iam_user")
)

type IAMUserResource struct{}

func NewIAMUserResourceTest() testsacc.TestACC {
	return &IAMUserResource{}
}

// GetResourceName returns the name of the resource.
func (r *IAMUserResource) GetResourceName() string {
	return IAMUserResourceName.String()
}

func (r *IAMUserResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *IAMUserResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.User)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_iam_user" "example" {
						name        = {{ generate . "name" }}
						role_name   = "Organization Administrator"
						password    = "Th!s1sSecur3P@ssword"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckNoResourceAttr(resourceName, "email"),
						resource.TestCheckNoResourceAttr(resourceName, "full_name"),
						resource.TestCheckNoResourceAttr(resourceName, "telephone"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "0"),
						resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "0"),
						resource.TestCheckResourceAttr(resourceName, "take_ownership", "true"),
						resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_iam_user" "example" {
							name             = {{ get . "name" }}
							full_name        = "Example User"
							password         = "Th!s1sSecur3P@ssword"
							role_name        = "Organization Administrator"
							enabled          = false
							email            = "foo@bar.org"
							telephone        = "1234567890"
							take_ownership   = false
							deployed_vm_quota = 10
							stored_vm_quota   = 5
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "full_name", "Example User"),
							resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "email", "foo@bar.org"),
							resource.TestCheckResourceAttr(resourceName, "telephone", "1234567890"),
							resource.TestCheckResourceAttr(resourceName, "take_ownership", "false"),
							resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "10"),
							resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "5"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder:    []string{"name"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"take_ownership", "password"},
					},
				},
			}
		},
		// This test is used by other tests. And is used by datasource because take_ownership is not tested here
		"example_2": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.User)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_iam_user" "example_2" {
						name        = {{ generate . "name" }}
						role_name   = "Organization Administrator"
						password    = "Th!s1sSecur3P@ssword"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckNoResourceAttr(resourceName, "email"),
						resource.TestCheckNoResourceAttr(resourceName, "full_name"),
						resource.TestCheckNoResourceAttr(resourceName, "telephone"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "0"),
						resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "0"),
						resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
					},
				},
			}
		},
	}
}

func TestAccIAMUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&IAMUserResource{}),
	})
}
