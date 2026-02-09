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
	"fmt"
	"log"
	"os"

	"github.com/iancoleman/strcase"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// * DependenciesConfigResponse
// Append appends the given dependencies config to the current one.
func (d *DependenciesConfigResponse) Append(tf func() map[string]TFData) {
	*d = append(*d, tf)
}

// * TFAdvanced
// IsEmpty returns true if the TFAdvanced is empty.
func (t *TFAdvanced) IsEmpty() bool {
	return t == nil
}

// * ResourceName
// String returns the string representation of the ResourceName.
func (r ResourceName) String() string {
	return string(r)
}

// * TestName
// Get return the name of the example formatted as a lower camel case string.
func (e TestName) Get() string {
	return strcase.ToSnake(e.String())
}

// String string representation of the example name.
func (e TestName) String() string {
	return string(e)
}

// ComputeResourceName returns the name of the resource under test.
// Return cloudavenue_<resource_name>.<example_name>.
func (e TestName) ComputeResourceName(resourceName string) string {
	return fmt.Sprintf("%s.%s", resourceName, e.Get())
}

// *TFData

// Get returns the Terraform configuration as a string.
func (t TFData) Get() string {
	return t.String()
}

// Set sets the Terraform configuration to the given string.
func (t *TFData) Set(s string) {
	*t = TFData(s)
}

// String returns the Terraform configuration as a string.
func (t TFData) String() string {
	return string(t)
}

// Append appends the given Terraform configuration to the current one.
func (t *TFData) Append(tf TFData) {
	t.append(tf)
}

// AppendWithoutResourceName appends the given Terraform configuration to the current one.
func (t *TFData) AppendWithoutResourceName(tf TFData) {
	t.appendWithoutResourceName(tf)
}

// is empty.
func (t *TFData) IsEmpty() bool {
	return t.String() == ""
}

// appendWithoutResourceName appends the given Terraform configuration to the current one.
func (t *TFData) appendWithoutResourceName(tf TFData) {
	t.append(tf)
}

// append appends the given Terraform configuration to the current one.
func (t *TFData) append(tf TFData) {
	*t = TFData(fmt.Sprintf("%s\n%s", t, tf.String()))
}

// The code below is commented out because it is not used in the current implementation.
// But it can be useful in the future.
// extractResourceName extracts the resource name and config name from the Terraform configuration.
// example: "resource "cloudavenue_catalog" "example" {}" => "cloudavenue_catalog.example"
// example: "data "cloudavenue_catalog" "example" {}" => "data.cloudavenue_catalog.example"
// func (t *TFData) extractResourceName() string {
// 	// find the first occurrence of "resource" or "data"
// 	re := regexp.MustCompile(`(resource|data) \"(.*)\" \"(.*)\"`)
// 	line := re.FindString(t.String())

// 	x := strings.Split(line, " ")
// 	// for each word remove the double quotes
// 	for i, v := range x {
// 		x[i] = strings.ReplaceAll(v, "\"", "")
// 	}

// 	// Result is <resource_name>.<config_name>
// 	if x[0] == "resource" {
// 		return fmt.Sprintf("%s.%s", x[1], x[2])
// 	}

// 	// Result is data.<resource_name>.<config_name>
// 	return fmt.Sprintf("%s.%s.%s", x[0], x[1], x[2])
// }

// *TFConfig

// Generate creates the Terraform configuration for the resource under test.
// It returns the Terraform configuration as a string.
// Concatenate the dependencies config and the resource config.
func (t TFConfig) Generate(_ context.Context, dependencies TFData) string {
	t.TFConfig.appendWithoutResourceName(dependencies)
	return t.TFConfig.Get()
}

// *Test

// initListOfDeps initializes the list of dependencies.
func (t *Test) initListOfDeps() {
	if t.listOfDeps == nil {
		t.listOfDeps = make(ListOfDependencies, 0)
	}
}

// ExistInListOfDeps.
func (t *Test) ExistInListOfDeps(configName string) bool {
	return t.listOfDeps.Exists(ConfigName(configName))
}

// Compute Dependencies config.
func (t *Test) ComputeDependenciesConfig(testACC TestACC) {
	t.initListOfDeps()
	log.Default().Print("Starting building dependencies config")
	for _, vs := range testACC.DependenciesConfig() {
		for configName, tfData := range vs() {
			t.listOfDeps.Append(ConfigName(configName), tfData)
		}
	}

	if t.CommonDependencies != nil {
		for _, vs := range t.CommonDependencies() {
			for configName, tfData := range vs() {
				t.listOfDeps.Append(ConfigName(configName), tfData)
			}
		}
	}
	log.Default().Printf("Finished building dependencies config. Found %d dependencies", len(t.listOfDeps))
}

// GetChecks returns the checks for the test.
func (t Test) GetChecks() []resource.TestCheckFunc {
	checks := make([]resource.TestCheckFunc, 0)
	checks = append(checks, t.CommonChecks...)
	checks = append(checks, t.Create.Checks...)
	return checks
}

