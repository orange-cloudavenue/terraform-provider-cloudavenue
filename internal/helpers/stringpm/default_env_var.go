package stringpm

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DefaultEnvVar(envVar string) planmodifier.String {
	return &stringDefaultEnvVarPlanModifier{envVar}
}

type stringDefaultEnvVarPlanModifier struct {
	EnvVar string
}

var _ planmodifier.String = (*stringDefaultEnvVarPlanModifier)(nil)

func (env *stringDefaultEnvVarPlanModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Default value: %s", env.EnvVar)
}

func (env *stringDefaultEnvVarPlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default value: `%s`", env.EnvVar)
}

func (env *stringDefaultEnvVarPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	// Load the value from the environment variable
	// If the environment variable is not set, we are done here
	v := types.StringValue(os.Getenv(env.EnvVar))
	if !v.IsNull() {
		resp.PlanValue = v
	} else {
		resp.Diagnostics.AddError("Environment variable not set", fmt.Sprintf("The environment variable %s is not set", env.EnvVar))
	}
}
