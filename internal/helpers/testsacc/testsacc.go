package testsacc

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type (
	ResourceName               string
	TFData                     string
	TestName                   string
	ListOfDependencies         []ResourceName
	DependenciesConfigResponse []func() TFData

	// TODO : Add TFACCLog management.
	TFACCLog struct {
		Level string `env:"LEVEL,default=info"`
	}

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
		// If EnableAutoCheck is true, these checks will be automatically added to the test.
		// If EnableAutoCheck is false, these checks will be the only checks run for the test.
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

// * DependenciesConfigResponse
// Append appends the given dependencies config to the current one.
func (d *DependenciesConfigResponse) Append(tf func() TFData) {
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

// * ListOfDependencies
// Append appends the given dependency to the list of dependencies.
// resourceName is a concatenation of the resource name and the name. For example, "cloudavenue_catalog.example".
func (l *ListOfDependencies) Append(resourceName ResourceName) {
	x := strings.Split(resourceName.String(), ".")

	if len(x) == 2 && !l.Exists(resourceName) {
		*l = append(*l, resourceName)
	}
}

// Exists checks if the given dependency exists in the list of dependencies.
// resourceName is a concatenation of the resource name and the name. For example, "cloudavenue_catalog.example".
func (l *ListOfDependencies) Exists(resourceName ResourceName) bool {
	for _, v := range *l {
		if v == resourceName {
			return true
		}
	}
	return false
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

// extractResourceName extracts the resource name and config name from the Terraform configuration.
// example: "resource "cloudavenue_catalog" "example" {}" => "cloudavenue_catalog.example"
func (t *TFData) extractResourceName() string {
	// find the first occurrence of "resource" or "data"
	re := regexp.MustCompile(`(resource|data) \"(.*)\" \"(.*)\"`)
	line := re.FindString(t.String())

	x := strings.Split(line, " ")
	// for each word remove the double quotes
	for i, v := range x {
		x[i] = strings.ReplaceAll(v, "\"", "")
	}
	return fmt.Sprintf("%s.%s", x[1], x[2])
}

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

// Compute Dependencies config.
func (t *Test) ComputeDependenciesConfig(testACC TestACC) {
	t.initListOfDeps()
	for _, v := range testACC.DependenciesConfig() {
		// tf contains terraform configuration for a dependency
		tf := v()
		if !tf.IsEmpty() && !t.listOfDeps.Exists(ResourceName(tf.extractResourceName())) {
			t.CacheDependenciesConfig.append(tf)
			t.listOfDeps.Append(ResourceName(tf.extractResourceName()))
		}
	}
}

// GenerateSteps generates the structure of the acceptance tests.
func (t Test) GenerateSteps(ctx context.Context, testName TestName, testACC TestACC) (steps []resource.TestStep) {
	// Init Slice
	steps = make([]resource.TestStep, 0)

	// listOfChecks is a concatenation of the common checks and the specific checks.
	listOfChecks := t.CommonChecks
	listOfChecks = append(listOfChecks, t.Create.Checks...)

	// lastConfigGenerated is the last Terraform configuration generated. (Used for destroy step)
	var lastConfigGenerated string

	// * Compute dependencies config
	t.ComputeDependenciesConfig(testACC)

	// * Create step
	lastConfigGenerated = t.Create.Generate(ctx, t.CacheDependenciesConfig)
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

			lastConfigGenerated = update.Generate(ctx, t.CacheDependenciesConfig)
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

		steps = append(steps, step(ctx, resourceName).GenerateSteps(ctx, testName, tacc)...)
	}

	return steps
}

// GenerateTestChecks Generate the checks for a specific test.
func GenerateTestChecks(ctx context.Context, tacc TestACC, resourceName string, testName TestName) []resource.TestCheckFunc {
	return tacc.Tests(ctx)[testName](ctx, resourceName).GenerateCheckWithCommonChecks()
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
