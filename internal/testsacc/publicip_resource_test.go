package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &PublicIPResource{}

const (
	PublicIPResourceName = testsacc.ResourceName("cloudavenue_publicip")
)

type PublicIPResource struct{}

func NewPublicIPResourceTest() testsacc.TestACC {
	return &PublicIPResource{}
}

// GetResourceName returns the name of the resource.
func (r *PublicIPResource) GetResourceName() string {
	return PublicIPResourceName.String()
}

func (r *PublicIPResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[EdgeGatewayResourceName]().GetDefaultConfig)
	return
}

func (r *PublicIPResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_publicip" "example" {
					  edge_gateway_id = cloudavenue_edgegateway.example.id
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
						resource.TestCheckResourceAttrWith(resourceName, "edge_gateway_id", uuid.TestIsType(uuid.Gateway)),
						resource.TestCheckNoResourceAttr(resourceName, "edge_gateway_name"),
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_with_edge_name": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
						resource "cloudavenue_publicip" "example_with_edge_name" {
							edge_gateway_name = cloudavenue_edgegateway.example.name
						}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
						resource.TestCheckNoResourceAttr(resourceName, "edge_gateway_id"),
						resource.TestCheckResourceAttrSet(resourceName, "edge_gateway_name"),
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"edge_gateway_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccPublicIPResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&PublicIPResource{}),
	})
}
