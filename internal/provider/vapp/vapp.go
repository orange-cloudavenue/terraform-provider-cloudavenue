// Package vapp provides a Terraform resource to manage vApps.
package vapp

// createVmFromTemplate is responsible for create VMs from template of two types:
// * Standalone VMs
// * VMs inside vApp (vApp VMs)
//
// Code flow has 3 layers:
// 1. Lookup common information, required for both types of VMs (Standalone and vApp child). Things such as
//   - Template to be used
//   - Network adapter configuration
//   - Storage profile configuration
//   - VM compute policy configuration
//
// 2. Perform VM creation operation based on type in separate switch/case
//   - standaloneVmType
//   - vAppVmType
//
// # This part includes defining initial structures for VM and also any explicitly required operations for that type of VM
//
// 3. Perform additional operations which are common for both types of VMs
//
// Note. VM Power ON (if it wasn't disabled in HCL configuration) occurs as last step after all configuration is done.
// func createVmFromTemplate(ctx context.Context, v *VappClient) (vm *govcd.VM, err error) {

// 	// If VDC is not defined at data source level, use the one defined at provider level
// 	if v.Plan.Vdc.IsNull() || v.Plan.Vdc.IsUnknown() {
// 		if v.Client.DefaultVdcExist() {
// 			v.Plan.Vdc = types.StringValue(v.Client.GetDefaultVdc())
// 		} else {
// 			err = errors.New("VDC is required when not defined at provider level")
// 			return nil, err
// 		}
// 	}

// 	// Step 1 - lookup common information
// 	org, vdc, err := v.Client.GetOrgAndVdc(v.Client.GetOrg(), v.Plan.Vdc.ValueString())
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving Org and VDC: %s", err)
// 	}

// 	// Look up VM template inside vApp template - either specified by `vm_name_in_template` or the
// 	// first one in vApp
// 	vmTemplate, err := lookupvAppTemplateforVm(v, org, vdc)
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding vApp template: %s", err)
// 	}

// 	// Look up vApp before setting up network configuration. Having a vApp set, will enable
// 	// additional network availability in vApp validations in `networksToConfig` function.
// 	// It is only possible for vApp VMs, as empty VMs will get their hidden vApps created after the
// 	// VM is created.
// 	var vapp *govcd.VApp
// 	if v.VmType == vappVmType {
// 		vapp, err = vdc.GetVAppByName(v.Plan.VappName.ValueString(), false)
// 		if err != nil {
// 			return nil, fmt.Errorf("[VM create] error finding vApp %s: %s", v.Plan.VappName.ValueString(), err)
// 		}
// 	}

// 	// Build up network configuration
// 	networkConnectionSection, err := networksToConfig(v, vapp)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to process network configuration: %s", err)
// 	}
// 	tflog.Info(ctx, fmt.Sprintf("[VM create] networkConnectionSection %# v", pretty.Formatter(networkConnectionSection)))

// 	// Lookup storage profile reference if it was specified
// 	storageProfilePtr, err := lookupStorageProfile(v.Plan.StorageProfile.ValueString(), vdc)
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding storage profile: %s", err)
// 	}

// 	// Look up compute policies
// 	sizingPolicy, err := lookupComputePolicy(v, v.Plan.SizingPolicyID.ValueString())
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding sizing policy: %s", err)
// 	}
// 	placementPolicy, err := lookupComputePolicy(v, v.Plan.PlacementPolicyID.ValueString())
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding placement policy: %s", err)
// 	}
// 	var vmComputePolicy *govcdtypes.ComputePolicy
// 	if sizingPolicy != nil || placementPolicy != nil {
// 		vmComputePolicy = &govcdtypes.ComputePolicy{}
// 		if sizingPolicy != nil {
// 			vmComputePolicy.VmSizingPolicy = &govcdtypes.Reference{HREF: sizingPolicy.Href}
// 		}
// 		if placementPolicy != nil {
// 			vmComputePolicy.VmPlacementPolicy = &govcdtypes.Reference{HREF: placementPolicy.Href}
// 		}
// 	}

