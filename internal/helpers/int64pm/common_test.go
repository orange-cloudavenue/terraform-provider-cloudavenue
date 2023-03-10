package int64pm_test

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func int64pmTestCases(expectedValue int64) map[string]struct {
	request  planmodifier.Int64Request
	expected *planmodifier.Int64Response
} {
	return map[string]struct {
		request  planmodifier.Int64Request
		expected *planmodifier.Int64Response
	}{
		"null-state": {
			// when we first create the resource, use the unknown
			// value
			request: planmodifier.Int64Request{
				StateValue:  types.Int64Null(),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			expected: &planmodifier.Int64Response{
				PlanValue: types.Int64Value(expectedValue),
			},
		},
		"known-plan": {
			// this would really only happen if we had a plan
			// modifier setting the value before this plan modifier
			// got to it
			//
			// but we still want to preserve that value, in this
			// case
			request: planmodifier.Int64Request{
				StateValue:  types.Int64Value(10),
				PlanValue:   types.Int64Value(11),
				ConfigValue: types.Int64Null(),
			},
			expected: &planmodifier.Int64Response{
				PlanValue: types.Int64Value(11),
			},
		},
		"non-null-state-unknown-plan": {
			// this is the situation we want to preserve the state
			// in
			request: planmodifier.Int64Request{
				StateValue:  types.Int64Value(10),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			expected: &planmodifier.Int64Response{
				PlanValue: types.Int64Value(expectedValue),
			},
		},
		"unknown-config": {
			// this is the situation in which a user is
			// interpolating into a field. We want that to still
			// show up as unknown, otherwise they'll get apply-time
			// errors for changing the value even though we knew it
			// was legitimately possible for it to change and the
			// provider can't prevent this from happening
			request: planmodifier.Int64Request{
				StateValue:  types.Int64Value(10),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Unknown(),
			},
			expected: &planmodifier.Int64Response{
				PlanValue: types.Int64Value(expectedValue),
			},
		},
	}
}
