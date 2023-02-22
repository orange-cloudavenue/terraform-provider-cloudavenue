package vm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

func updateVM(v *VMClient) (*govcd.VM, error) { //nolint:gocyclo
	var (
		vdc  *govcd.Vdc
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
	_, vdc, err = v.Client.GetOrgAndVDC(v.Client.GetOrg(), v.Plan.VDC.ValueString())
	if err != nil {
		return nil, fmt.Errorf("error retrieving VDC %s: %s", v.Plan.VDC.ValueString(), err)
	}

	// Get vApp
	vapp, err = vdc.GetVAppByName(v.Plan.VappName.ValueString(), false)
	if err != nil {
		return nil, fmt.Errorf("error retrieving vApp %s: %s", v.Plan.VappName.ValueString(), err)
	}

	// Get VM
	vm, err = vapp.GetVMByNameOrId(v.Plan.VMName.ValueString(), false)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VM %s: %s", v.Plan.VMName.ValueString(), err)
	}

	if !v.Plan.Resource.MemoryHotAddEnabled.Equal(v.State.Resource.MemoryHotAddEnabled) && v.Plan.Resource.Memory.Equal(v.State.Resource.Memory) {
		if err := vm.ChangeMemory(v.Plan.Resource.Memory.ValueInt64()); err != nil {
			return nil, fmt.Errorf("error changing memory: %s", err)
		}
	}

	if !v.Plan.Resource.CPUHotAddEnabled.Equal(v.State.Resource.CPUHotAddEnabled) && v.Plan.Resource.CPUs.Equal(v.State.Resource.CPUs) {
		if err := vm.ChangeCPU(int(v.Plan.Resource.CPUs.ValueInt64()), int(v.Plan.Resource.CPUCores.ValueInt64())); err != nil {
			return nil, fmt.Errorf("error changing CPU: %s", err)
		}
	}

	if len(v.Plan.Networks) != len(v.State.Networks) {
		requireUpdate := false
		foundPrimaryNic := false

		for i, network := range v.Plan.Networks {
			if network.IsPrimary.ValueBool() {
				foundPrimaryNic = true
			}

			if !network.AdapterType.Equal(v.State.Networks[i].AdapterType) {
				requireUpdate = true
				break
			}

			if !network.Type.Equal(v.State.Networks[i].Type) {
				requireUpdate = true
				break
			}

			if !network.IPAllocationMode.Equal(v.State.Networks[i].IPAllocationMode) {
				requireUpdate = true
				break
			}

			if !network.Name.Equal(v.State.Networks[i].Name) {
				requireUpdate = true
				break
			}

			if !network.IsPrimary.Equal(v.State.Networks[i].IsPrimary) {
				requireUpdate = true
				break
			}

			if !network.IP.Equal(v.State.Networks[i].IP) {
				requireUpdate = true
				break
			}

			if !network.Mac.Equal(v.State.Networks[i].Mac) {
				requireUpdate = true
				break
			}

			if !network.Connected.Equal(v.State.Networks[i].Connected) {
				requireUpdate = true
				break
			}
		}

		// * Primary NIC cannot be removed on a powered on VM
		if requireUpdate && foundPrimaryNic {
			networkConnectionSection, err := networksToConfig(v, vapp)
			if err != nil {
				return nil, fmt.Errorf("unable to setup network configuration for update: %s", err)
			}
			err = vm.UpdateNetworkConnectionSection(&networkConnectionSection)
			if err != nil {
				return nil, fmt.Errorf("unable to update network configuration: %s", err)
			}
		}

		err = addRemoveGuestProperties(v, vm)
		if err != nil {
			return nil, fmt.Errorf("unable to update guest properties: %s", err)
		}

		if !v.Plan.SizingPolicyID.Equal(v.State.SizingPolicyID) || !v.Plan.PlacementPolicyID.Equal(v.State.PlacementPolicyID) {
			_, err = vm.UpdateComputePolicyV2(v.Plan.SizingPolicyID.ValueString(), v.Plan.PlacementPolicyID.ValueString(), "")
			if err != nil {
				return nil, fmt.Errorf("unable to update compute policy: %s", err)
			}
		}

		if !v.Plan.StorageProfile.Equal(v.State.StorageProfile) {
			if v.Plan.StorageProfile.ValueString() != "" {
				storageProfile, err := vdc.FindStorageProfileReference(v.Plan.StorageProfile.ValueString())
				if err != nil {
					return nil, fmt.Errorf("error retrieving storage profile %s : %s", v.Plan.StorageProfile.ValueString(), err)
				}
				_, err = vm.UpdateStorageProfile(storageProfile.HREF)
				if err != nil {
					return nil, fmt.Errorf("error updating changing storage profile to %s: %s", v.Plan.StorageProfile.ValueString(), err)
				}
			}
		}
	}

	return vm, nil
}
