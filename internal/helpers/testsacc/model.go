/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type (
	ResourceName string // Resource name is cloudavenue_<resource_name>
	ConfigName   string // Config name is cloudavenue_<resource_name>.<config_name>

	TFData                     string
	TestName                   string
	ListOfDependencies         map[ConfigName]TFData
	DependenciesConfigResponse []func() map[string]TFData

	TestACC interface {
		// GetResourceName returns the name of the resource under test.
		GetResourceName() string

		// DependenciesConfig returns the Terraform configuration used to create any dependencies of the resource under test.
		DependenciesConfig() DependenciesConfigResponse

		// Tests returns the acceptance tests to run for the resource under test.
		// resourceName is a concatenation of the resource name and the example name. For example, "cloudavenue_catalog.example".
		Tests(context.Context) map[TestName]func(ctx context.Context, resourceName string) Test
	}

	Test struct {
		// CommonChecks is a list of common Terraform checks applied to all tests.
		CommonChecks []resource.TestCheckFunc

		// CommonDependencies is a list of common Terraform configurations used to create any dependencies of the resource under test.
		CommonDependencies func() DependenciesConfigResponse

		// Create returns the Terraform configuration to use for the test.
		// This should be a valid Terraform configuration that can be used to create, update, and destroy resources.
		Create TFConfig

		// Update returns the Terraform configurations to use for the update test.
		// This should be a valid Terraform configuration that can be used to update the resource under test.
		Updates []TFConfig

		// Import returns the Terraform configurations to use for the import test.
		Imports []TFImport

		// Destroy will create a destroy plan if set to true.
		Destroy bool

		// CacheDependenciesConfig is used to cache the dependencies config.
		CacheDependenciesConfig TFData

		// listOfDeps is a list of dependencies.
		listOfDeps ListOfDependencies
	}

	TFConfig struct {
		// Checks is a Terraform checks to run for checking the resource under test.
		Checks []resource.TestCheckFunc

		// TFCongig is the Terraform configuration to use for the test.
		TFConfig TFData

		// TFAdvanced is the Terraform advanced configuration to use for the test.
		TFAdvanced TFAdvanced
	}

	TFAdvanced struct {
		// PreConfig is called before the Config is applied to perform any per-step
		// setup that needs to happen. This is called regardless of "test mode"
		// below.
		PreConfig func()

		// Taint is a list of resource addresses to taint prior to the execution of
		// the step. Be sure to only include this at a step where the referenced
		// address will be present in state, as it will fail the test if the resource
		// is missing.
		//
		// This option is ignored on ImportState tests, and currently only works for
		// resources in the root module path.
		Taint []string

		// Destroy will create a destroy plan if set to true.
		Destroy bool

		// ExpectNonEmptyPlan can be set to true for specific types of tests that are
		// looking to verify that a diff occurs
		ExpectNonEmptyPlan bool

		// ExpectError allows the construction of test cases that we expect to fail
		// with an error. The specified regexp must match against the error for the
		// test to pass.
		ExpectError *regexp.Regexp

		// PlanOnly can be set to only run `plan` with this configuration, and not
		// actually apply it. This is useful for ensuring config changes result in
		// no-op plans
		PlanOnly bool

		// PreventDiskCleanup can be set to true for testing terraform modules which
		// require access to disk at runtime. Note that this will leave files in the
		// temp folder
		PreventDiskCleanup bool

		// PreventPostDestroyRefresh can be set to true for cases where data sources
		// are tested alongside real resources
		PreventPostDestroyRefresh bool
	}

	TFImport struct {
		// ImportStateId is the ID to perform an ImportState operation with.
		// This is optional. If it isn't set, then the resource ID is automatically
		// determined by inspecting the state for ResourceName's ID.
		ImportStateID string

		// ImportStateIdFunc is a function that can be used to dynamically generate
		// the ID for the ImportState tests. It is sent the state, which can be
		// checked to derive the attributes necessary and generate the string in the
		// desired format.
		ImportStateIDFunc resource.ImportStateIdFunc

		// ImportStateIDBuilder is a function that can be used to dynamically generate
		// the ID for the ImportState tests. It is sent the state, which can be
		// checked to derive the attributes necessary and generate the string in the
		// desired format.
		// Specifie the list of attribute names to use to build the ID.
		// Example of use: []string{"vdc_name", "edgegateway_name"} => "vdcExample.edgegatewayExample"
		ImportStateIDBuilder []string

		// ImportStateVerifyIgnore is a list of prefixes of fields that should
		// not be verified to be equal. These can be set to ephemeral fields or
		// fields that can't be refreshed and don't matter.
		ImportStateVerifyIgnore []string

		// ImportStatePersist, if true, will update the persisted state with the
		// state generated by the import operation (i.e., terraform import). When
		// false (default) the state generated by the import operation is discarded
		// at the end of the test step that is verifying import behavior.
		ImportStatePersist bool

		// ImportStateVerify, if true, will also check that the state values
		// that are finally put into the state after import match for all the
		// IDs returned by the Import.
		ImportStateVerify bool

		// ImportState, if true, will test the functionality of ImportState
		// by importing the resource with ID of that resource.
		ImportState bool
	}
)
