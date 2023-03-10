package vm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// createVmFromTemplate is responsible for create vApp VMs from template :
// * VMs inside vApp (vApp VMs)
//
// Code flow has 3 layers:
// 1. Lookup common information, required for both types of VMs (Standalone and vApp child). Things such as
//   - Template to be used
//   - Network adapter configuration
//   - Storage profile configuration
//   - VM compute policy configuration
//
// # This part includes defining initial structures for VM and also any explicitly required operations for that type of VM
//
// 3. Perform additional operations which are common for both types of VMs
//
// Note. VM Power ON (if it wasn't disabled in HCL configuration) occurs as last step after all configuration is done.
func createVM(_ context.Context, v *Client) (vm *govcd.VM, err error) { //nolint:gocyclo
	var (
		vdc        *govcd.Vdc
		org        *govcd.Org
		vapp       *govcd.VApp
		vmTemplate govcd.VAppTemplate
	)

	// If VDC is not defined at resource level, use the one defined at provider level
	if v.Plan.VDC.IsNull() || v.Plan.VDC.IsUnknown() {
		if v.Client.DefaultVDCExist() {
			v.Plan.VDC = types.StringValue(v.Client.GetDefaultVDC())
		} else {
			return nil, fmt.Errorf("VDC is required when not defined at provider level")
		}
	}

	// Get vcd object
	org, vdc, err = v.Client.GetOrgAndVDC(v.Client.GetOrg(), v.Plan.VDC.ValueString())
	if err != nil {
		return nil, fmt.Errorf("error retrieving VDC %s: %w", v.Plan.VDC.ValueString(), err)
	}

	if !v.Plan.VappTemplateID.IsNull() && !v.Plan.VappTemplateID.IsUnknown() {
		// Look up VM template inside vApp template - either specified by `vm_name_in_template` or the
		// first one in vApp
		vmTemplate, err = lookupvAppTemplateforVM(v, org, vdc)
		if err != nil {
			return nil, fmt.Errorf("error finding vApp template: %w", err)
		}
	}

	// Get vApp
	vapp, err = vdc.GetVAppByName(v.Plan.VappName.ValueString(), false)
	if err != nil {
		return nil, fmt.Errorf("[VM create] error finding vApp %s: %w", v.Plan.VappName.ValueString(), err)
	}

	// Build up network configuration
	networkConnectionSection, err := networksToConfig(v, vapp)
	if err != nil {
		return nil, fmt.Errorf("unable to process network configuration: %w", err)
	}

	// Lookup storage profile reference if it was specified
	storageProfilePtr, err := lookupStorageProfile(v.Plan.StorageProfile.ValueString(), vdc)
	if err != nil {
		return nil, fmt.Errorf("error finding storage profile: %w", err)
	}

	// Look up compute policies
	sizingPolicy, err := lookupComputePolicy(v, v.Plan.SizingPolicyID.ValueString())
	if err != nil {
		return nil, fmt.Errorf("error finding sizing policy: %w", err)
	}
	placementPolicy, err := lookupComputePolicy(v, v.Plan.PlacementPolicyID.ValueString())
	if err != nil {
		return nil, fmt.Errorf("error finding placement policy: %w", err)
	}
	var vmComputePolicy *govcdtypes.ComputePolicy
	if sizingPolicy != nil || placementPolicy != nil {
		vmComputePolicy = &govcdtypes.ComputePolicy{}
		if sizingPolicy != nil {
			vmComputePolicy.VmSizingPolicy = &govcdtypes.Reference{HREF: sizingPolicy.Href}
		}
		if placementPolicy != nil {
			vmComputePolicy.VmPlacementPolicy = &govcdtypes.Reference{HREF: placementPolicy.Href}
		}
	}

	if !v.Plan.VappTemplateID.IsNull() && !v.Plan.VappTemplateID.IsUnknown() {
		vmFromTemplateParams := &govcdtypes.ReComposeVAppParams{
			Ovf:              govcdtypes.XMLNamespaceOVF,
			Xsi:              govcdtypes.XMLNamespaceXSI,
			Xmlns:            govcdtypes.XMLNamespaceVCloud,
			AllEULAsAccepted: v.Plan.AcceptAllEulas.ValueBool(),
			Name:             vapp.VApp.Name,
			PowerOn:          false, // VM will be powered on after all configuration is done
			SourcedItem: &govcdtypes.SourcedCompositionItemParam{
				Source: &govcdtypes.Reference{
					HREF: vmTemplate.VAppTemplate.HREF,
					Name: v.Plan.VMName.ValueString(), // This VM name defines the VM name after creation
				},
				VMGeneralParams: &govcdtypes.VMGeneralParams{
					Description: v.Plan.Description.ValueString(),
				},
				InstantiationParams: &govcdtypes.InstantiationParams{
					// If a MAC address is specified for NIC - it does not get set with this call,
					// therefore an additional `vm.UpdateNetworkConnectionSection` is required.
					NetworkConnectionSection: &networkConnectionSection,
				},
				ComputePolicy:  vmComputePolicy,
				StorageProfile: storageProfilePtr,
			},
		}

		vm, err = vapp.AddRawVM(vmFromTemplateParams)
		if err != nil {
			return nil, errRemoveResource
		}
	} else {
		var bootImage *govcdtypes.Media

		storageProfilePtr, err := lookupStorageProfile(v.Plan.StorageProfile.ValueString(), vdc)
		if err != nil {
			return nil, fmt.Errorf("error finding storage profile: %w", err)
		}

		customizationSection := &govcdtypes.GuestCustomizationSection{}

		// Process parameters from 'customization' block
		updateCustomizationSection(v, customizationSection)

		isVirtualCPUType64 := strings.Contains(v.Plan.OsType.ValueString(), "64")
		virtualCPUType := "VM32"
		if isVirtualCPUType64 {
			virtualCPUType = "VM64"
		}

		if !v.Plan.BootImageID.IsNull() && !v.Plan.BootImageID.IsUnknown() {
			var (
				bootMediaIdentifier string
				mediaRecord         *govcd.MediaRecord
				err                 error
			)

			bootMediaIdentifier = v.Plan.BootImageID.ValueString()
			mediaRecord, err = v.Client.Vmware.QueryMediaById(bootMediaIdentifier)
			if err != nil {
				return nil, fmt.Errorf("[VM creation] error getting boot image %s: %w", bootMediaIdentifier, err)
			}

			// This workaround is to check that the Media file is synchronized in catalog, even if it isn't an iso
			// file. It's not officially documented that IsIso==true means that, but it's the only way we have at the moment.
			if !mediaRecord.MediaRecord.IsIso {
				return nil, fmt.Errorf("[VM creation] error getting boot image %s: Media is not synchronized in the catalog", bootMediaIdentifier)
			}

			bootImage = &govcdtypes.Media{HREF: mediaRecord.MediaRecord.HREF}
		}

		vmParams := &govcdtypes.RecomposeVAppParamsForEmptyVm{
			XmlnsVcloud: govcdtypes.XMLNamespaceVCloud,
			XmlnsOvf:    govcdtypes.XMLNamespaceOVF,
			PowerOn:     false, // Power on is handled at the end of VM creation process
			CreateItem: &govcdtypes.CreateItem{
				Name: v.Plan.VMName.ValueString(),
				// BUG in VCD, do not allow empty NetworkConnectionSection, so we pass simplest
				// network configuration and after VM created update with real config
				NetworkConnectionSection: &govcdtypes.NetworkConnectionSection{
					PrimaryNetworkConnectionIndex: 0,
					NetworkConnection: []*govcdtypes.NetworkConnection{
						{Network: "none", NetworkConnectionIndex: 0, IPAddress: "any", IsConnected: false, IPAddressAllocationMode: "NONE"},
					},
				},
				StorageProfile:            storageProfilePtr,
				ComputePolicy:             vmComputePolicy,
				Description:               v.Plan.Description.ValueString(),
				GuestCustomizationSection: customizationSection,
				VmSpecSection: &govcdtypes.VmSpecSection{
					Modified: utils.TakeBoolPointer(true),
					Info:     "Virtual Machine specification",
					OsType:   v.Plan.OsType.ValueString(),
					// can be created with resource internal_disk
					DiskSection:    &govcdtypes.DiskSection{DiskSettings: []*govcdtypes.DiskSettings{}},
					VirtualCpuType: virtualCPUType,
				},
				BootImage: bootImage,
			},
		}

		vm, err = vapp.AddEmptyVm(vmParams)
		if err != nil {
			return nil, fmt.Errorf("[VM creation] error creating VM %s : %w", v.Plan.VMName.ValueString(), err)
		}
	}

	// If a MAC address is specified for NIC - it does not get set with initial create call therefore
	// running additional update call to make sure it is set correctly

	err = vm.UpdateNetworkConnectionSection(&networkConnectionSection)
	if err != nil {
		return nil, fmt.Errorf("unable to setup network configuration for empty VM %w", err)
	}

	// Refresh VM to have the latest structure
	if err := vm.Refresh(); err != nil {
		return nil, fmt.Errorf("error refreshing VM %s : %w", v.Plan.VMName.ValueString(), err)
	}

	// OS Type is set only for VMs created from template
	// * os_type
	err = updateOsType(v, vm)
	if err != nil {
		return nil, fmt.Errorf("error updating hardware version and OS type : %w", err)
	}

	if err := vm.Refresh(); err != nil {
		return nil, fmt.Errorf("error refreshing VM %s : %w", v.Plan.VMName.ValueString(), err)
	}

	// Template VMs require CPU/Memory setting
	// Lookup CPU values either from schema or from sizing policy. If nothing is set - it will be
	// inherited from template
	var cpuCores, cpuCoresPerSocket *int
	var memory *int64
	if sizingPolicy != nil {
		cpuCores, cpuCoresPerSocket, memory, err = getCPUMemoryValues(v, sizingPolicy.VdcComputePolicyV2)
	} else {
		cpuCores, cpuCoresPerSocket, memory, err = getCPUMemoryValues(v, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting CPU/Memory compute values: %w", err)
	}

	if cpuCores != nil && cpuCoresPerSocket != nil {
		err = vm.ChangeCPUAndCoreCount(cpuCores, cpuCoresPerSocket)
		if err != nil {
			return nil, fmt.Errorf("error changing CPU settings: %w", err)
		}

		if err := vm.Refresh(); err != nil {
			return nil, fmt.Errorf("error refreshing VM %s : %w", v.Plan.VMName.ValueString(), err)
		}
	}

	if memory != nil {
		err = vm.ChangeMemory(*memory)
		if err != nil {
			return nil, fmt.Errorf("error setting memory size from schema for VM from template: %w", err)
		}

		if err := vm.Refresh(); err != nil {
			return nil, fmt.Errorf("error refreshing VM %s : %w", v.Plan.VMName.ValueString(), err)
		}
	}

	return vm, nil
}
