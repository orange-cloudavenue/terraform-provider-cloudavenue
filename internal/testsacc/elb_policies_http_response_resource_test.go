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

var _ testsacc.TestACC = &ELBPoliciesHTTPResponseResource{}

const (
	ELBPoliciesHTTPResponseResourceName = testsacc.ResourceName("cloudavenue_elb_policies_http_response")
)

type ELBPoliciesHTTPResponseResource struct{}

func NewELBPoliciesHTTPResponseResourceTest() testsacc.TestACC {
	return &ELBPoliciesHTTPResponseResource{}
}

// GetResourceName returns the name of the resource.
func (r *ELBPoliciesHTTPResponseResource) GetResourceName() string {
	return ELBPoliciesHTTPResponseResourceName.String()
}

func (r *ELBPoliciesHTTPResponseResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[ELBVirtualServiceResourceName]().GetDefaultConfig)
	return
}

func (r *ELBPoliciesHTTPResponseResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
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
					resource "cloudavenue_elb_policies_http_response" "example" {
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
								response_headers = [
									{
										criteria = "CONTAINS"
										name     = "X-RESPONSE"
										values    = ["example"]
									},
									{
										criteria = "BEGINS_WITH"
										name     = "X-RESPONSE-CUSTOM"
										values    = ["example"]
									}
								]
								service_ports = {
									criteria = "IS_IN"
									ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
								}
								location = {
									criteria = "BEGINS_WITH"
									values = [
										"example.com"
									]
								}
								status_code = {
									criteria = "IS_IN"
									codes    = ["200","302"]
								}
							}

							// Define the action to take when the criteria is met
							actions = {
								
								location_rewrite = {
									host      = "example.org"
									protocol  = "HTTPS"
									keep_query = true
									port      = 443
								}
								
								modify_headers = [
									{
										action = "ADD"
										name   = "X-FROM-OLD-DOMAIN"
										value  = "example.com"
									}
								]
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
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.response_headers.#", "2"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.response_headers.*", map[string]string{
							"criteria": "CONTAINS",
							"name":     "X-RESPONSE",
							"values.#": "1",
						}),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.response_headers.*", map[string]string{
							"criteria": "BEGINS_WITH",
							"name":     "X-RESPONSE-CUSTOM",
							"values.#": "1",
						}),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.location.criteria", "BEGINS_WITH"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.location.values.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.status_code.criteria", "IS_IN"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.status_code.codes.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.location_rewrite.host", "example.org"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.location_rewrite.protocol", "HTTPS"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.location_rewrite.keep_query", "true"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.location_rewrite.port", "443"),
						resource.TestCheckResourceAttr(resourceName, "policies.0.actions.modify_headers.#", "1"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.actions.modify_headers.*", map[string]string{
							"action": "ADD",
							"name":   "X-FROM-OLD-DOMAIN",
							"value":  "example.com",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Test error
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_policies_http_response" "example" {
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
								response_headers = [
									{
										criteria = "CONTAINS"
										name     = "X-RESPONSE"
										values    = ["example"]
									},
									{
										criteria = "BEGINS_WITH"
										name     = "X-RESPONSE-CUSTOM"
										values    = ["example"]
									}
								]
								service_ports = {
									criteria = "IS_IN"
									ports    = ["80"] // Only 80 because only port 80 is set in the virtual service
								}
								location = {
									criteria = "BEGINS_WITH"
									values = [
										"example.com"
									]
								}
								status_code = {
									criteria = "IS_IN"
									codes    = ["200","302"]
								}
							}

							// Define the action to take when the criteria is met
							actions = {
								
								location_rewrite = {
									host      = "example.org"
									protocol  = "HTTPS"
									keep_query = true
									port      = 443
								}
								
								modify_headers = [
									{
										action = "ADD"
										name   = "X-FROM-OLD-DOMAIN"
										value  = "example.com"
									},
																		{
										action = "ADD"
										name   = "X-FROM-OLD-DOMAIN2"
										value  = "example.org"
									},
								]
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
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.response_headers.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.response_headers.*", map[string]string{
								"criteria": "CONTAINS",
								"name":     "X-RESPONSE",
								"values.#": "1",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.criteria.response_headers.*", map[string]string{
								"criteria": "BEGINS_WITH",
								"name":     "X-RESPONSE-CUSTOM",
								"values.#": "1",
							}),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.service_ports.ports.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.location.criteria", "BEGINS_WITH"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.location.values.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.status_code.criteria", "IS_IN"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.status_code.codes.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.criteria.query.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.location_rewrite.host", "example.org"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.location_rewrite.protocol", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.location_rewrite.keep_query", "true"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.location_rewrite.port", "443"),
							resource.TestCheckResourceAttr(resourceName, "policies.0.actions.modify_headers.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.actions.modify_headers.*", map[string]string{
								"action": "ADD",
								"name":   "X-FROM-OLD-DOMAIN",
								"value":  "example.com",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "policies.0.actions.modify_headers.*", map[string]string{
								"action": "ADD",
								"name":   "X-FROM-OLD-DOMAIN2",
								"value":  "example.org",
							}),
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

func TestAccELBPoliciesHTTPResponseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ELBPoliciesHTTPResponseResource{}),
	})
}
