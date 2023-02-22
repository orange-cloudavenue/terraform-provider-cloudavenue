package vm

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type VMClient struct {
	Client *client.CloudAvenue
	Plan   *vmResourceModel
	State  *vmResourceModel
}

const vmUnknownStatus = "-unknown-status-"

var errRemoveResource = errors.New("resource is being removed")

func addRemoveGuestProperties(v *VMClient, vm *govcd.VM) error {
	if !v.Plan.GuestProperties.IsNull() || !v.Plan.GuestProperties.IsUnknown() {
		vmProperties, err := getGuestProperties(v.Plan.GuestProperties)
		if err != nil {
			return fmt.Errorf("unable to convert guest properties to data structure")
		}

		_, err = vm.SetProductSectionList(vmProperties)
		if err != nil {
			return fmt.Errorf("error setting guest properties: %s", err)
		}
	}
	return nil
}

// getGuestProperties returns a struct for setting guest properties
func getGuestProperties(guestProperties types.Map) (*govcdtypes.ProductSectionList, error) {
	vmProperties := &govcdtypes.ProductSectionList{
		ProductSection: &govcdtypes.ProductSection{
			Info:     "Custom properties",
			Property: []*govcdtypes.Property{},
		},
	}
	for key, value := range guestProperties.Elements() {
		log.Printf("[TRACE] Adding guest property: key=%s, value=%s to object", key, value)
		oneProp := &govcdtypes.Property{
			UserConfigurable: true,
			Type:             "string",
			Key:              key,
			Label:            key,
			Value:            &govcdtypes.Value{Value: value.String()},
		}
		vmProperties.ProductSection.Property = append(vmProperties.ProductSection.Property, oneProp)
	}

	return vmProperties, nil
}

// updateGuestCustomizationSetting is responsible for setting all the data related to VM customization
func updateGuestCustomizationSetting(v *VMClient, vm *govcd.VM) error {
	// Retrieve existing customization section to only customize what was throughout this function
	customizationSection, err := vm.GetGuestCustomizationSection()
	if err != nil {
		return fmt.Errorf("error getting existing customization section before changing: %s", err)
	}

	// for back compatibility we allow to set computer name from `name` if computer_name isn't provided
	var computerName string

	if !v.Plan.ComputerName.IsNull() && !v.Plan.ComputerName.IsUnknown() {
		customizationSection.ComputerName = computerName
	}

	// Process parameters from 'customization' block
	updateCustomizationSection(v, customizationSection)

	// Apply any of the settings we have set
	if _, err = vm.SetGuestCustomizationSection(customizationSection); err != nil {
		return fmt.Errorf("error applying guest customization details: %s", err)
	}

	return nil
}

