package client

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type VM struct {
	*govcd.VM
}

// * Guest properties
// GetGuestProperties returns the guest properties of a VM.
func (v VM) GetGuestProperties() (guestProperties []*govcdtypes.Property, err error) {
	x, err := v.GetProductSectionList()
	if err != nil {
		return nil, err
	}

	guestProperties = append(guestProperties, x.ProductSection.Property...)

	return
}

// SetGuestProperties sets the guest properties of a VM
// If the guest property already exists, it will be updated.
func (v VM) SetGuestProperties(guestProperties map[string]string) (err error) {
	listGuestProperties := make([]*govcdtypes.Property, 0)

	for key, value := range guestProperties {
		listGuestProperties = append(listGuestProperties, &govcdtypes.Property{
			UserConfigurable: true,
			Type:             "string",
			Key:              key,
			Label:            key,
			Value:            &govcdtypes.Value{Value: value},
		})
	}

	_, err = v.SetProductSectionList(&govcdtypes.ProductSectionList{
		ProductSection: &govcdtypes.ProductSection{
			Info:     "Custom properties",
			Property: listGuestProperties,
		},
	})

	return
}

// * Customization

// GetCustomization returns the customization of a VM.
func (v VM) GetCustomization() (guestCustomization *govcdtypes.GuestCustomizationSection, err error) {
	return v.GetGuestCustomizationSection()
}

// SetCustomization sets the customization of a VM.
func (v VM) SetCustomization(guestCustomization *govcdtypes.GuestCustomizationSection) (err error) {
	_, err = v.SetGuestCustomizationSection(guestCustomization)
	return
}

// * OS type

// SetOSType sets the OS type of a VM.
func (v VM) SetOSType(osType string) (err error) {
	updateOsType := v.VM.VM.VmSpecSection

	updateOsType.OsType = osType

	_, err = v.UpdateVmSpecSection(updateOsType, v.VM.VM.Description)
	return
}

// SetExposeHardwareVirtualization sets the expose hardware virtualization of a VM.
func (v VM) SetExposeHardwareVirtualization(isEnabled bool) (err error) {
	task, err := v.ToggleHardwareVirtualization(isEnabled)
	if err != nil {
		return err
	}

	return task.WaitTaskCompletion()
}

// GetExposeHardwareVirtualization returns the expose hardware virtualization of a VM.
func (v VM) GetExposeHardwareVirtualization() bool {
	return v.VM.VM.NestedHypervisorEnabled
}

// NetworkConnectionIsDefined checks if a network connection is defined.
func (v VM) NetworkConnectionIsDefined() bool {
	return v.VM.VM.NetworkConnectionSection.NetworkConnection != nil
}

// GetNetworkConnection returns the network connection of a VM.
func (v VM) GetNetworkConnection() (networkConnections []*govcdtypes.NetworkConnection) {
	if v.VM.VM.NetworkConnectionSection == nil {
		return nil
	}

	return v.VM.VM.NetworkConnectionSection.NetworkConnection
}

// GetOSType returns the OS type of a VM.
func (v VM) GetOSType() string {
	return v.VM.VM.VmSpecSection.OsType
}

// GetAffinityRuleID returns the affinity rule ID of a VM.
func (v VM) GetAffinityRuleID() string {
	if v.VM.VM.ComputePolicy == nil || v.VM.VM.ComputePolicy.VmPlacementPolicy == nil {
		return ""
	}

	return v.VM.VM.ComputePolicy.VmPlacementPolicy.ID
}

// GetDefaultAffinityRuleID returns the default affinity rule ID of a VM.
func (v VM) GetDefaultAffinityRuleID() (string, error) {
	vdc, err := v.GetParentVdc()
	if err != nil {
		return "", err
	}

	return vdc.Vdc.DefaultComputePolicy.ID, nil
}

// GetAffinityRuleIDOrDefault returns the affinity rule ID of a VM or the default affinity rule ID if the VM has no affinity rule ID.
func (v VM) GetAffinityRuleIDOrDefault() (string, error) {
	affinityRuleID := v.GetAffinityRuleID()
	if affinityRuleID != "" {
		return affinityRuleID, nil
	}

	return v.GetDefaultAffinityRuleID()
}

// GetStorageProfileName returns the storage profile name of a VM.
func (v VM) GetStorageProfileName() string {
	if v.VM.VM.StorageProfile == nil {
		return ""
	}

	return v.VM.VM.StorageProfile.Name
}

// IsCpusIsDefined returns true if the number of CPUs of a VM is defined.
func (v VM) CpusIsDefined() bool {
	return v.VM.VM.VmSpecSection.NumCpus != nil
}

// IsCpusCoresIsDefined returns true if the number of cores per CPU of a VM is defined.
func (v VM) CpusCoresIsDefined() bool {
	return v.VM.VM.VmSpecSection.NumCoresPerSocket != nil
}

// GetCpus returns the number of CPUs of a VM.
func (v VM) GetCpus() int {
	if !v.CpusIsDefined() {
		return 0
	}

	return *v.VM.VM.VmSpecSection.NumCpus
}

// GetCpusCores returns the number of cores per CPU of a VM.
func (v VM) GetCpusCores() int {
	if !v.CpusCoresIsDefined() {
		return 0
	}

	return *v.VM.VM.VmSpecSection.NumCoresPerSocket
}

// MemoryIsDefined returns true if the memory of a VM is defined.
func (v VM) MemoryIsDefined() bool {
	return v.VM.VM.VmSpecSection.MemoryResourceMb != nil
}

// GetMemory returns the memory of a VM.
func (v VM) GetMemory() int64 {
	if !v.MemoryIsDefined() {
		return 0
	}

	return v.VM.VM.VmSpecSection.MemoryResourceMb.Configured
}

// HotAddIsDefined returns true if the hot add of a VM is defined.
func (v VM) HotAddIsDefined() bool {
	return v.VM.VM.VMCapabilities != nil
}

// GetCpuHotAdd returns the hot add of a VM.
func (v VM) GetCpuHotAddEnabled() bool {
	if !v.HotAddIsDefined() {
		return false
	}

	return v.VM.VM.VMCapabilities.CPUHotAddEnabled
}

// GetMemoryHotAdd returns the hot add of a VM.
func (v VM) GetMemoryHotAddEnabled() bool {
	if !v.HotAddIsDefined() {
		return false
	}

	return v.VM.VM.VMCapabilities.MemoryHotAddEnabled
}
