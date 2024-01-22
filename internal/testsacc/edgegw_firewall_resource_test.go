package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &EdgeGatewayFirewallResource{}

const (
	EdgeGatewayFirewallResourceName = testsacc.ResourceName("cloudavenue_edgegateway_firewall")
)

type EdgeGatewayFirewallResource struct{}

func NewEdgeGatewayFirewallResourceTest() testsacc.TestACC {
	return &EdgeGatewayFirewallResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayFirewallResource) GetResourceName() string {
	return EdgeGatewayFirewallResourceName.String()
}

func (r *EdgeGatewayFirewallResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewaySecurityGroupResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayFirewallResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_edgegateway_firewall" "example" {
					  edge_gateway_id = cloudavenue_edgegateway.example.id
					  rules = [
					    {
					      action      = "ALLOW"
					      name        = "allow all IPv4 traffic"
					      direction   = "IN_OUT"
					      ip_protocol = "IPV4"
					    }
					  ]
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.action", "ALLOW"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.name", "allow all IPv4 traffic"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.direction", "IN_OUT"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.ip_protocol", "IPV4"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_edgegateway_firewall" "example" {
							edge_gateway_id = cloudavenue_edgegateway.example.id
							rules = [
							  {
								action      = "ALLOW"
								name        = "allow all IPv4 traffic"
								direction   = "IN"
								ip_protocol = "IPV4"
							  },
							  {
								action      = "ALLOW"
								name        = "allow OUT IPv4 traffic"
								direction   = "OUT"
								ip_protocol = "IPV4"
							  }
							]
						  }`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.action", "ALLOW"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.name", "allow all IPv4 traffic"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.direction", "IN"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.ip_protocol", "IPV4"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.action", "ALLOW"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.name", "allow OUT IPv4 traffic"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.direction", "OUT"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.ip_protocol", "IPV4"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_with_ids": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
						resource "cloudavenue_edgegateway_firewall" "example_with_ids" {
							edge_gateway_id = cloudavenue_edgegateway.example.id
							rules = [{
								action      = "ALLOW"
								name        = "To Internet"
								direction   = "OUT"
								ip_protocol = "IPV4"
							},
							{
								action      = "ALLOW"
								name        = "From Internet to HTTP"
								direction   = "IN"
								ip_protocol = "IPV4"

								destination_ids = [cloudavenue_edgegateway_security_group.example.id]
								app_port_profile_ids = ["urn:vcloud:applicationPortProfile:4d8cc407-fe83-3a9f-af20-95dfe3a1e9a2"]
							},
							{
								action      = "ALLOW"
								name        = "From Internet to HTTPS"
								direction   = "IN"
								ip_protocol = "IPV4"

								destination_ids = [cloudavenue_edgegateway_security_group.example.id]
								app_port_profile_ids = ["urn:vcloud:applicationPortProfile:9c8049b5-9820-36f9-b90c-ab8f462df3c6"]
							}]
						}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),

						resource.TestCheckResourceAttr(resourceName, "rules.#", "3"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.action", "ALLOW"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.name", "To Internet"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.direction", "OUT"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.ip_protocol", "IPV4"),

						resource.TestCheckResourceAttr(resourceName, "rules.1.action", "ALLOW"),
						resource.TestCheckResourceAttr(resourceName, "rules.1.name", "From Internet to HTTP"),
						resource.TestCheckResourceAttr(resourceName, "rules.1.direction", "IN"),
						resource.TestCheckResourceAttr(resourceName, "rules.1.ip_protocol", "IPV4"),
						resource.TestCheckResourceAttr(resourceName, "rules.1.destination_ids.#", "1"),
						resource.TestCheckResourceAttrWith(resourceName, "rules.1.destination_ids.0", uuid.TestIsType(uuid.SecurityGroup)),
						resource.TestCheckResourceAttr(resourceName, "rules.1.app_port_profile_ids.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "rules.1.app_port_profile_ids.0", "urn:vcloud:applicationPortProfile:4d8cc407-fe83-3a9f-af20-95dfe3a1e9a2"),

						resource.TestCheckResourceAttr(resourceName, "rules.2.action", "ALLOW"),
						resource.TestCheckResourceAttr(resourceName, "rules.2.name", "From Internet to HTTPS"),
						resource.TestCheckResourceAttr(resourceName, "rules.2.direction", "IN"),
						resource.TestCheckResourceAttr(resourceName, "rules.2.ip_protocol", "IPV4"),
						resource.TestCheckResourceAttr(resourceName, "rules.2.destination_ids.#", "1"),
						resource.TestCheckResourceAttrWith(resourceName, "rules.2.destination_ids.0", uuid.TestIsType(uuid.SecurityGroup)),
						resource.TestCheckResourceAttr(resourceName, "rules.2.app_port_profile_ids.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "rules.2.app_port_profile_ids.0", "urn:vcloud:applicationPortProfile:9c8049b5-9820-36f9-b90c-ab8f462df3c6"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_edgegateway_firewall" "example_with_ids" {
							edge_gateway_id = cloudavenue_edgegateway.example.id
							rules = [{
								action      = "ALLOW"
								name        = "To Internet"
								direction   = "OUT"
								ip_protocol = "IPV4"
							},
							{
								action      = "ALLOW"
								name        = "From Internet to HTTP"
								direction   = "IN"
								ip_protocol = "IPV4"

								destination_ids = [cloudavenue_edgegateway_security_group.example.id]
								app_port_profile_ids = ["urn:vcloud:applicationPortProfile:4d8cc407-fe83-3a9f-af20-95dfe3a1e9a2"]
							},
							{
								action      = "ALLOW"
								name        = "From Internet to HTTPS"
								direction   = "IN"
								ip_protocol = "IPV4"

								destination_ids = [cloudavenue_edgegateway_security_group.example.id]
								app_port_profile_ids = ["urn:vcloud:applicationPortProfile:9c8049b5-9820-36f9-b90c-ab8f462df3c6"]
							},
							{
								action      = "ALLOW"
								name        = "Allow local traffic"
								direction   = "IN_OUT"
								ip_protocol = "IPV4"

								source_ids = [cloudavenue_edgegateway_security_group.example.id]
								destination_ids = [cloudavenue_edgegateway_security_group.example.id]
							}]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),

							resource.TestCheckResourceAttr(resourceName, "rules.#", "4"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.action", "ALLOW"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.name", "To Internet"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.direction", "OUT"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.ip_protocol", "IPV4"),

							resource.TestCheckResourceAttr(resourceName, "rules.1.action", "ALLOW"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.name", "From Internet to HTTP"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.direction", "IN"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.ip_protocol", "IPV4"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.destination_ids.#", "1"),
							resource.TestCheckResourceAttrWith(resourceName, "rules.1.destination_ids.0", uuid.TestIsType(uuid.SecurityGroup)),
							resource.TestCheckResourceAttr(resourceName, "rules.1.app_port_profile_ids.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.app_port_profile_ids.0", "urn:vcloud:applicationPortProfile:4d8cc407-fe83-3a9f-af20-95dfe3a1e9a2"),

							resource.TestCheckResourceAttr(resourceName, "rules.2.action", "ALLOW"),
							resource.TestCheckResourceAttr(resourceName, "rules.2.name", "From Internet to HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "rules.2.direction", "IN"),
							resource.TestCheckResourceAttr(resourceName, "rules.2.ip_protocol", "IPV4"),
							resource.TestCheckResourceAttr(resourceName, "rules.2.destination_ids.#", "1"),
							resource.TestCheckResourceAttrWith(resourceName, "rules.2.destination_ids.0", uuid.TestIsType(uuid.SecurityGroup)),
							resource.TestCheckResourceAttr(resourceName, "rules.2.app_port_profile_ids.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "rules.2.app_port_profile_ids.0", "urn:vcloud:applicationPortProfile:9c8049b5-9820-36f9-b90c-ab8f462df3c6"),

							resource.TestCheckResourceAttr(resourceName, "rules.3.action", "ALLOW"),
							resource.TestCheckResourceAttr(resourceName, "rules.3.name", "Allow local traffic"),
							resource.TestCheckResourceAttr(resourceName, "rules.3.direction", "IN_OUT"),
							resource.TestCheckResourceAttr(resourceName, "rules.3.ip_protocol", "IPV4"),
							resource.TestCheckResourceAttr(resourceName, "rules.3.source_ids.#", "1"),
							resource.TestCheckResourceAttrWith(resourceName, "rules.3.source_ids.0", uuid.TestIsType(uuid.SecurityGroup)),
							resource.TestCheckResourceAttr(resourceName, "rules.3.destination_ids.#", "1"),
							resource.TestCheckResourceAttrWith(resourceName, "rules.3.destination_ids.0", uuid.TestIsType(uuid.SecurityGroup)),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccEdgeGatewayFirewallResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayFirewallResource{}),
	})
}
