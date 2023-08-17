package testsacc

import (
	"context"
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TODO : ADD Generic func ImportStateIDFunc to generate the ID for the ImportState tests

type (
	TFData   string
	TestName string

	TestACC interface {
		// GetResourceName returns the name of the resource under test.
		GetResourceName() string

		// DependenciesConfig returns the Terraform configuration used to create any dependencies of the resource under test.
		DependenciesConfig() TFData

		// Tests returns the acceptance tests to run for the resource under test.
		// resourceName is a concatenation of the resource name and the example name. For example, "cloudavenue_catalog.example".
		Tests(context.Context) map[TestName]func(ctx context.Context, resourceName string) Test
	}

	Test struct {
		// CommonChecks is a list of common Terraform checks applied to all tests.
		CommonChecks []resource.TestCheckFunc

		// Create returns the Terraform configuration to use for the test.
		// This should be a valid Terraform configuration that can be used to create, update, and destroy resources.
		Create TFConfig

		// Update returns the Terraform configurations to use for the update test.
		// This should be a valid Terraform configuration that can be used to update the resource under test.
		Updates []TFConfig

		// Import returns the Terraform configurations to use for the import test.
		Imports []TFImport
	}

	TFConfig struct {
		// Checks is a Terraform checks to run for checking the resource under test.
		// If EnableAutoCheck is true, these checks will be automatically added to the test.
		// If EnableAutoCheck is false, these checks will be the only checks run for the test.
		Checks []resource.TestCheckFunc

		// TFCongig is the Terraform configuration to use for the test.
		TFConfig TFData
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
	*t = TFData(fmt.Sprintf("%s\n%s", t, tf))
}

// *TFConfig

// Generate creates the Terraform configuration for the resource under test.
// It returns the Terraform configuration as a string.
// Concatenate the dependencies config and the resource config.
func (t TFConfig) Generate(ctx context.Context, dependencies TFData) string {
	dependencies.Append(t.TFConfig)
	return dependencies.Get()
}

// *Test

// GenerateSteps generates the structure of the acceptance tests.
func (t Test) GenerateSteps(ctx context.Context, testName TestName, testACC TestACC) (steps []resource.TestStep) {
	// Init Slice
	steps = make([]resource.TestStep, 0)

	// listOfChecks is a concatenation of the common checks and the specific checks.
	listOfChecks := t.CommonChecks
	listOfChecks = append(listOfChecks, t.Create.Checks...)

	// * Create step
	steps = append(steps, resource.TestStep{
		Config: t.Create.Generate(ctx, testACC.DependenciesConfig()),
		Check: resource.ComposeAggregateTestCheckFunc(
			listOfChecks...,
		),
	})

	// * Update steps
	if len(t.Updates) > 0 {
		for _, update := range t.Updates {
			listOfChecks := t.CommonChecks
			listOfChecks = append(listOfChecks, update.Checks...)

			steps = append(steps, resource.TestStep{
				Config: update.Generate(ctx, testACC.DependenciesConfig()),
				Check: resource.ComposeAggregateTestCheckFunc(
					listOfChecks...,
				),
			})
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

	return
}

// GenerateCheckWithCommonChecks concatenates the common checks and the specific checks.
func (t Test) GenerateCheckWithCommonChecks() []resource.TestCheckFunc {
	listOfChecks := t.CommonChecks
	listOfChecks = append(listOfChecks, t.Create.Checks...)

	return listOfChecks
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

		steps = step(ctx, resourceName).GenerateSteps(ctx, testName, tacc)
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
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		// Build the ID
		id := ""
		for _, attributeName := range attributeNames {
			id += rs.Primary.Attributes[attributeName]
			if attributeName != attributeNames[len(attributeNames)-1] {
				id += "."
			}
		}

		return id, nil
	}
}
