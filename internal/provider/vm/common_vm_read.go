package vm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

// readVM reads the vApp VM configuration from the vApp
func readVM(v *VMClient) (*govcd.VM, error) {
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
		if govcd.IsNotFound(err) {
			return nil, errRemoveResource
		}
		return nil, fmt.Errorf("error retrieving vApp %s: %s", v.Plan.VappName.ValueString(), err)
	}

	// Get VM
	vm, err = vapp.GetVMByNameOrId(v.Plan.VMName.ValueString(), false)
	if err != nil {
		if govcd.IsNotFound(err) {
			return nil, errRemoveResource
		}
		return nil, fmt.Errorf("error retrieving VM %s: %s", v.Plan.VMName.ValueString(), err)
	}

	return vm, nil
}