// 	// Step 2 - perform VM creation operation based on type
// 	// VM creation uses different structure depending on if it is a standaloneVmType or vappVmType
// 	// These structures differ and one might accept all required parameters, while other
// 	switch v.VmType {
// 	case standaloneVmType:
// 		standaloneVmParams := govcdtypes.InstantiateVmTemplateParams{
// 			Xmlns:            govcdtypes.XMLNamespaceVCloud,
// 			Name:             v.Plan.VMName.ValueString(), // VM name post creation
// 			PowerOn:          false,                       // VM will be powered on after all configuration is done
// 			AllEULAsAccepted: v.Plan.AcceptAllEulas.ValueBool(),
// 			ComputePolicy:    vmComputePolicy,
// 			SourcedVmTemplateItem: &govcdtypes.SourcedVmTemplateParams{
// 				Source: &govcdtypes.Reference{
// 					HREF: vmTemplate.VAppTemplate.HREF,
// 					ID:   vmTemplate.VAppTemplate.ID,
// 					Type: vmTemplate.VAppTemplate.Type,
// 					Name: vmTemplate.VAppTemplate.Name,
// 				},
// 				VmGeneralParams: &govcdtypes.VMGeneralParams{
// 					Description: v.Plan.Description.ValueString(),
// 				},
// 				VmTemplateInstantiationParams: &govcdtypes.InstantiationParams{
// 					// If a MAC address is specified for NIC - it does not get set with this call,
// 					// therefore an additional `vm.UpdateNetworkConnectionSection` is required.
// 					NetworkConnectionSection: &networkConnectionSection,
// 				},
// 				StorageProfile: storageProfilePtr,
// 			},
// 		}

// 		util.Logger.Printf("%# v", pretty.Formatter(standaloneVmParams))
// 		vm, err = vdc.CreateStandaloneVMFromTemplate(&standaloneVmParams)
// 		if err != nil {
// 			return nil, removeResource
// 		}

// 		// d.SetId(vm.VM.ID)

// 		vapp, err = vm.GetParentVApp()
// 		if err != nil {
// 			return nil, removeResource
// 		}
// 		// util.Logger.Printf("[VM create] vApp after creation %# v", pretty.Formatter(vapp.VApp))
// 		// dSet(d, "vapp_name", vapp.VApp.Name)
// 		// dSet(d, "vm_type", string(standaloneVmType))

// 	////////////////////////////////////////////////////////////////////////////////////////////
// 	// This part of code handles additional VM create operations, which can not be set during
// 	// initial VM creation.
// 	// __Explicitly__ template based Standalone VMs are addressed here.
// 	////////////////////////////////////////////////////////////////////////////////////////////

// 	case vappVmType:
// 		vappName := v.Plan.VappName.ValueString()
// 		vapp, err = vdc.GetVAppByName(vappName, false)
// 		if err != nil {
// 			return nil, fmt.Errorf("[VM create] error finding vApp %s: %s", vappName, err)
// 		}

// 		vappVmParams := &govcdtypes.ReComposeVAppParams{
// 			Ovf:              govcdtypes.XMLNamespaceOVF,
// 			Xsi:              govcdtypes.XMLNamespaceXSI,
// 			Xmlns:            govcdtypes.XMLNamespaceVCloud,
// 			AllEULAsAccepted: v.Plan.AcceptAllEulas.ValueBool(),
// 			Name:             vapp.VApp.Name,
// 			PowerOn:          false, // VM will be powered on after all configuration is done
// 			SourcedItem: &govcdtypes.SourcedCompositionItemParam{
// 				Source: &govcdtypes.Reference{
// 					HREF: vmTemplate.VAppTemplate.HREF,
// 					Name: v.Plan.VMName.ValueString(), // This VM name defines the VM name after creation
// 				},
// 				VMGeneralParams: &govcdtypes.VMGeneralParams{
// 					Description: v.Plan.Description.ValueString(),
// 				},
// 				InstantiationParams: &govcdtypes.InstantiationParams{
// 					// If a MAC address is specified for NIC - it does not get set with this call,
// 					// therefore an additional `vm.UpdateNetworkConnectionSection` is required.
// 					NetworkConnectionSection: &networkConnectionSection,
// 				},
// 				ComputePolicy:  vmComputePolicy,
// 				StorageProfile: storageProfilePtr,
// 			},
// 		}