// GenerateSteps generates the structure of the acceptance tests.
func (t Test) GenerateSteps(ctx context.Context, testName TestName, testACC TestACC) (steps []resource.TestStep) {
	// Init Slice
	steps = make([]resource.TestStep, 0)

	// * Compute dependencies config
	t.ComputeDependenciesConfig(testACC)

	// * Create step
	lastConfigGenerated := t.Create.Generate(ctx, t.listOfDeps.ToTFData())

	log.Default().Print(lastConfigGenerated)
	// if env TF_ACC_ONLY_PRINT is set to true, print the config and exit.
	if os.Getenv("TF_ACC_ONLY_PRINT") == "true" {
		os.Exit(0)
	}

	// listOfChecks is a concatenation of the common checks and the specific checks.
	listOfChecks := t.CommonChecks
	listOfChecks = append(listOfChecks, t.Create.Checks...)
	createTestStep := resource.TestStep{
		Config: lastConfigGenerated,
		Check: resource.ComposeAggregateTestCheckFunc(
			listOfChecks...,
		),
	}

	if !t.Create.TFAdvanced.IsEmpty() {
		createTestStep.PreConfig = t.Create.TFAdvanced.PreConfig
		createTestStep.Taint = t.Create.TFAdvanced.Taint
		createTestStep.Destroy = t.Create.TFAdvanced.Destroy
		createTestStep.ExpectNonEmptyPlan = t.Create.TFAdvanced.ExpectNonEmptyPlan
		createTestStep.ExpectError = t.Create.TFAdvanced.ExpectError
		createTestStep.PlanOnly = t.Create.TFAdvanced.PlanOnly
		createTestStep.PreventDiskCleanup = t.Create.TFAdvanced.PreventDiskCleanup
		createTestStep.PreventPostDestroyRefresh = t.Create.TFAdvanced.PreventPostDestroyRefresh
	}

	steps = append(steps, createTestStep)

	// * Update steps
	if len(t.Updates) > 0 {
		for _, update := range t.Updates {
			listOfChecks := t.CommonChecks
			listOfChecks = append(listOfChecks, update.Checks...)

			lastConfigGenerated = update.Generate(ctx, t.listOfDeps.ToTFData())
			updateTestStep := resource.TestStep{
				Config: lastConfigGenerated,
				Check: resource.ComposeAggregateTestCheckFunc(
					listOfChecks...,
				),
			}

			if !update.TFAdvanced.IsEmpty() {
				updateTestStep.PreConfig = update.TFAdvanced.PreConfig
				updateTestStep.Taint = update.TFAdvanced.Taint
				updateTestStep.Destroy = update.TFAdvanced.Destroy
				updateTestStep.ExpectNonEmptyPlan = update.TFAdvanced.ExpectNonEmptyPlan
				updateTestStep.ExpectError = update.TFAdvanced.ExpectError
				updateTestStep.PlanOnly = update.TFAdvanced.PlanOnly
				updateTestStep.PreventDiskCleanup = update.TFAdvanced.PreventDiskCleanup
				updateTestStep.PreventPostDestroyRefresh = update.TFAdvanced.PreventPostDestroyRefresh
			}

			steps = append(steps, updateTestStep)
		}
	}

	// * Import steps
	if len(t.Imports) > 0 {
		for _, importStep := range t.Imports {
			importTest := resource.TestStep{
				ResourceName:            testName.ComputeResourceName(testACC.GetResourceName()),
				ImportState:             importStep.ImportState,
				ImportStateId:           importStep.ImportStateID,
				ImportStateVerify:       importStep.ImportStateVerify,
				ImportStateVerifyIgnore: importStep.ImportStateVerifyIgnore,
				ImportStatePersist:      importStep.ImportStatePersist,
				ImportStateIdFunc:       importStep.ImportStateIDFunc,
			}

			if len(importStep.ImportStateIDBuilder) > 0 {
				importTest.ImportStateIdFunc = ImportStateIDBuilder(testName.ComputeResourceName(testACC.GetResourceName()), importStep.ImportStateIDBuilder)
			}

			steps = append(steps, importTest)
		}
	}

	// * Destroy step
	if t.Destroy {
		destroyTestStep := resource.TestStep{
			Config:  lastConfigGenerated,
			Destroy: true,
		}

		steps = append(steps, destroyTestStep)
	}

	return steps
}

// *TestACC
// GenerateTests generates the acceptance tests for the resource under test.
func GenerateTests(tacc TestACC) []resource.TestStep {
	var (
		ctx   = context.Background()
		steps = make([]resource.TestStep, 0)
	)

	// For each test
	for testName, step := range tacc.Tests(ctx) {
		// resourceName is a concatenation of the resource name and the example name. For example, "cloudavenue_catalog.example".
		resourceName := testName.ComputeResourceName(tacc.GetResourceName())
		if envvar, ok := os.LookupEnv("TF_ACC_RUN_TEST"); ok {
			if envvar != testName.String() {
				continue
			}
		}
		steps = append(steps, step(ctx, resourceName).GenerateSteps(ctx, testName, tacc)...)
	}

	return steps
}

// * Other
// ImportStateIDBuilder is a function that can be used to dynamically generate
// the ID for the ImportState tests. It is sent the state, which can be
// checked to derive the attributes necessary and generate the string in the
// desired format.
// Specifie the list of attribute names to use to build the ID.
// Example of use: []string{"vdc_name", "edgegateway_name"} => "vdcExample.edgegatewayExample".
func ImportStateIDBuilder(resourceName string, attributeNames []string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ImportStateIDBuilder : Resource %s not found", resourceName)
		}

		// Build the ID
		id := ""
		for _, attributeName := range attributeNames {
			// Catch attribute not found
			i, ok := rs.Primary.Attributes[attributeName]
			if !ok {
				return "", fmt.Errorf("ImportStateIDBuilder : Attribute %s not found", attributeName)
			}

			id += i
			if attributeName != attributeNames[len(attributeNames)-1] {
				id += "."
			}
		}

		if id == "" {
			return "", fmt.Errorf("ImportStateIDBuilder : ID is empty after building")
		}

		return id, nil
	}
}
