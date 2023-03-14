package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

/*
Disks is a list of Disk.
*/
type Disks []Disk

/*
Disk is independent disk attached to a VM.
This disk is always detached disk type.
*/
type Disk struct {
	ID             types.String `tfsdk:"id"`
	VAppName       types.String `tfsdk:"vapp_name"`
	VAppID         types.String `tfsdk:"vapp_id"`
	VMName         types.String `tfsdk:"vm_name"`
	VMID           types.String `tfsdk:"vm_id"`
	VDC            types.String `tfsdk:"vdc"`
	Name           types.String `tfsdk:"name"`
	SizeInMb       types.Int64  `tfsdk:"size_in_mb"`
	StorageProfile types.String `tfsdk:"storage_profile"`
	IsDetachable   types.Bool   `tfsdk:"is_detachable"`

	BusType    types.String `tfsdk:"bus_type"`
	BusNumber  types.Int64  `tfsdk:"bus_number"`
	UnitNumber types.Int64  `tfsdk:"unit_number"`
}

/*
ToAttrValue

converts the Disk struct to a terraform plan

attr.Value is a map[string]attr.Value
  - "id"
  - "vdc"
  - "vapp_name"
  - "vapp_id"
  - "vm_name"
  - "vm_id"
  - "name"
  - "size_in_mb"
  - "storage_profile"
  - "is_detachable"
  - "bus_type"
  - "bus_number"
  - "unit_number"
*/
func (d *Disk) ToAttrValue() map[string]attr.Value {
	return map[string]attr.Value{
		"id":              d.ID,
		"vapp_name":       d.VAppName,
		"vapp_id":         d.VAppID,
		"vm_name":         d.VMName,
		"vm_id":           d.VMID,
		"vdc":             d.VDC,
		"name":            d.Name,
		"size_in_mb":      d.SizeInMb,
		"storage_profile": d.StorageProfile,
		"is_detachable":   d.IsDetachable,
		"bus_type":        d.BusType,
		"bus_number":      d.BusNumber,
		"unit_number":     d.UnitNumber,
	}
}

/*
DiskAttrType

returns the attr.Type map for the disk

attr.Type is a map[string]attr.Type
  - "id" 				(types.StringType)
  - "vapp_name" 		(types.StringType)
  - "vapp_id" 			(types.StringType)
  - "vm_name" 			(types.StringType)
  - "vm_id" 			(types.StringType)
  - "vdc" 				(types.StringType)
  - "name" 				(types.StringType)
  - "size_in_mb" 		(types.Int64Type)
  - "storage_profile" 	(types.StringType)
  - "bus_type" 			(types.StringType)
  - "bus_number" 		(types.Int64Type)
  - "unit_number" 		(types.Int64Type)
*/
func DiskAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"vapp_name":       types.StringType,
		"vapp_id":         types.StringType,
		"vm_name":         types.StringType,
		"vm_id":           types.StringType,
		"vdc":             types.StringType,
		"name":            types.StringType,
		"bus_type":        types.StringType,
		"size_in_mb":      types.Int64Type,
		"storage_profile": types.StringType,
		"is_detachable":   types.BoolType,
		"bus_number":      types.Int64Type,
		"unit_number":     types.Int64Type,
	}
}

// DisksFromPlan converts the terraform plan to a OLDDisk struct.
func DisksFromPlan(ctx context.Context, x types.Set) (*Disk, diag.Diagnostics) {
	if x.IsNull() || x.IsUnknown() {
		return &Disk{}, diag.Diagnostics{}
	}

	c := &Disk{}

	d := x.ElementsAs(ctx, c, false)

	return c, d
}

/*
ElementType

return the attr.Type for the disk.
*/
func (d *Disk) ElementType() attr.Type {
	return types.ObjectType{AttrTypes: DiskAttrType()}
}

func (d *Disks) ElementType() attr.Type {
	return types.ObjectType{AttrTypes: DiskAttrType()}
}