// 		vm, err = vapp.AddRawVM(vappVmParams)
// 		if err != nil {
// 			return nil, removeResource
// 		}

// 		// d.SetId(vm.VM.ID)
// 		// dSet(d, "vm_type", string(vappVmType))

// 	////////////////////////////////////////////////////////////////////////////////////////////
// 	// This part of code handles additional VM create operations, which can not be set during
// 	// initial VM creation.
// 	// __Explicitly__ template based vApp VMs are addressed here.
// 	////////////////////////////////////////////////////////////////////////////////////////////

// 	default:
// 		return nil, fmt.Errorf("unknown VM type %s", v.VmType)
// 	}

// 	////////////////////////////////////////////////////////////////////////////////////////////////
// 	// This part of code handles additional VM create operations, which can not be set during
// 	// initial VM creation.
// 	// __Only__ template based VMs are addressed here.
// 	////////////////////////////////////////////////////////////////////////////////////////////////

// 	// If a MAC address is specified for NIC - it does not get set with initial create call therefore
// 	// running additional update call to make sure it is set correctly

// 	err = vm.UpdateNetworkConnectionSection(&networkConnectionSection)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to setup network configuration for empty VM %s", err)
// 	}

// 	// Refresh VM to have the latest structure
// 	if err := vm.Refresh(); err != nil {
// 		return nil, fmt.Errorf("error refreshing VM %s : %s", v.Plan.VMName.ValueString(), err)
// 	}

// 	// TODO Wait remi
// 	// update existing internal disks in template (it is only applicable to VMs created
// 	// Such fields are processed:
// 	// * override_template_disk
// 	err = updateTemplateInternalDisks(d, meta, *vm)
// 	if err != nil {
// 		dSet(d, "override_template_disk", nil)
// 		return nil, fmt.Errorf("error managing internal disks : %s", err)
// 	}
// 	// ! End

// 	if err := vm.Refresh(); err != nil {
// 		return nil, fmt.Errorf("error refreshing VM %s : %s", vmName, err)
// 	}

// 	// OS Type and Hardware version should only be changed if specified. (Only applying to VMs from
// 	// templates as empty VMs require this by default)
// 	// Such fields are processed:
// 	// * os_type
// 	// * hardware_version
// 	err = updateOsType(v, vm)
// 	if err != nil {
// 		return nil, fmt.Errorf("error updating hardware version and OS type : %s", err)
// 	}

// 	if err := vm.Refresh(); err != nil {
// 		return nil, fmt.Errorf("error refreshing VM %s : %s", v.Plan.VMName.ValueString(), err)
// 	}

// 	// Template VMs require CPU/Memory setting
// 	// Lookup CPU values either from schema or from sizing policy. If nothing is set - it will be
// 	// inherited from template
// 	var cpuCores, cpuCoresPerSocket *int
// 	var memory *int64
// 	if sizingPolicy != nil {
// 		cpuCores, cpuCoresPerSocket, memory, err = getCpuMemoryValues(v, sizingPolicy.VdcComputePolicyV2)
// 	} else {
// 		cpuCores, cpuCoresPerSocket, memory, err = getCpuMemoryValues(v, nil)
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting CPU/Memory compute values: %s", err)
// 	}

// 	if cpuCores != nil || cpuCoresPerSocket != nil {
// 		err = vm.ChangeCPUAndCoreCount(cpuCores, cpuCoresPerSocket)
// 		if err != nil {
// 			return nil, fmt.Errorf("error changing CPU settings: %s", err)
// 		}

// 		if err := vm.Refresh(); err != nil {
// 			return nil, fmt.Errorf("error refreshing VM %s : %s", v.Plan.VMName.ValueString(), err)
// 		}
// 	}

// 	if memory != nil {
// 		err = vm.ChangeMemory(*memory)
// 		if err != nil {
// 			return nil, fmt.Errorf("error setting memory size from schema for VM from template: %s", err)
// 		}

