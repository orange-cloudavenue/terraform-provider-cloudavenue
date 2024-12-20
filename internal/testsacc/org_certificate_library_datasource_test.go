package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &OrgCertificateLibraryDatasource{}

const (
	OrgCertificateLibraryDatasourceName = testsacc.ResourceName("data.cloudavenue_org_certificate_library_datasources_go")
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
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_org_certificate_library_datasources_go" "example" {
						foo_id = cloudavenue_foo_bar.example.id
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					// Checks: GetResourceConfig()[org_CertificateLibraryDatasourcesGoResourceName]().GetDefaultChecks()
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
