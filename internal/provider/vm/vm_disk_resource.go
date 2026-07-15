/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vm provides a Terraform resource.
package vm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &diskResource{}
	_ resource.ResourceWithConfigure   = &diskResource{}
	_ resource.ResourceWithImportState = &diskResource{}
	_ resource.ResourceWithModifyPlan  = &diskResource{}
)

// NewDiskResource is a helper function to simplify the provider implementation.
func NewDiskResource() resource.Resource {
	return &diskResource{}
}

// diskResource is the resource implementation.
type diskResource struct {
	client *client.CloudAvenue
	vapp   vapp.VAPP
	vdc    vdc.VDC
	org    org.Org
	vm     vm.VM
}

type diskResourceModel vm.Disk

// Metadata returns the resource type name.
func (r *diskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "disk"
}

// Schema defines the schema for the resource.
func (r *diskResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DiskSuperSchema().GetResource(ctx)
}

func (r *diskResource) Init(_ context.Context, rm *vm.Disk) (diags diag.Diagnostics) {
	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return diags
	}

	r.vdc, diags = vdc.Init(r.client, rm.VDC)
	if diags.HasError() {
		return diags
	}

	r.vapp, diags = vapp.Init(r.client, r.vdc, rm.VAppID, rm.VAppName)
	if diags.HasError() {
		return diags
	}

	if rm.VMName.ValueString() != "" || rm.VMID.ValueString() != "" {
		r.vm, diags = vm.Get(r.vapp, vm.GetVMOpts{
			ID:   rm.VMID,
			Name: rm.VMName,
		})
	}
	return diags
}