// 		if err := vm.Refresh(); err != nil {
// 			return nil, fmt.Errorf("error refreshing VM %s : %s", v.Plan.VMName.ValueString(), err)
// 		}
// 	}

// 	return vm, nil
// }

// createVmEmpty is responsible for creating empty VMs of two types:
// * Standalone VMs
// * VMs inside vApp (vApp VMs)
//
// Code flow has 3 layers:
// 1. Lookup common information, required for both types of VMs (Standalone and vApp child). Things such as
//   - OS Type
//   - Hardware version
//   - Storage profile configuration
//   - VM compute policy configuration
//   - Boot image
//
// 2. Perform VM creation operation based on type in separate switch/case
//   - standaloneVmType
//   - vAppVmType
//
// # This part includes defining initial structures for VM and also any explicitly required operations for that type of VM
//
// 3. Perform additional operations which are common for both types of VMs
//
// Note. VM Power ON (if it wasn't disabled in HCL configuration) occurs as last step after all configuration is done.
// func createVmEmpty(ctx context.Context, v *VappClient) (*govcd.VM, error) {

// 	_, vdc, err := v.Client.GetOrgAndVdc(v.Client.GetOrg(), v.Plan.Vdc.ValueString())
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving Org and VDC: %s", err)
// 	}
// 	var vapp *govcd.VApp

// 	if v.VmType == vappVmType {
// 		vapp, err = vdc.GetVAppByName(v.Plan.VappName.ValueString(), false)
// 		if err != nil {
// 			return nil, fmt.Errorf("[VM create] error finding vApp for empty VM %s: %s", v.Plan.VappName.ValueString(), err)
// 		}
// 	}

// 	var (
// 		ok           bool
// 		osType       string
// 		computerName string
// 	)

// 	if !v.Plan.OsType.IsNull() && !v.Plan.OsType.IsUnknown() {
// 		osType = v.Plan.OsType.ValueString()
// 	} else {
// 		return nil, fmt.Errorf("`os_type` is required when creating empty VM")
// 	}

// 	if !v.Plan.ComputerName.IsNull() && !v.Plan.ComputerName.IsUnknown() {
// 		computerName = v.Plan.ComputerName.ValueString()
// 	} else {
// 		return nil, fmt.Errorf("`computer_name` is required when creating empty VM")
// 	}

// 	var bootImage *govcdtypes.Media
// 	if !v.Plan.BootImageID.IsNull() && !v.Plan.BootImageID.IsUnknown() {

// 		var (
// 			bootMediaIdentifier string
// 			mediaRecord         *govcd.MediaRecord
// 			err                 error
// 		)

// 		bootMediaIdentifier = v.Plan.BootImageID.ValueString()
// 		mediaRecord, err = v.Client.Vmware.QueryMediaById(bootMediaIdentifier)
// 		if err != nil {
// 			return nil, fmt.Errorf("[VM creation] error getting boot image %s: %s", bootMediaIdentifier, err)
// 		}

// 		// This workaround is to check that the Media file is synchronized in catalog, even if it isn't an iso
// 		// file. It's not officially documented that IsIso==true means that, but it's the only way we have at the moment.
// 		if !mediaRecord.MediaRecord.IsIso {
// 			return nil, fmt.Errorf("[VM creation] error getting boot image %s: Media is not synchronized in the catalog", bootMediaIdentifier)
// 		}

// 		bootImage = &govcdtypes.Media{HREF: mediaRecord.MediaRecord.HREF}
// 	}

// 	storageProfilePtr, err := lookupStorageProfile(v.Plan.StorageProfile.ValueString(), vdc)
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding storage profile: %s", err)
// 	}

// 	customizationSection := &govcdtypes.GuestCustomizationSection{}
// 	customizationSection.ComputerName = v.Plan.ComputerName.ValueString()

// 	// Process parameters from 'customization' block
// 	updateCustomizationSection(v, customizationSection)

// 	isVirtualCpuType64 := strings.Contains(v.Plan.OsType.ValueString(), "64")
// 	virtualCpuType := "VM32"
// 	if isVirtualCpuType64 {
// 		virtualCpuType = "VM64"
// 	}