// ToPlan converts the disk struct to a terraform plan.
func (d *Disks) ToPlan(ctx context.Context) (basetypes.SetValue, diag.Diagnostics) {
	if d == nil {
		return types.SetNull(d.ElementType()), diag.Diagnostics{}
	}

	return types.SetValueFrom(ctx, d.ElementType(), d)
}

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
			Computed: true,
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
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				requireReplaceIfNotDetachable(),
				removeStateIfConfigIsUnset(),
			},
		},

		"is_detachable": schema.BoolAttribute{
			MarkdownDescription: "If the disk is detachable or not. If set to `true`, the disk will be attached to any VM " +
				"that is created from the vApp. If set to `false`, the disk will be attached to the VM specified in `vm_name` or `vm_id`. " +
				"Change this field requires a replacement of the disk.",
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

/*
DiskCreate

creates a detachable disk.
*/
func DiskCreate(ctx context.Context, org org.Org, vdc vdc.VDC, vm *govcd.VM, disk *Disk, inVapp vapp.VAPP) (*Disk, diag.Diagnostics) {
	d := diag.Diagnostics{}

	if inVapp.VApp == nil || org.Org == nil || vdc.VDC == nil {
		d.AddError("Error creating disk", "Empty vApp, org or vdc")
		return nil, d
	}

	// Lock vApp
	d.Append(inVapp.LockVAPP(ctx)...)
	if d.HasError() {
		return nil, d
	}
	defer d.Append(inVapp.UnlockVAPP(ctx)...)

	// Checking if the disk name is already existing in the vDC
	existingDisk, err := vdc.QueryDisk(disk.Name.ValueString())
	if existingDisk != (govcd.DiskRecord{}) || err == nil {
		d.AddError("already exists in the vDC", fmt.Sprintf("The disk %s already exists in the vDC", disk.Name.ValueString()))
		return nil, d
	}

	// Init struct for creating a disk
	diskCreateParams := &govcdtypes.DiskCreateParams{
		Disk: &govcdtypes.Disk{
			Name:        disk.Name.ValueString(),
			SizeMb:      disk.SizeInMb.ValueInt64(),
			SharingType: "None",
		},
	}

	// If the bus type is set checking if it exists and setting it
	if !disk.BusType.IsNull() && !disk.BusType.IsUnknown() {
		diskCreateParams.Disk.BusType = diskparams.GetBusTypeByName(disk.BusType.ValueString()).Code()
		diskCreateParams.Disk.BusSubType = diskparams.GetBusTypeByName(disk.BusType.ValueString()).SubType()
	}

	// If the storage profile is set checking if it exists and setting it
	if !disk.StorageProfile.IsNull() && !disk.StorageProfile.IsUnknown() {
		storageReference, err := vdc.FindStorageProfileReference(disk.StorageProfile.ValueString())
		if err != nil {
			d.AddError("storage profile not found", fmt.Sprintf("The storage profile %s does not exist in the vDC", disk.StorageProfile.ValueString()))
			return nil, d
		}
		diskCreateParams.Disk.StorageProfile = &govcdtypes.Reference{HREF: storageReference.HREF}
	}

	// Create the disk
	task, err := vdc.CreateDisk(diskCreateParams)
	if err != nil {
		d.AddError("error creating disk", err.Error())
		return nil, d
	}

	// Wait for the task to finish
	err = task.WaitTaskCompletion()
	if err != nil {
		d.AddError("error on creating disk", err.Error())
		return nil, d
	}

	// Get the disk by the Href
	x, err := vdc.GetDiskByHref(task.Task.Owner.HREF)
	if err != nil {
		d.AddError("unable to find disk after creating", fmt.Sprintf("unable to find disk with href %s: %s", task.Task.HREF, err))
		return nil, d
	}

	if x.Disk == nil {
		d.AddError("unable to find disk after creating", fmt.Sprintf("unable to find disk with href %s", task.Task.HREF))
		return nil, d
	}

	// Try to attach the disk to the VM if it is detachable and the VM is not nil (vm is nil if VMName or VMID is not set)
	if disk.IsDetachable.ValueBool() && vm != nil {
		var busNumber, unitNumber types.Int64

		// If the bus number or the unit number are not set, compute them
		if disk.BusNumber.IsNull() || disk.BusNumber.IsUnknown() {
			b, u := diskparams.ComputeBusAndUnitNumber(vm.VM.VmSpecSection.DiskSection.DiskSettings)
			busNumber = types.Int64Value(int64(b))
			unitNumber = types.Int64Value(int64(u))
		} else {
			busNumber = disk.BusNumber
			unitNumber = disk.UnitNumber
		}

		disk.ID = types.StringValue(x.Disk.Id)

		// Attach the disk to the VM
		d.Append(DiskAttach(ctx, vdc, &Disk{
			ID:         disk.ID,
			Name:       disk.Name,
			BusNumber:  busNumber,
			UnitNumber: unitNumber,
		}, vm)...)
		if d.HasError() {
			return nil, d
		}
	}

	dCreated := disk
	dCreated.ID = types.StringValue(x.Disk.Id)
	dCreated.Name = types.StringValue(x.Disk.Name)
	dCreated.BusType = types.StringValue(diskparams.GetBusTypeByCode(x.Disk.BusType, x.Disk.BusSubType).Name())
	dCreated.StorageProfile = types.StringValue(x.Disk.StorageProfile.Name)

	return dCreated, d
}

/*
DiskRead

reads a detachable disk
If the disk is not found, returns nil

if disk and diag.Diagnostics are nil the disk is
not found and the resource should be removed from state.
*/
func DiskRead(ctx context.Context, client *client.CloudAvenue, org org.Org, vdc vdc.VDC, disk *Disk, inVapp vapp.VAPP) (*Disk, diag.Diagnostics) {
	d := diag.Diagnostics{}

	var (
		x   *govcd.Disk
		err error
	)

	if inVapp.VApp == nil || org.Org == nil || vdc.VDC == nil {
		d.AddError("Error read disk", "Empty vApp, org or vdc")
		return nil, d
	}

	// Lock vApp
	d.Append(inVapp.LockVAPP(ctx)...)
	if d.HasError() {
		return nil, d
	}
	defer d.Append(inVapp.UnlockVAPP(ctx)...)

	var r *govcd.VM
	// VMName is required if VMID is set
	if disk.VMID.ValueString() != "" {
		r, err = inVapp.GetVMById(disk.VMID.ValueString(), true)
		if err != nil {
			d.AddError("unable to find vm", fmt.Sprintf("unable to find vm with id %s: %s", disk.VMID.ValueString(), err))
			return nil, d
		}
		disk.VMName = types.StringValue(r.VM.Name)
	} else if disk.VMName.ValueString() != "" {
		r, err = inVapp.GetVMByName(disk.VMName.ValueString(), true)
		if err != nil {
			d.AddError("unable to find vm", fmt.Sprintf("unable to find vm with name %s: %s", disk.VMName.ValueString(), err))
			return nil, d
		}
		disk.VMID = types.StringValue(r.VM.ID)
	}

	dRead := disk

	if disk.IsDetachable.ValueBool() {
		if !disk.ID.IsNull() && !disk.ID.IsUnknown() {
			// Get the disk by the ID
			x, err = vdc.GetDiskById(disk.ID.ValueString(), true)
			if err != nil {
				if govcd.IsNotFound(err) {
					return nil, nil
				}
				d.AddError("unable to find disk", fmt.Sprintf("unable to find disk with id %s: %s", disk.ID.ValueString(), err))
				return nil, d
			}
		} else {
			// Get the disk by the Name
			disks, err := vdc.GetDisksByName(disk.Name.ValueString(), true)
			if err != nil {
				if govcd.IsNotFound(err) {
					return nil, nil
				}
				d.AddError("unable to find disk", fmt.Sprintf("unable to find disk with name %s: %s", disk.Name.ValueString(), err))
				return nil, d
			}

			if len(*disks) > 1 {
				d.AddError("multiple disks found", fmt.Sprintf("multiple disks found with name %s", disk.Name.ValueString()))
				return nil, d
			}

			x = &(*disks)[0]
		}

		attachedVmsHrefs, err := x.GetAttachedVmsHrefs()
		if err != nil {
			d.AddError("unable to find attached VM", fmt.Sprintf("unable to find attached VM for disk %s: %s", disk.Name.ValueString(), err))
			return nil, d
		}

		// Normally a disk can be attached to only one VM
		if len(attachedVmsHrefs) > 1 {
			d.AddError("multiple VMs attached", fmt.Sprintf("multiple VMs attached to disk %s", disk.Name.ValueString()))
			return nil, d
		}

		if len(attachedVmsHrefs) == 1 {
			r, err = client.Vmware.Client.GetVMByHref(attachedVmsHrefs[0])
			if err != nil {
				d.AddError("unable to find attached VM", fmt.Sprintf("unable to find attached VM for disk %s: %s", disk.Name.ValueString(), err))
				return nil, d
			}

			if r.VM == nil {
				d.AddError("unable to find attached VM", fmt.Sprintf("unable to find attached VM for disk %s: %s", disk.Name.ValueString(), err))
				return nil, d
			}

			err = r.Refresh()
			if err != nil {
				d.AddError("unable to refresh attached VM", fmt.Sprintf("unable to refresh attached VM for disk %s: %s", disk.Name.ValueString(), err))
				return nil, d
			}
		}

		dRead.ID = types.StringValue(x.Disk.Id)
		dRead.Name = types.StringValue(x.Disk.Name)
		dRead.BusType = types.StringValue(diskparams.GetBusTypeByCode(x.Disk.BusType, x.Disk.BusSubType).Name())
		dRead.SizeInMb = types.Int64Value(x.Disk.SizeMb)
		dRead.StorageProfile = types.StringValue(x.Disk.StorageProfile.Name)
	} else {
		// Internal Disk
		internalDisk, diagErr := InternalDiskRead(ctx, client, &InternalDisk{
			ID: disk.ID,
		}, r)
		if internalDisk == nil && diagErr == nil {
			return nil, nil
		}
		d.Append(diagErr...)
		if diagErr.HasError() {
			return nil, d
		}

		dRead.ID = internalDisk.ID
		dRead.BusType = internalDisk.BusType
		dRead.SizeInMb = internalDisk.SizeInMb
		dRead.StorageProfile = internalDisk.StorageProfile
	}

	if r != nil {
		dRead.VMName = types.StringValue(r.VM.Name)
		dRead.VMID = types.StringValue(r.VM.ID)
	} else {
		if disk.VMName.ValueString() != "" {
			dRead.VMName = types.StringValue("")
		}
		if disk.VMID.ValueString() != "" {
			dRead.VMID = types.StringValue("")
		}
	}

	return dRead, d
}

/*
DiskUpdate

updates a detachable disk

List of attributes require attach/detach disk to be updated if the disk is detachable and attached to a VM:
  - size_in_mb
  - storage_profile
*/
func DiskUpdate(ctx context.Context, client *client.CloudAvenue, diskPlan, diskState *Disk, org org.Org, vdc vdc.VDC, inVapp vapp.VAPP) (*Disk, diag.Diagnostics) { //nolint:gocyclo
	d := diag.Diagnostics{}

	// Preventing nil pointer
	if diskPlan == nil || diskState == nil {
		d.AddError("disk is nil", "disk is nil")
		return nil, d
	}

	if inVapp.VApp == nil || org.Org == nil || vdc.VDC == nil {
		d.AddError("Error read disk", "Empty vApp, org or vdc")
		return nil, d
	}

	// Lock vApp
	d.Append(inVapp.LockVAPP(ctx)...)
	if d.HasError() {
		return nil, d
	}
	defer d.Append(inVapp.UnlockVAPP(ctx)...)

	// If VDC is not defined at resource level, use the one defined at provider level
	if diskPlan.VDC.IsNull() || diskPlan.VDC.IsUnknown() {
		if client.DefaultVDCExist() {
			diskPlan.VDC = types.StringValue(client.GetDefaultVDC())
		} else {
			d.AddError("VDC is required", "VDC is required when not defined at provider level")
			return nil, d
		}
	}

	if diskPlan.IsDetachable.ValueBool() {
		// Get the disk by the ID
		x, err := vdc.GetDiskById(diskState.ID.ValueString(), true)
		if err != nil {
			d.AddError("unable to find disk", fmt.Sprintf("unable to find disk with id %s: %s", diskState.ID.ValueString(), err))
			return nil, d
		}

		// Check if size or storage profile has changed
		if !diskPlan.SizeInMb.Equal(diskState.SizeInMb) ||
			!diskPlan.StorageProfile.Equal(diskState.StorageProfile) ||
			!diskPlan.VMID.Equal(diskState.VMID) ||
			!diskPlan.VMName.Equal(diskState.VMName) {
			var (
				storageProfileRef *govcdtypes.Reference
				vm                *govcd.VM
				diskIsDetached    bool
				vmDiskDetached    *govcd.VM
			)

			// Get VM object
			vapp, err := vdc.GetVAppById(diskState.VAppID.ValueString(), true)
			if err != nil {
				d.AddError("unable to find vapp", fmt.Sprintf("unable to find vapp with id %s: %s", diskState.VAppID.ValueString(), err))
				return nil, d
			}

			// if VMName or VMID is not define is not necessary to detach the disk
			// Use diskState to get possible OLD VMID
			if diskState.VMName.ValueString() != "" || diskState.VMID.ValueString() != "" {
				if diskState.VMID.ValueString() != "" {
					vm, err = vapp.GetVMById(diskState.VMID.ValueString(), true)
				} else {
					vm, err = vapp.GetVMByName(diskState.VMName.ValueString(), true)
				}
				if err != nil {
					d.AddError("unable to find vm", fmt.Sprintf("unable to find vm with id %s: %s", diskState.VMID.ValueString(), err))
					return nil, d
				}

				// Use diskState to get possible OLD VMID
				d.Append(DiskDetach(ctx, vdc, diskState, vm)...)
				if d.HasError() {
					return nil, d
				}

				diskIsDetached = true
				vmDiskDetached = vm

				err = x.Refresh()
				if err != nil {
					d.AddError("unable to refresh disk", fmt.Sprintf("unable to refresh disk %s (id:%s): %s", diskPlan.Name.ValueString(), diskPlan.ID.ValueString(), err))
					return nil, d
				}
			}

			if !diskPlan.SizeInMb.Equal(diskState.SizeInMb) ||
				!diskPlan.StorageProfile.Equal(diskState.StorageProfile) {
				// If the storage profile is set checking if it exists and setting it
				if !diskPlan.StorageProfile.Equal(diskState.StorageProfile) {
					storageReference, err := vdc.FindStorageProfileReference(diskPlan.StorageProfile.ValueString())
					if err != nil {
						d.AddError("storage profile not found", fmt.Sprintf("The storage profile %s does not exist in the vDC", diskPlan.StorageProfile.ValueString()))
						return nil, d
					}
					storageProfileRef = &govcdtypes.Reference{HREF: storageReference.HREF, Name: storageReference.Name}
				}

				x.Disk.SizeMb = diskPlan.SizeInMb.ValueInt64()
				if storageProfileRef != nil {
					x.Disk.StorageProfile = storageProfileRef
				}

				// Updating the disk
				task, err := x.Update(x.Disk)
				if err != nil {
					d.AddError("unable to update disk", fmt.Sprintf("unable to update disk %s (id:%s): %s", diskPlan.Name.ValueString(), diskPlan.ID.ValueString(), err))
					return nil, d
				}

				err = task.WaitTaskCompletion()
				if err != nil {
					d.AddError("unable to update disk", fmt.Sprintf("unable to update disk %s (id:%s): %s", diskPlan.Name.ValueString(), diskPlan.ID.ValueString(), err))
					return nil, d
				}
			}

			if diskPlan.VMName.ValueString() != "" ||
				diskPlan.VMID.ValueString() != "" {
				var vm *govcd.VM
				if diskPlan.VMID.ValueString() != "" {
					vm, err = vapp.GetVMById(diskPlan.VMID.ValueString(), true)
				} else {
					vm, err = vapp.GetVMByName(diskPlan.VMName.ValueString(), true)
				}
				if err != nil {
					d.AddError("unable to find vm", fmt.Sprintf("unable to find vm with id %s: %s", diskPlan.VMID.ValueString(), err))
					return nil, d
				}

				var busNumber, unitNumber types.Int64

				if diskIsDetached {
					for _, disk := range vmDiskDetached.VM.VmSpecSection.DiskSection.DiskSettings {
						if disk.DiskId == x.Disk.Id {
							busNumber = types.Int64Value(int64(disk.BusNumber))
							unitNumber = types.Int64Value(int64(disk.UnitNumber))
							break
						}
					}
				} else if diskPlan.BusNumber.IsNull() || diskPlan.UnitNumber.IsNull() {
					b, u := diskparams.ComputeBusAndUnitNumber(vm.VM.VmSpecSection.DiskSection.DiskSettings)
					busNumber = types.Int64Value(int64(b))
					unitNumber = types.Int64Value(int64(u))
				}

				d.Append(DiskAttach(ctx, vdc, &Disk{
					ID:         diskPlan.ID,
					BusNumber:  busNumber,
					UnitNumber: unitNumber,
					Name:       diskPlan.Name,
				}, vm)...)
				if d.HasError() {
					return nil, d
				}

				diskPlan.VMID = types.StringValue(vm.VM.ID)
				diskPlan.VMName = types.StringValue(vm.VM.Name)
			}
		}

		err = x.Refresh()
		if err != nil {
			d.AddError("unable to refresh disk", fmt.Sprintf("unable to refresh disk %s (id:%s): %s", diskPlan.Name.ValueString(), diskPlan.ID.ValueString(), err))
			return nil, d
		}
	} else { //nolint:gocritic
		if !diskPlan.SizeInMb.Equal(diskState.SizeInMb) ||
			!diskPlan.StorageProfile.Equal(diskState.StorageProfile) {
			_, diagerr := InternalDiskUpdate(ctx, client, InternalDisk{
				ID:             diskPlan.ID,
				BusType:        diskPlan.BusType,
				SizeInMb:       diskPlan.SizeInMb,
				StorageProfile: diskPlan.StorageProfile,
				BusNumber:      diskPlan.BusNumber,
				UnitNumber:     diskPlan.UnitNumber,
			}, diskPlan.VAppName, diskPlan.VMName, diskPlan.VDC)
			d.Append(diagerr...)
			if d.HasError() {
				return nil, d
			}
		}
	}

	diskPlan.ID = diskState.ID

	return diskPlan, d
}

/*
DiskDelete

delete a disk

if the disk is attached to a VM, it will return an error.
*/
func DiskDelete(ctx context.Context, client *client.CloudAvenue, disk *Disk, org org.Org, vdc vdc.VDC, inVapp vapp.VAPP) diag.Diagnostics {
	d := diag.Diagnostics{}

	if inVapp.VApp == nil || org.Org == nil || vdc.VDC == nil {
		d.AddError("Error read disk", "Empty vApp, org or vdc")
		return d
	}

	// Lock vApp
	d.Append(inVapp.LockVAPP(ctx)...)
	if d.HasError() {
		return d
	}
	defer d.Append(inVapp.UnlockVAPP(ctx)...)

	if disk.IsDetachable.ValueBool() {
		diskRecord, err := vdc.QueryDisk(disk.Name.ValueString())
		if err != nil {
			return d
		}

		if diskRecord.Disk.IsAttached {
			var vmByNameOrID types.String
			if disk.VMName.ValueString() != "" {
				vmByNameOrID = disk.VMName
			} else {
				vmByNameOrID = disk.VMID
			}
			myVM, err := inVapp.GetVMByNameOrId(vmByNameOrID.ValueString(), false)
			if err != nil {
				d.AddError("Error retrieving VM", err.Error())
				return d
			}
			if myVM.VM == nil {
				d.AddError("Error retrieving VM", "VM not found")
				return d
			}

			d.Append(DiskDetach(ctx, vdc, disk, myVM)...)
			if d.HasError() {
				return d
			}
		}

		// Get disk object
		x, err := vdc.GetDiskByHref(diskRecord.Disk.HREF)
		if err != nil {
			d.AddError("error retrieving disk", fmt.Sprintf("error retrieving disk %s: %s", disk.Name.ValueString(), err))
			return d
		}

		// Delete disk
		task, err := x.Delete()
		if err != nil {
			d.AddError("error deleting disk", fmt.Sprintf("error deleting disk %s: %s", disk.Name.ValueString(), err))
			return d
		}

		err = task.WaitTaskCompletion()
		if err != nil {
			d.AddError("error deleting disk", fmt.Sprintf("error deleting disk %s: %s", disk.Name.ValueString(), err))
			return d
		}
	} else {
		var vmByNameOrID types.String
		if disk.VMName.ValueString() != "" {
			vmByNameOrID = disk.VMName
		} else {
			vmByNameOrID = disk.VMID
		}
		myVM, err := inVapp.GetVMByNameOrId(vmByNameOrID.ValueString(), false)
		if err != nil {
			d.AddError("Error retrieving VM", err.Error())
			return d
		}
		if myVM.VM == nil {
			d.AddError("Error retrieving VM", "VM not found")
			return d
		}

		d.Append(InternalDiskDelete(ctx, &InternalDisk{
			ID: disk.ID,
		}, myVM)...)
	}

	return d
}

/*
DiskAttach

attach a disk to a VM.
*/
func DiskAttach(ctx context.Context, vdc vdc.VDC, disk *Disk, vm *govcd.VM) diag.Diagnostics {
	d := diag.Diagnostics{}

	// Get disk object
	x, err := vdc.GetDiskById(disk.ID.ValueString(), true)
	if err != nil {
		d.AddError("error retrieving disk", fmt.Sprintf("error retrieving disk %s (ID:%s): %s", disk.Name.ValueString(), disk.ID.ValueString(), err))
		return d
	}

	attachParams := &govcdtypes.DiskAttachOrDetachParams{
		Disk:       &govcdtypes.Reference{HREF: x.Disk.HREF},
		BusNumber:  utils.TakeIntPointer(int(disk.BusNumber.ValueInt64())),
		UnitNumber: utils.TakeIntPointer(int(disk.UnitNumber.ValueInt64())),
	}

	// Attach disk
	task, err := vm.AttachDisk(attachParams)
	if err != nil {
		d.AddError("error attaching disk", fmt.Sprintf("error attaching disk %s: %s", disk.Name.ValueString(), err))
		return d
	}

	err = task.WaitTaskCompletion()
	if err != nil {
		d.AddError("error attaching disk", fmt.Sprintf("error attaching disk %s: %s", disk.Name.ValueString(), err))
		return d
	}

	return d
}

/*
DiskDetach

detach a disk from a VM.
*/
func DiskDetach(ctx context.Context, vdc vdc.VDC, disk *Disk, vm *govcd.VM) diag.Diagnostics {
	d := diag.Diagnostics{}

	// Get disk object
	x, err := vdc.GetDiskById(disk.ID.ValueString(), true)
	if err != nil {
		d.AddError("error retrieving disk", fmt.Sprintf("error retrieving disk %s: %s", disk.Name.ValueString(), err))
		return d
	}

	detachParams := &govcdtypes.DiskAttachOrDetachParams{
		Disk: &govcdtypes.Reference{HREF: x.Disk.HREF},
	}

	// Detach disk
	task, err := vm.DetachDisk(detachParams)
	if err != nil {
		d.AddError("error detaching disk", fmt.Sprintf("error detaching disk %s: %s", disk.Name.ValueString(), err))
		return d
	}

	err = task.WaitTaskCompletion()
	if err != nil {
		d.AddError("error detaching disk", fmt.Sprintf("error detaching disk %s: %s", disk.Name.ValueString(), err))
		return d
	}

	return d
}
