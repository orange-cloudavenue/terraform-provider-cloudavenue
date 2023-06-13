package vm

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type VM struct {
	*client.VM
	vApp vapp.VAPP
}

type GetVMOpts struct {
	ID   types.String
	Name types.String
}

// vmIDOrName returns the ID or name of the VM.
func (v GetVMOpts) vmIDOrName() string {
	if v.ID.IsNull() || v.ID.IsUnknown() {
		return v.Name.ValueString()
	}
	return v.ID.ValueString()
}

// ConstructObject is a special function that is used to construct the VM object from the govcd.VM.
func ConstructObject(vApp vapp.VAPP, vm *govcd.VM) VM {
	return VM{VM: &client.VM{VM: vm}, vApp: vApp}
}

/*
Init

Initializes a VM struct with a VM and a vApp.
*/
func Init(_ *client.CloudAvenue, vApp vapp.VAPP, vmInfo GetVMOpts) (vm VM, d diag.Diagnostics) {
	vmOut, err := vApp.GetVMByNameOrId(vmInfo.vmIDOrName(), true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			d.AddError("VM not found", err.Error())
			return VM{}, d
		}
		d.AddError("Error retrieving VM", err.Error())
		return VM{}, d
	}
	return VM{VM: &client.VM{VM: vmOut}, vApp: vApp}, nil
}

func Get(vApp vapp.VAPP, vmInfo GetVMOpts) (vm VM, d diag.Diagnostics) {
	return Init(nil, vApp, vmInfo)
}

func (v VM) constructLockKey() string {
	return fmt.Sprintf("vm:%s", v.GetID())
}

// LockVM locks VM.
func (v VM) LockVM(ctx context.Context) (d diag.Diagnostics) {
	if v.GetID() == "" || ctx == nil {
		d.AddError("Incorrect lock args", "VM: "+v.GetID())
		return
	}

	mutex.GlobalMutex.KvLock(ctx, v.constructLockKey())
	return
}

// UnlockVM unlocks VM.
func (v VM) UnlockVM(ctx context.Context) (d diag.Diagnostics) {
	if v.GetID() == "" || ctx == nil {
		d.AddError("Incorrect Unlock args", "VM: "+v.GetID())
		return
	}

	mutex.GlobalMutex.KvUnlock(ctx, v.constructLockKey())
	return
}

// GetName returns the name of the VM.
func (v VM) GetName() string {
	return v.VM.VM.VM.Name
}

// GetID returns the ID of the VM.
func (v VM) GetID() string {
	return v.VM.VM.VM.ID
}

// GetDescription returns the description of the VM.
func (v VM) GetDescription() string {
	return v.VM.VM.VM.Description
}

// GetDiskSettings returns the disk settings of the VM.
func (v VM) GetDiskSettings() []*govcdtypes.DiskSettings {
	return v.VM.VM.VM.VmSpecSection.DiskSection.DiskSettings
}

// IsPoweredOn returns true if the VM is powered on.
func (v VM) IsPoweredON() bool {
	status, _ := v.GetStatus()
	return status == powerON
}

// GetCPUCount returns the number of CPUs of the VM.
func (v VM) GetCPUCount() int64 {
	return int64(*v.VM.VM.VM.VmSpecSection.NumCpus)
}

// GetCPUCoresCount returns the number of CPU cores of the VM.
func (v VM) GetCPUCoresCount() int64 {
	return int64(*v.VM.VM.VM.VmSpecSection.NumCoresPerSocket)
}

// IsCPUHotAddEnabled returns true if CPU hot add is enabled.
func (v VM) IsCPUHotAddEnabled() bool {
	return v.VM.VM.VM.VMCapabilities.CPUHotAddEnabled
}

// AttachDiskSettings represents the settings for attaching a disk to a VM.
func (v VM) AttachDiskSettings(busNumber, unitNumber types.Int64, diskHREF string) *govcdtypes.DiskAttachOrDetachParams {
	var b, u int

	if busNumber.IsNull() || unitNumber.IsNull() {
		b, u = diskparams.ComputeBusAndUnitNumber(v.GetDiskSettings())
	} else {
		b = int(busNumber.ValueInt64())
		u = int(unitNumber.ValueInt64())
	}

	return &govcdtypes.DiskAttachOrDetachParams{
		Disk:       &govcdtypes.Reference{HREF: diskHREF},
		BusNumber:  utils.TakeIntPointer(b),
		UnitNumber: utils.TakeIntPointer(u),
	}
}

type NetworkConnection struct {
	Name             types.String
	Connected        types.Bool
	IPAllocationMode types.String
	IP               types.String
	Type             types.String
	Mac              types.String
	AdapterType      types.String
	IsPrimary        types.Bool
}

