package vm

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	commonvm "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

func updateVM(ctx context.Context, v *Client) (*govcd.VM, error) { //nolint:gocyclo
	var (
		vapp *govcd.VApp
		vm   *govcd.VM
		err  error
	)

	// If VDC is not defined at resource/data source level, use the one defined at provider level
	if v.Plan.VDC.IsNull() || v.Plan.VDC.IsUnknown() {
		if v.Client.DefaultVDCExist() {
			v.Plan.VDC = types.StringValue(v.Client.GetDefaultVDC())
		} else {
			return nil, fmt.Errorf("VDC is required when not defined at provider level")
		}
	}

	// Get vcd object
	vdc, err := v.Client.GetVDC(client.WithVDCName(v.State.VDC.ValueString()))
	if err != nil {
		return nil, fmt.Errorf("error retrieving VDC %s: %w", v.Plan.VDC.ValueString(), err)
	}

	// Get vApp
	vapp, err = vdc.GetVAppByName(v.Plan.VappName.ValueString(), false)
	if err != nil {
		return nil, fmt.Errorf("error retrieving vApp %s: %w", v.Plan.VappName.ValueString(), err)
	}

	// Get VM
	vm, err = vapp.GetVMByNameOrId(v.Plan.VMName.ValueString(), false)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VM %s: %w", v.Plan.VMName.ValueString(), err)
	}

	customizationNeeded := isForcedCustomization(v)

	// * End of Init *

	// ! Hot Update (VM must be powered on or off)
	var (
		cpuNeedsColdChange      = false
		memoryNeedsColdChange   = false
		networksNeedsColdChange = false
		resource, resourceState *commonvm.Resource

		d diag.Diagnostics
	)

	resource, d = commonvm.ResourceFromPlan(ctx, v.Plan.Resource)
	if d.HasError() {
		return nil, fmt.Errorf("error retrieving resource from plan: %s", d)
	}

	resourceState, d = commonvm.ResourceFromPlan(ctx, v.State.Resource)
	if d.HasError() {
		return nil, fmt.Errorf("error retrieving resource from state: %s", d)
	}

	// Update Resource
	if !v.Plan.Resource.Equal(v.State.Resource) {
		// * Update Memory
		if resource.MemoryHotAddEnabled.ValueBool() && !resource.Memory.Equal(resourceState.Memory) {
			if err := vm.ChangeMemory(resource.Memory.ValueInt64()); err != nil {
				return nil, fmt.Errorf("error changing memory: %w", err)
			}
		} else if !resource.MemoryHotAddEnabled.ValueBool() && !resource.Memory.Equal(resourceState.Memory) {
			memoryNeedsColdChange = true
		}

		// * Update CPU
		if resource.CPUHotAddEnabled.ValueBool() && resource.CPUs.Equal(resourceState.CPUs) {
			if err := vm.ChangeCPUAndCoreCount(utils.TakeIntPointer(int(resource.CPUs.ValueInt64())), utils.TakeIntPointer(int(resource.CPUCores.ValueInt64()))); err != nil {
				return nil, fmt.Errorf("error changing CPU: %w", err)
			}
		} else if !resource.CPUHotAddEnabled.ValueBool() && !resource.CPUs.Equal(resourceState.CPUs) {
			cpuNeedsColdChange = true
		}
	}
	// End of Update Resource

	// Update Networks
	if v.Plan.Networks.Equal(v.State.Networks) {
		var (
			requireUpdate = false
			isPrimaryNic  = false
		)

		var networks, networksState commonvm.Networks
		diag := v.Plan.Networks.ElementsAs(context.Background(), &networks, false)
		if diag.HasError() {
			return nil, fmt.Errorf("error retrieving networks: %s", diag)
		}

		diag = v.State.Networks.ElementsAs(context.Background(), &networksState, false)
		if diag.HasError() {
			return nil, fmt.Errorf("error retrieving networks: %s", diag)
		}

		for i, network := range networks {
			if network.IsPrimary.ValueBool() {
				isPrimaryNic = true
			}

			if !network.AdapterType.Equal(networksState[i].AdapterType) {
				requireUpdate = true
				break
			}

			if !network.Type.Equal(networksState[i].Type) {
				requireUpdate = true
				break
			}

			if !network.IPAllocationMode.Equal(networksState[i].IPAllocationMode) {
				requireUpdate = true
				break
			}

			if !network.Name.Equal(networksState[i].Name) {
				requireUpdate = true
				break
			}

			if !network.IsPrimary.Equal(networksState[i].IsPrimary) {
				requireUpdate = true
				break
			}

			if !network.IP.Equal(networksState[i].IP) {
				requireUpdate = true
				break
			}

			if !network.Mac.Equal(networksState[i].Mac) {
				requireUpdate = true
				break
			}

			if !network.Connected.Equal(networksState[i].Connected) {
				requireUpdate = true
				break
			}
		}

		// * Primary NIC cannot be removed on a powered on VM
		if requireUpdate && !isPrimaryNic {
			networkConnectionSection, err := networksToConfig(v, vapp)
			if err != nil {
				return nil, fmt.Errorf("unable to setup network configuration for update: %w", err)
			}
			err = vm.UpdateNetworkConnectionSection(&networkConnectionSection)
			if err != nil {
				return nil, fmt.Errorf("unable to update network configuration: %w", err)
			}
		} else if requireUpdate && isPrimaryNic {
			networksNeedsColdChange = true
		}

		err = addRemoveGuestProperties(v, vm)
		if err != nil {
			return nil, fmt.Errorf("unable to update guest properties: %w", err)
		}

		if !v.Plan.SizingPolicyID.Equal(v.State.SizingPolicyID) || !v.Plan.PlacementPolicyID.Equal(v.State.PlacementPolicyID) {
			_, err = vm.UpdateComputePolicyV2(v.Plan.SizingPolicyID.ValueString(), v.Plan.PlacementPolicyID.ValueString(), "")
			if err != nil {
				return nil, fmt.Errorf("unable to update compute policy: %w", err)
			}
		}

		if !v.Plan.StorageProfile.Equal(v.State.StorageProfile) {
			if v.Plan.StorageProfile.ValueString() != "" {
				storageProfile, err := vdc.FindStorageProfileReference(v.Plan.StorageProfile.ValueString())
				if err != nil {
					return nil, fmt.Errorf("error retrieving storage profile %s : %w", v.Plan.StorageProfile.ValueString(), err)
				}
				_, err = vm.UpdateStorageProfile(storageProfile.HREF)
				if err != nil {
					return nil, fmt.Errorf("error updating changing storage profile to %s: %w", v.Plan.StorageProfile.ValueString(), err)
				}
			}
		}
	}
	// End of update networks

	// Update Guest Properties
	if !v.Plan.GuestProperties.Equal(v.State.GuestProperties) {
		err = addRemoveGuestProperties(v, vm)
		if err != nil {
			return nil, fmt.Errorf("unable to update guest properties: %w", err)
		}
	}
	// End of update guest properties

	// Update Sizing Policy
	if !v.Plan.SizingPolicyID.Equal(v.State.SizingPolicyID) || !v.Plan.PlacementPolicyID.Equal(v.State.PlacementPolicyID) {
		_, err = vm.UpdateComputePolicyV2(v.Plan.SizingPolicyID.ValueString(), v.Plan.PlacementPolicyID.ValueString(), "")
		if err != nil {
			return nil, fmt.Errorf("unable to update compute policy: %w", err)
		}
	}
	// End of update sizing policy

	// Update Storage Profile
	if !v.Plan.StorageProfile.Equal(v.State.StorageProfile) {
		if !v.Plan.StorageProfile.IsUnknown() && v.Plan.StorageProfile.ValueString() != "" {
			sP, err := vdc.FindStorageProfileReference(v.Plan.StorageProfile.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error retrieving storage profile %s : %w", v.Plan.StorageProfile.ValueString(), err)
			}

			_, err = vm.UpdateStorageProfile(sP.HREF)
			if err != nil {
				return nil, fmt.Errorf("error updating changing storage profile to %s: %w", v.Plan.StorageProfile.ValueString(), err)
			}
		}
	}
	// End of update storage profile

	// Update Customization and Computer Name
	if !v.Plan.Customization.Equal(v.State.Customization) || !v.Plan.ComputerName.Equal(v.State.ComputerName) {
		err = updateGuestCustomizationSetting(v, vm)
		if err != nil {
			return nil, fmt.Errorf("unable to update guest customization: %w", err)
		}
	}
	// End of update customization and computer name

	// ! End of Hot Update

	// ! Cold Update (VM must be powered off)
	vmStatusBeforeUpdate, err := vm.GetStatus()
	if err != nil {
		return nil, fmt.Errorf("unable to get VM status before update: %w", err)
	}

	// if power_on, boot_image, expose_hardware_virtualization, os_type, description, cpu_hot_add_enabled, memory_hot_add_enabled haschange
	if !v.Plan.PowerON.Equal(v.State.PowerON) ||
		!v.Plan.ExposeHardwareVirtualization.Equal(v.State.ExposeHardwareVirtualization) ||
		!v.Plan.OsType.Equal(v.State.OsType) ||
		!v.Plan.Description.Equal(v.State.Description) ||
		!resource.CPUHotAddEnabled.Equal(resourceState.CPUHotAddEnabled) ||
		!resource.MemoryHotAddEnabled.Equal(resourceState.MemoryHotAddEnabled) ||
		!resource.CPUCores.Equal(resourceState.CPUCores) ||
		networksNeedsColdChange ||
		cpuNeedsColdChange ||
		memoryNeedsColdChange {
		if vmStatusBeforeUpdate != "POWERED_OFF" {
			if v.Plan.PreventUpdatePowerOff.IsNull() || v.Plan.PreventUpdatePowerOff.IsUnknown() || v.Plan.PreventUpdatePowerOff.ValueBool() {
				return nil, fmt.Errorf("update stopped: VM needs to power off to change properties, but `prevent_update_power_off` is `true`")
			}

			task, err := vm.Undeploy()
			if err != nil {
				return nil, fmt.Errorf("error triggering undeploy for VM %s: %w", vm.VM.Name, err)
			}

			err = task.WaitTaskCompletion()
			if err != nil {
				return nil, fmt.Errorf("error waiting for undeploy task for VM %s: %w", vm.VM.Name, err)
			}
		}

		// Update CPUCores
		if !resource.CPUCores.Equal(resourceState.CPUCores) || cpuNeedsColdChange {
			err := vm.ChangeCPUAndCoreCount(utils.TakeIntPointer(int(resource.CPUs.ValueInt64())), utils.TakeIntPointer(int(resource.CPUCores.ValueInt64())))
			if err != nil {
				return nil, fmt.Errorf("unable to update CPU cores: %w", err)
			}
		}
		// End of update CPUCores

		// Update Memory
		if memoryNeedsColdChange {
			err := vm.ChangeMemory(resource.Memory.ValueInt64())
			if err != nil {
				return nil, fmt.Errorf("unable to update memory: %w", err)
			}
		}
		// End of update Memory

		// Update Network
		if networksNeedsColdChange {
			networkConnectionSection, err := networksToConfig(v, vapp)
			if err != nil {
				return nil, fmt.Errorf("unable to update network: %w", err)
			}

			err = vm.UpdateNetworkConnectionSection(&networkConnectionSection)
			if err != nil {
				return nil, fmt.Errorf("unable to update network: %w", err)
			}
		}
		// End of update Network

		// Update ExposeHardwareVirtualization
		if !v.Plan.ExposeHardwareVirtualization.Equal(v.State.ExposeHardwareVirtualization) {
			task, err := vm.ToggleHardwareVirtualization(v.Plan.ExposeHardwareVirtualization.ValueBool())
			if err != nil {
				return nil, fmt.Errorf("error changing hardware assisted virtualization: %w", err)
			}

			err = task.WaitTaskCompletion()
			if err != nil {
				return nil, err
			}
		}
		// End of update ExposeHardwareVirtualization

		// Update VM spec section and description
		var (
			vmSpecSectionUpdate = false
			vmSpecSection       = vm.VM.VmSpecSection
			description         = vm.VM.Description
		)
		if !v.Plan.Description.Equal(v.State.Description) {
			vmSpecSectionUpdate = true
			description = v.Plan.Description.ValueString()
		}

		if !v.Plan.OsType.Equal(v.State.OsType) {
			vmSpecSectionUpdate = true
			vmSpecSection.OsType = v.Plan.OsType.ValueString()
		}

		if vmSpecSectionUpdate {
			task, err := vm.UpdateVmSpecSectionAsync(vmSpecSection, description)
			if err != nil {
				return nil, fmt.Errorf("unable to update VM spec section: %w", err)
			}

			err = task.WaitTaskCompletion()
			if err != nil {
				return nil, fmt.Errorf("unable to update VM spec section: %w", err)
			}
		}
		// End of update VM spec section and description

		// Update CPU/Memory Hot Add
		if !resource.CPUHotAddEnabled.Equal(resourceState.CPUHotAddEnabled) ||
			!resource.MemoryHotAddEnabled.Equal(resourceState.MemoryHotAddEnabled) {
			_, err := vm.UpdateVmCpuAndMemoryHotAdd(resource.CPUHotAddEnabled.ValueBool(), resource.MemoryHotAddEnabled.ValueBool())
			if err != nil {
				return nil, fmt.Errorf("unable to update CPU/Memory Hot Add: %w", err)
			}
		}
		// End of update CPU/Memory Hot Add
	}
	// ! End of Cold Update

	// Update Power ON/OFF
	if v.Plan.PowerON.ValueBool() {
		vmStatus, err := vm.GetStatus()
		if err != nil {
			return nil, fmt.Errorf("unable to get VM status: %w", err)
		}

		if !customizationNeeded && vmStatus != "POWERED_ON" {
			task, err := vm.PowerOn()
			if err != nil {
				return nil, fmt.Errorf("unable to power on VM: %w", err)
			}

			err = task.WaitTaskCompletion()
			if err != nil {
				return nil, fmt.Errorf("unable to power on VM: %w", err)
			}
		}

		if customizationNeeded {
			if vmStatus != "POWERED_OFF" {
				task, err := vm.Undeploy()
				if err != nil {
					return nil, fmt.Errorf("unable to undeploy VM: %w", err)
				}

				err = task.WaitTaskCompletion()
				if err != nil {
					return nil, fmt.Errorf("unable to undeploy VM: %w", err)
				}
			}

			err = vm.PowerOnAndForceCustomization()
			if err != nil {
				return nil, fmt.Errorf("unable to power on and force customization: %w", err)
			}
		}
	}
	// End of update power on/off

	return vm, nil
}
