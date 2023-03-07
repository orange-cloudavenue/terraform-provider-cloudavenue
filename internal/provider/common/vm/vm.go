package vm

import (
	"fmt"
	"log"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

const (
	// Used when a task fails. The placeholder is for the error
	errorCompletingTask = "error completing tasks: %s"
)

// GetVM returns a VM by name.
func GetVM(vdc *govcd.Vdc, vappNameOrID, vmNameOrID string) (*govcd.VM, error) {
	vapp, err := vdc.GetVAppByNameOrId(vappNameOrID, true)
	if err != nil {
		return nil, fmt.Errorf("[getVm] failed to get vApp: %s", err)
	}
	vm, err := vapp.GetVMByNameOrId(vmNameOrID, false)
	if err != nil {
		return nil, fmt.Errorf("[getVm] failed to get VM: %s", err)
	}
	return vm, err
}

// PowerOnIfNeeded powers on a VM if it was powered on before and the bus type is IDE.
func PowerOnIfNeeded(vm *govcd.VM, busType string, allowVMReboot bool, vmStatusBefore string) error {
	vmStatus, err := vm.GetStatus()
	if err != nil {
		return fmt.Errorf("error getting VM status before ensuring it is powered on: %s", err)
	}

	if vmStatusBefore == "POWERED_ON" && vmStatus != "POWERED_ON" && busType == "ide" && allowVMReboot {
		log.Printf("[DEBUG] Powering on VM %s after adding internal disk.", vm.VM.Name)

		task, err := vm.PowerOn()
		if err != nil {
			return fmt.Errorf("error powering on VM for adding/updating internal disk: %s", err)
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
		return "", fmt.Errorf("error getting VM status before ensuring it is powered off: %s", err)
	}
	vmStatusBefore := vmStatus

	if vmStatus != "POWERED_OFF" && busType == "ide" && allowVMReboot {
		log.Printf("[DEBUG] Powering off VM %s for adding/updating internal disk.", vm.VM.Name)

		task, err := vm.PowerOff()
		if err != nil {
			return vmStatusBefore, fmt.Errorf("error powering off VM for adding internal disk: %s", err)
		}
		err = task.WaitTaskCompletion()
		if err != nil {
			return vmStatusBefore, fmt.Errorf(errorCompletingTask, err)
		}
	}
	return vmStatusBefore, nil
}

var (
	vcdMutexKV = mutex.NewKV()
)
