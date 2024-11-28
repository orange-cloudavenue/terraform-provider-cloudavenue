package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ALBPoolResource{}

const (
	ALBPoolResourceName = testsacc.ResourceName("cloudavenue_alb_pool")
)

type ALBPoolResource struct{}

func NewALBPoolResourceTest() testsacc.TestACC {
	return &ALBPoolResource{}
}

func (r *ALBPoolResource) GetResourceName() string {
	return ALBPoolResourceName.String()
}

func (r *ALBPoolResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetDataSourceConfig()[EdgeGatewayDataSourceName]().GetSpecificConfig("example_with_id"))
	return
}

func (r *ALBPoolResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.LoadBalancerPool)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_alb_pool" "example" {
						edge_gateway_id = data.cloudavenue_edgegateway.example_with_id.id
						name              = {{ generate . "name" }}
						persistence_profile = {
							type = "CLIENT_IP"
						}
						members = [
							{
								ip_address = "192.168.1.1"
								port       = "80"
							},
							{
								ip_address = "192.168.1.2"
								port       = "80"
							},
							{
								ip_address = "192.168.1.3"
								port       = "80"
							}
						]
						health_monitors = ["UDP", "TCP"]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_id"),
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "persistence_profile.type", "CLIENT_IP"),
						resource.TestCheckResourceAttr(resourceName, "members.0.ip_address", "192.168.1.1"),
						resource.TestCheckResourceAttr(resourceName, "members.0.port", "80"),
						resource.TestCheckResourceAttr(resourceName, "members.1.ip_address", "192.168.1.2"),
						resource.TestCheckResourceAttr(resourceName, "members.1.port", "80"),
						resource.TestCheckResourceAttr(resourceName, "members.2.ip_address", "192.168.1.3"),
						resource.TestCheckResourceAttr(resourceName, "members.2.port", "80"),
						resource.TestCheckTypeSetElemAttr(resourceName, "health_monitors.*", "UDP"),
						resource.TestCheckTypeSetElemAttr(resourceName, "health_monitors.*", "TCP"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
							resource "cloudavenue_alb_pool" "example" {
								edge_gateway_name = data.cloudavenue_edgegateway.example_with_id.name
								name              = {{ get . "name" }}
								persistence_profile = {
									type = "HTTP_COOKIE"
									value = {{ generate . "httpcookie" }}
								}
								members = [
									{
										ip_address = "192.168.1.1"
										port       = "80"
									},
									{
										ip_address = "192.168.1.2"
										port       = "80"
									},
									{
										ip_address = "192.168.1.4"
										port       = "8080"
									}
								]
								health_monitors = ["TCP"]
								algorithm = "ROUND_ROBIN"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "persistence_profile.type", "HTTP_COOKIE"),
							resource.TestCheckResourceAttr(resourceName, "persistence_profile.value", testsacc.GetValueFromTemplate(resourceName, "httpcookie")),
							resource.TestCheckResourceAttr(resourceName, "members.0.ip_address", "192.168.1.1"),
							resource.TestCheckResourceAttr(resourceName, "members.0.port", "80"),
							resource.TestCheckResourceAttr(resourceName, "members.1.ip_address", "192.168.1.2"),
							resource.TestCheckResourceAttr(resourceName, "members.1.port", "80"),
							resource.TestCheckResourceAttr(resourceName, "members.2.ip_address", "192.168.1.4"),
							resource.TestCheckResourceAttr(resourceName, "members.2.port", "8080"),
							resource.TestCheckTypeSetElemAttr(resourceName, "health_monitors.*", "TCP"),
							resource.TestCheckResourceAttr(resourceName, "algorithm", "ROUND_ROBIN"),
						},
					},
				},
				// ! Import testing
				Imports: []testsacc.TFImport{
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
	}
}

func TestAccALBPoolResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ALBPoolResource{}),
	})
}
