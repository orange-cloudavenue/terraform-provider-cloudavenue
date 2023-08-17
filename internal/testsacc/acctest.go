// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider"
)

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
	if v := os.Getenv("CLOUDAVENUE_URL"); v == "" {
		t.Fatal("CLOUDAVENUE_URL must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_USER"); v == "" {
		t.Fatal("CLOUDAVENUE_USER must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_PASSWORD"); v == "" {
		t.Fatal("CLOUDAVENUE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_ORG"); v == "" {
		t.Fatal("CLOUDAVENUE_ORG must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_VDC"); v == "" {
		t.Fatal("CLOUDAVENUE_VDC must be set for acceptance tests")
	}
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

type ResourceName string

// String returns the string representation of the ResourceName.
func (r ResourceName) String() string {
	return string(r)
}

type resourceConfig struct {
	testsacc.TestACC
}

// GetDefaultConfig returns the create configuration for the test named "example".
func (r resourceConfig) GetDefaultConfig() testsacc.TFData {
	return r.Tests(context.Background())["example"](
		context.Background(),
		r.GetResourceName(),
	).Create.TFConfig
}

// GetSpecificConfig returns the create configuration for the test named "example".
func (r resourceConfig) GetSpecificConfig(testName string) testsacc.TFData {
	return r.Tests(context.Background())[testsacc.TestName(testName)](
		context.Background(),
		r.GetResourceName(),
	).Create.TFConfig
}

// AddConstantConfig returns the create configuration from constant.
func AddConstantConfig(config string) testsacc.TFData {
	return testsacc.TFData(config)
}

func NewResourceConfig(data testsacc.TestACC) func() resourceConfig {
	return func() resourceConfig {
		return resourceConfig{
			TestACC: data,
		}
	}
}
