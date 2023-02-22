// Package vm contains the common code for the VM resource and the VM datasource.
package vm

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
)

var (
	// InternalDiskBusTypes is a map of internal disk bus types.
	InternalDiskBusTypes = map[string]string{
		"ide":         "1",
		"parallel":    "3",
		"sas":         "4",
		"paravirtual": "5",
		"sata":        "6",
		"nvme":        "7",
	}
	// InternalDiskBusTypesFromValues is a map of internal disk bus types.
	InternalDiskBusTypesFromValues = map[string]string{
		"1": "ide",
		"3": "parallel",
		"4": "sas",
		"5": "paravirtual",
		"6": "sata",
		"7": "nvme",
	}
)

type InternalDiskModel struct {
	ID             types.String `tfsdk:"id"`
	BusType        types.String `tfsdk:"bus_type"`
	BusNumber      types.Int64  `tfsdk:"bus_number"`
	SizeInMb       types.Int64  `tfsdk:"size_in_mb"`
	UnitNumber     types.Int64  `tfsdk:"unit_number"`
	StorageProfile types.String `tfsdk:"storage_profile"`
}

func InternalDiskSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The ID of the internal disk.",
		},
		"bus_type":        diskparams.BusTypeAttribute(),
		"size_in_mb":      diskparams.SizeInMBAttribute(),
		"bus_number":      diskparams.BusNumberAttribute(),
		"unit_number":     diskparams.UnitNumberAttribute(),
		"storage_profile": diskparams.StorageProfileAttribute(),
	}
}