func updateCustomizationSection(v *VMClient, customizationSection *govcdtypes.GuestCustomizationSection) {
	customizationSlice := v.Plan.Customization
	if len(customizationSlice) == 1 {
		customization := customizationSlice[0]

		if !customization.Enabled.IsNull() && !customization.Enabled.IsUnknown() {
			customizationSection.Enabled = utils.TakeBoolPointer(customization.Enabled.ValueBool())
		}
		// init script
		if !customization.InitScript.IsNull() && !customization.InitScript.IsUnknown() {
			customizationSection.CustomizationScript = customization.InitScript.ValueString()
		}
		// change SID
		if !customization.ChangeSid.IsNull() && !customization.ChangeSid.IsUnknown() {
			customizationSection.ChangeSid = utils.TakeBoolPointer(customization.ChangeSid.ValueBool())
		}
		// allow local admin password
		if !customization.AllowLocalAdminPassword.IsNull() && !customization.AllowLocalAdminPassword.IsUnknown() {
			customizationSection.AdminPasswordEnabled = utils.TakeBoolPointer(customization.AllowLocalAdminPassword.ValueBool())
		}
		// must change password on first login
		if !customization.MustChangePasswordOnFirstLogin.IsNull() && !customization.MustChangePasswordOnFirstLogin.IsUnknown() {
			customizationSection.ResetPasswordRequired = utils.TakeBoolPointer(customization.MustChangePasswordOnFirstLogin.ValueBool())
		}
		// auto generate password
		if !customization.AutoGeneratePassword.IsNull() && !customization.AutoGeneratePassword.IsUnknown() {
			customizationSection.AdminPasswordAuto = utils.TakeBoolPointer(customization.AutoGeneratePassword.ValueBool())
		}
		// admin password
		if !customization.AdminPassword.IsNull() && !customization.AdminPassword.IsUnknown() {
			customizationSection.AdminPassword = customization.AdminPassword.ValueString()
		}
		// number of auto logins
		if !customization.NumberOfAutoLogons.IsNull() && !customization.NumberOfAutoLogons.IsUnknown() {
			// The AdminAutoLogonEnabled is "hidden" from direct user input to behave exactly like UI does. UI sets
			// the value of this field behind the scenes based on number_of_auto_logons count.
			// AdminAutoLogonEnabled=false if number_of_auto_logons == 0
			// AdminAutoLogonEnabled=true if number_of_auto_logons > 0
			customizationSection.AdminAutoLogonEnabled = utils.TakeBoolPointer(customization.NumberOfAutoLogons.ValueInt64() > 0)
			customizationSection.AdminAutoLogonCount = int(customization.NumberOfAutoLogons.ValueInt64())
		}
		// join domain
		if !customization.JoinDomain.IsNull() && !customization.JoinDomain.IsUnknown() {
			customizationSection.JoinDomainEnabled = utils.TakeBoolPointer(customization.JoinDomain.ValueBool())
		}
		// join org domain
		if !customization.JoinOrgDomain.IsNull() && !customization.JoinOrgDomain.IsUnknown() {
			customizationSection.UseOrgSettings = utils.TakeBoolPointer(customization.JoinOrgDomain.ValueBool())
		}
		// domain name
		if !customization.JoinDomainName.IsNull() && !customization.JoinDomainName.IsUnknown() {
			customizationSection.DomainName = customization.JoinDomainName.ValueString()
		}
		// domain user
		if !customization.JoinDomainUser.IsNull() && !customization.JoinDomainUser.IsUnknown() {
			customizationSection.DomainUserName = customization.JoinDomainUser.ValueString()
		}
		// domain domain password
		if !customization.JoinDomainPassword.IsNull() && !customization.JoinDomainPassword.IsUnknown() {
			customizationSection.DomainUserPassword = customization.JoinDomainPassword.ValueString()
		}
		// domain account ou
		if !customization.JoinDomainAccountOU.IsNull() && !customization.JoinDomainAccountOU.IsUnknown() {
			customizationSection.MachineObjectOU = customization.JoinDomainAccountOU.ValueString()
		}
	}
}

// isForcedCustomization checks "customization" block in resource and checks if the value of field "force"
// is set to "true". It returns false if the value is not set or is set to false
func isForcedCustomization(v *VMClient) bool {
	if len(v.Plan.Customization) != 1 {
		return false
	}

	cust := v.Plan.Customization[0]

	if !cust.Force.IsNull() && !cust.Force.IsUnknown() {
		return cust.Force.ValueBool()
	} else {
		return false
	}
}

// lookupvAppTemplateforVm will do the following
// evaluate if optional parameter `vm_name_in_template` was specified.
//
// If `vm_name_in_template` was specified
// * It will look up the exact VM with given `vm_name_in_template` inside `vapp_template_id`
//
// If `vm_name_in_template` was not specified:
// * Return error
func lookupvAppTemplateforVM(v *VMClient, org *govcd.Org, vdc *govcd.Vdc) (govcd.VAppTemplate, error) {
	if !v.Plan.VappTemplateID.IsNull() && !v.Plan.VappName.IsUnknown() {
		// Lookup of vApp Template using URN

		vAppTemplate, err := v.Client.Vmware.GetVAppTemplateById(v.Plan.VappTemplateID.ValueString())
		if err != nil {
			return govcd.VAppTemplate{}, fmt.Errorf("error finding vApp Template with URN %s: %s", v.Plan.VappTemplateID.ValueString(), err)
		}

		if !v.Plan.VMNameInTemplate.IsNull() && !v.Plan.VMNameInTemplate.IsUnknown() {
			vmInTemplateRecord, err := v.Client.Vmware.QuerySynchronizedVmInVAppTemplateByHref(vAppTemplate.VAppTemplate.HREF, v.Plan.VMNameInTemplate.ValueString())
			if err != nil {
				return govcd.VAppTemplate{}, fmt.Errorf("error obtaining VM '%s' inside vApp Template: %s", v.Plan.VMNameInTemplate.ValueString(), err)
			}

			returnedVAppTemplate, err := v.Client.Vmware.GetVAppTemplateByHref(vmInTemplateRecord.HREF)
			if err != nil {
				return govcd.VAppTemplate{}, fmt.Errorf("error getting vApp template from inner VM %s: %s", v.Plan.VMNameInTemplate.ValueString(), err)
			}

			return *returnedVAppTemplate, err
		} else {
			// If no VM name was specified - we will pick the first VM inside the vApp Template

			if vAppTemplate.VAppTemplate == nil || vAppTemplate.VAppTemplate.Children == nil || len(vAppTemplate.VAppTemplate.Children.VM) == 0 {
				return govcd.VAppTemplate{}, fmt.Errorf("vApp Template %s does not contain any VMs", v.Plan.VappTemplateID.ValueString())
			}
			returnedVAppTemplate := govcd.NewVAppTemplate(&v.Client.Vmware.Client)
			returnedVAppTemplate.VAppTemplate = vAppTemplate.VAppTemplate.Children.VM[0]
			return *returnedVAppTemplate, nil
		}
	} else {
		return govcd.VAppTemplate{}, fmt.Errorf("vApp Template ID is not specified")
	}
}