// 	// Look up compute policies
// 	sizingPolicy, err := lookupComputePolicy(v, "sizing_policy_id")
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding sizing policy: %s", err)
// 	}
// 	placementPolicy, err := lookupComputePolicy(v, "placement_policy_id")
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding placement policy: %s", err)
// 	}
// 	var vmComputePolicy *govcdtypes.ComputePolicy
// 	if sizingPolicy != nil || placementPolicy != nil {
// 		vmComputePolicy = &govcdtypes.ComputePolicy{}
// 		if sizingPolicy != nil {
// 			vmComputePolicy.VmSizingPolicy = &govcdtypes.Reference{HREF: sizingPolicy.Href}
// 		}
// 		if placementPolicy != nil {
// 			vmComputePolicy.VmPlacementPolicy = &govcdtypes.Reference{HREF: placementPolicy.Href}
// 		}
// 	}

// 	// Lookup CPU/Memory parameters
// 	var cpuCores, cpuCoresPerSocket *int
// 	var memory *int64
// 	if sizingPolicy != nil {
// 		cpuCores, cpuCoresPerSocket, memory, err = getCpuMemoryValues(v, sizingPolicy.VdcComputePolicyV2)
// 	} else {
// 		cpuCores, cpuCoresPerSocket, memory, err = getCpuMemoryValues(v, nil)
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting CPU/Memory compute values: %s", err)
// 	}
// 	// Wrap memory definition into a suitable construct if it is set
// 	var memoryResourceMb *govcdtypes.MemoryResourceMb
// 	if memory != nil {
// 		memoryResourceMb = &govcdtypes.MemoryResourceMb{Configured: *memory}
// 	}

// 	vmName := v.Plan.VMName.ValueString()
// 	var newVm *govcd.VM

// 	switch v.VmType {
// 	case standaloneVmType:
// 		var mediaReference *govcdtypes.Reference
// 		if bootImage != nil {
// 			mediaReference = &govcdtypes.Reference{
// 				HREF: bootImage.HREF,
// 				ID:   bootImage.ID,
// 				Type: bootImage.Type,
// 				Name: bootImage.Name,
// 			}
// 		}
// 		params := govcdtypes.CreateVmParams{
// 			Xmlns:       govcdtypes.XMLNamespaceVCloud,
// 			Name:        vmName,
// 			PowerOn:     false, // Power on is handled at the end of VM creation process
// 			Description: v.Plan.Description.ValueString(),
// 			CreateVm: &govcdtypes.Vm{
// 				Name:          vmName,
// 				ComputePolicy: vmComputePolicy,
// 				// BUG in VCD, do not allow empty NetworkConnectionSection, so we pass simplest
// 				// network configuration and after VM created update with real config
// 				NetworkConnectionSection: &govcdtypes.NetworkConnectionSection{
// 					PrimaryNetworkConnectionIndex: 0,
// 					NetworkConnection: []*govcdtypes.NetworkConnection{
// 						{Network: "none", NetworkConnectionIndex: 0, IPAddress: "any", IsConnected: false, IPAddressAllocationMode: "NONE"}},
// 				},
// 				VmSpecSection: &govcdtypes.VmSpecSection{
// 					Modified:          utils.TakeBoolPointer(true),
// 					Info:              "Virtual Machine specification",
// 					OsType:            v.Plan.OsType.ValueString(),
// 					CpuResourceMhz:    &govcdtypes.CpuResourceMhz{Configured: 0},
// 					NumCpus:           cpuCores,
// 					NumCoresPerSocket: cpuCoresPerSocket,
// 					MemoryResourceMb:  memoryResourceMb,

// 					// can be created with resource internal_disk
// 					DiskSection:    &govcdtypes.DiskSection{DiskSettings: []*govcdtypes.DiskSettings{}},
// 					VirtualCpuType: virtualCpuType,
// 				},
// 				GuestCustomizationSection: customizationSection,
// 				StorageProfile:            storageProfilePtr,
// 			},
// 			Media: mediaReference,
// 		}

// 		newVm, err = vdc.CreateStandaloneVm(&params)
// 		if err != nil {
// 			return nil, err
// 		}

