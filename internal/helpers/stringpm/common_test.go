package stringpm_test

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func stringpmTestCases(expectedValue string) map[string]struct {
	request  planmodifier.StringRequest
	expected *planmodifier.StringResponse
} {
	return map[string]struct {
		request  planmodifier.StringRequest
		expected *planmodifier.StringResponse
	}{
		"null-state": {
			// when we first create the resource, use the unknown
			// value
			request: planmodifier.StringRequest{
				StateValue:  types.StringNull(),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue(expectedValue),
			},
		},
		"known-plan": {
			// this would really only happen if we had a plan
			// modifier setting the value before this plan modifier
			// got to it
			//
			// but we still want to preserve that value, in this
			// case
			request: planmodifier.StringRequest{
				StateValue:  types.StringValue("other"),
				PlanValue:   types.StringValue("test"),
				ConfigValue: types.StringNull(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("test"),
			},
		},
		"non-null-state-unknown-plan": {
			// this is the situation we want to preserve the state
			// in
			request: planmodifier.StringRequest{
				StateValue:  types.StringValue("test"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue(expectedValue),
			},
		},
		"unknown-config": {
			// this is the situation in which a user is
			// interpolating into a field. We want that to still
			// show up as unknown, otherwise they'll get apply-time
			// errors for changing the value even though we knew it
			// was legitimately possible for it to change and the
			// provider can't prevent this from happening
			request: planmodifier.StringRequest{
				StateValue:  types.StringValue("test"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringUnknown(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue(expectedValue),
			},
		},
	}
}
