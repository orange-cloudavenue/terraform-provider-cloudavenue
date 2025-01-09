package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &OrgCertificateLibraryDatasource{}

const (
	OrgCertificateLibraryDatasourceName = testsacc.ResourceName("data.cloudavenue_org_certificate_library")
)

type OrgCertificateLibraryDatasource struct{}

func NewOrgCertificateLibraryDatasourceTest() testsacc.TestACC {
	return &OrgCertificateLibraryDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *OrgCertificateLibraryDatasource) GetResourceName() string {
	return OrgCertificateLibraryDatasourceName.String()
}

func (r *OrgCertificateLibraryDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	// resp.Append(GetResourceConfig()[OrgCertificateLibraryDatasourcesGoResourceName]().GetDefaultConfig),
	return
}

func (r *OrgCertificateLibraryDatasource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_org_certificate_library" "example" {
						name = "cert-auto-self-sign"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.CertificateLibraryItem)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrSet(resourceName, "certificate"),
					},
				},
			}
		},
	}
}

func TestAccOrgCertificateLibraryDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&OrgCertificateLibraryDatasource{}),
	})
}
