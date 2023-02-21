// Package vm contains the common code for the VM resource and the VM datasource.
package vm

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	SizeInMb       types.Int64  `tfsdk:"size_in_mb"`
	BusNumber      types.Int64  `tfsdk:"bus_number"`
	UnitNumber     types.Int64  `tfsdk:"unit_number"`
	StorageProfile types.String `tfsdk:"storage_profile"`
}

func InternalDiskSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The ID of the internal disk.",
		},
		"bus_type": schema.StringAttribute{
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			MarkdownDescription: "The type of disk controller. Possible values: `ide`, `parallel` (LSI Logic Parallel SCSI), `sas` (LSI Logic SAS SCSI), `paravirtual` (Paravirtual SCSI), `sata`, `nvme`.",
			Validators: []validator.String{
				stringvalidator.OneOf("ide", "parallel", "sas", "paravirtual", "sata", "nvme"),
			},
		},
		"size_in_mb": schema.Int64Attribute{
			Required:            true,
			MarkdownDescription: "The size of the disk in MB.",
		},
		"bus_number": schema.Int64Attribute{
			Required: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
			MarkdownDescription: "The number of the `SCSI` or `IDE` controller itself.",
		},
		"unit_number": schema.Int64Attribute{
			Required: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
			MarkdownDescription: "The device number on the `SCSI` or `IDE` controller of the disk.",
		},
		"storage_profile": schema.StringAttribute{
			Computed:            true,
			Optional:            true,
			MarkdownDescription: "Storage profile to override the VM default one.",
		},
	}
}
