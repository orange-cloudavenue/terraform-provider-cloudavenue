package vm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

// readVM reads the vApp VM configuration from the vApp.
func readVM(v *Client) (*govcd.VM, error) {
	var (
		vdc  *govcd.Vdc
		vapp *govcd.VApp
		vm   *govcd.VM
		err  error
	)

	// If VDC is not defined at resource/data source level, use the one defined at provider level
	if v.State.VDC.IsNull() || v.State.VDC.IsUnknown() {
		if v.Client.DefaultVDCExist() {
			v.State.VDC = types.StringValue(v.Client.GetDefaultVDC())
		} else {
			return nil, fmt.Errorf("VDC is required when not defined at provider level")
		}
	}

	// Get vcd object
	_, vdcHandler, err := v.Client.GetOrgAndVDC(v.Client.GetOrg(), v.State.VDC.ValueString())
	if err != nil {
		return nil, fmt.Errorf("error retrieving VDC %s: %w", v.State.VDC.ValueString(), err)
	}

	vdc, isVDC := vdcHandler.(*govcd.Vdc)
	if !isVDC {
		return nil, fmt.Errorf("expected *govcd.Vdc type, have %T", vdcHandler)
	}

	// Get vApp
	vapp, err = vdc.GetVAppByName(v.State.VappName.ValueString(), true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving vApp %s: %w", v.State.VappName.ValueString(), err)
	}

	// Get VM
	vm, err = vapp.GetVMByNameOrId(v.State.ID.ValueString(), true)
	if err != nil {
		if govcd.IsNotFound(err) {
			return nil, errRemoveResource
		}
		return nil, fmt.Errorf("error retrieving VM %s (ID:%s): %w", v.State.VMName.ValueString(), v.State.ID.ValueString(), err)
	}

	return vm, nil
}
