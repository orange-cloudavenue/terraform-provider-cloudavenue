package vm

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/int64pm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
)

type TemplateDiskModel struct {
	BusType        types.String `tfsdk:"bus_type"`
	BusNumber      types.Int64  `tfsdk:"bus_number"`
	SizeInMb       types.Int64  `tfsdk:"size_in_mb"`
	Iops           types.Int64  `tfsdk:"iops"`
	StorageProfile types.String `tfsdk:"storage_profile"`
	UnitNumber     types.Int64  `tfsdk:"unit_number"`
}

// TemplateDiskSchema returns the schema for the template disk
func TemplateDiskSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"bus_type":        diskparams.BusTypeAttribute(),
		"size_in_mb":      diskparams.SizeInMBAttribute(),
		"bus_number":      diskparams.BusNumberAttribute(),
		"unit_number":     diskparams.UnitNumberAttribute(),
		"storage_profile": diskparams.StorageProfileAttribute(),
		"iops": schema.Int64Attribute{
			MarkdownDescription: "Specifies the IOPS for the disk. Default is 0.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				int64pm.SetDefault(0),
				int64planmodifier.RequiresReplace(),
			},
		},
	}
}
