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
				// ? Create a minimalist resource (1)
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
									connection = "CLOSE"
								}
							} // End policy 1
						] // End policies
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "CLOSE"),
					},
				},
				// ! Update criteria and action
				Updates: []testsacc.TFConfig{
					// ? Update criteria path (2)
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
											criteria = "EQUALS"
											paths    = ["/example/updated"]
										}
									}
						
									// Define the action to take when the criteria is met
									actions = {
										connection = "ALLOW"
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "EQUALS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "ALLOW"),
						},
					},
					// ? Update criteria to client_ip (3)
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
										client_ip = {
											criteria = "IS_IN"
											ip_addresses = [
												"12.13.14.15",
												"12.13.14.0/24",
												"12.13.14.1-12.13.14.15"
											]
										}	
									}
						
									// Define the action to take when the criteria is met
									actions = {
										connection = "ALLOW"
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.ip_addresses.#", "3"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.criteria.path.criteria"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "ALLOW"),
						},
					},
					// ? Update criteria to cookie (4)
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
										cookie = {
											criteria = "BEGINS_WITH"
											name     = "example"
											value    = "example"
										}
									}
						
									// Define the action to take when the criteria is met
									actions = {
										connection = "ALLOW"
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.criteria.client_ip.criteria"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.criteria", "BEGINS_WITH"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.value", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "ALLOW"),
						},
					},
					// ? Update criteria to http_methods (5)
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
										http_methods = {
											criteria = "IS_IN"
												methods  = ["GET", "POST"]
											}
										}
						
										// Define the action to take when the criteria is met
									actions = {
										connection = "ALLOW"
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.criteria.cookie.criteria"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.methods.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "ALLOW"),
						},
					},
					// ? Update criteria to protocol (6)
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
										protocol = "HTTPS"
									}
						
									// Define the action to take when the criteria is met
									actions = {
										connection = "ALLOW"
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.criteria.http_methods.criteria"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "ALLOW"),
						},
					},
					// ? Update criteria to query (7)
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
										query = [
											"example=example"
										]
									}
						
									// Define the action to take when the criteria is met
									actions = {
										connection = "ALLOW"
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.criteria.protocol"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.0", "example=example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "ALLOW"),
						},
					},
					// ? Update criteria to request_headers (8)
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
										request_headers = [
											{
												criteria = "EXISTS"
												name     = "X-EXAMPLE"
											},
											{
												criteria = "CONTAINS"
												name     = "X-EXAMPLE"
												values    = ["example"]
											},
											{
												criteria = "BEGINS_WITH"
												name     = "X-CUSTOM"
												values    = ["example"]
											}
										]
									}
						
									// Define the action to take when the criteria is met
									actions = {
										connection = "ALLOW"
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.criteria.query"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.request_headers.#", "3"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.request_headers.*", map[string]string{
								"criteria": "CONTAINS",
								"name":     "X-EXAMPLE",
								"values.#": "1",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.request_headers.*", map[string]string{
								"criteria": "BEGINS_WITH",
								"name":     "X-CUSTOM",
								"values.#": "1",
							}),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "ALLOW"),
						},
					},
					// ? Update crto service_ports (9)
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
										service_ports = {
											criteria = "IS_IN"
											ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
										}
									}
						
									// Define the action to take when the criteria is met
									actions = {
										connection = "ALLOW"
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.criteria.request_headers"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "ALLOW"),
						},
					},
					// ? Update Action to redirect_to_https (10)
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
										service_ports = {
											criteria = "IS_IN"
											ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
										}
									}

									// Define the action to take when the criteria is met
									actions = {
										redirect_to_https = 8443
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.actions.connection"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.redirect_to_https", "8443"),
						},
					},
					// ? Update Action to send_response (11)
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
										service_ports = {
											criteria = "IS_IN"
											ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
										}
									}

									// Define the action to take when the criteria is met
									actions = {
										send_response = {
											status_code = 200
											content     = base64encode("example")
											content_type = "text/plain"
										}
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.actions.redirect_to_https"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.send_response.status_code", "200"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.send_response.content", "ZXhhbXBsZQ=="),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.send_response.content_type", "text/plain"),
						},
					},
					// ? Update Action to rate_limit with close_connection (12)
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
										service_ports = {
											criteria = "IS_IN"
											ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
										}
									}

									// Define the action to take when the criteria is met
									actions = {
										rate_limit = {
											count  = 100
											period = 10
											close_connection = true
										}
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.actions.send_response"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.count", "100"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.period", "10"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.close_connection", "true"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.actions.rate_limit.local_response"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.actions.rate_limit.redirect"),
						},
					},
					// ? Update Action to rate_limit with local_response (13)
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
										service_ports = {
											criteria = "IS_IN"
											ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
										}
									}

									// Define the action to take when the criteria is met
									actions = {
										rate_limit = {
											count  = 100
											period = 10
											local_response = {
												status_code  = 200
												content 	 = base64encode("example")
												content_type = "text/plain"
											}
										}
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.count", "100"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.period", "10"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.actions.rate_limit.close_connection"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.local_response.status_code", "200"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.local_response.content", "ZXhhbXBsZQ=="),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.local_response.content_type", "text/plain"),
						},
					},
					// ? Update Action to rate_limit with redirect (14)
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
										service_ports = {
											criteria = "IS_IN"
											ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
										}
									}

									// Define the action to take when the criteria is met
									actions = {
										rate_limit = {
											count  = 100
											period = 10
											redirect = {
												port  = 8443
												protocol = "HTTPS"
											}
										}
									}
								}
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.count", "100"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.period", "10"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.actions.rate_limit.local_response"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.redirect.port", "8443"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rate_limit.redirect.protocol", "HTTPS"),
						},
					},
				},
				// ! Imports testing
				// ? Import the resource (15)
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_multiple_criteria": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					return
				},
				// ! Create testing
				// ? Create a resource with multiple criteria (16)
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_policies_http_security" "example_multiple_criteria" {
						virtual_service_id = cloudavenue_elb_virtual_service.example.id
						policies = [
							// Policy 1
							{
							name = "example"

							criteria = {
								cookie = {
									criteria = "BEGINS_WITH"
									name     = "example"
									value    = "example"
								}
								http_methods = {
									criteria = "IS_IN"
									methods  = ["GET", "POST"]
								}
								path = {
									criteria = "CONTAINS"
									paths    = ["/example"]
								}
								protocol = "HTTPS"
								query = [
									"example=example"
								]
								request_headers = [
									{
										criteria = "CONTAINS"
										name     = "X-EXAMPLE"
										values    = ["example"]
									}
								]
								service_ports = {
									criteria = "IS_IN"
									ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
								}
							}

							// Define the action to take when the criteria is met
							actions = {
								connection = "CLOSE"
							}
						} // End policy 1
					] // End policies
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.connection", "CLOSE"),
					},
				},
				// ! Update criteria and action
				Updates: []testsacc.TFConfig{
					// ? Update criteria values (17)
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_policies_http_security" "example_multiple_criteria" {
							virtual_service_id = cloudavenue_elb_virtual_service.example.id
							policies = [
								// Policy 1
								{
								name = "example"

								criteria = {
									cookie = {
										criteria = "BEGINS_WITH"
										name     = "example"
										value    = "example_updated"
									}
									http_methods = {
										criteria = "IS_IN"
										methods  = ["GET", "POST", "PUT"]
									}
									path = {
										criteria = "CONTAINS"
										paths    = ["/example_updated"]
									}
									protocol = "HTTPS"
									query = [
										"example=example_updated"
									]
									request_headers = [
										{
											criteria = "CONTAINS"
											name     = "X-EXAMPLE"
											values    = ["example_updated"]
										},
										{
											criteria = "BEGINS_WITH"
											name     = "X-CUSTOM"
											values    = ["example_updated"]
										}
									]
									service_ports = {
										criteria = "IS_IN"
										ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
									}
								}

								// Define the action to take when the criteria is met
								actions = {
										redirect_to_https = 8443
								}
							} // End policy 1
						] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.criteria", "BEGINS_WITH"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.value", "example_updated"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.methods.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.0", "example=example_updated"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.request_headers.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.request_headers.*", map[string]string{
								"criteria": "CONTAINS",
								"name":     "X-EXAMPLE",
								"values.#": "1",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.request_headers.*", map[string]string{
								"criteria": "BEGINS_WITH",
								"name":     "X-CUSTOM",
								"values.#": "1",
							}),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.0", "80"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.actions.connection"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.redirect_to_https", "8443"),
						},
					},
					// ? Update criteria to 2 policies (18)
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_policies_http_security" "example_multiple_criteria" {
							virtual_service_id = cloudavenue_elb_virtual_service.example.id
							policies = [
								// Policy 1
								{
									name = "example"
									criteria = {
										client_ip = {
											criteria = "IS_IN"
											ip_addresses = [
												"12.13.14.15",
												"12.13.14.0/24",
												"12.13.14.1-12.13.14.15"
											]
										}
										cookie = {
											criteria = "BEGINS_WITH"
											name     = "example"
											value    = "example_updated"
										}
										http_methods = {
											criteria = "IS_IN"
											methods  = ["GET", "POST", "PUT"]
										}
										path = {
											criteria = "CONTAINS"
											paths    = ["/example_updated"]
										}
										protocol = "HTTPS"
										query = [
											"example=example_updated"
										]
										request_headers = [
											{
												criteria = "CONTAINS"
												name     = "X-EXAMPLE"
												values    = ["example_updated"]
											},
											{
												criteria = "BEGINS_WITH"
												name     = "X-CUSTOM"
												values    = ["example_updated"]
											}
										]
										service_ports = {
											criteria = "IS_IN"
											ports    = ["80"]
										}
									}
									// Define the action to take when the criteria is met
									actions = {
										redirect_to_https = 8443
									}
								}, // End policy 1
								// Policy 2
								{
									name = "example2"
									criteria = {
										cookie = {
											criteria = "BEGINS_WITH"
											name     = "example"
											value    = "example2"
										}
										http_methods = {
											criteria = "IS_IN"
											methods  = ["GET", "POST", "PUT"]
										}
										path = {
											criteria = "CONTAINS"
											paths    = ["/example2"]
										}
										protocol = "HTTP"
										query = [
											"example=example2"
										]
										request_headers = [
											{
												criteria = "CONTAINS"
												name     = "X-EXAMPLE"
												values    = ["example2"]
											}
										]
										service_ports = {
											criteria = "IS_IN"
											ports    = ["80"]
										}
									}
									// Define the action to take when the criteria is met
									actions = {
										rate_limit = {
											count  = 100
											period = 10
											local_response = {
												status_code  = 200
												content 	 = base64encode("example")
												content_type = "text/plain"
											}
										}	
									}
								} // End policy 2
							] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.ip_addresses.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.criteria", "BEGINS_WITH"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.value", "example_updated"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.methods.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.0", "example=example_updated"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.request_headers.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.request_headers.*", map[string]string{
								"criteria": "CONTAINS",
								"name":     "X-EXAMPLE",
								"values.#": "1",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.request_headers.*", map[string]string{
								"criteria": "BEGINS_WITH",
								"name":     "X-CUSTOM",
								"values.#": "1",
							}),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.0", "80"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.0.actions.connection"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.redirect_to_https", "8443"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.name", "example2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.cookie.criteria", "BEGINS_WITH"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.cookie.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.cookie.value", "example2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.http_methods.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.http_methods.methods.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.protocol", "HTTP"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.query.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.query.0", "example=example2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.request_headers.#", "1"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.1.criteria.request_headers.*", map[string]string{
								"criteria": "CONTAINS",
								"name":     "X-EXAMPLE",
								"values.#": "1",
							}),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.service_ports.ports.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.service_ports.ports.0", "80"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.1.actions.redirect_to_https"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.rate_limit.count", "100"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.rate_limit.period", "10"),
							resource.TestCheckNoResourceAttr(resourceName, "policies.1.actions.rate_limit.close_connection"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.rate_limit.local_response.status_code", "200"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.rate_limit.local_response.content", "ZXhhbXBsZQ=="),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.rate_limit.local_response.content_type", "text/plain"),
						},
					},
				},
				// ! Delete testing
				Destroy: true,
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
