package testsacc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// localCache is used to store Test rendered.
var localCache = make(map[string]Test)

type ResourceConfig struct {
	TestACC
}

// GetDefaultConfig returns the create configuration for the test named "example".
func (r ResourceConfig) GetDefaultConfig() map[string]TFData {
	return r.GetSpecificConfig("example")()
}

// GetSpecificConfig returns the create configuration for the test named "example".
func (r *ResourceConfig) GetSpecificConfig(testName string) func() map[string]TFData {
	return func() map[string]TFData {
		x := make(map[string]TFData)

		// Add the dependencies config (scope : Tests)
		for _, v := range r.DependenciesConfig() {
			for k, v := range v() {
				x[k] = v
			}
		}

		if _, ok := localCache[r.GetResourceName()+"."+testName]; !ok {
			localCache[r.GetResourceName()+"."+testName] = r.Tests(context.Background())[TestName(testName)](
				context.Background(),
				r.GetResourceName()+"."+testName,
			)
		}

		// Add the config
		x[r.GetResourceName()+"."+testName] = localCache[r.GetResourceName()+"."+testName].Create.TFConfig

		// Add the Test dependencies config (scope : Test)
		if localCache[r.GetResourceName()+"."+testName].CommonDependencies != nil {
			for _, v := range localCache[r.GetResourceName()+"."+testName].CommonDependencies() {
				if v != nil {
					for k, v := range v() {
						x[k] = v
					}
				}
			}
		}

		return x
	}
}

// GetDefaultChecks returns the checks for the test named "example".
func (r ResourceConfig) GetDefaultChecks() []resource.TestCheckFunc {
	return r.GetSpecificChecks("example")
}

// GetSpecificChecks returns the checks for the test named.
func (r *ResourceConfig) GetSpecificChecks(testName string) []resource.TestCheckFunc {
	if _, ok := localCache[r.GetResourceName()+"."+testName]; !ok {
		localCache[r.GetResourceName()+"."+testName] = r.Tests(context.Background())[TestName(testName)](
			context.Background(),
			r.GetResourceName()+"."+testName,
		)
	}
	return localCache[r.GetResourceName()+"."+testName].GetChecks()
}

func NewResourceConfig(data TestACC) func() *ResourceConfig {
	return func() *ResourceConfig {
		return &ResourceConfig{
			TestACC: data,
		}
	}
}
