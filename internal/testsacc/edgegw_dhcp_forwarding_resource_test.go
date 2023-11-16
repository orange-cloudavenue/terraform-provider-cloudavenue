package testsacc

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccDhcpForwardingResourceConfig = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10"
	]
}
`

const testAccDhcpForwardingResourceConfigUpdate = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10",
		"192.168.10.11"
	]
}
`

const testAccDhcpForwardingResourceConfigUpdateError = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
	enabled = false
	dhcp_servers = [
		"192.168.10.10"
	]
}
`

const testAccDhcpForwardingResourceConfigWithVDCGroup = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_group.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10",
	]
}
`

const testAccDhcpForwardingResourceConfigWithVDCGroupUpdate = `
resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
	edge_gateway_id = cloudavenue_edgegateway.example_with_group.id
	enabled = true
	dhcp_servers = [
		"192.168.10.10",
		"192.168.10.11"
	]
}
`

func dhcpForwardingTestCheck(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "1"),
	)
}

func dhcpForwardingTestCheckWithVDCGroup(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
		resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "2"),
	)
}

func TestAccDhcpForwardingResource(t *testing.T) {
	resourceName := "cloudavenue_edgegateway_dhcp_forwarding.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// * Test with VDC
			{
				// Apply
				Config: ConcatTests(testAccEdgeGatewayResourceConfig, testAccDhcpForwardingResourceConfig),
				Check:  dhcpForwardingTestCheck(resourceName),
			},
			{
				// Update
				Config: ConcatTests(testAccEdgeGatewayResourceConfig, testAccDhcpForwardingResourceConfigUpdate),
				Check:  dhcpForwardingTestCheckWithVDCGroup(resourceName),
			},
			{
				// Import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// check error when updating dhcp_servers if enabled is false
				Config:             ConcatTests(testAccEdgeGatewayResourceConfig, testAccDhcpForwardingResourceConfigUpdateError),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile("DHCP Servers cannot be edited"),
			},
			// Destroy test with VDC
			{
				Destroy: true,
				Config:  ConcatTests(testAccEdgeGatewayResourceConfig, testAccDhcpForwardingResourceConfigUpdate),
			},

			// * Test with VDC group
			{
				// Apply
				Config: ConcatTests(testAccEdgeGatewayGroupResourceConfig, testAccDhcpForwardingResourceConfigWithVDCGroup),
				Check:  dhcpForwardingTestCheck(resourceName),
			},
			{
				// Update
				Config: ConcatTests(testAccEdgeGatewayGroupResourceConfig, testAccDhcpForwardingResourceConfigWithVDCGroupUpdate),
				Check:  dhcpForwardingTestCheckWithVDCGroup(resourceName),
			},
			{
				// Import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var _ testsacc.TestACC = &EdgeGatewayDhcpForwardingResource{}

const (
	EdgeGatewayDhcpForwardingResourceName = testsacc.ResourceName("data.cloudavenue_backup")
)

type EdgeGatewayDhcpForwardingResource struct{}

func NewEdgeGatewayDhcpForwardingResourceTest() testsacc.TestACC {
	return &EdgeGatewayDhcpForwardingResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayDhcpForwardingResource) GetResourceName() string {
	return EdgeGatewayDhcpForwardingResourceName.String()
}

func (r *EdgeGatewayDhcpForwardingResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayDhcpForwardingResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (backup vdc example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
						edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
						enabled = true
						dhcp_servers = [
							"192.168.10.10"
						]
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.10"),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
							enabled = true
							dhcp_servers = [
								"192.168.10.10",
								"192.168.10.11"
							]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "2"),
							resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.10"),
							resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.11"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
							enabled = false
							dhcp_servers = [
								"192.168.10.10",
								"192.168.10.11"
							]
						}`,
						TFAdvanced: testsacc.TFAdvanced{
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
							ExpectError:        regexp.MustCompile("DHCP Servers cannot be edited"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			}
		},
		"example_with_vdc_group": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Gateway)),
					resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
				},
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_edgegateway_dhcp_forwarding" "example_with_vdc_group" {
						edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
						enabled = true
						dhcp_servers = [
							"192.168.10.10"
						]
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.10"),
					},
				},
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
							enabled = true
							dhcp_servers = [
								"192.168.10.10",
								"192.168.10.11"
							]
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "dhcp_servers.#", "2"),
							resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.10"),
							resource.TestCheckTypeSetElemAttr(resourceName, "dhcp_servers.*", "192.168.10.11"),
						},
					},
					{
						TFConfig: `
						resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
							edge_gateway_id = cloudavenue_edgegateway.example_with_vdc.id
							enabled = false
							dhcp_servers = [
								"192.168.10.10",
								"192.168.10.11"
							]
						}`,
						TFAdvanced: testsacc.TFAdvanced{
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
							ExpectError:        regexp.MustCompile("DHCP Servers cannot be edited"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			}
		},
	}
}

func TestAccEdgeGatewayDhcpForwardingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayDhcpForwardingResource{}),
	})
}