// networksToConfig converts terraform schema for 'network' to types.NetworkConnectionSection
// which is used for creating new VM
//
// The `vapp` parameter does not play critical role in the code, but adds additional validations:
// * `org` type of networks will be checked if they are already attached to the vApp
// * `vapp` type networks will be checked for existence inside the vApp
func networksToConfig(v *VMClient, vapp *govcd.VApp) (govcdtypes.NetworkConnectionSection, error) {
	networkConnectionSection := govcdtypes.NetworkConnectionSection{}

	// sets existing primary network connection index. Further code changes index only if change is
	// found
	for index, singleNetwork := range v.Plan.Networks {
		if singleNetwork.IsPrimary.ValueBool() {
			networkConnectionSection.PrimaryNetworkConnectionIndex = index
		}
	}

	for index, singleNetwork := range v.Plan.Networks {
		netConn := &govcdtypes.NetworkConnection{}

		networkName := singleNetwork.Name.ValueString()
		ipAllocationMode := singleNetwork.IPAllocationMode.ValueString()
		ip := singleNetwork.IP.ValueString()

		if v.State != nil && !v.Plan.Networks[index].IsPrimary.Equal(v.State.Networks[index].IsPrimary) && singleNetwork.IsPrimary.ValueBool() {
			networkConnectionSection.PrimaryNetworkConnectionIndex = index
		}

		switch singleNetwork.Type.ValueString() {
		case "org":
			isVappOrgNetwork, err := isItVappOrgNetwork(networkName, *vapp)
			if err != nil {
				return govcdtypes.NetworkConnectionSection{}, err
			}
			if !isVappOrgNetwork {
				return govcdtypes.NetworkConnectionSection{}, fmt.Errorf("vApp Org network : %s is not found", networkName)
			}
		case "vapp":
			isVappNetwork, err := isItVappNetwork(networkName, *vapp)
			if err != nil {
				return govcdtypes.NetworkConnectionSection{}, fmt.Errorf("unable to find vApp network %s: %s", networkName, err)
			}
			if !isVappNetwork {
				return govcdtypes.NetworkConnectionSection{}, fmt.Errorf("vApp network : %s is not found", networkName)
			}
		}

		netConn.IsConnected = singleNetwork.Connected.ValueBool()
		netConn.IPAddressAllocationMode = ipAllocationMode
		netConn.NetworkConnectionIndex = index
		netConn.Network = networkName

		if !singleNetwork.Mac.IsNull() && !singleNetwork.Mac.IsUnknown() {
			netConn.MACAddress = singleNetwork.Mac.ValueString()
		}

		if ipAllocationMode == govcdtypes.IPAllocationModeNone {
			netConn.Network = govcdtypes.NoneNetwork
		}

		if net.ParseIP(ip) != nil {
			netConn.IPAddress = ip
		}

		if !singleNetwork.AdapterType.IsNull() && !singleNetwork.AdapterType.IsUnknown() {
			netConn.NetworkAdapterType = singleNetwork.AdapterType.ValueString()
		}

		networkConnectionSection.NetworkConnection = append(networkConnectionSection.NetworkConnection, netConn)
	}
	return networkConnectionSection, nil
}

