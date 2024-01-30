// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider"
)

var localCacheResource = map[string]testsacc.Test{}

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cloudavenue": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// You can add code here to run prior to any test case execution, for example assertions
// about the appropriate environment variables being set are common to see in a pre-check
// function.
func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("CLOUDAVENUE_USERNAME"); v == "" {
		t.Fatal("CLOUDAVENUE_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_PASSWORD"); v == "" {
		t.Fatal("CLOUDAVENUE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_ORG"); v == "" {
		t.Fatal("CLOUDAVENUE_ORG must be set for acceptance tests")
	}

	// Generate a new execution ID for this run.
	// Not error checking here because it's not critical.
	x, _ := uuid.NewUUID()
	metrics.GlobalExecutionID = "testacc_" + x.String()
	log.Default().Printf("TestACC: execution ID is %s", metrics.GlobalExecutionID)
}

// Deprecated: Use ContactConfigs instead.
func ConcatTests(tests ...string) string {
	return ContactConfigs(tests...)
}

// ContactConfigs concatenates the given configs into a single string.
func ContactConfigs(configs ...string) string {
	var result string
	for _, config := range configs {
		result += config + "\n"
	}
	return result
}

type resourceConfig struct {
	testsacc.TestACC
}

// GetDefaultConfig returns the create configuration for the test named "example".
func (r resourceConfig) GetDefaultConfig() testsacc.TFData {
	return r.GetSpecificConfig("example")()
}

// GetSpecificConfig returns the create configuration for the test named "example".
func (r resourceConfig) GetSpecificConfig(testName string) func() testsacc.TFData {
	return func() testsacc.TFData {
		return r.fromCache(testName).CacheDependenciesConfig
	}
}

func (r resourceConfig) fromCache(testName string) testsacc.Test {
	t, ok := localCacheResource[r.GetResourceName()+"."+testName]
	if !ok {
		t = r.Tests(context.Background())[testsacc.TestName(testName)](
			context.Background(),
			r.GetResourceName()+"."+testName,
		)
		t.ComputeDependenciesConfig(r.TestACC)
		t.CacheDependenciesConfig.Append(t.Create.TFConfig)
		localCacheResource[r.GetResourceName()+"."+testName] = t
	}

	return t
}

// GetDefaultChecks returns the checks for the test named "example".
func (r resourceConfig) GetDefaultChecks() []resource.TestCheckFunc {
	return r.GetSpecificChecks("example")
}

// GetSpecificChecks returns the checks for the test named.
func (r resourceConfig) GetSpecificChecks(testName string) []resource.TestCheckFunc {
	return r.fromCache(testName).Create.Checks
}

// AddConstantConfig returns the create configuration from constant.
func AddConstantConfig(config string) func() testsacc.TFData {
	return func() testsacc.TFData {
		return testsacc.TFData(config)
	}
}

func NewResourceConfig(data testsacc.TestACC) func() resourceConfig {
	return func() resourceConfig {
		return resourceConfig{
			TestACC: data,
		}
	}
}
