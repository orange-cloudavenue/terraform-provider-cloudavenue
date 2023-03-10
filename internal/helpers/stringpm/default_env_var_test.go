// Package stringpm provides a plan modifier for string values.
package stringpm_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/stringpm"
)

func TestDefaultEnvVarModifierPlanModifyString(t *testing.T) {
	const (
		envVarName = "TEST_VAR"
		envValue   = "testFromEnvVar"
	)

	for name, testCase := range stringpmTestCases(envValue) {
		name, testCase := name, testCase
		// set environnement variable
		t.Setenv(envVarName, envValue)

		t.Run(name, func(t *testing.T) {
			// t.Parallel()

			resp := &planmodifier.StringResponse{
				PlanValue: testCase.request.PlanValue,
			}

			stringpm.SetDefaultEnvVar(envVarName).PlanModifyString(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
