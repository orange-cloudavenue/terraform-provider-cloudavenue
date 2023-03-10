// Package boolpm provides a plan modifier for boolean values.
package boolpm_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/boolpm"
)

func TestDefaultModifierPlanModifyBool(t *testing.T) {
	expectedValue := true

	for name, testCase := range boolpmTestCases(expectedValue) {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			// t.Parallel()

			resp := &planmodifier.BoolResponse{
				PlanValue: testCase.request.PlanValue,
			}

			boolpm.SetDefault(expectedValue).PlanModifyBool(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
