package vm

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/stringpm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
)

type vmResourceModelDisks []vmResourceModelDisk

type vmResourceModelDisk struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	BusType        types.String `tfsdk:"bus_type"`
	SizeInMb       types.Int64  `tfsdk:"size_in_mb"`
	StorageProfile types.String `tfsdk:"storage_profile"`
	BusNumber      types.Int64  `tfsdk:"bus_number"`
	UnitNumber     types.Int64  `tfsdk:"unit_number"`
}

func vmResourceModelDiskAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"name":            types.StringType,
		"bus_type":        types.StringType,
		"size_in_mb":      types.Int64Type,
		"storage_profile": types.StringType,
		"bus_number":      types.Int64Type,
		"unit_number":     types.Int64Type,
	}
}

func (d *vmResourceModelDisks) ObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: vmResourceModelDiskAttrType(),
	}
}

// DiskInternalExternalSchema returns schema for internal and external disks.
func diskInternalExternalSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The ID of the disk.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the disk.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringpm.SetDefaultEmptyString(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"bus_type":        diskparams.BusTypeAttribute(),
		"size_in_mb":      diskparams.SizeInMBAttribute(),
		"storage_profile": diskparams.StorageProfileAttribute(),
		"bus_number":      diskparams.BusNumberAttribute(),
		"unit_number":     diskparams.UnitNumberAttribute(),
	}
}

// ToPlan converts the vmResourceModelDisks struct to a terraform plan.
func (d *vmResourceModelDisks) ToPlan(ctx context.Context) (types.Set, diag.Diagnostics) {
	if d == nil {
		return types.SetNull(d.ObjectType()), diag.Diagnostics{}
	}

	return types.SetValueFrom(ctx, d.ObjectType(), d)
}

// DisksFromPlan converts the terraform plan to a OLDDisk struct.
func DisksFromPlan(ctx context.Context, x types.Set) ([]vmResourceModelDisk, diag.Diagnostics) {
	if x.IsNull() || x.IsUnknown() {
		return []vmResourceModelDisk{}, diag.Diagnostics{}
	}

	c := make([]vmResourceModelDisk, 0)

	d := x.ElementsAs(ctx, c, false)

	return c, d
}

// IsExternal checks if the disk is external.
func (d *vmResourceModelDisk) IsExternal() bool {
	return !d.ID.IsNull()
}

// DiskFromInternalDisk converts an internal disk to a vmResourceModelDisk.
func DiskFromInternalDisk(internalDisk vm.InternalDisk) vmResourceModelDisk {
	return vmResourceModelDisk{
		ID:             internalDisk.ID,
		BusType:        internalDisk.BusType,
		SizeInMb:       internalDisk.SizeInMb,
		StorageProfile: internalDisk.StorageProfile,
		BusNumber:      internalDisk.BusNumber,
		UnitNumber:     internalDisk.UnitNumber,
	}
}

// DiskFromExternalDisk converts an external disk to a vmResourceModelDisk.
func DiskFromExternalDisk(externalDisk vm.Disk, busNumber, unitNumber types.Int64) vmResourceModelDisk {
	return vmResourceModelDisk{
		ID:             externalDisk.ID,
		Name:           externalDisk.Name,
		BusType:        externalDisk.BusType,
		SizeInMb:       externalDisk.SizeInMb,
		StorageProfile: externalDisk.StorageProfile,
		BusNumber:      busNumber,
		UnitNumber:     unitNumber,
	}
}

// DiskFromGovcdDiskSettings converts a govcd disk settings to a vmResourceModelDisk.
func DiskFromGovcdDiskSettings(diskSettings *govcdtypes.DiskSettings) vmResourceModelDisk {
	d := vmResourceModelDisk{
		BusType:    types.StringValue(vm.GetBusTypeByCode(diskSettings.AdapterType).Name()),
		SizeInMb:   types.Int64Value(diskSettings.SizeMb),
		BusNumber:  types.Int64Value(int64(diskSettings.BusNumber)),
		UnitNumber: types.Int64Value(int64(diskSettings.UnitNumber)),
	}

	if diskSettings.Disk != nil {
		d.ID = types.StringValue(diskSettings.Disk.ID)
		d.Name = types.StringValue(diskSettings.Disk.Name)
	}

	if diskSettings.StorageProfile != nil {
		d.StorageProfile = types.StringValue(diskSettings.StorageProfile.Name)
	}

	return d
}

// Append appends a disk to the list of disks.
func (d *vmResourceModelDisks) Append(disk ...vmResourceModelDisk) {
	*d = append(*d, disk...)
}

/*
DiskRead

Reads the disks of a VM.
*/
func DisksRead(vm *govcd.VM) (disks vmResourceModelDisks, err error) {
	if vm.VM != nil && vm.VM.VmSpecSection != nil && vm.VM.VmSpecSection.DiskSection != nil && vm.VM.VmSpecSection.DiskSection.DiskSettings != nil {
		disks := vmResourceModelDisks{}

		for _, disk := range vm.VM.VmSpecSection.DiskSection.DiskSettings {
			vmResourceModelDisk := DiskFromGovcdDiskSettings(disk)
			disks.Append(vmResourceModelDisk)
		}

		return disks, nil
	}
	return vmResourceModelDisks{}, errors.New("no disks found")
}
