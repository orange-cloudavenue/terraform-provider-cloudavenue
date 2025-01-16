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

var _ testsacc.TestACC = &IAMUserSAMLResource{}

const (
	IAMUserSAMLResourceName = testsacc.ResourceName("cloudavenue_iam_user_saml")
)

type IAMUserSAMLResource struct{}

func NewIAMUserSAMLResourceTest() testsacc.TestACC {
	return &IAMUserSAMLResource{}
}

// GetResourceName returns the name of the resource.
func (r *IAMUserSAMLResource) GetResourceName() string {
	return IAMUserSAMLResourceName.String()
}

func (r *IAMUserSAMLResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *IAMUserSAMLResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.User)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_iam_user_saml" "example" {
						user_name = "mickael.stanislas.ext"
						role_name = "Organization Administrator"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "user_name", "mickael.stanislas.ext"),
						resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "0"),
						resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "0"),
						resource.TestCheckResourceAttr(resourceName, "take_ownership", "true"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Disable the user
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_iam_user_saml" "example" {
							user_name = "mickael.stanislas.ext"
							role_name = "Organization Administrator"
							enabled = false
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "user_name", "mickael.stanislas.ext"),
							resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "0"),
							resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "0"),
							resource.TestCheckResourceAttr(resourceName, "take_ownership", "true"),
						},
					},
					// * Re-enable the user
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_iam_user_saml" "example" {
							user_name = "mickael.stanislas.ext"
							role_name = "Organization Administrator"
							enabled = true
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "user_name", "mickael.stanislas.ext"),
							resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "0"),
							resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "0"),
							resource.TestCheckResourceAttr(resourceName, "take_ownership", "true"),
						},
					},
					// * Change Quotas
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_iam_user_saml" "example" {
							user_name = "mickael.stanislas.ext"
							role_name = "Organization Administrator"
							enabled = true
							deployed_vm_quota = 10
							stored_vm_quota = 5
					}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "user_name", "mickael.stanislas.ext"),
							resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "10"),
							resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "5"),
							resource.TestCheckResourceAttr(resourceName, "take_ownership", "true"),
						},
					},
					// * Change Take Ownership
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_iam_user_saml" "example" {
							user_name = "mickael.stanislas.ext"
							role_name = "Organization Administrator"
							enabled = true
							deployed_vm_quota = 10
							stored_vm_quota = 5
							take_ownership = false
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "user_name", "mickael.stanislas.ext"),
							resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "10"),
							resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "5"),
							resource.TestCheckResourceAttr(resourceName, "take_ownership", "false"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder:    []string{"user_name"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"take_ownership"},
					},
				},
				Destroy: true,
			}
		},
		"example_quota_on_creation": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.User)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_iam_user_saml" "example_quota_on_creation" {
						user_name = "mickael.stanislas.ext"
						role_name = "Organization Administrator"
						deployed_vm_quota = 10
						stored_vm_quota = 5
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "user_name", "mickael.stanislas.ext"),
						resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "10"),
						resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "5"),
						resource.TestCheckResourceAttr(resourceName, "take_ownership", "true"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * unset quotas
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_iam_user_saml" "example_quota_on_creation" {
							user_name = "mickael.stanislas.ext"
							role_name = "Organization Administrator"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "user_name", "mickael.stanislas.ext"),
							resource.TestCheckResourceAttr(resourceName, "role_name", "Organization Administrator"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "deployed_vm_quota", "0"),
							resource.TestCheckResourceAttr(resourceName, "stored_vm_quota", "0"),
							resource.TestCheckResourceAttr(resourceName, "take_ownership", "true"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder:    []string{"user_name"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"take_ownership"},
					},
				},
			}
		},
	}
}

func TestAccIAMUserSAMLResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&IAMUserSAMLResource{}),
	})
}
