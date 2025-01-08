/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vm contains the common code for the VM resource and the VM datasource.
package vm

import (
	"context"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
)

type busType struct {
	key  string
	name string
	code string
}

func (b busType) Name() string {
	return strings.ToUpper(b.name)
}

func (b busType) Code() string {
	return b.code
}

// INTERNAL DISK BUS TYPES
// <option _ngcontent-ldg-c121="" value="4" class="ng-star-inserted"> LSI Logic SAS (SCSI) </option>
// <option _ngcontent-ldg-c121="" value="1" class="ng-star-inserted"> IDE </option>
// <option _ngcontent-ldg-c121="" value="6" class="ng-star-inserted"> SATA </option>
// <option _ngcontent-ldg-c121="" value="7" class="ng-star-inserted"> NVME </option>

var (
	busTypeSATA = busType{key: "sata", name: "sata", code: "6"} // Bus type SATA
	busTypeSCSI = busType{key: "sas", name: "scsi", code: "4"}  // Bus type SCSI
	busTypeNVME = busType{key: "nvme", name: "nvme", code: "7"} // Bus type NVME
)

func GetBusTypeByCode(code string) busType {
	switch code {
	case busTypeSATA.code:
		return busTypeSATA
	case busTypeSCSI.code:
		return busTypeSCSI
	case busTypeNVME.code:
		return busTypeNVME
	default:
		return busTypeSATA
	}
}

func GetBusTypeByKey(key string) busType {
	switch strings.ToLower(key) {
	case busTypeSATA.name:
		return busTypeSATA
	case busTypeSCSI.name:
		return busTypeSCSI
	case busTypeNVME.name:
		return busTypeNVME
	default:
		return busTypeSATA
	}
}

type InternalDisks []InternalDisk

type InternalDisk struct {
	ID             types.String `tfsdk:"id"`
	BusType        types.String `tfsdk:"bus_type"`
	BusNumber      types.Int64  `tfsdk:"bus_number"`
	SizeInMb       types.Int64  `tfsdk:"size_in_mb"`
	UnitNumber     types.Int64  `tfsdk:"unit_number"`
	StorageProfile types.String `tfsdk:"storage_profile"`
}

// InternalDiskAttrType returns the type map for the internal disk.
func InternalDiskAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"bus_type":        types.StringType,
		"bus_number":      types.Int64Type,
		"size_in_mb":      types.Int64Type,
		"unit_number":     types.Int64Type,
		"storage_profile": types.StringType,
	}
}

func InternalDiskSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The ID of the internal disk.",
			Computed:            true,
		},
		"bus_type":        diskparams.BusTypeAttribute(),
		"size_in_mb":      diskparams.SizeInMBAttribute(),
		"bus_number":      diskparams.BusNumberAttribute(),
		"unit_number":     diskparams.UnitNumberAttribute(),
		"storage_profile": diskparams.StorageProfileAttribute(),
	}
}

func InternalDiskSchemaComputed() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The ID of the internal disk.",
			Computed:            true,
		},
		"bus_type":        diskparams.BusTypeAttributeComputed(),
		"size_in_mb":      diskparams.SizeInMBAttributeComputed(),
		"bus_number":      diskparams.BusNumberAttributeComputed(),
		"unit_number":     diskparams.UnitNumberAttributeComputed(),
		"storage_profile": diskparams.StorageProfileAttributeComputed(),
	}
}

// ToPlan converts a InternalDisks struct to a terraform plan.
func (i *InternalDisks) ToPlan(ctx context.Context) (basetypes.SetValue, diag.Diagnostics) {
	if i == nil {
		return types.SetUnknown(types.ObjectType{AttrTypes: InternalDiskAttrType()}), diag.Diagnostics{}
	}

	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: InternalDiskAttrType()}, i)
}

// InternalDiskFromPlan creates a InternalDisks from a plan.
func InternalDiskFromPlan(ctx context.Context, x types.Set) (*InternalDisks, diag.Diagnostics) {
	if x.IsNull() || x.IsUnknown() {
		return &InternalDisks{}, diag.Diagnostics{}
	}

	i := &InternalDisks{}

	d := x.ElementsAs(ctx, i, false)

	return i, d
}

