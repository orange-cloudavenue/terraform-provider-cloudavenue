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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ELBVirtualServiceResource{}

const (
	ELBVirtualServiceResourceName = testsacc.ResourceName("cloudavenue_elb_virtual_service")
)

type ELBVirtualServiceResource struct{}

func NewELBVirtualServiceResourceTest() testsacc.TestACC {
	return &ELBVirtualServiceResource{}
}

// GetResourceName returns the name of the resource.
func (r *ELBVirtualServiceResource) GetResourceName() string {
	return ELBVirtualServiceResourceName.String()
}

func (r *ELBVirtualServiceResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_for_elb"))
	resp.Append(GetResourceConfig()[ELBPoolResourceName]().GetDefaultConfig)
	return resp
}

func (r *ELBVirtualServiceResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Example with service_type: HTTP and Simple Ports
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
					resource.TestCheckResourceAttrWith(resourceName, "pool_id", urn.TestIsType(urn.LoadBalancerPool)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					resource.TestCheckResourceAttrSet(resourceName, "pool_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_virtual_service" "example" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
						enabled = true
						pool_id = cloudavenue_elb_pool.example.id
						virtual_ip = {{ generate . "virtual_ip" "private-ipv4" }}
						service_type = "HTTP"
						service_ports = [
							{
								start = 80
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

						resource.TestCheckResourceAttr(resourceName, "service_type", "HTTP"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "80"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "80"), // port end = port start if not specified
						resource.TestCheckNoResourceAttr(resourceName, "certificate_id"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Test error (bad service_type)
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_virtual_service" "example" {
							name = {{ get . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							pool_id = cloudavenue_elb_pool.example.id
							virtual_ip = {{ get . "virtual_ip" }}
							service_type = "HTTE"
							service_ports = [
								{
									start = 80
								}
							]
						}`),
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							PlanOnly:           true,
							ExpectError:        regexp.MustCompile(`Attribute service_type value must be one of`),
						},
					},
					// * Test Update (Add a port & change virtual_ip)
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_virtual_service" "example" {
							name = {{ get . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							pool_id = cloudavenue_elb_pool.example.id
							virtual_ip = {{ generate . "virtual_ip" "private-ipv4" }}
							service_type = "HTTP"
							service_ports = [
								{
									start = 80
								},
								{
									start = 8080
									end   = 8090
								},
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),
							resource.TestCheckResourceAttr(resourceName, "service_type", "HTTP"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "80"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "80"), // port end = port start if not specified
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.start", "8080"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.end", "8090"),
							resource.TestCheckNoResourceAttr(resourceName, "certificate_id"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		// * Example with service_type: HTTPS / Certificate / PublicIP
		"example_https": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
					resource.TestCheckResourceAttrWith(resourceName, "pool_id", urn.TestIsType(urn.LoadBalancerPool)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					resource.TestCheckResourceAttrSet(resourceName, "pool_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[ORGCertificateLibraryResourceName]().GetDefaultConfig)
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_virtual_service" "example_https" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
						enabled = true
						pool_id = cloudavenue_elb_pool.example.id
						virtual_ip = {{ generate . "virtual_ip" "public-ipv4" }}
						service_type = "HTTPS"
						certificate_id = cloudavenue_org_certificate_library.example.id
						service_ports = [
							{
								start = 443
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

						resource.TestCheckResourceAttr(resourceName, "service_type", "HTTPS"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "443"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "443"), // port end = port start if not specified
						resource.TestCheckResourceAttrSet(resourceName, "certificate_id"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Update description and add a port
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_virtual_service" "example_https" {
							name = {{ get . "name" }}
							description = {{ generate . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							pool_id = cloudavenue_elb_pool.example.id
							virtual_ip = {{ get . "virtual_ip" }}

							service_type = "HTTPS"
							certificate_id = cloudavenue_org_certificate_library.example.id
							service_ports = [
								{
									start = 443
								},
								{
									start = 8443
									end   = 8446
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

							resource.TestCheckResourceAttr(resourceName, "service_type", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "443"), // port end = port start if not specified
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.start", "8443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.end", "8446"),
							resource.TestCheckResourceAttrSet(resourceName, "certificate_id"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
				Destroy: true,
			}
		},

		"example_l_4_tls": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
					resource.TestCheckResourceAttrWith(resourceName, "pool_id", urn.TestIsType(urn.LoadBalancerPool)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					resource.TestCheckResourceAttrSet(resourceName, "pool_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[ORGCertificateLibraryResourceName]().GetDefaultConfig)
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_virtual_service" "example_l_4_tls" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
						enabled = true
						pool_id = cloudavenue_elb_pool.example.id
						virtual_ip = {{ generate . "virtual_ip" "public-ipv4" }}
						service_type = "L4_TLS"
						certificate_id = cloudavenue_org_certificate_library.example.id
						service_ports = [
							{
								start = 443
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

						resource.TestCheckResourceAttr(resourceName, "service_type", "L4_TLS"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "443"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "443"), // port end = port start if not specified
						resource.TestCheckResourceAttrSet(resourceName, "certificate_id"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Update name and add a port
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_virtual_service" "example_l_4_tls" {
							name = {{ generate . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							pool_id = cloudavenue_elb_pool.example.id
							virtual_ip = {{ get . "virtual_ip" }}

							service_type = "L4_TLS"
							certificate_id = cloudavenue_org_certificate_library.example.id
							service_ports = [
								{
									start = 443
								},
								{
									start = 8443
									end   = 8446
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

							resource.TestCheckResourceAttr(resourceName, "service_type", "L4_TLS"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "443"), // port end = port start if not specified
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.start", "8443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.end", "8446"),
							resource.TestCheckResourceAttrSet(resourceName, "certificate_id"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_l_4_tcp": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
					resource.TestCheckResourceAttrWith(resourceName, "pool_id", urn.TestIsType(urn.LoadBalancerPool)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					resource.TestCheckResourceAttrSet(resourceName, "pool_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[ORGCertificateLibraryResourceName]().GetDefaultConfig)
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_virtual_service" "example_l_4_tcp" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
						enabled = true
						pool_id = cloudavenue_elb_pool.example.id
						virtual_ip = {{ generate . "virtual_ip" "public-ipv4" }}
						service_type = "L4_TCP"
						service_ports = [
							{
								start = 443
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

						resource.TestCheckResourceAttr(resourceName, "service_type", "L4_TCP"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "443"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "443"), // port end = port start if not specified
						resource.TestCheckNoResourceAttr(resourceName, "certificate_id"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Update name and add a port
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_virtual_service" "example_l_4_tcp" {
							name = {{ generate . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							pool_id = cloudavenue_elb_pool.example.id
							virtual_ip = {{ get . "virtual_ip" }}

							service_type = "L4_TCP"
							service_ports = [
								{
									start = 443
								},
								{
									start = 8443
									end   = 8446
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

							resource.TestCheckResourceAttr(resourceName, "service_type", "L4_TCP"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "443"), // port end = port start if not specified
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.start", "8443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.end", "8446"),
							resource.TestCheckNoResourceAttr(resourceName, "certificate_id"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_l_4_udp": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
					resource.TestCheckResourceAttrWith(resourceName, "pool_id", urn.TestIsType(urn.LoadBalancerPool)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					resource.TestCheckResourceAttrSet(resourceName, "pool_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[ORGCertificateLibraryResourceName]().GetDefaultConfig)
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_virtual_service" "example_l_4_udp" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
						enabled = true
						pool_id = cloudavenue_elb_pool.example.id
						virtual_ip = {{ generate . "virtual_ip" "public-ipv4" }}
						service_type = "L4_UDP"
						service_ports = [
							{
								start = 443
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

						resource.TestCheckResourceAttr(resourceName, "service_type", "L4_UDP"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "443"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "443"), // port end = port start if not specified
						resource.TestCheckNoResourceAttr(resourceName, "certificate_id"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Update name and add a port
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_elb_virtual_service" "example_l_4_udp" {
							name = {{ generate . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
							enabled = true
							pool_id = cloudavenue_elb_pool.example.id
							virtual_ip = {{ get . "virtual_ip" }}

							service_type = "L4_UDP"
							service_ports = [
								{
									start = 443
								},
								{
									start = 8443
									end   = 8446
								}
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

							resource.TestCheckResourceAttr(resourceName, "service_type", "L4_UDP"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "443"), // port end = port start if not specified
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.start", "8443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.end", "8446"),
							resource.TestCheckNoResourceAttr(resourceName, "certificate_id"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_with_service_engine_group": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
					resource.TestCheckResourceAttrWith(resourceName, "pool_id", urn.TestIsType(urn.LoadBalancerPool)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					resource.TestCheckResourceAttrSet(resourceName, "pool_name"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetDataSourceConfig()[ELBServiceEngineGroupsDataSourceName]().GetDefaultConfig)
					return resp
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_elb_virtual_service" "example_with_service_engine_group" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_id = data.cloudavenue_edgegateway.example_for_elb.id
						enabled = true
						pool_id = cloudavenue_elb_pool.example.id
						virtual_ip = {{ generate . "virtual_ip" "public-ipv4" }}

						service_type = "L4_UDP"
						service_ports = [
							{
								start = 443
							}
						]

						service_engine_group_name = data.cloudavenue_elb_service_engine_groups.example.service_engine_groups[0].name
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "virtual_ip", testsacc.GetValueFromTemplate(resourceName, "virtual_ip")),

						resource.TestCheckResourceAttr(resourceName, "service_type", "L4_UDP"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.start", "443"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.end", "443"), // port end = port start if not specified
						resource.TestCheckNoResourceAttr(resourceName, "certificate_id"),
					},
				},
			}
		},
	}
}

func TestAccELBVirtualServiceResource(t *testing.T) {
	cleanup := orgCertificateLibraryResourcePreCheck()
	defer cleanup()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ELBVirtualServiceResource{}),
	})
}