// isItVappOrgNetwork checks if it is a vApp Org network (not vApp Network)
func isItVappOrgNetwork(vAppNetworkName string, vapp govcd.VApp) (bool, error) {
	vAppNetworkConfig, err := vapp.GetNetworkConfig()
	if err != nil {
		return false, fmt.Errorf("error getting vApp networks: %s", err)
	}

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.NetworkName == vAppNetworkName &&
			!govcd.IsVappNetwork(networkConfig.Configuration) {
			return true, nil
		}
	}

	return false, fmt.Errorf("configured vApp Org network isn't found: %s", vAppNetworkName)
}

// isItVappNetwork checks if it is a vApp network (not vApp Org Network)
func isItVappNetwork(vAppNetworkName string, vapp govcd.VApp) (bool, error) {
	vAppNetworkConfig, err := vapp.GetNetworkConfig()
	if err != nil {
		return false, fmt.Errorf("error getting vApp networks: %s", err)
	}

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.NetworkName == vAppNetworkName &&
			govcd.IsVappNetwork(networkConfig.Configuration) {
			return true, nil
		}
	}

	return false, fmt.Errorf("configured vApp network isn't found: %s", vAppNetworkName)
}

func lookupStorageProfile(storageProfileName string, vdc *govcd.Vdc) (*govcdtypes.Reference, error) {
	// If no storage profile lookup was requested - bail out early and return nil reference
	if storageProfileName == "" {
		return nil, nil
	}

	storageProfile, err := vdc.FindStorageProfileReference(storageProfileName)
	if err != nil {
		return nil, fmt.Errorf("[vm creation] error retrieving storage profile %s : %s", storageProfileName, err)
	}

	return &storageProfile, nil
}

// lookupComputePolicy returns the Compute Policy associated to the value of the given Compute Policy attribute. If the
// attribute is not set, the returned policy will be nil. If the obtained policy is incorrect, it will return an error.
func lookupComputePolicy(v *VMClient, value string) (*govcd.VdcComputePolicyV2, error) {
	if value == "" {
		return nil, nil
	}

	computePolicy, err := v.Client.Vmware.GetVdcComputePolicyV2ById(value)
	if err != nil {
		return nil, fmt.Errorf("error getting compute policy %s: %s", value, err)
	}
	if computePolicy.Href == "" {
		return nil, fmt.Errorf("empty compute policy HREF detected")
	}
	return computePolicy, nil
}

// lookupComputePolicySizingPolicyID returns the Compute Policy Sizing Policy ID associated to the value of the given
// Compute Policy Sizing Policy ID attribute. If the attribute is not set, the returned policy will be nil. If the
// obtained policy is incorrect, it will return an error.
// func lookupComputePolicySizingPolicyID(v *VMClient, sizingPolicyID string) (*govcd.VdcComputePolicyV2, error) {
// 	if sizingPolicyID == "" {
// 		return nil, nil
// 	}

// 	return lookupComputePolicy(v, sizingPolicyID)
// }

// lookupComputePolicyPlacementPolicyID returns the Compute Policy Placement Policy ID associated to the value of the
// given Compute Policy Placement Policy ID attribute. If the attribute is not set, the returned policy will be nil.
// If the obtained policy is incorrect, it will return an error.
// func lookupComputePolicyPlacementPolicyID(v *VMClient, placementPolicyID string) (*govcd.VdcComputePolicyV2, error) {
// 	if placementPolicyID == "" {
// 		return nil, nil
// 	}

// 	return lookupComputePolicy(v, placementPolicyID)
// }

