package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ALBVirtualServiceResource{}

const (
	ALBVirtualServiceResourceName = testsacc.ResourceName("cloudavenue_alb_virtual_service")
)

type ALBVirtualServiceResource struct{}

func NewALBVirtualServiceResourceTest() testsacc.TestACC {
	return &ALBVirtualServiceResource{}
}

// GetResourceName returns the name of the resource.
func (r *ALBVirtualServiceResource) GetResourceName() string {
	return ALBVirtualServiceResourceName.String()
}

func (r *ALBVirtualServiceResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_with_id"))
	resp.Append(GetResourceConfig()[ALBPoolResourceName]().GetDefaultConfig)
	// resp.Append(GetDataSourceConfig()[PublicIPResourceName]().GetDefaultConfig)
	return
}

func (r *ALBVirtualServiceResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// TODO : Complete tests
		// * First with service_type: HTTP and Simple Ports
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "pool_id"),
				},
				// ! Create testing with service_type: HTTP
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_alb_virtual_service" "example" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						edge_gateway_id = data.cloudavenue_edgegateway.example_with_id.id
						enabled = true
						pool_id = cloudavenue_alb_pool.example.id
						virtual_ip = "192.168.10.10"
						service_type = "HTTP"
						service_ports = [
							{
								port_start = 80
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "virtual_ip", "192.168.10.10"),
						resource.TestCheckResourceAttr(resourceName, "service_type", "HTTP"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "service_ports.*", map[string]string{
							"port_start": "80",
						}),
					},
				},
				// ! Updates testing with service_type: HTTPS
				// * Add a new port 443 with SSL and another port 8443 without SSL with service_type: HTTPS
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
							resource "cloudavenue_alb_virtual_service" "example" {
								name = {{ get . "name" }}
								description = {{ generate . "description" }}
								edge_gateway_id = data.cloudavenue_edgegateway.example_with_id.id
								enabled = true
								pool_id = cloudavenue_alb_pool.example.id
								virtual_ip = "192.168.10.10"
								certificate_id = "urn:vcloud:certificateLibraryItem:f9caac3a-2555-477e-ae58-0740687d4daf"
								service_type = "HTTPS"
								service_ports = [
									{
										port_start = 443
										port_type = "TCP_PROXY"
										port_ssl = "true"
									},
									{
										port_start = 8443
									}
								]
							}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "certificate_id", "urn:vcloud:certificateLibraryItem:f9caac3a-2555-477e-ae58-0740687d4daf"),
							resource.TestCheckResourceAttr(resourceName, "virtual_ip", "192.168.10.10"),
							resource.TestCheckResourceAttr(resourceName, "service_type", "HTTPS"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_start", "443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_type", "TCP_PROXY"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_ssl", "true"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.port_start", "8443"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.1.port_ssl", "false"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
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
				// ! Destroy testing
				Destroy: true,
			}
		},
		"example_with_public_ip": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				},
				// ! Create testing
				// * Create a new Virtual Service with a public IP and service_type: L4 (TCP)
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_alb_virtual_service" "example_with_public_ip" {
							name = {{ generate . "name" }}
							description = {{ generate . "description" }}
							edge_gateway_name = data.cloudavenue_edgegateway.example_with_id.name
							enabled = true
							pool_id = cloudavenue_alb_pool.example.id
							virtual_ip = "62.161.20.174"
							service_type = "L4"
							service_ports = [
								{
									port_start = 80
									port_type = "TCP_PROXY"
									port_ssl = false
								}
							]
						}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrSet(resourceName, "pool_id"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttrSet(resourceName, "virtual_ip"),
						resource.TestCheckResourceAttr(resourceName, "service_type", "L4"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "service_ports.*", map[string]string{
							"port_start": "80",
							"port_type":  "TCP_PROXY",
							"port_ssl":   "false",
						}),
					},
				},
				// ! Updates testing
				// * Update service_ports with a range (UDP_FAST_PROXY) and service_type: L4
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
							resource "cloudavenue_alb_virtual_service" "example_with_public_ip" {
								name = {{ get . "name" }}
								description = {{ generate . "description" }}
								edge_gateway_name = data.cloudavenue_edgegateway.example_with_id.name
								enabled = true
								pool_id = cloudavenue_alb_pool.example.id
								virtual_ip = "62.161.20.174"
								service_type = "L4"
								service_ports = [
									{
										port_start = 80
										port_end = 90
										port_type = "UDP_FAST_PATH"
									}
								]
							}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
							resource.TestCheckResourceAttrSet(resourceName, "pool_id"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttrSet(resourceName, "virtual_ip"),
							resource.TestCheckResourceAttr(resourceName, "service_type", "L4"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "service_ports.*", map[string]string{
								"port_start": "80",
								"port_end":   "90",
								"port_type":  "UDP_FAST_PATH",
								"port_ssl":   "false",
							}),
						},
					},
				},
			}
		},
	}
}

func TestAccALBVirtualServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ALBVirtualServiceResource{}),
	})
}
