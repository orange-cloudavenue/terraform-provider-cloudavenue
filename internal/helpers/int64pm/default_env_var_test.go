// Package int64pm provides a plan modifier for int64 values.
package int64pm_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/int64pm"
)

func TestDefaultEnvVarModifierPlanModifyInt64(t *testing.T) {
	const (
		envVarName = "TEST_VAR"
		envValue   = 123
	)

	for name, testCase := range int64pmTestCases(envValue) {
		name, testCase := name, testCase
		// set environnement variable
		t.Setenv(envVarName, "123")

		t.Run(name, func(t *testing.T) {
			resp := &planmodifier.Int64Response{
				PlanValue: testCase.request.PlanValue,
			}

			int64pm.SetDefaultEnvVar(envVarName).PlanModifyInt64(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
