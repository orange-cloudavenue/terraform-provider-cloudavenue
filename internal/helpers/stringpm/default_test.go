// Package stringpm provides a plan modifier for string values.
package stringpm_test

import (
	"context"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/stringpm"
)

func TestDefaultModifierPlanModifyString(t *testing.T) {
	expectedValue := uniuri.New()

	for name, testCase := range stringpmTestCases(expectedValue) {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			// t.Parallel()

			resp := &planmodifier.StringResponse{
				PlanValue: testCase.request.PlanValue,
			}

			stringpm.SetDefault(expectedValue).PlanModifyString(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
