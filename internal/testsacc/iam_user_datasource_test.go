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

var _ testsacc.TestACC = &IAMUserDataSource{}

const (
	IAMUserDataSourceName = testsacc.ResourceName("data.cloudavenue_iam_user")
)

type IAMUserDataSource struct{}

func NewIAMUserDataSourceTest() testsacc.TestACC {
	return &IAMUserDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *IAMUserDataSource) GetResourceName() string {
	return IAMUserDataSourceName.String()
}

func (r *IAMUserDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *IAMUserDataSource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[IAMUserResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_iam_user" "example" {
						name = cloudavenue_iam_user.example.name
					}`,
					Checks: GetResourceConfig()[IAMUserResourceName]().GetDefaultChecks(),
				},
			}
		},
		"example_saml_user": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[IAMUserSAMLResourceName]().GetDefaultConfig)
					return
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_iam_user" "example_saml_user" {
						name = cloudavenue_iam_user_saml.example.user_name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.User)),
						resource.TestCheckResourceAttr(resourceName, "name", "mickael.stanislas.ext"),
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

func TestAccIAMUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&IAMUserDataSource{}),
	})
}
