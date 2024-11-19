package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var _ testsacc.TestACC = &iamUserSAMLResource{}

const (
	iamUserSAMLResourceName = testsacc.ResourceName("cloudavenue_iam_user_saml")
)

type iamUserSAMLResource struct{}

func NewiamUserSAMLResourceTest() testsacc.TestACC {
	return &iamUserSAMLResource{}
}

// GetResourceName returns the name of the resource.
func (r *iamUserSAMLResource) GetResourceName() string {
	return iamUserSAMLResourceName.String()
}

func (r *iamUserSAMLResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// TODO : Add dependencies config
	// resp.Append(GetResourceConfig()[CatalogResourceName]().GetDefaultConfig)

	// This is method for add dependencies legacy config
	// resp.Append(AddConstantConfig(constantName))
	return
}

/*
	Unit tests not working for now. Because the phone number imported from the SAML provider is valid and the update is not working
*/

func (r *iamUserSAMLResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// TODO : Complete tests
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", uuid.TestIsType(uuid.User)), // TODO : Change type
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_iam_user_saml" "example" {
						user_name = "mickael.stanislas.ext"
						role_name = "Organization Administrator"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "foo", "bar"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_iam_user_saml" "example" {
							foo = "barUpdated"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "foo", "barUpdated"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_iam_user_saml" "example" {
							foo = "barUpdated"
							bar = "foo"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "foo", "barUpdated"),
							resource.TestCheckResourceAttr(resourceName, "bar", "foo"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAcciamUserSAMLResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&iamUserSAMLResource{}),
	})
}