func ConstructNetworksConnectionWithoutVM(vApp vapp.VAPP, networks []NetworkConnection) (networkConnection govcdtypes.NetworkConnectionSection, err error) {
	x := VM{
		vApp: vApp,
	}
	return x.ConstructNetworksConnection(networks)
}

// ConstructNetworksConnection constructs a NetworkConnectionSection from a list of NetworkConnection.
func (v VM) ConstructNetworksConnection(networks []NetworkConnection) (networkConnection govcdtypes.NetworkConnectionSection, err error) {
	for index, network := range networks {
		netCon := &govcdtypes.NetworkConnection{
			Network:                 network.Name.ValueString(),
			IsConnected:             network.Connected.ValueBool(),
			IPAddressAllocationMode: network.IPAllocationMode.ValueString(),
			IPAddress:               network.IP.ValueString(),
			NetworkConnectionIndex:  index,
		}

		if v.vApp.VAPP == nil {
			return govcdtypes.NetworkConnectionSection{}, fmt.Errorf("vApp is not initialized")
		}

		switch network.Type.ValueString() {
		case "vapp":
			if ok, err := v.vApp.IsVAPPNetwork(network.Name.ValueString()); err != nil {
				return govcdtypes.NetworkConnectionSection{}, err
			} else if !ok {
				return govcdtypes.NetworkConnectionSection{}, fmt.Errorf("vApp network : %s is not found", network.Name.ValueString())
			}
		case "org":
			if ok, err := v.vApp.IsVAPPOrgNetwork(network.Name.ValueString()); err != nil {
				return govcdtypes.NetworkConnectionSection{}, err
			} else if !ok {
				return govcdtypes.NetworkConnectionSection{}, fmt.Errorf("org network : %s is not found", network.Name.ValueString())
			}
		}

		if network.Mac.ValueString() != "" {
			netCon.MACAddress = network.Mac.ValueString()
		}

		if network.AdapterType.ValueString() != "" {
			netCon.NetworkAdapterType = network.AdapterType.ValueString()
		}

		networkConnection.NetworkConnection = append(networkConnection.NetworkConnection, netCon)

		if network.IsPrimary.ValueBool() {
			networkConnection.PrimaryNetworkConnectionIndex = index
		}
	}

	return networkConnection, nil
}

// ! LEGACY

const (
	// Used when a task fails. The placeholder is for the error.
	errorCompletingTask = "error completing tasks: %s"
)

// GetVM returns a VM by name.
func GetVMOLD(vdc *govcd.Vdc, vappNameOrID, vmNameOrID string) (*govcd.VM, error) {
	vapp, err := vdc.GetVAppByNameOrId(vappNameOrID, true)
	if err != nil {
		return nil, fmt.Errorf("[getVm] failed to get vApp: %w", err)
	}
	vm, err := vapp.GetVMByNameOrId(vmNameOrID, false)
	if err != nil {
		return nil, fmt.Errorf("[getVm] failed to get VM: %w", err)
	}
	return vm, err
}

// PowerOnIfNeeded powers on a VM if it was powered on before and the bus type is IDE.
func PowerOnIfNeeded(vm *govcd.VM, busType string, allowVMReboot bool, vmStatusBefore string) error {
	vmStatus, err := vm.GetStatus()
	if err != nil {
		return fmt.Errorf("error getting VM status before ensuring it is powered on: %w", err)
	}

	if vmStatusBefore == "POWERED_ON" && vmStatus != "POWERED_ON" && busType == "ide" && allowVMReboot {
		log.Printf("[DEBUG] Powering on VM %s after adding internal disk.", vm.VM.Name)

		task, err := vm.PowerOn()
		if err != nil {
			return fmt.Errorf("error powering on VM for adding/updating internal disk: %w", err)
		}
		err = task.WaitTaskCompletion()
		if err != nil {
			return fmt.Errorf(errorCompletingTask, err)
		}
	}
	return nil
}

// PowerOffIfNeeded powers off a VM if it was powered off before and the bus type is IDE.
func PowerOffIfNeeded(vm *govcd.VM, busType string, allowVMReboot bool) (string, error) {
	vmStatus, err := vm.GetStatus()
	if err != nil {
		return "", fmt.Errorf("error getting VM status before ensuring it is powered off: %w", err)
	}
	vmStatusBefore := vmStatus

	if vmStatus != "POWERED_OFF" && busType == "ide" && allowVMReboot {
		log.Printf("[DEBUG] Powering off VM %s for adding/updating internal disk.", vm.VM.Name)

		task, err := vm.PowerOff()
		if err != nil {
			return vmStatusBefore, fmt.Errorf("error powering off VM for adding internal disk: %w", err)
		}
		err = task.WaitTaskCompletion()
		if err != nil {
			return vmStatusBefore, fmt.Errorf(errorCompletingTask, err)
		}
	}
	return vmStatusBefore, nil
}

// var vcdMutexKV = mutex.NewKV()