func (r *diskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// ModifyPlan is called before Create, Update, and Delete to modify the plan.
func (r *diskResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var d diag.Diagnostics

	diskPlan := &diskResourceModel{}
	diskState := &diskResourceModel{}

	d = req.Plan.Get(ctx, diskPlan)
	if d.HasError() {
		// Plan is not available, so we can't validate the plan.
		return
	}

	// if disk is not detachable, vm_id or vm_name is required
	if diskPlan.IsDetachable.IsNull() || diskPlan.IsDetachable.IsUnknown() || !diskPlan.IsDetachable.ValueBool() {
		if (diskPlan.VMID.IsNull() || diskPlan.VMID.IsUnknown()) && (diskPlan.VMName.IsNull() || diskPlan.VMName.IsUnknown()) {
			resp.Diagnostics.AddError(
				"VM is required",
				"if \"is_detachable\" attribute is false \"vm_id\" or \"vm_name\" is required to attach disk to a VM",
			)
			return
		}
	}

	d = req.State.Get(ctx, diskState)
	if d.HasError() {
		// State is not available, so we can't validate the plan.
		return
	}

	if diskPlan.IsDetachable.ValueBool() {
		if !diskPlan.SizeInMb.Equal(diskState.SizeInMb) {
			resp.Diagnostics.AddWarning(
				"Warning detach/attach disk is required",
				"Disk size cannot be changed when disk is detachable. Detach/attach disk is required. \n"+
					"If you apply this change, the disk will be detached and attached again.",
			)
		}
	} else {
		if diskPlan.BusType.ValueString() == diskparams.BusTypeIDE.Name() {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("bus_type"),
				"Warning IDE bus type require VM power off to be applied",
				"IDE bus type require VM power off to be applied. \n"+
					"If you apply this change, power off before apply and power on after apply will be required.",
			)
		}

		// if disk is not detachable, vm_id or vm_name is not editable
		// This setting is not in schema definition because if is_detachable is true, vm_id and vm_name is editable
		if !diskPlan.VMID.Equal(diskState.VMID) {
			resp.RequiresReplace.Append(path.Root("vm_id"))
		}
		if !diskPlan.VMName.Equal(diskState.VMName) {
			resp.RequiresReplace.Append(path.Root("vm_name"))
		}
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *diskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { //nolint:gocyclo
	defer metrics.New("cloudavenue_vm_disk", r.client.GetOrgName(), metrics.Create)()

	plan := &vm.Disk{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newPlan := *plan
	newPlan.VDC = types.StringValue(r.vdc.GetName())

	if plan.IsDetachable.ValueBool() {
		resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
		if resp.Diagnostics.HasError() {
			return
		}

		defer r.vapp.UnlockVAPP(ctx)

		if ok, err := r.vdc.DiskExist(plan.Name.ValueString()); ok {
			resp.Diagnostics.AddError("Disk already exists", "Disk with name "+plan.Name.ValueString()+" already exists")
			return
		} else if err != nil {
			resp.Diagnostics.AddError("Error checking disk", "Error checking if disk with name "+plan.Name.ValueString()+" already exists. Error : "+err.Error())
			return
		}

		// Init struct for creating a disk
		diskCreateParams := &govcdtypes.DiskCreateParams{
			Disk: &govcdtypes.Disk{
				Name:        plan.Name.ValueString(),
				SizeMb:      plan.SizeInMb.ValueInt64(),
				SharingType: "None",
			},
		}

		// If the bus type is set checking if it exists and setting it
		if !plan.BusType.IsNull() && !plan.BusType.IsUnknown() {
			diskCreateParams.Disk.BusType = diskparams.GetBusTypeByName(plan.BusType.ValueString()).Code()
			diskCreateParams.Disk.BusSubType = diskparams.GetBusTypeByName(plan.BusType.ValueString()).SubType()
		}

		// If the storage profile is set checking if it exists and setting it
		if !plan.StorageProfile.IsNull() && !plan.StorageProfile.IsUnknown() {
			storageReference, err := r.vdc.FindStorageProfileReference(plan.StorageProfile.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("storage profile not found", fmt.Sprintf("The storage profile %s does not exist in the vDC", plan.StorageProfile.ValueString()))
				return
			}
			diskCreateParams.Disk.StorageProfile = &govcdtypes.Reference{HREF: storageReference.HREF}
		}

		// Create the disk
		task, err := r.vdc.CreateDisk(diskCreateParams)
		if err != nil {
			resp.Diagnostics.AddError("error creating disk", err.Error())
			return
		}

		// Wait for the task to finish
		if err = task.WaitTaskCompletion(); err != nil {
			resp.Diagnostics.AddError("error on creating disk", err.Error())
			return
		}

		disk, err := r.vdc.GetDiskByHref(task.Task.Owner.HREF)
		if err != nil {
			resp.Diagnostics.AddError("error on getting disk", err.Error())
			return
		}

		newPlan.ID = types.StringValue(disk.Disk.Id)
		newPlan.StorageProfile = types.StringValue(disk.Disk.StorageProfile.Name)

		if r.vm != (vm.VM{}) {
			resp.Diagnostics.Append(r.vm.LockVM(ctx)...)
			if resp.Diagnostics.HasError() {
				return
			}

			defer r.vm.UnlockVM(ctx)

			// loop with Timeout
			diskRefreshTimeout := 20 * time.Second
			diskRefreshTicker := time.NewTicker(2 * time.Second)
			defer diskRefreshTicker.Stop()
			timeout := time.After(diskRefreshTimeout)
			refreshEnded := false

			for !refreshEnded {
				select {
				case <-diskRefreshTicker.C:
					if err := r.vm.Refresh(); err != nil {
						resp.Diagnostics.AddError("error refreshing vm", fmt.Errorf("error refreshing vm: %w", err).Error())
						return
					}
					if r.vm.VM == nil || r.vm.VM.VM == nil || r.vm.VM.VM.VM == nil || r.vm.VM.VM.VM.Link == nil {
						continue
					}
					if len(r.vm.VM.VM.VM.Link) > 0 {
						refreshEnded = true
					}
				case <-timeout:
					resp.Diagnostics.AddError("error refreshing disk", "timeout refreshing disk")
					return
				}
			}

			if err := r.vm.Refresh(); err != nil {
				resp.Diagnostics.AddError("error refreshing vm", err.Error())
				return
			}

			// Attach disk
			task, err = r.vm.AttachDisk(r.vm.AttachDiskSettings(plan.BusNumber, plan.UnitNumber, task.Task.Owner.HREF))
			if err != nil {
				resp.Diagnostics.AddError("error attaching disk", fmt.Errorf("error attaching disk %s: %w", plan.Name.ValueString(), err).Error())
				return
			}

			if err = task.WaitTaskCompletion(); err != nil {
				resp.Diagnostics.AddError("error attaching disk", fmt.Errorf("error attaching disk %s: %w", plan.Name.ValueString(), err).Error())
				return
			}

			if err = disk.Refresh(); err != nil {
				resp.Diagnostics.AddError("error refreshing disk", fmt.Errorf("error refreshing disk %s: %w", plan.Name.ValueString(), err).Error())
				return
			}
			newPlan.BusType = types.StringValue(diskparams.GetBusTypeByCode(disk.Disk.BusType, disk.Disk.BusSubType).Name())

			if err := r.vm.Refresh(); err != nil {
				resp.Diagnostics.AddError("error refreshing vm", fmt.Errorf("error refreshing vm: %w", err).Error())
				return
			}
			var diskSettings []*govcdtypes.DiskSettings
			if r.vm.VM != nil && r.vm.VM.VM != nil && r.vm.VM.VM.VM != nil && r.vm.VM.VM.VM.VmSpecSection != nil && r.vm.VM.VM.VM.VmSpecSection.DiskSection != nil {
				diskSettings = r.vm.VM.VM.VM.VmSpecSection.DiskSection.DiskSettings
			}

			found := false
			for _, diskSetting := range diskSettings {
				if diskSetting.DiskId == disk.Disk.Id {
					newPlan.BusNumber = types.Int64Value(int64(diskSetting.BusNumber))
					newPlan.UnitNumber = types.Int64Value(int64(diskSetting.UnitNumber))
					found = true
					break
				}
			}
			if !found {
				resp.Diagnostics.AddError("error reading disk settings", fmt.Sprintf("unable to find disk %s in VM disk settings", disk.Disk.Id))
				return
			}
		} else {
			newPlan.BusNumber = types.Int64Null()
			newPlan.UnitNumber = types.Int64Null()
		} // End if r.vm != (vm.VM{})
	} else { // Disk not detachable it's an internal disk
		if r.vm.VM == nil || r.vm.VM.VM == nil || r.vm.VM.VM.VM == nil {
			resp.Diagnostics.AddError("VM not found", "VM is not available to create internal disk")
			return
		}

		// storage profile
		var (
			storageProfilePrt *govcdtypes.Reference
			overrideVMDefault bool
		)

		if plan.StorageProfile.IsNull() || plan.StorageProfile.IsUnknown() {
			storageProfilePrt = r.vm.VM.VM.VM.StorageProfile
			overrideVMDefault = false
		} else {
			storageProfile, errFindStorage := r.vdc.FindStorageProfileReference(plan.StorageProfile.ValueString())
			if errFindStorage != nil {
				resp.Diagnostics.AddError("Error retrieving storage profile", errFindStorage.Error())
				return
			}
			storageProfilePrt = &storageProfile
			overrideVMDefault = true
		}

		// value is required but not treated.
		isThinProvisioned := true

		var busNumber, unitNumber types.Int64

		var computedBus, computedUnit int
		if plan.BusNumber.IsNull() || plan.BusNumber.IsUnknown() || plan.UnitNumber.IsNull() || plan.UnitNumber.IsUnknown() {
			var diskSettings []*govcdtypes.DiskSettings
			if r.vm.VM != nil && r.vm.VM.VM != nil && r.vm.VM.VM.VM != nil && r.vm.VM.VM.VM.VmSpecSection != nil && r.vm.VM.VM.VM.VmSpecSection.DiskSection != nil {
				diskSettings = r.vm.VM.VM.VM.VmSpecSection.DiskSection.DiskSettings
			}
			computedBus, computedUnit = diskparams.ComputeBusAndUnitNumber(diskSettings)
		}

		if plan.BusNumber.IsNull() || plan.BusNumber.IsUnknown() {
			busNumber = types.Int64Value(int64(computedBus))
		} else {
			busNumber = plan.BusNumber
		}

		if plan.UnitNumber.IsNull() || plan.UnitNumber.IsUnknown() {
			unitNumber = types.Int64Value(int64(computedUnit))
		} else {
			unitNumber = plan.UnitNumber
		}

		diskSetting := &govcdtypes.DiskSettings{
			SizeMb:              plan.SizeInMb.ValueInt64(),
			UnitNumber:          int(unitNumber.ValueInt64()),
			BusNumber:           int(busNumber.ValueInt64()),
			AdapterType:         vm.GetBusTypeByKey(plan.BusType.ValueString()).Code(),
			ThinProvisioned:     &isThinProvisioned,
			StorageProfile:      storageProfilePrt,
			VirtualQuantityUnit: "byte",
			OverrideVmDefault:   overrideVMDefault,
		}

		resp.Diagnostics.Append(r.vm.LockVM(ctx)...)
		if resp.Diagnostics.HasError() {
			return
		}

		defer r.vm.UnlockVM(ctx)

		diskID, err := r.vm.AddInternalDisk(diskSetting)
		if err != nil {
			resp.Diagnostics.AddError("Error creating disk", err.Error())
			return
		}

		newPlan.ID = types.StringValue(diskID)
		newPlan.BusType = types.StringValue(strings.ToUpper(vm.GetBusTypeByCode(diskSetting.AdapterType).Name()))
		newPlan.SizeInMb = types.Int64Value(diskSetting.SizeMb)
		newPlan.StorageProfile = types.StringValue(storageProfilePrt.Name)
		newPlan.BusNumber = types.Int64Value(int64(diskSetting.BusNumber))
		newPlan.UnitNumber = types.Int64Value(int64(diskSetting.UnitNumber))
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, newPlan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *diskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vm_disk", r.client.GetOrgName(), metrics.Read)()

	state := &vm.Disk{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedState := *state

	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defer r.vapp.UnlockVAPP(ctx)

	// * Detachable disk
	if state.IsDetachable.ValueBool() {
		// Get the disk by the ID
		x, err := r.vdc.GetDiskById(state.ID.ValueString(), true)
		if err != nil {
			if govcd.IsNotFound(err) {
				// Disk not found, remove from state
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("unable to find disk", fmt.Sprintf("unable to find disk with id %s: %v", state.ID.ValueString(), err))
			return
		}

		attachedVmsHrefs, err := x.GetAttachedVmsHrefs()
		if err != nil {
			resp.Diagnostics.AddError("unable to find attached VM", fmt.Errorf("unable to find attached VM for disk %s: %w", state.Name.ValueString(), err).Error())
			return
		}

		// Normally a disk can be attached to only one VM
		if len(attachedVmsHrefs) > 1 {
			resp.Diagnostics.AddError("multiple VMs attached", fmt.Sprintf("multiple VMs attached to disk %s.", state.Name.ValueString()))
			return
		}

		updatedState.ID = types.StringValue(x.Disk.Id)
		updatedState.Name = types.StringValue(x.Disk.Name)
		updatedState.SizeInMb = types.Int64Value(x.Disk.SizeMb)
		updatedState.BusType = types.StringValue(strings.ToUpper(diskparams.GetBusTypeByCode(x.Disk.BusType, x.Disk.BusSubType).Name()))
		updatedState.StorageProfile = types.StringValue(x.Disk.StorageProfile.Name)

		// Normally a disk can be attached to only one VM
		if len(attachedVmsHrefs) == 1 {
			govcdVM, err := r.client.Vmware.Client.GetVMByHref(attachedVmsHrefs[0])
			if err != nil {
				resp.Diagnostics.AddError("unable to find attached VM", fmt.Errorf("unable to find attached VM for disk %s: %w", state.Name.ValueString(), err).Error())
				return
			}

			var found bool
			if govcdVM != nil && govcdVM.VM != nil && govcdVM.VM.VmSpecSection != nil && govcdVM.VM.VmSpecSection.DiskSection != nil {
				for _, diskSetting := range govcdVM.VM.VmSpecSection.DiskSection.DiskSettings {
					if diskSetting.DiskId == state.ID.ValueString() {
						updatedState.BusNumber = types.Int64Value(int64(diskSetting.BusNumber))
						updatedState.UnitNumber = types.Int64Value(int64(diskSetting.UnitNumber))
						found = true
						break
					}
				}
			}
			if !found {
				updatedState.BusNumber = types.Int64Null()
				updatedState.UnitNumber = types.Int64Null()
				resp.Diagnostics.AddWarning("disk bus information not found", fmt.Sprintf("could not find bus_number/unit_number for disk %s in VM disk settings", state.ID.ValueString()))
			}
		} else if len(attachedVmsHrefs) == 0 {
			updatedState.BusNumber = types.Int64Null()
			updatedState.UnitNumber = types.Int64Null()
		}
	} else {
		// * Internal disk
		internalDisk, err := r.vm.GetInternalDiskById(state.ID.ValueString(), true)
		if err != nil {
			if govcd.IsNotFound(err) {
				// Disk not found, remove from state
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("unable to find disk", fmt.Sprintf("unable to find disk with id %s: %v", state.ID.ValueString(), err))
			return
		}

		updatedState.ID = types.StringValue(internalDisk.DiskId)
		updatedState.Name = types.StringNull()
		updatedState.SizeInMb = types.Int64Value(internalDisk.SizeMb)
		updatedState.StorageProfile = types.StringValue(internalDisk.StorageProfile.Name)
		updatedState.BusType = types.StringValue(strings.ToUpper(vm.GetBusTypeByCode(internalDisk.AdapterType).Name()))
		updatedState.BusNumber = types.Int64Value(int64(internalDisk.BusNumber))
		updatedState.UnitNumber = types.Int64Value(int64(internalDisk.UnitNumber))
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *diskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { //nolint:gocyclo
	defer metrics.New("cloudavenue_vm_disk", r.client.GetOrgName(), metrics.Update)()

	plan := &vm.Disk{}
	state := &vm.Disk{}

	// Get current plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defer r.vapp.UnlockVAPP(ctx)

	updatedState := *plan

	if state.IsDetachable.ValueBool() {
		// Get the disk by the ID
		disk, err := r.vdc.GetDiskById(state.ID.ValueString(), true)
		if err != nil {
			resp.Diagnostics.AddError("unable to find disk", fmt.Sprintf("unable to find disk with id %s: %v", state.ID.ValueString(), err))
			return
		}

		// Check if size, storage profile, vm id, vm name, bus number or unit number has changed
		if !plan.SizeInMb.Equal(state.SizeInMb) ||
			!plan.StorageProfile.Equal(state.StorageProfile) ||
			!plan.VMID.Equal(state.VMID) ||
			!plan.VMName.Equal(state.VMName) ||
			(!plan.BusNumber.IsUnknown() && !plan.BusNumber.Equal(state.BusNumber)) ||
			(!plan.UnitNumber.IsUnknown() && !plan.UnitNumber.Equal(state.UnitNumber)) {
			// Check if disk is attached to a VM
			if !(state.VMID.IsNull() && state.VMName.IsNull()) {
				// Detach disk from VM

				vmOld, diag := vm.Get(r.vapp, vm.GetVMOpts{
					ID:   state.VMID,
					Name: state.VMName,
				})
				resp.Diagnostics.Append(diag...)
				if resp.Diagnostics.HasError() {
					return
				}

				resp.Diagnostics.Append(vmOld.LockVM(ctx)...)
				if resp.Diagnostics.HasError() {
					return
				}

				// Detach disk
				task, err := vmOld.DetachDisk(&govcdtypes.DiskAttachOrDetachParams{
					Disk: &govcdtypes.Reference{HREF: disk.Disk.HREF},
				})
				if err != nil {
					vmOld.UnlockVM(ctx)
					resp.Diagnostics.AddError("error detaching disk", fmt.Errorf("error detaching disk %s: %w", state.Name.ValueString(), err).Error())
					return
				}
				if err = task.WaitTaskCompletion(); err != nil {
					vmOld.UnlockVM(ctx)
					resp.Diagnostics.AddError("error detaching disk", fmt.Errorf("error detaching disk %s: %w", state.Name.ValueString(), err).Error())
					return
				}
				vmOld.UnlockVM(ctx)
			}

			if err = disk.Refresh(); err != nil {
				resp.Diagnostics.AddError("unable to refresh disk", fmt.Sprintf("unable to refresh disk %s (id:%s): %s", state.Name.ValueString(), state.ID.ValueString(), err))
				return
			}

			if !plan.SizeInMb.Equal(state.SizeInMb) ||
				!plan.StorageProfile.Equal(state.StorageProfile) {
				// If the storage profile is set checking if it exists and setting it
				if !plan.StorageProfile.Equal(state.StorageProfile) {
					storageReference, err := r.vdc.FindStorageProfileReference(plan.StorageProfile.ValueString())
					if err != nil {
						resp.Diagnostics.AddError("storage profile not found", fmt.Sprintf("The storage profile %s does not exist in the vDC", plan.StorageProfile.ValueString()))
						return
					}
					disk.Disk.StorageProfile = &govcdtypes.Reference{HREF: storageReference.HREF, Name: storageReference.Name}
				}

				disk.Disk.SizeMb = plan.SizeInMb.ValueInt64()

				// Updating the disk
				task, err := disk.Update(disk.Disk)
				if err != nil {
					resp.Diagnostics.AddError("unable to update disk", fmt.Sprintf("unable to update disk %s (id:%s): %s", plan.Name.ValueString(), plan.ID.ValueString(), err))
					return
				}

				if err = task.WaitTaskCompletion(); err != nil {
					resp.Diagnostics.AddError("unable to update disk", fmt.Sprintf("unable to update disk %s (id:%s): %s", plan.Name.ValueString(), plan.ID.ValueString(), err))
					return
				}
			}

			if plan.VMName.ValueString() != "" ||
				plan.VMID.ValueString() != "" {
				var (
					vmNew vm.VM
					d     diag.Diagnostics
				)
				vmNew, d = vm.Get(r.vapp, vm.GetVMOpts{
					ID:   plan.VMID,
					Name: plan.VMName,
				})
				if d.HasError() {
					resp.Diagnostics.Append(d...)
					return
				}

				var busNumber, unitNumber types.Int64

				resp.Diagnostics.Append(vmNew.LockVM(ctx)...)
				if resp.Diagnostics.HasError() {
					return
				}

				defer vmNew.UnlockVM(ctx)

				if err := vmNew.Refresh(); err != nil {
					resp.Diagnostics.AddError("error refreshing vm", err.Error())
					return
				}

				var computedBus, computedUnit int
				if plan.BusNumber.IsNull() || plan.BusNumber.IsUnknown() || plan.UnitNumber.IsNull() || plan.UnitNumber.IsUnknown() {
					var diskSettings []*govcdtypes.DiskSettings
					if vmNew.VM != nil && vmNew.VM.VM != nil && vmNew.VM.VM.VM != nil && vmNew.VM.VM.VM.VmSpecSection != nil && vmNew.VM.VM.VM.VmSpecSection.DiskSection != nil {
						diskSettings = vmNew.VM.VM.VM.VmSpecSection.DiskSection.DiskSettings
					}
					computedBus, computedUnit = diskparams.ComputeBusAndUnitNumber(diskSettings)
				}

				if plan.BusNumber.IsNull() || plan.BusNumber.IsUnknown() {
					busNumber = types.Int64Value(int64(computedBus))
				} else {
					busNumber = plan.BusNumber
				}

				if plan.UnitNumber.IsNull() || plan.UnitNumber.IsUnknown() {
					unitNumber = types.Int64Value(int64(computedUnit))
				} else {
					unitNumber = plan.UnitNumber
				}

				// Attach disk
				task, err := vmNew.AttachDisk(&govcdtypes.DiskAttachOrDetachParams{
					Disk:       &govcdtypes.Reference{HREF: disk.Disk.HREF},
					BusNumber:  utils.TakeIntPointer(int(busNumber.ValueInt64())),
					UnitNumber: utils.TakeIntPointer(int(unitNumber.ValueInt64())),
				})
				if err != nil {
					resp.Diagnostics.AddError("error attaching disk", fmt.Errorf("error attaching disk %s: %w", plan.Name.ValueString(), err).Error())
					return
				}

				if err = task.WaitTaskCompletion(); err != nil {
					resp.Diagnostics.AddError("error attaching disk", fmt.Errorf("error attaching disk %s: %w", plan.Name.ValueString(), err).Error())
					return
				}

				// Read actual bus_number and unit_number from the VM
				if err := vmNew.Refresh(); err != nil {
					resp.Diagnostics.AddError("error refreshing vm", fmt.Errorf("error refreshing vm: %w", err).Error())
					return
				}
				var diskSettings []*govcdtypes.DiskSettings
				if vmNew.VM != nil && vmNew.VM.VM != nil && vmNew.VM.VM.VM != nil && vmNew.VM.VM.VM.VmSpecSection != nil && vmNew.VM.VM.VM.VmSpecSection.DiskSection != nil {
					diskSettings = vmNew.VM.VM.VM.VmSpecSection.DiskSection.DiskSettings
				}

				found := false
				for _, x := range diskSettings {
					if x.DiskId == disk.Disk.Id {
						updatedState.BusNumber = types.Int64Value(int64(x.BusNumber))
						updatedState.UnitNumber = types.Int64Value(int64(x.UnitNumber))
						found = true
						break
					}
				}
				if !found {
					resp.Diagnostics.AddError("error reading disk settings", fmt.Sprintf("unable to find disk %s in VM disk settings", disk.Disk.Id))
					return
				}
			} else {
				updatedState.BusNumber = types.Int64Null()
				updatedState.UnitNumber = types.Int64Null()
			}
		}

		// If the detachable disk is not attached to any VM and is not being
		// attached to one, bus_number and unit_number must be null. Otherwise a
		// change to bus/unit in config while the disk is detached would leave
		// stale non-null plan values that a subsequent Read resets to null,
		// causing a permadiff.
		if state.VMID.IsNull() && state.VMName.IsNull() && plan.VMID.IsNull() && plan.VMName.IsNull() {
			updatedState.BusNumber = types.Int64Null()
			updatedState.UnitNumber = types.Int64Null()
		}
	} else {
		// * Internal disk
		internalDisk, err := r.vm.GetInternalDiskById(state.ID.ValueString(), true)
		if err != nil {
			resp.Diagnostics.AddError("unable to find disk", fmt.Sprintf("unable to find disk with id %s: %v", state.ID.ValueString(), err))
			return
		}

		resp.Diagnostics.Append(r.vm.LockVM(ctx)...)
		if resp.Diagnostics.HasError() {
			return
		}

		defer r.vm.UnlockVM(ctx)

		if r.vm.VM == nil || r.vm.VM.VM == nil || r.vm.VM.VM.VM == nil {
			resp.Diagnostics.AddError("VM not found", "VM is not available to update internal disk")
			return
		}

		internalDisk.SizeMb = plan.SizeInMb.ValueInt64()

		var (
			storageProfilePrt *govcdtypes.Reference
			overrideVMDefault bool
		)

		if plan.StorageProfile.IsNull() || plan.StorageProfile.IsUnknown() {
			storageProfilePrt = r.vm.VM.VM.VM.StorageProfile
			overrideVMDefault = false
		} else {
			storageProfile, errFindStorage := r.vdc.FindStorageProfileReference(plan.StorageProfile.ValueString())
			if errFindStorage != nil {
				resp.Diagnostics.AddError("Error retrieving storage profile", errFindStorage.Error())
				return
			}
			storageProfilePrt = &storageProfile
			overrideVMDefault = true
		}

		internalDisk.StorageProfile = storageProfilePrt
		internalDisk.OverrideVmDefault = overrideVMDefault

		if _, err := r.vm.UpdateInternalDisks(r.vm.VM.VM.VM.VmSpecSection); err != nil {
			resp.Diagnostics.AddError("error updating internal disk", err.Error())
			return
		}

		internalDisk, err = r.vm.GetInternalDiskById(state.ID.ValueString(), true)
		if err != nil {
			resp.Diagnostics.AddError("unable to find disk", fmt.Errorf("unable to find disk with id %s: %w", state.ID.ValueString(), err).Error())
			return
		}

		updatedState.BusNumber = types.Int64Value(int64(internalDisk.BusNumber))
		updatedState.UnitNumber = types.Int64Value(int64(internalDisk.UnitNumber))
	}
	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *diskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vm_disk", r.client.GetOrgName(), metrics.Delete)()

	state := &vm.Disk{}
	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defer r.vapp.UnlockVAPP(ctx)

	if state.IsDetachable.ValueBool() {
		// Get the disk by the ID
		x, err := r.vdc.GetDiskById(state.ID.ValueString(), true)
		if err != nil {
			if govcd.IsNotFound(err) {
				// Disk not found, remove from state
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("unable to find disk", fmt.Sprintf("unable to find disk with id %s: %v", state.ID.ValueString(), err))
			return
		}

		attached, err := x.AttachedVM()
		if err != nil {
			resp.Diagnostics.AddError("unable to find disk", fmt.Sprintf("unable to find disk with id %s: %v", state.ID.ValueString(), err))
			return
		}

		if attached != nil {
			resp.Diagnostics.Append(r.vm.LockVM(ctx)...)
			if resp.Diagnostics.HasError() {
				return
			}

			defer r.vm.UnlockVM(ctx)

			task, err := r.vm.DetachDisk(&govcdtypes.DiskAttachOrDetachParams{
				Disk: &govcdtypes.Reference{
					HREF: x.Disk.HREF,
				},
			})
			if err != nil {
				resp.Diagnostics.AddError("error detaching disk", fmt.Errorf("error detaching disk %s: %w", state.Name.ValueString(), err).Error())
				return
			}

			if err = task.WaitTaskCompletion(); err != nil {
				resp.Diagnostics.AddError("error detaching disk", fmt.Errorf("error detaching disk %s: %w", state.Name.ValueString(), err).Error())
				return
			}
		}

		if err := x.Refresh(); err != nil {
			resp.Diagnostics.AddError("error refreshing disk", fmt.Sprintf("error refreshing disk %s: %v", state.Name.ValueString(), err))
			return
		}

		// Delete disk
		task, err := x.Delete()
		if err != nil {
			resp.Diagnostics.AddError("error deleting disk", fmt.Sprintf("error deleting disk %s: %v", state.Name.ValueString(), err))
			return
		}

		if err = task.WaitTaskCompletion(); err != nil {
			resp.Diagnostics.AddError("error deleting disk", fmt.Sprintf("error deleting disk %s: %v", state.Name.ValueString(), err))
			return
		}
	} else {
		// Delete disk
		if err := r.vm.DeleteInternalDisk(state.ID.ValueString()); err != nil {
			resp.Diagnostics.AddError("error deleting disk", fmt.Sprintf("error deleting disk %s: %v", state.Name.ValueString(), err))
			return
		}
	}
}

func (r *diskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vm_disk", r.client.GetOrgName(), metrics.Import)()

	idParts := strings.Split(req.ID, ".")

	var (
		vdcName          string
		vAppID, vAppName string
		vmID, vmName     string
		diskID           string
		isDetachable     bool

		diags diag.Diagnostics
	)

	switch len(idParts) {
	// Case 2 : vAppIDOrName.DiskID
	case 2:
		if urn.IsVAPP(idParts[0]) {
			vAppID = idParts[0]
		} else {
			vAppName = idParts[0]
		}

		diskID = idParts[1]
		if urn.IsDisk(idParts[1]) {
			isDetachable = true
		}

	// Case 3 : vAppIDOrName.VmIDOrName.DiskID or vdcName.vAppIDOrName.DiskID
	case 3:
		diskID = idParts[2]
		if urn.IsDisk(diskID) {
			isDetachable = true
		}

		if isDetachable {
			_, diags = vdc.Init(r.client, types.StringValue(idParts[0]))
			if !diags.HasError() {
				// FORMAT : vdcName.vAppIDOrName.DiskID
				vdcName = idParts[0]
				if urn.IsVAPP(idParts[1]) {
					vAppID = idParts[1]
				} else {
					vAppName = idParts[1]
				}
				goto next
			}
		}

		// FORMAT : vAppIDOrName.VmIDOrName.DiskID
		if urn.IsVAPP(idParts[0]) {
			vAppID = idParts[0]
		} else {
			vAppName = idParts[0]
		}

		if urn.IsVM(idParts[1]) {
			vmID = idParts[1]
		} else {
			vmName = idParts[1]
		}

	// Case 4 : vdcName.vAppIDOrName.VmIDOrName.DiskID
	case 4:
		vdcName = idParts[0]
		if urn.IsVAPP(idParts[1]) {
			vAppID = idParts[1]
		} else {
			vAppName = idParts[1]
		}

		if urn.IsVM(idParts[2]) {
			vmID = idParts[2]
		} else {
			vmName = idParts[2]
		}

		diskID = idParts[3]
		if urn.IsDisk(diskID) {
			isDetachable = true
		}
	default:
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with formats: vAppIDOrName.DiskID or vAppIDOrName.VmIDOrName.DiskID or vdcName.vAppIDOrName.DiskID or vdcName.vAppIDOrName.VmIDOrName.DiskID Got: %q", req.ID),
		)
		return
	}

next:

	r.org, diags = org.Init(r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.vdc, diags = vdc.Init(r.client, utils.StringValueOrNull(vdcName))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.vapp, diags = vapp.Init(r.client, r.vdc, utils.StringValueOrNull(vAppID), utils.StringValueOrNull(vAppName))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if vmID != "" || vmName != "" {
		r.vm, diags = vm.Get(r.vapp, vm.GetVMOpts{
			ID:   utils.StringValueOrNull(vmID),
			Name: utils.StringValueOrNull(vmName),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), diskID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vdc"), r.vdc.GetName())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("is_detachable"), isDetachable)...)

	if vAppID != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_id"), r.vapp.GetID())...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), r.vapp.GetName())...)
	}

	if vmID != "" || vmName != "" {
		if vmID != "" {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vm_id"), r.vm.GetID())...)
		} else {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vm_name"), r.vm.GetName())...)
		}
	}
}
