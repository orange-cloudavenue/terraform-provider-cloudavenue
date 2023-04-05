package vm

import (
	"fmt"
	"log"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type VM struct {
	*client.VM
}

// GetName returns the name of the VM.
func (v VM) GetName() string {
	return v.VM.VM.VM.Name
}

// GetID returns the ID of the VM.
func (v VM) GetID() string {
	return v.VM.VM.VM.ID
}

// AttachDiskSettings represents the settings for attaching a disk to a VM.
func (v VM) AttachDiskSettings(busNumber, unitNumber types.Int64, diskHREF string) *govcdtypes.DiskAttachOrDetachParams {
	var b, u int

	if busNumber.IsNull() || unitNumber.IsNull() {
		b, u = diskparams.ComputeBusAndUnitNumber(v.VM.VM.VM.VmSpecSection.DiskSection.DiskSettings)
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

// ! LEGACY

const (
	// Used when a task fails. The placeholder is for the error.
	errorCompletingTask = "error completing tasks: %s"
)

// GetVM returns a VM by name.
func GetVM(vdc *govcd.Vdc, vappNameOrID, vmNameOrID string) (*govcd.VM, error) {
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
