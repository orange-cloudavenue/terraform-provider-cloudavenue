package vm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
)

// Specific planmodifier
// if is_detachable is false the VMName/VMID is not modifiable.
func requireReplaceIfNotDetachable() planmodifier.String {
	description := "Attribute requires replacement if `is_detachable` is false"

	return stringplanmodifier.RequiresReplaceIf(stringplanmodifier.RequiresReplaceIfFunc(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
		isDetachable := &types.Bool{}

		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("is_detachable"), isDetachable)...)

		if !isDetachable.ValueBool() {
			resp.RequiresReplace = true
		}
	}), description, description)
}

// DiskSchema returns the schema for the OLDdisk.
func DiskSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The ID of the disk.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"vdc": vdc.Schema(),

		"vapp_id": vapp.Schema()["vapp_id"],

		"vapp_name": vapp.Schema()["vapp_name"],

		"vm_name": schema.StringAttribute{
			MarkdownDescription: "The name of the VM. If `vm_id` is not set and `ìs_detachable` is set to `true`, " +
				"the disk will be attached to any VM. This field is required if `is_detachable` is set to `false`.",
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				requireReplaceIfNotDetachable(),
				removeStateIfConfigIsUnset(),
			},
		},
		"vm_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the VM. If `vm_name` is not set and `ìs_detachable` is set to `true`, " +
				"the disk will be attached to any VM. This field is required if `is_detachable` is set to `false`.",
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				requireReplaceIfNotDetachable(),
				removeStateIfConfigIsUnset(),
			},
		},

		"is_detachable": schema.BoolAttribute{
			MarkdownDescription: "This field specifies whether the disk is detachable. " +
				"If set to true, the disk can be attached to any VM created from the vApp. " +
				"If set to false, the disk will be attached only to the VM specified in `vm_name` or `vm_id`. " +
				"Changing this field will require replacing the disk.",
			Required: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
				boolplanmodifier.UseStateForUnknown(),
			},
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the disk.",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"bus_type":        diskparams.BusTypeAttribute(),
		"size_in_mb":      diskparams.SizeInMBAttribute(),
		"storage_profile": diskparams.StorageProfileAttribute(),
		"bus_number":      diskparams.BusNumberAttribute(),
		"unit_number":     diskparams.UnitNumberAttribute(),
	}
}