/*
InternalDiskCreate

Creates a new internal disk associated with a VM.
*/
func InternalDiskCreate(ctx context.Context, c *client.CloudAvenue, disk InternalDisk, vAppName, vmName, vdcName types.String) (newDisk *InternalDisk, d diag.Diagnostics) {
	vdc, err := c.GetVDC(vdcName.ValueString())
	if err != nil {
		d.AddError("Error retrieving VDC", err.Error())
		return nil, d
	}

	myVM, err := GetVMOLD(vdc.Vdc, vAppName.ValueString(), vmName.ValueString())
	if err != nil {
		d.AddError("Error retrieving VM", err.Error())
		return nil, d
	}

	// storage profile
	var (
		storageProfilePrt *govcdtypes.Reference
		overrideVMDefault bool
	)

	if disk.StorageProfile.IsNull() || disk.StorageProfile.IsUnknown() {
		storageProfilePrt = myVM.VM.StorageProfile
		overrideVMDefault = false
	} else {
		storageProfile, errFindStorage := vdc.FindStorageProfileReference(disk.StorageProfile.ValueString())
		if errFindStorage != nil {
			d.AddError("Error retrieving storage profile", errFindStorage.Error())
			return nil, d
		}
		storageProfilePrt = &storageProfile
		overrideVMDefault = true
	}

	// value is required but not treated.
	isThinProvisioned := true

	var busNumber, unitNumber types.Int64

	if disk.BusNumber.IsNull() || disk.BusNumber.IsUnknown() {
		b, u := diskparams.ComputeBusAndUnitNumber(myVM.VM.VmSpecSection.DiskSection.DiskSettings)
		busNumber = types.Int64Value(int64(b))
		unitNumber = types.Int64Value(int64(u))
	} else {
		busNumber = disk.BusNumber
		unitNumber = disk.UnitNumber
	}

	diskSetting := &govcdtypes.DiskSettings{
		SizeMb:              disk.SizeInMb.ValueInt64(),
		UnitNumber:          int(busNumber.ValueInt64()),
		BusNumber:           int(unitNumber.ValueInt64()),
		AdapterType:         GetBusTypeByKey(disk.BusType.ValueString()).Code(),
		ThinProvisioned:     &isThinProvisioned,
		StorageProfile:      storageProfilePrt,
		VirtualQuantityUnit: "byte",
		OverrideVmDefault:   overrideVMDefault,
	}

	diskID, err := myVM.AddInternalDisk(diskSetting)
	if err != nil {
		d.AddError("Error creating disk", err.Error())
		return nil, d
	}

	newDisk = &disk
	newDisk.ID = types.StringValue(diskID)
	newDisk.BusType = types.StringValue(GetBusTypeByCode(diskSetting.AdapterType).Name())
	newDisk.SizeInMb = types.Int64Value(diskSetting.SizeMb)
	newDisk.StorageProfile = types.StringValue(storageProfilePrt.Name)

	return newDisk, d
}

/*
InternalDiskRead

Reads an internal disk associated with a VM.
*/
func InternalDiskRead(ctx context.Context, client *client.CloudAvenue, disk *InternalDisk, vm *govcd.VM) (readDisk *InternalDisk, d diag.Diagnostics) {
	diskSettings, err := vm.GetInternalDiskById(disk.ID.ValueString(), true)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// If the disk is not found, then return nil so that we can show
			return nil, nil
		}
		d.AddError("Error retrieving disk with id "+disk.ID.ValueString(), err.Error())
		return
	}

	readDisk = disk
	readDisk.ID = types.StringValue(diskSettings.DiskId)
	readDisk.BusType = types.StringValue(GetBusTypeByCode(diskSettings.AdapterType).Name())
	readDisk.SizeInMb = types.Int64Value(diskSettings.SizeMb)
	readDisk.StorageProfile = types.StringValue(diskSettings.StorageProfile.Name)

	return
}

/*
InternalDiskUpdate

Updates an internal disk associated with a VM.
*/
func InternalDiskUpdate(ctx context.Context, c *client.CloudAvenue, disk InternalDisk, vAppName, vmName, vdcName types.String) (updatedDisk *InternalDisk, d diag.Diagnostics) {
	vdc, err := c.GetVDC(vdcName.ValueString())
	if err != nil {
		d.AddError("Error retrieving VDC", err.Error())
		return
	}

	myVM, err := GetVMOLD(vdc.Vdc, vAppName.ValueString(), vmName.ValueString())
	if err != nil {
		d.AddError("Error retrieving VM", err.Error())
		return
	}

	diskSettingsToUpdate, err := myVM.GetInternalDiskById(disk.ID.ValueString(), false)
	if err != nil {
		d.AddError("Error retrieving disk", err.Error())
		return
	}

	diskSettingsToUpdate.SizeMb = disk.SizeInMb.ValueInt64()
	// Note can't change adapter type, bus number, unit number as vSphere changes diskId

	var (
		storageProfilePrt *govcdtypes.Reference
		overrideVMDefault bool
	)

	storageProfileName := disk.StorageProfile.ValueString()
	if storageProfileName != "" {
		storageProfile, errFindStorage := vdc.FindStorageProfileReference(storageProfileName)
		if errFindStorage != nil {
			d.AddError("Error retrieving storage profile", errFindStorage.Error())
			return
		}
		storageProfilePrt = &storageProfile
		overrideVMDefault = true
	} else {
		storageProfilePrt = myVM.VM.StorageProfile
		overrideVMDefault = false
	}

	diskSettingsToUpdate.StorageProfile = storageProfilePrt
	diskSettingsToUpdate.OverrideVmDefault = overrideVMDefault

	_, err = myVM.UpdateInternalDisks(myVM.VM.VmSpecSection)
	if err != nil {
		d.AddError("Error updating disk", err.Error())
		return
	}

	updatedDisk = &disk
	updatedDisk.ID = types.StringValue(diskSettingsToUpdate.DiskId)
	updatedDisk.BusType = types.StringValue(GetBusTypeByCode(diskSettingsToUpdate.AdapterType).Name())
	updatedDisk.SizeInMb = types.Int64Value(diskSettingsToUpdate.SizeMb)
	updatedDisk.StorageProfile = types.StringValue(storageProfilePrt.Name)

	return updatedDisk, nil
}
