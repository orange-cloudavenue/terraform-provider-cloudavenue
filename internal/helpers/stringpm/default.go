package stringpm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Default(str string) planmodifier.String {
	return &defaultPlanModifier{str}
}

type defaultPlanModifier struct {
	value string
}

var _ planmodifier.String = (*defaultPlanModifier)(nil)

func (str *defaultPlanModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Default value: %s", str)
}

func (str *defaultPlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default value: `%s`", str)
}

func (str *defaultPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	resp.PlanValue = types.StringValue(str.value)
}