func updateTemplateInternalDisks(v *VMClient, vm govcd.VM) error {
	// Get vcd object
	_, vdc, err := v.Client.GetOrgAndVDC(v.Client.GetOrg(), v.Plan.VDC.ValueString())
	if err != nil {
		return fmt.Errorf("error retrieving VDC %s: %s", v.Plan.VDC.ValueString(), err)
	}

	if vm.VM.VmSpecSection == nil || vm.VM.VmSpecSection.DiskSection == nil {
		return fmt.Errorf("[updateTemplateInternalDisks] VmSpecSection part is missing")
	}

	diskSettings := vm.VM.VmSpecSection.DiskSection.DiskSettings

	var storageProfilePrt *govcdtypes.Reference
	var overrideVMDefault bool

	if len(v.Plan.OverrideTemplateDisks) == 0 {
		return nil
	}

	for _, internalDiskProvidedConfig := range v.Plan.OverrideTemplateDisks {
		diskCreatedByTemplate := getMatchedDisk(internalDiskProvidedConfig, diskSettings)

		storageProfileName := internalDiskProvidedConfig.StorageProfile.ValueString()
		if storageProfileName != "" {
			storageProfile, err := vdc.FindStorageProfileReference(storageProfileName)
			if err != nil {
				return fmt.Errorf("[vm creation] error retrieving storage profile %s : %s", storageProfileName, err)
			}
			storageProfilePrt = &storageProfile
			overrideVMDefault = true
		} else {
			storageProfilePrt = vm.VM.StorageProfile
			overrideVMDefault = false
		}

		if diskCreatedByTemplate == nil {
			return fmt.Errorf("[vm creation] disk with bus type %s, bus number %d and unit number %d not found",
				internalDiskProvidedConfig.BusType.ValueString(), internalDiskProvidedConfig.BusNumber.ValueInt64(), internalDiskProvidedConfig.UnitNumber.ValueInt64())
		}

		// Update details of internal disk for disk existing in template
		if !internalDiskProvidedConfig.Iops.IsNull() && !internalDiskProvidedConfig.Iops.IsUnknown() {
			diskCreatedByTemplate.Iops = utils.TakeInt64Pointer(internalDiskProvidedConfig.Iops.ValueInt64())
		}

		// value is required but not treated.
		isThinProvisioned := true
		diskCreatedByTemplate.ThinProvisioned = &isThinProvisioned

		diskCreatedByTemplate.SizeMb = internalDiskProvidedConfig.SizeInMb.ValueInt64()
		diskCreatedByTemplate.StorageProfile = storageProfilePrt
		diskCreatedByTemplate.OverrideVmDefault = overrideVMDefault
	}

	vmSpecSection := vm.VM.VmSpecSection
	vmSpecSection.DiskSection.DiskSettings = diskSettings
	_, err = vm.UpdateInternalDisks(vmSpecSection)
	if err != nil {
		return fmt.Errorf("error updating VM disks: %s", err)
	}

	return nil
}

// getMatchedDisk returns matched disk by adapter type, bus number and unit number
func getMatchedDisk(internalDiskProvidedConfig vm.TemplateDiskModel, diskSettings []*govcdtypes.DiskSettings) *govcdtypes.DiskSettings {
	for _, diskSetting := range diskSettings {
		if diskSetting.AdapterType == vm.InternalDiskBusTypes[internalDiskProvidedConfig.BusType.ValueString()] &&
			diskSetting.BusNumber == int(internalDiskProvidedConfig.BusNumber.ValueInt64()) &&
			diskSetting.UnitNumber == int(internalDiskProvidedConfig.UnitNumber.ValueInt64()) {
			return diskSetting
		}
	}
	return nil
}

func updateOsType(v *VMClient, vm *govcd.VM) error {
	var err error

	vmSpecSection := vm.VM.VmSpecSection

	if !v.Plan.OsType.IsNull() && !v.Plan.OsType.IsUnknown() {
		vmSpecSection.OsType = v.Plan.OsType.ValueString()
		_, err = vm.UpdateVmSpecSection(vmSpecSection, v.Plan.Description.ValueString())
		if err != nil {
			return fmt.Errorf("error changing VM spec section: %s", err)
		}
	}

	return nil
}