// 	////////////////////////////////////////////////////////////////////////////////////////////
// 	// This part of code handles additional VM create operations, which can not be set during
// 	// initial VM creation.
// 	// __Explicitly__ empty  Standalone VMs are addressed here.
// 	////////////////////////////////////////////////////////////////////////////////////////////

// 	case vappVmType:
// 		recomposeVAppParamsForEmptyVm := &govcdtypes.RecomposeVAppParamsForEmptyVm{
// 			XmlnsVcloud: govcdtypes.XMLNamespaceVCloud,
// 			XmlnsOvf:    govcdtypes.XMLNamespaceOVF,
// 			PowerOn:     false, // Power on is handled at the end of VM creation process
// 			CreateItem: &govcdtypes.CreateItem{
// 				Name: vmName,
// 				// BUG in VCD, do not allow empty NetworkConnectionSection, so we pass simplest
// 				// network configuration and after VM created update with real config
// 				NetworkConnectionSection: &govcdtypes.NetworkConnectionSection{
// 					PrimaryNetworkConnectionIndex: 0,
// 					NetworkConnection: []*govcdtypes.NetworkConnection{
// 						{Network: "none", NetworkConnectionIndex: 0, IPAddress: "any", IsConnected: false, IPAddressAllocationMode: "NONE"}},
// 				},
// 				StorageProfile:            storageProfilePtr,
// 				ComputePolicy:             vmComputePolicy,
// 				Description:               v.Plan.Description.ValueString(),
// 				GuestCustomizationSection: customizationSection,
// 				VmSpecSection: &govcdtypes.VmSpecSection{
// 					Modified:          utils.TakeBoolPointer(true),
// 					Info:              "Virtual Machine specification",
// 					OsType:            v.Plan.OsType.ValueString(),
// 					NumCpus:           cpuCores,
// 					NumCoresPerSocket: cpuCoresPerSocket,
// 					MemoryResourceMb:  memoryResourceMb,
// 					// can be created with resource internal_disk
// 					DiskSection:    &govcdtypes.DiskSection{DiskSettings: []*govcdtypes.DiskSettings{}},
// 					VirtualCpuType: virtualCpuType,
// 				},
// 				BootImage: bootImage,
// 			},
// 		}

// 		newVm, err = vapp.AddEmptyVm(recomposeVAppParamsForEmptyVm)
// 		if err != nil {
// 			return nil, fmt.Errorf("[VM creation] error creating VM %s : %s", vmName, err)
// 		}

// 	////////////////////////////////////////////////////////////////////////////////////////////
// 	// This part of code handles additional VM create operations, which can not be set during
// 	// initial VM creation.
// 	// __Explicitly__ empty vApp VMs are addressed here.
// 	////////////////////////////////////////////////////////////////////////////////////////////

// 	default:
// 		return nil, fmt.Errorf("unknown VM type %s", v.VmType)
// 	}

// 	////////////////////////////////////////////////////////////////////////////////////////////////
// 	// This part of code handles additional VM create operations, which can not be set during
// 	// initial VM creation.
// 	// __Only__ empty VMs are addressed here.
// 	////////////////////////////////////////////////////////////////////////////////////////////////

// 	vapp, err = newVm.GetParentVApp()
// 	if err != nil {
// 		return nil, fmt.Errorf("[VM creation] error retrieving vApp from standalone VM %s : %s", v.Plan.VMName.ValueString(), err)
// 	}

// 	// Due to the Bug in VCD, VM creation works only with Org VDC networks, not vApp networks - we
// 	// setup network configuration with update.

// 	// firstly cleanup dummy network as network adapter type can't be changed
// 	err = newVm.UpdateNetworkConnectionSection(&govcdtypes.NetworkConnectionSection{})
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to setup network configuration for empty VM %s", err)
// 	}

// 	networkConnectionSection, err := networksToConfig(v, vapp)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to setup network configuration for empty VM: %s", err)
// 	}

// 	// add real network configuration
// 	err = newVm.UpdateNetworkConnectionSection(&networkConnectionSection)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to setup network configuration for empty VM %s", err)
// 	}

// 	return newVm, nil
// }
