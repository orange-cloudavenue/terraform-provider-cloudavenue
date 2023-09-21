package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &NetworkRoutedResource{}

const (
	NetworkRoutedResourceName = testsacc.ResourceName("cloudavenue_network_routed")
)

type NetworkRoutedResource struct{}

func NewNetworkRoutedResourceTest() testsacc.TestACC {
	return &NetworkRoutedResource{}
}

// GetResourceName returns the name of the resource.
func (r *NetworkRoutedResource) GetResourceName() string {
	return NetworkRoutedResourceName.String()
}

func (r *NetworkRoutedResource) DependenciesConfig() (configs testsacc.TFData) {
	configs.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig())
	return
}

func (r *NetworkRoutedResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.Network)),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_network_routed" "example" {
						name        = {{ generate . "name" }}
						description = {{ generate . "description" }}
					  
						edge_gateway_id = cloudavenue_edgegateway.example.id
					  
						gateway       = "192.168.1.254"
						prefix_length = 24
					  
						dns1 = "1.1.1.1"
						dns2 = "8.8.8.8"
					  
						dns_suffix = "example"
					  
						static_ip_pool = [
						  {
							start_address = "192.168.1.10"
							end_address   = "192.168.1.20"
						  }
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.1.254"),
						resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.1"),
						resource.TestCheckResourceAttr(resourceName, "dns2", "8.8.8.8"),
						resource.TestCheckResourceAttr(resourceName, "dns_suffix", "example"),
						resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.start_address", "192.168.1.10"),
						resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.end_address", "192.168.1.20"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_network_routed" "example" {
							name        = {{ get . "name" }}
							description = {{ get . "description" }}
						  
							edge_gateway_id = cloudavenue_edgegateway.example.id
						  
							gateway       = "192.168.1.250"
							prefix_length = 24
						  
							dns1 = "1.1.1.2"
							dns2 = "8.8.8.9"
						  
							dns_suffix = "exampleupdated"
						  
							static_ip_pool = [
							  {
								start_address = "192.168.1.1"
								end_address   = "192.168.1.30"
							  }
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.1.250"),
							resource.TestCheckResourceAttr(resourceName, "dns1", "1.1.1.2"),
							resource.TestCheckResourceAttr(resourceName, "dns2", "8.8.8.9"),
							resource.TestCheckResourceAttr(resourceName, "dns_suffix", "exampleupdated"),
							resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.start_address", "192.168.1.1"),
							resource.TestCheckResourceAttr(resourceName, "static_ip_pool.0.end_address", "192.168.1.30"),
						},
					},
				},
				// ! Imports testing
				// TODO : Add import test after resolving this issue https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/526
				// Imports: []testsacc.TFImport{
				// 	{
				// 		ImportStateIDBuilder: []string{"edge_gateway_name", "name"},
				// 		ImportState:          true,
				// 		ImportStateVerify:    true,
				// 	},
				// },
			}
		},
	}
	// TODO: ADD Test with VDC Group
}

func TestAccNetworkRoutedResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&NetworkRoutedResource{}),
	})
}
