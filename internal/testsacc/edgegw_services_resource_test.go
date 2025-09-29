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

	"github.com/orange-cloudavenue/common-go/regex"
	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/helpers"
)

var _ testsacc.TestACC = &EdgeGatewayServicesResource{}

const (
	EdgeGatewayServicesResourceName = testsacc.ResourceName("cloudavenue_edgegateway_services")
)

type EdgeGatewayServicesResource struct{}

func NewEdgeGatewayServicesResourceTest() testsacc.TestACC {
	return &EdgeGatewayServicesResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayServicesResource) GetResourceName() string {
	return EdgeGatewayServicesResourceName.String()
}

func (r *EdgeGatewayServicesResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return resp
}

func (r *EdgeGatewayServicesResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_services" "example" {
						edge_gateway_name = cloudavenue_edgegateway.example.name
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", helpers.TestIsType(urn.EdgeGateway)),
						resource.TestMatchResourceAttr(resourceName, "edge_gateway_name", regex.EdgeGatewayNameRegex()),
						resource.TestCheckResourceAttrSet(resourceName, "network"),
						resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
						resource.TestCheckResourceAttr(resourceName, "services.%", "2"),
						// CAV service Administration
						resource.TestCheckResourceAttr(resourceName, "services.administration.%", "3"),
						resource.TestCheckResourceAttrSet(resourceName, "services.administration.network"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.category", "administration"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.%", "8"),
						// DNS
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.dns-authoritative.%", "5"),
						resource.TestCheckResourceAttrSet(resourceName, "services.administration.services.dns-authoritative.description"),
						resource.TestCheckNoResourceAttr(resourceName, "services.administration.services.dns-authoritative.fqdns"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.dns-authoritative.ips.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.dns-authoritative.name", "dns-authoritative"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.dns-authoritative.ports.#", "2"),
						// DNS Resolver
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.dns-resolver.%", "5"),
						resource.TestCheckResourceAttrSet(resourceName, "services.administration.services.dns-resolver.description"),
						resource.TestCheckNoResourceAttr(resourceName, "services.administration.services.dns-resolver.fqdns"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.dns-resolver.ips.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.dns-resolver.name", "dns-resolver"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.dns-resolver.ports.#", "2"),
						// HTTP Proxy
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.linux-repository.%", "5"),
						resource.TestCheckResourceAttrSet(resourceName, "services.administration.services.linux-repository.description"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.linux-repository.fqdns.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.linux-repository.fqdns.0", "repo.service.cav"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.linux-repository.ips.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.linux-repository.name", "linux-repository"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.linux-repository.ports.#", "1"),
						// NTP
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.ntp.%", "5"),
						resource.TestCheckResourceAttrSet(resourceName, "services.administration.services.ntp.description"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.ntp.fqdns.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.ntp.ips.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.ntp.name", "ntp"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.ntp.ports.#", "1"),
						// RHUI Repository
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.rhui-repository.%", "5"),
						resource.TestCheckResourceAttrSet(resourceName, "services.administration.services.rhui-repository.description"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.rhui-repository.fqdns.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.rhui-repository.fqdns.0", "rhui.service.cav"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.rhui-repository.ips.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.rhui-repository.name", "rhui-repository"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.rhui-repository.ports.#", "1"),
						// SMTP
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.smtp.%", "5"),
						resource.TestCheckResourceAttrSet(resourceName, "services.administration.services.smtp.description"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.smtp.fqdns.0", "smtp.service.cav"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.smtp.ips.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.smtp.name", "smtp"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.smtp.ports.#", "1"),
						// Windows KMS
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-kms.%", "5"),
						resource.TestCheckResourceAttrSet(resourceName, "services.administration.services.windows-kms.description"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-kms.fqdns.0", "kms.service.cav"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-kms.ips.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-kms.name", "windows-kms"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-kms.ports.#", "1"),
						// Windows Repository
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-repository.%", "5"),
						resource.TestCheckResourceAttrSet(resourceName, "services.administration.services.windows-repository.description"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-repository.fqdns.0", "wsus.service.cav"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-repository.ips.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-repository.name", "windows-repository"),
						resource.TestCheckResourceAttr(resourceName, "services.administration.services.windows-repository.ports.#", "2"),
						// CAV service S3 Object Storage
						resource.TestCheckResourceAttr(resourceName, "services.s3.%", "3"),
						resource.TestCheckResourceAttr(resourceName, "services.s3.network", "194.206.55.5/32"),
						resource.TestCheckResourceAttr(resourceName, "services.s3.category", "s3"),
						resource.TestCheckResourceAttr(resourceName, "services.s3.services.%", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.s3.services.s3-internal.%", "5"),
						resource.TestCheckResourceAttrSet(resourceName, "services.s3.services.s3-internal.description"),
						resource.TestCheckResourceAttr(resourceName, "services.s3.services.s3-internal.fqdns.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.s3.services.s3-internal.ips.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "services.s3.services.s3-internal.name", "s3-internal"),
						resource.TestCheckResourceAttr(resourceName, "services.s3.services.s3-internal.ports.#", "1"),
					},
				},
				// ! Update normally is not supported, but we check to enable twice the resource on the same edge gateway
				// This should not fail and leave the resource as is.
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_services" "example" {
							edge_gateway_name = cloudavenue_edgegateway.example.name
						}
						resource "cloudavenue_edgegateway_services" "example_duplicate" {
							edge_gateway_name = cloudavenue_edgegateway.example.name
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "id"),
							resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", helpers.TestIsType(urn.EdgeGateway)),
							resource.TestMatchResourceAttr(resourceName, "edge_gateway_name", regex.EdgeGatewayNameRegex()),
						},
					},
				},
				// ! Import is not supported. Create a new resource instead.
			}
		},
	}
}

func TestAccEdgeGatewayServicesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayServicesResource{}),
	})
}
