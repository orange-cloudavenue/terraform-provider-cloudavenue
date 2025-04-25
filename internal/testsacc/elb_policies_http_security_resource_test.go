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

var _ testsacc.TestACC = &ELBPoliciesHTTPSecurityResource{}

const (
	ELBPoliciesHTTPSecurityResourceName = testsacc.ResourceName("cloudavenue_elb_policies_http_security")
)

type ELBPoliciesHTTPSecurityResource struct{}

func NewELBPoliciesHTTPSecurityResourceTest() testsacc.TestACC {
	return &ELBPoliciesHTTPSecurityResource{}
}

// GetResourceName returns the name of the resource.
func (r *ELBPoliciesHTTPSecurityResource) GetResourceName() string {
	return ELBPoliciesHTTPSecurityResourceName.String()
}

func (r *ELBPoliciesHTTPSecurityResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[ELBVirtualServiceResourceName]().GetDefaultConfig)
	return
}

func (r *ELBPoliciesHTTPSecurityResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_policies_http_security" "example" {
						virtual_service_id = cloudavenue_elb_virtual_service.example.id
						policies = [
							// Policy 1
							{
							name = "example"

							// Define the criteria for the policy
							criteria = {
								path = {
									criteria = "CONTAINS"
									paths    = ["/example"]
								}
							}

							// Define the action to take when the criteria is met
							actions = {
								connection = "DENY"
							}
						} // End policy 1
					] // End policies
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rewrite_url.path", "/example"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "DENY"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Add new policy
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_policies_http_security" "example" {
						virtual_service_id = cloudavenue_elb_virtual_service.example.id
						policies = [
							// Policy 1
							{
								name = "example"

								// Define the criteria for the policy
								criteria = {
									path = {
										criteria = "CONTAINS"
										paths    = ["/example"]
									}
								}

								// Define the action to take when the criteria is met
								actions = {
									connection = "ALLOW"
								}
							}, // End policy 1
							// Policy 2
							{
								name = "example2"
								// Define the criteria for the policy
								criteria = {
									protocol = "HTTPS"
								}
								// Define the action to take when the criteria is met
								actions = {
									connection = "DENY"
								}
							} // End policy 2
					] // End policies
					}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.connection", "ALLOW"),

							resource.TestCheckResourceAttr(resourceName, "policies.1.name", "example2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.connection", "DENY"),
						},
					},
					// * Test Update
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_policies_http_security" "example" {
						virtual_service_id = cloudavenue_elb_virtual_service.example.id
						policies = [
							// Policy 1
							{
								name = "example"

								// Define the criteria for the policy
								criteria = {
									path = {
										criteria = "CONTAINS"
										paths    = ["/example/updated"]
									}
								}

								// Define the action to take when the criteria is met
								actions = {
									connection = "DENY"
								}
							}, // End policy 1
							// Policy 2
							{
								name = "example2"
								// Define the criteria for the policy
								criteria = {
									protocol = "HTTP"
								}
								// Define the action to take when the criteria is met
								actions = {
									connection = "ALLOW"
								}
							} // End policy 2
					] // End policies
					}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "DENY"),

							resource.TestCheckResourceAttr(resourceName, "policies.1.name", "example2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.connection", "DENY"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccELBPoliciesHTTPSecurityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ELBPoliciesHTTPSecurityResource{}),
	})
}
