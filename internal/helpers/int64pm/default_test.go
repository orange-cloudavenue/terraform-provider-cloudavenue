// Package int64pm provides a plan modifier for int64 values.
package int64pm_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/int64pm"
)

func TestDefaultModifierPlanModifyInt64(t *testing.T) {
	const expectedValue = 123

	for name, testCase := range int64pmTestCases(expectedValue) {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			// t.Parallel()

			resp := &planmodifier.Int64Response{
				PlanValue: testCase.request.PlanValue,
			}

			int64pm.SetDefault(expectedValue).PlanModifyInt64(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
