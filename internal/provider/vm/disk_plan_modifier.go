package vm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// removeStateIfConfigIsUnset returns a plan modifier that removes the state value
// if the config value is unset.
func removeStateIfConfigIsUnset() planmodifier.String {
	return removeStateIfConfigIsUnsetModifier{}
}

// removeStateIfConfigIsUnsetModifier implements the plan modifier.
type removeStateIfConfigIsUnsetModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m removeStateIfConfigIsUnsetModifier) Description(_ context.Context) string {
	return ""
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m removeStateIfConfigIsUnsetModifier) MarkdownDescription(_ context.Context) string {
	return ""
}

// PlanModifyString implements the plan modification logic.
func (m removeStateIfConfigIsUnsetModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	vmName := &types.String{}
	vmID := &types.String{}

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("vm_name"), vmName)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("vm_id"), vmID)...)

	// Fields are empty in Config
	if vmName.ValueString() == "" && vmID.ValueString() == "" {
		resp.PlanValue = types.StringNull()
		return
	}

	// Actual field is vm_name and vm_id is not empty
	if req.Path.Equal(path.Root("vm_name")) && vmID.ValueString() != "" {
		if req.StateValue.ValueString() == "" {
			resp.PlanValue = types.StringUnknown()
		}
		return
	}

	// Actual field is vm_id and vm_name is not empty
	if req.Path.Equal(path.Root("vm_id")) && vmName.ValueString() != "" {
		if req.StateValue.ValueString() == "" {
			resp.PlanValue = types.StringUnknown()
		}
		return
	}
}
