package testsacc

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketResource{}

const (
	S3BucketResourceName = testsacc.ResourceName("cloudavenue_s3_bucket")
)

type S3BucketResource struct{}

func NewS3BucketResourceTest() testsacc.TestACC {
	return &S3BucketResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketResource) GetResourceName() string {
	return S3BucketResourceName.String()
}

func (r *S3BucketResource) DependenciesConfig() (configs testsacc.TFData) {
	return
}

func (r *S3BucketResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// TODO : Complete tests
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_bucket" "example" {
						name = {{ generate . "name" }}
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "object_lock", "false"),
						resource.TestCheckResourceAttr(resourceName, "endpoint", fmt.Sprintf("https://%s.s3-region01.cloudavenue.orange-business.com", testsacc.GetValueFromTemplate(resourceName, "name"))),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// No updates
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"examplewithobjectlock": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_bucket" "examplewithobjectlock" {
						name = {{ generate . "name" }}
						object_lock = true
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "object_lock", "true"),
						resource.TestCheckResourceAttr(resourceName, "endpoint", fmt.Sprintf("https://%s.s3-region01.cloudavenue.orange-business.com", testsacc.GetValueFromTemplate(resourceName, "name"))),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// No updates
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccS3BucketResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketResource{}),
	})
}
