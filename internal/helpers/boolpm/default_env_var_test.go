// Package boolpm provides a plan modifier for boolean values.
package boolpm_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/boolpm"
)

func TestDefaultEnvVarModifierPlanModifyBool(t *testing.T) {
	const (
		envVarName = "TEST_VAR"
		envValue   = true
	)

	for name, testCase := range boolpmTestCases(envValue) {
		name, testCase := name, testCase
		// set environnement variable
		t.Setenv(envVarName, "true")

		t.Run(name, func(t *testing.T) {
			resp := &planmodifier.BoolResponse{
				PlanValue: testCase.request.PlanValue,
			}

			boolpm.SetDefaultEnvVar(envVarName).PlanModifyBool(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
