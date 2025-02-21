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

var _ testsacc.TestACC = &ELBPoliciesHTTPRequestResource{}

const (
	ELBPoliciesHTTPRequestResourceName = testsacc.ResourceName("cloudavenue_elb_policies_http_request")
)

type ELBPoliciesHTTPRequestResource struct{}

func NewELBPoliciesHTTPRequestResourceTest() testsacc.TestACC {
	return &ELBPoliciesHTTPRequestResource{}
}

// GetResourceName returns the name of the resource.
func (r *ELBPoliciesHTTPRequestResource) GetResourceName() string {
	return ELBPoliciesHTTPRequestResourceName.String()
}

func (r *ELBPoliciesHTTPRequestResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[ELBVirtualServiceResourceName]().GetDefaultConfig)
	return
}

func (r *ELBPoliciesHTTPRequestResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
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
					resource "cloudavenue_elb_policies_http_request" "example" {
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
								},
								{
									criteria = "BEGINS_WITH"
									name     = "X-CUSTOM"
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
								modify_headers = [
									{
										action = "ADD"
										name   = "X-SECURE"
										value  = "example"
									},
									{
										action = "REMOVE"
										name   = "X-EXAMPLE"
									}
								]
								rewrite_url = {
									host = "example.com"
									path = "/example"
								}
							}
						} // End policy 1
					] // End policies
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.criteria", "IS_IN"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.ip_addresses.#", "3"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.criteria", "BEGINS_WITH"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.name", "example"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.value", "example"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.criteria", "IS_IN"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.methods.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.protocol", "HTTPS"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.#", "1"),
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
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rewrite_url.host", "example.com"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rewrite_url.path", "/example"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.modify_headers.#", "2"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.actions.modify_headers.*", map[string]string{
							"action": "ADD",
							"name":   "X-SECURE",
							"value":  "example",
						}),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.actions.modify_headers.*", map[string]string{
							"action": "REMOVE",
							"name":   "X-EXAMPLE",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Add new policy
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_policies_http_request" "example" {
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
								},
								{
									criteria = "BEGINS_WITH"
									name     = "X-CUSTOM"
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
								modify_headers = [
									{
										action = "ADD"
										name   = "X-SECURE"
										value  = "example"
									},
									{
										action = "REMOVE"
										name   = "X-EXAMPLE"
									}
								]
								rewrite_url = {
									host = "example.com"
									path = "/example"
								}
							}
						}, // End policy 1
						// Policy 2
							{
							name = "example2"

							// Define the criteria for the policy
							criteria = {
								client_ip = {
									criteria = "IS_IN"
									ip_addresses = [
										"12.13.14.15",
										"12.13.14.0/24"
									]
								}
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
								},
								{
									criteria = "BEGINS_WITH"
									name     = "X-CUSTOM"
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
								modify_headers = [
									{
										action = "ADD"
										name   = "X-SECURE"
										value  = "example"
									},
									{
										action = "REMOVE"
										name   = "X-EXAMPLE"
									}
								]
								rewrite_url = {
									host = "example.com"
									path = "/example"
								}
							}
						} // End policy 2
					] // End policies
					}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.ip_addresses.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.criteria", "BEGINS_WITH"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.value", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.methods.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.#", "1"),
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
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rewrite_url.host", "example.com"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rewrite_url.path", "/example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.modify_headers.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.actions.modify_headers.*", map[string]string{
								"action": "ADD",
								"name":   "X-SECURE",
								"value":  "example",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.actions.modify_headers.*", map[string]string{
								"action": "REMOVE",
								"name":   "X-EXAMPLE",
							}),

							resource.TestCheckResourceAttr(resourceName, "policies.1.name", "example2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.client_ip.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.client_ip.ip_addresses.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.cookie.criteria", "BEGINS_WITH"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.cookie.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.cookie.value", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.http_methods.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.http_methods.methods.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.query.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.request_headers.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.1.criteria.request_headers.*", map[string]string{
								"criteria": "CONTAINS",
								"name":     "X-EXAMPLE",
								"values.#": "1",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.1.criteria.request_headers.*", map[string]string{
								"criteria": "BEGINS_WITH",
								"name":     "X-CUSTOM",
								"values.#": "1",
							}),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.service_ports.ports.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.rewrite_url.host", "example.com"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.rewrite_url.path", "/example"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.modify_headers.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.1.actions.modify_headers.*", map[string]string{
								"action": "ADD",
								"name":   "X-SECURE",
								"value":  "example",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.1.actions.modify_headers.*", map[string]string{
								"action": "REMOVE",
								"name":   "X-EXAMPLE",
							}),
						},
					},
					// * Test Update
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_policies_http_request" "example" {
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
											"12.13.14.0/24"
											// "12.13.14.1-12.13.14.15"
										]
									}
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
									},
									{
										criteria = "BEGINS_WITH"
										name     = "X-CUSTOM"
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
									modify_headers = [
										{
											action = "ADD"
											name   = "X-SECURE"
											value  = "example"
										},
										{
											action = "REMOVE"
											name   = "X-EXAMPLE"
										}
									]
									rewrite_url = {
										host = "example.com"
										path = "/example"
									}
								}
							}, // End policy 1
							// Policy 2
								{
								name = "example2"
	
								// Define the criteria for the policy
								criteria = {
									client_ip = {
										criteria = "IS_IN"
										ip_addresses = [
											"12.13.14.15",
											"12.13.14.0/24"
											// "12.13.14.1-12.13.14.15"
										]
									}
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
									},
									{
										criteria = "BEGINS_WITH"
										name     = "X-CUSTOM"
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
									redirect = {
										port = 80
										protocol = "HTTP"
									}
								}
							} // End policy 2
						] // End policies
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "policies.0.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.client_ip.ip_addresses.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.criteria", "BEGINS_WITH"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.cookie.value", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.http_methods.methods.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.#", "1"),
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
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rewrite_url.host", "example.com"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.rewrite_url.path", "/example"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.modify_headers.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.actions.modify_headers.*", map[string]string{
								"action": "ADD",
								"name":   "X-SECURE",
								"value":  "example",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.actions.modify_headers.*", map[string]string{
								"action": "REMOVE",
								"name":   "X-EXAMPLE",
							}),

							resource.TestCheckResourceAttr(resourceName, "policies.1.name", "example2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.client_ip.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.client_ip.ip_addresses.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.cookie.criteria", "BEGINS_WITH"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.cookie.name", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.cookie.value", "example"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.http_methods.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.http_methods.methods.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.path.criteria", "CONTAINS"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.path.paths.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.query.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.request_headers.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.1.criteria.request_headers.*", map[string]string{
								"criteria": "CONTAINS",
								"name":     "X-EXAMPLE",
								"values.#": "1",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.1.criteria.request_headers.*", map[string]string{
								"criteria": "BEGINS_WITH",
								"name":     "X-CUSTOM",
								"values.#": "1",
							}),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.criteria.service_ports.ports.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.redirect.port", "80"),
							resource.TestCheckResourceAttr(resourceName, "policies.1.actions.redirect.protocol", "HTTP"),
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

func TestAccELBPoliciesHTTPRequestResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ELBPoliciesHTTPRequestResource{}),
	})
}