// getCpuMemoryValues returns CPU, CPU core count and Memory variables. Priority comes from HCL
// schema configuration and then whatever is present in compute policy (if it was specified at all)
func getCPUMemoryValues(v *VMClient, vdcComputePolicy *govcdtypes.VdcComputePolicyV2) (cpu, cores *int, memory *int64, err error) {
	var (
		setCPU    int
		setCores  int
		setMemory int64
	)

	if !v.Plan.Resource.Memory.IsNull() && !v.Plan.Resource.Memory.IsUnknown() {
		setMemory = v.Plan.Resource.Memory.ValueInt64()
	}

	if !v.Plan.Resource.CPUs.IsNull() && !v.Plan.Resource.CPUs.IsUnknown() {
		setCPU = int(v.Plan.Resource.CPUs.ValueInt64())
	}

	if !v.Plan.Resource.CPUCores.IsNull() && !v.Plan.Resource.CPUCores.IsUnknown() {
		setCores = int(v.Plan.Resource.CPUCores.ValueInt64())
	}

	// Check if sizing policy has any settings settings and override VM configuration with it
	if vdcComputePolicy != nil {
		if vdcComputePolicy.Memory != nil {
			mem := int64(*vdcComputePolicy.Memory)
			setMemory = mem
		}

		if vdcComputePolicy.CPUCount != nil {
			setCPU = *vdcComputePolicy.CPUCount
		}

		if vdcComputePolicy.CoresPerSocket != nil {
			setCores = *vdcComputePolicy.CoresPerSocket
		}
	}

	return &setCPU, &setCores, &setMemory, nil
}

// attachDetachIndependentDisks updates attached disks to latest state, removes not needed, and adds
// new ones
func attachDetachIndependentDisks(v *VMClient, gvcdvm govcd.VM, vdc *govcd.Vdc) error {
	if reflect.DeepEqual(v.Plan.Disks, v.State.Disks) {
		// No changes in disks
		return nil
	}

	var (
		attachDisks []vm.DiskModel
		detachDisks []vm.DiskModel
	)

	for i := range v.Plan.Disks {
		if reflect.DeepEqual(v.Plan.Disks[i], v.State.Disks[i]) {
			// No changes in disk
			continue
		} else {
			// Disk changed

			// Determine if disk was added or removed
			if v.State.Disks[i].Name.IsNull() || v.State.Disks[i].Name.IsUnknown() {
				// Disk was added
				attachDisks = append(attachDisks, v.Plan.Disks[i])
			} else {
				// Disk was removed
				detachDisks = append(detachDisks, v.State.Disks[i])
			}
		}
	}

	for _, diskData := range detachDisks {
		disk, err := vdc.QueryDisk(diskData.Name.ValueString())
		if err != nil {
			return fmt.Errorf("did not find disk `%s`: %s", diskData.Name.ValueString(), err)
		}

		attachParams := &govcdtypes.DiskAttachOrDetachParams{
			Disk:       &govcdtypes.Reference{HREF: disk.Disk.HREF},
			UnitNumber: utils.TakeIntPointer(int(diskData.UnitNumber.ValueInt64())),
			BusNumber:  utils.TakeIntPointer(int(diskData.BusNumber.ValueInt64())),
		}

		task, err := gvcdvm.DetachDisk(attachParams)
		if err != nil {
			return fmt.Errorf("error detaching disk `%s` to vm %s", diskData.Name.ValueString(), err)
		}
		err = task.WaitTaskCompletion()
		if err != nil {
			return fmt.Errorf("error waiting for task to complete detaching disk `%s` to vm %s", diskData.Name.ValueString(), err)
		}
	}

	sort.SliceStable(attachDisks, func(i, j int) bool {
		if attachDisks[i].BusNumber.Equal(attachDisks[j].BusNumber) {
			return attachDisks[i].UnitNumber.ValueInt64() > attachDisks[j].UnitNumber.ValueInt64()
		}
		return attachDisks[i].BusNumber.ValueInt64() > attachDisks[j].BusNumber.ValueInt64()
	})

	for _, diskData := range attachDisks {
		disk, err := vdc.QueryDisk(diskData.Name.ValueString())
		if err != nil {
			return fmt.Errorf("did not find disk `%s`: %s", diskData.Name.ValueString(), err)
		}

		attachParams := &govcdtypes.DiskAttachOrDetachParams{
			Disk:       &govcdtypes.Reference{HREF: disk.Disk.HREF},
			UnitNumber: utils.TakeIntPointer(int(diskData.UnitNumber.ValueInt64())),
			BusNumber:  utils.TakeIntPointer(int(diskData.BusNumber.ValueInt64())),
		}

		task, err := gvcdvm.AttachDisk(attachParams)
		if err != nil {
			return fmt.Errorf("error attaching disk `%s` to vm %s", diskData.Name.ValueString(), err)
		}
		err = task.WaitTaskCompletion()
		if err != nil {
			return fmt.Errorf("error waiting for task to complete attaching disk `%s` to vm %s", diskData.Name.ValueString(), err)
		}
	}
	return nil
}
