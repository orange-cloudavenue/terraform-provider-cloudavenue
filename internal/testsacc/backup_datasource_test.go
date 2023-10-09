package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &BackupDataSource{}

const (
	BackupDataSourceName = testsacc.ResourceName("data.cloudavenue_backup")
)

type BackupDataSource struct{}

func NewBackupDataSourceTest() testsacc.TestACC {
	return &BackupDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *BackupDataSource) GetResourceName() string {
	return BackupDataSourceName.String()
}

func (r *BackupDataSource) DependenciesConfig() (configs testsacc.TFData) {
	// Add dependencies config to the resource
	configs.Append(GetResourceConfig()[BackupResourceName]().GetDefaultConfig())
	configs.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig())
	return
}

func (r *BackupDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (backup vdc example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_backup" "example" {
						type = "vdc"
						target_name = cloudavenue_backup.example.target_name
					}`,
					Checks: NewBackupResourceTest().Tests(ctx)["example"](ctx, resourceName).GenerateCheckWithCommonChecks(),
				},
			}
		},
	}
}

func TestAccBackupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&BackupDataSource{}),
	})
}
