package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ORGCertificateLibraryResource{}

const (
	ORGCertificateLibraryResourceName = testsacc.ResourceName("cloudavenue_org_certificate_library")
)

type ORGCertificateLibraryResource struct{}

func NewORGCertificateLibraryResourceTest() testsacc.TestACC {
	return &ORGCertificateLibraryResource{}
}

// GetResourceName returns the name of the resource.
func (r *ORGCertificateLibraryResource) GetResourceName() string {
	return ORGCertificateLibraryResourceName.String()
}

func (r *ORGCertificateLibraryResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *ORGCertificateLibraryResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.CertificateLibraryItem)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					resource "cloudavenue_org_certificate_library" "example" {
						name = "example"
						description = "This is a certificate"
						certificate = file("/Users/micheneaudavid/cav-cert.pem")
						private_key = file("/Users/micheneaudavid/cav-key.pem")
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", "example"),
						resource.TestCheckResourceAttr(resourceName, "description", "This is a certificate"),
						resource.TestCheckResourceAttrSet(resourceName, "certificate"),
						resource.TestCheckResourceAttrSet(resourceName, "private_key"),
						resource.TestCheckNoResourceAttr(resourceName, "passphrase"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: `
						resource "cloudavenue_org_certificate_library" "example" {
							name = "example updated"
							description = "This is a certificate updated"
							certificate = file("/Users/micheneaudavid/cav-cert.pem")
							private_key = file("/Users/micheneaudavid/cav-key.pem")
						}`,
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", "example updated"),
							resource.TestCheckResourceAttr(resourceName, "description", "This is a certificate updated"),
							resource.TestCheckResourceAttrSet(resourceName, "certificate"),
							resource.TestCheckResourceAttrSet(resourceName, "private_key"),
							resource.TestCheckNoResourceAttr(resourceName, "passphrase"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_org_certificate_library" "example" {
							name = "example updated"
							description = "This is a certificate updated"
							certificate = file("/Users/micheneaudavid/cav-cert.pem")
							private_key = file("/Users/micheneaudavid/cav-key.pem")
							passphrase = "password"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", "example updated"),
							resource.TestCheckResourceAttr(resourceName, "description", "This is a certificate updated"),
							resource.TestCheckResourceAttrSet(resourceName, "certificate"),
							resource.TestCheckResourceAttrSet(resourceName, "private_key"),
							resource.TestCheckResourceAttrSet(resourceName, "passphrase"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder:    []string{"id"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"passphrase", "private_key"},
					},
					{
						ImportStateIDBuilder:    []string{"name"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"passphrase", "private_key"},
					},
				},
			}
		},
	}
}

func TestAccORGCertificateLibraryResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ORGCertificateLibraryResource{}),
	})
}
