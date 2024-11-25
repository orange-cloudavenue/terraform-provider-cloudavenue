package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &EdgeGatewayAppPortProfileResource{}

const (
	EdgeGatewayAppPortProfileResourceName = testsacc.ResourceName("cloudavenue_edgegateway_app_port_profile")
)

type EdgeGatewayAppPortProfileResource struct{}

func NewEdgeGatewayAppPortProfileResourceTest() testsacc.TestACC {
	return &EdgeGatewayAppPortProfileResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayAppPortProfileResource) GetResourceName() string {
	return EdgeGatewayAppPortProfileResourceName.String()
}

func (r *EdgeGatewayAppPortProfileResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayAppPortProfileResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.AppPortProfile)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway_app_port_profile" "example" {
					  name = {{ generate . "name" }}
					  description = {{ generate . "description" }}
					  edge_gateway_id = cloudavenue_edgegateway.example.id
					  app_ports = [
					    {
					    	protocol = "ICMPv4"
					    }
					  ]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.AppPortProfile)),
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "app_ports.*", map[string]string{
							"protocol": "ICMPv4",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway_app_port_profile" "example" {
							name = {{ get . "name" }}
							description = {{ get . "description" }}
							edge_gateway_id = cloudavenue_edgegateway.example.id
							app_ports = [
							  {
							  	protocol = "ICMPv4"
							  },
							  {
								protocol = "TCP"
								ports = [
									"80",
									"443",
									"8080",
								]
							  },
							  {
								protocol = "UDP"
								ports = [
									"53",
								]
							  }
							]
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.AppPortProfile)),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "app_ports.*", map[string]string{
								"protocol": "ICMPv4",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "app_ports.*", map[string]string{
								"protocol": "TCP",
								"ports.#":  "3",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "app_ports.*", map[string]string{
								"protocol": "UDP",
								"ports.#":  "1",
							}),
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
						ImportStateIDBuilder: []string{"edge_gateway_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
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
	}
}

func TestAccEdgeGatewayAppPortProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayAppPortProfileResource{}),
	})
}
