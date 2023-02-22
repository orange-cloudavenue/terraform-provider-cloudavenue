package vm

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DiskModel struct {
	Name       types.String `tfsdk:"name"`
	BusNumber  types.Int64  `tfsdk:"bus_number"`
	SizeInMb   types.Int64  `tfsdk:"size_in_mb"`
	UnitNumber types.Int64  `tfsdk:"unit_number"`
}

// DiskSchema returns the schema for the disk
func DiskSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			MarkdownDescription: "Independent disk name.",
			Required:            true,
		},
		"bus_number": schema.StringAttribute{
			MarkdownDescription: "Bus number on which to place the disk controller.",
			Required:            true,
		},
		"size_in_mb": schema.Int64Attribute{
			MarkdownDescription: "The size of the disk in MB.",
			Computed:            true,
		},
		"unit_number": schema.Int64Attribute{
			MarkdownDescription: "Unit number (slot) on the bus specified by BusNumber",
			Required:            true,
		},
	}
}
