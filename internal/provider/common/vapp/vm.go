package vapp

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type VM struct {
	vApp VAPP
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
