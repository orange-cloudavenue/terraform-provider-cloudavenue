// Package stringpm provides a plan modifier for string values.
package stringpm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RequiresReplaceIfFunc is a conditional function used in the RequiresReplaceIf
// plan modifier to determine whether the attribute requires replacement.
type DefaultFunc func(context.Context, planmodifier.StringRequest, *DefaultFuncResponse)

// RequiresReplaceIfFuncResponse is the response type for a RequiresReplaceIfFunc.
type DefaultFuncResponse struct {
	// Diagnostics report errors or warnings related to this logic. An empty
	// or unset slice indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics

	// Value is the value to use by default if the attribute is not configured.
	Value string
}

// setDefaultFunc returns a plan modifier that conditionally requires
// resource replacement if:
//
//   - The resource is planned for update.
//   - The plan and state values are not equal.
//   - The given function returns true. Returning false will not unset any
//     prior resource replacement.
func setDefaultFunc(f DefaultFunc, description, markdownDescription string) planmodifier.String {
	return defaultFuncPlanModifier{
		f:                   f,
		description:         description,
		markdownDescription: markdownDescription,
	}
}

// defaultFuncPlanModifier is an plan modifier that sets RequiresReplace
// on the attribute if a given function is true.
type defaultFuncPlanModifier struct {
	f                   DefaultFunc
	description         string
	markdownDescription string
}

// Description returns a human-readable description of the plan modifier.
func (m defaultFuncPlanModifier) Description(_ context.Context) string {
	return m.description
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m defaultFuncPlanModifier) MarkdownDescription(_ context.Context) string {
	return m.markdownDescription
}

// PlanModifyString implements the plan modification logic.
func (m defaultFuncPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do not replace if the plan and state values are equal.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	funcResp := &DefaultFuncResponse{}

	m.f(ctx, req, funcResp)

	resp.Diagnostics.Append(funcResp.Diagnostics...)
	resp.PlanValue = types.StringValue(funcResp.Value)
}
