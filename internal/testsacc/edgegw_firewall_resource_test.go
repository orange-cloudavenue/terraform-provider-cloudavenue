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
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayFirewallResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First Test For a VDC Backup named "example"
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
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
							"action":      "ALLOW",
							"name":        "allow all IPv4 traffic",
							"direction":   "IN_OUT",
							"ip_protocol": "IPV4",
						}),
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
								direction   = "IN_OUT"
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
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "ALLOW",
								"name":        "allow all IPv4 traffic",
								"direction":   "IN_OUT",
								"ip_protocol": "IPV4",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rules.*", map[string]string{
								"action":      "ALLOW",
								"name":        "allow OUT IPv4 traffic",
								"direction":   "OUT",
								"ip_protocol": "IPV4",
							}),
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
