package testsacc

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
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
	// resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_with_id"))
	resp.Append(GetResourceConfig()[ALBPoolResourceName]().GetDefaultConfig)
	// resp.Append(GetDataSourceConfig()[PublicIPResourceName]().GetDefaultConfig)
	return
}

func (r *ALBVirtualServiceResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First with service_type: HTTP and Simple Ports
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
					resource.TestCheckResourceAttrWith(resourceName, "pool_id", urn.TestIsType(urn.LoadBalancerPool)),
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
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestMatchResourceAttr(resourceName, "service_engine_group_name", regexp.MustCompile(`^v[0-9]{3}w[0-9]{2}[e,i]pr[a-z]{12}[0-9]{2}$`)),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "virtual_ip", "192.168.10.10"),
						resource.TestCheckResourceAttr(resourceName, "service_type", "HTTP"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_start", "80"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_end", "80"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_type", "TCP_PROXY"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_ssl", "false"),
					},
				},
				// ! Updates testing with service_type: HTTPS
				// * Replace a port 80 to 443 with SSL and another port 8443 without SSL with service_type: HTTPS
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
								// This is a temporary solution with a real certificate id for this test. A datasource for certificates doesn't exist yet. Will be updated soon.
								certificate_id = "urn:vcloud:certificateLibraryItem:f9caac3a-2555-477e-ae58-0740687d4daf"
								service_type = "HTTPS"
								service_ports = [
									{
										port_start = 443
										port_type = "TCP_PROXY"
										port_ssl = true
									},
									{
										port_start = 8443
									}
								]
							}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
							resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttrWith(resourceName, "certificate_id", urn.TestIsType(urn.CertificateLibraryItem)),
							resource.TestCheckResourceAttr(resourceName, "virtual_ip", "192.168.10.10"),
							resource.TestCheckResourceAttr(resourceName, "service_type", "HTTPS"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "service_ports.*", map[string]string{
								"port_start": "443",
								"port_end":   "443",
								"port_type":  "TCP_PROXY",
								"port_ssl":   "true",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "service_ports.*", map[string]string{
								"port_start": "8443",
								"port_end":   "8443",
								"port_type":  "TCP_PROXY",
								"port_ssl":   "false",
							}),
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
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerVirtualService)),
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
							pool_name = cloudavenue_alb_pool.example.name
							# The datasource for a simple publicip doesn't exist yet: We use a temporary solution with a real Public IP for this test.
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
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttrSet(resourceName, "virtual_ip"),
						resource.TestCheckResourceAttr(resourceName, "service_type", "L4"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_start", "80"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_type", "TCP_PROXY"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_end", "80"),
						resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_ssl", "false"),
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
								pool_name = cloudavenue_alb_pool.example.name
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
							resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", urn.TestIsType(urn.Gateway)),
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttrSet(resourceName, "virtual_ip"),
							resource.TestCheckResourceAttr(resourceName, "service_type", "L4"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_start", "80"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_type", "UDP_FAST_PATH"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_end", "90"),
							resource.TestCheckResourceAttr(resourceName, "service_ports.0.port_ssl", "false"),
						},
					},
				},
				// ! Destroy testing
				Destroy: true,
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
