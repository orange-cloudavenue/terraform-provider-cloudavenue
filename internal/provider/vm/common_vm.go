package vm

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	commonvm "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type VMClient struct {
	Client *client.CloudAvenue
	Plan   *vmResourceModel
	State  *vmResourceModel
}

const vmUnknownStatus = "-unknown-status-"

var errRemoveResource = errors.New("resource is being removed")

/*
	* Guest properties

	Following code is responsible for setting guest properties on the VM
	3 functions are available:
		- addRemoveGuestProperties
		- getGuestProperties
		- updateGuestCustomizationSetting

*/

// addRemoveGuestProperties is responsible for setting guest properties on the VM.
func addRemoveGuestProperties(v *VMClient, vm *govcd.VM) error {
	// * GuestPropertiers is Optional Value in Terraform Schema.
	// * If it is not set, we don't need to do anything and return `nil`

	if !v.Plan.GuestProperties.IsNull() || !v.Plan.GuestProperties.IsUnknown() {
		vmProperties, err := getGuestProperties(v.Plan.GuestProperties)
		if err != nil {
			return fmt.Errorf("unable to convert guest properties to data structure")
		}

		_, err = vm.SetProductSectionList(vmProperties)
		if err != nil {
			return fmt.Errorf("error setting guest properties: %w", err)
		}
	}
	return nil
}

// getGuestProperties returns a struct for setting guest properties.
func getGuestProperties(guestProperties types.Map) (*govcdtypes.ProductSectionList, error) {
	// Init Struct
	vmProperties := &govcdtypes.ProductSectionList{
		ProductSection: &govcdtypes.ProductSection{
			Info:     "Custom properties",
			Property: []*govcdtypes.Property{},
		},
	}

	// For each key/value pair, add it to the struct
	for key, value := range guestProperties.Elements() {
		oneProp := &govcdtypes.Property{
			UserConfigurable: true,
			Type:             "string",
			Key:              key,
			Label:            key,
			Value:            &govcdtypes.Value{Value: value.String()},
		}
		vmProperties.ProductSection.Property = append(vmProperties.ProductSection.Property, oneProp)
	}

	if len(guestProperties.Elements()) != len(vmProperties.ProductSection.Property) {
		return nil, fmt.Errorf("unable to convert guest properties to data structure")
	}

	return vmProperties, nil
}

// * End of guest properties

/*
	* Customization

	Following code is responsible for setting VM customization
	3 functions are available:
		- updateGuestCustomizationSetting
		- updateCustomizationSection
		- isForcedCustomization

*/

// updateGuestCustomizationSetting is responsible for setting all the data related to VM customization.
func updateGuestCustomizationSetting(v *VMClient, vm *govcd.VM) error {
	// Retrieve existing customization section to only customize what was throughout this function
	customizationSection, err := vm.GetGuestCustomizationSection()
	if err != nil {
		return fmt.Errorf("error getting existing customization section before changing: %w", err)
	}

	// Process parameters from 'customization' block
	updateCustomizationSection(v, customizationSection)

	// Apply any of the settings we have set
	if _, err = vm.SetGuestCustomizationSection(customizationSection); err != nil {
		return fmt.Errorf("error applying guest customization details: %w", err)
	}

	return nil
}

// updateCustomizationSection is responsible for setting all the data related to VM customization.
func updateCustomizationSection(v *VMClient, customizationSection *govcdtypes.GuestCustomizationSection) {
	if v.Plan.ComputerName.IsNull() {
		// for back compatibility we allow to set computer name from `name` if computer_name isn't provided
		customizationSection.ComputerName = v.Plan.VMName.ValueString()
	} else {
		customizationSection.ComputerName = v.Plan.ComputerName.ValueString()
	}

	var customization commonvm.Customization
	v.Plan.Customization.As(context.Background(), customization, basetypes.ObjectAsOptions{})

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

// isForcedCustomization checks "customization" block in resource and checks if the value of field "force"
// is set to "true". It returns false if the value is not set or is set to false.
func isForcedCustomization(v *VMClient) bool {
	if v.Plan.Customization.IsNull() {
		return false
	}

	var customization commonvm.Customization
	v.Plan.Customization.As(context.Background(), customization, basetypes.ObjectAsOptions{})

	if !customization.Force.IsNull() && !customization.Force.IsUnknown() {
		return customization.Force.ValueBool()
	}

	return false
}

// * End of customization section

// lookupvAppTemplateforVm will do the following
// evaluate if optional parameter `vm_name_in_template` was specified.
//
// If `vm_name_in_template` was specified
// * It will look up the exact VM with given `vm_name_in_template` inside `vapp_template_id`
//
// If `vm_name_in_template` was not specified:
// * Return error.
func lookupvAppTemplateforVM(v *VMClient, org *govcd.Org, vdc *govcd.Vdc) (govcd.VAppTemplate, error) {
	if !v.Plan.VappTemplateID.IsNull() && !v.Plan.VappTemplateID.IsUnknown() {
		// Lookup of vApp Template using URN

		vAppTemplate, err := v.Client.Vmware.GetVAppTemplateById(v.Plan.VappTemplateID.ValueString())
		if err != nil {
			return govcd.VAppTemplate{}, fmt.Errorf("error finding vApp Template with URN %s: %w", v.Plan.VappTemplateID.ValueString(), err)
		}

		if !v.Plan.VMNameInTemplate.IsNull() && !v.Plan.VMNameInTemplate.IsUnknown() {
			vmInTemplateRecord, err := v.Client.Vmware.QuerySynchronizedVmInVAppTemplateByHref(vAppTemplate.VAppTemplate.HREF, v.Plan.VMNameInTemplate.ValueString())
			if err != nil {
				return govcd.VAppTemplate{}, fmt.Errorf("error obtaining VM '%s' inside vApp Template: %w", v.Plan.VMNameInTemplate.ValueString(), err)
			}

			returnedVAppTemplate, err := v.Client.Vmware.GetVAppTemplateByHref(vmInTemplateRecord.HREF)
			if err != nil {
				return govcd.VAppTemplate{}, fmt.Errorf("error getting vApp template from inner VM %s: %w", v.Plan.VMNameInTemplate.ValueString(), err)
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
// * `vapp` type networks will be checked for existence inside the vApp.
func networksToConfig(v *VMClient, vapp *govcd.VApp) (govcdtypes.NetworkConnectionSection, error) {
	networkConnectionSection := govcdtypes.NetworkConnectionSection{}

	if v.Plan.Networks.IsNull() {
		return networkConnectionSection, nil
	}

	networks, err := commonvm.NetworksFromPlan(v.Plan.Networks)
	if err != nil {
		return networkConnectionSection, err
	}

	// sets existing primary network connection index. Further code changes index only if change is
	// found
	for index, singleNetwork := range *networks {
		if singleNetwork.IsPrimary.ValueBool() {
			networkConnectionSection.PrimaryNetworkConnectionIndex = index
		}
	}

	for index, singleNetwork := range *networks {
		netConn := &govcdtypes.NetworkConnection{}

		networkName := singleNetwork.Name.ValueString()
		ipAllocationMode := singleNetwork.IPAllocationMode.ValueString()
		ip := singleNetwork.IP.ValueString()

		if v.State != nil {
			networksState, err := commonvm.NetworksFromPlan(v.State.Networks)
			if err != nil {
				return govcdtypes.NetworkConnectionSection{}, err
			}

			if !singleNetwork.IsPrimary.Equal((*networksState)[index].IsPrimary) && singleNetwork.IsPrimary.ValueBool() {
				networkConnectionSection.PrimaryNetworkConnectionIndex = index
			}
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
				return govcdtypes.NetworkConnectionSection{}, fmt.Errorf("unable to find vApp network %s: %w", networkName, err)
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

// isItVappOrgNetwork checks if it is a vApp Org network (not vApp Network).
func isItVappOrgNetwork(vAppNetworkName string, vapp govcd.VApp) (bool, error) {
	vAppNetworkConfig, err := vapp.GetNetworkConfig()
	if err != nil {
		return false, fmt.Errorf("error getting vApp networks: %w", err)
	}

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.NetworkName == vAppNetworkName &&
			!govcd.IsVappNetwork(networkConfig.Configuration) {
			return true, nil
		}
	}

	return false, fmt.Errorf("configured vApp Org network isn't found: %s", vAppNetworkName)
}

// isItVappNetwork checks if it is a vApp network (not vApp Org Network).
func isItVappNetwork(vAppNetworkName string, vapp govcd.VApp) (bool, error) {
	vAppNetworkConfig, err := vapp.GetNetworkConfig()
	if err != nil {
		return false, fmt.Errorf("error getting vApp networks: %w", err)
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
		return nil, errors.New("storageProfileName is an empty string")
	}

	storageProfile, err := vdc.FindStorageProfileReference(storageProfileName)
	if err != nil {
		return nil, fmt.Errorf("[vm creation] error retrieving storage profile %s : %w", storageProfileName, err)
	}

	return &storageProfile, nil
}

// lookupComputePolicy returns the Compute Policy associated to the value of the given Compute Policy attribute. If the
// attribute is not set, the returned policy will be nil. If the obtained policy is incorrect, it will return an error.
func lookupComputePolicy(v *VMClient, value string) (*govcd.VdcComputePolicyV2, error) {
	if value == "" {
		return nil, errors.New("value is an empty string")
	}

	computePolicy, err := v.Client.Vmware.GetVdcComputePolicyV2ById(value)
	if err != nil {
		return nil, fmt.Errorf("error getting compute policy %s: %w", value, err)
	}
	if computePolicy.Href == "" {
		return nil, fmt.Errorf("empty compute policy HREF detected")
	}
	return computePolicy, nil
}

func updateOsType(v *VMClient, vm *govcd.VM) error {
	var err error

	vmSpecSection := vm.VM.VmSpecSection

	if !v.Plan.OsType.IsNull() && !v.Plan.OsType.IsUnknown() {
		vmSpecSection.OsType = v.Plan.OsType.ValueString()
		_, err = vm.UpdateVmSpecSection(vmSpecSection, v.Plan.Description.ValueString())
		if err != nil {
			return fmt.Errorf("error changing VM spec section: %w", err)
		}
	}

	return nil
}

// getCpuMemoryValues returns CPU, CPU core count and Memory variables. Priority comes from HCL
// schema configuration and then whatever is present in compute policy (if it was specified at all).
func getCPUMemoryValues(v *VMClient, vdcComputePolicy *govcdtypes.VdcComputePolicyV2) (cpu, cores *int, memory *int64, err error) {
	var (
		setCPU    *int
		setCores  *int
		setMemory *int64
	)

	var resource commonvm.Resource
	if !v.Plan.Resource.IsNull() && !v.Plan.Resource.IsUnknown() {
		d := v.Plan.Resource.As(context.Background(), resource, basetypes.ObjectAsOptions{})
		if d.HasError() {
			return nil, nil, nil, fmt.Errorf("error retrieving resource: %s", d)
		}

		if !resource.Memory.IsNull() && !resource.Memory.IsUnknown() {
			setMemory = utils.TakeInt64Pointer(resource.Memory.ValueInt64())
		}

		if !resource.CPUs.IsNull() && !resource.CPUs.IsUnknown() {
			setCPU = utils.TakeIntPointer(int(resource.CPUs.ValueInt64()))
		}

		if !resource.CPUCores.IsNull() && !resource.CPUCores.IsUnknown() {
			setCores = utils.TakeIntPointer(int(resource.CPUCores.ValueInt64()))
		}
	}

	return setCPU, setCores, setMemory, nil
}

func createPlan(ctx context.Context, gvm *govcd.VM, x *vmResourceModel) (newPlan *vmResourceModel, err error) {
	vdc, err := gvm.GetParentVdc()
	if err != nil {
		return
	}

	networks, err := commonvm.NetworksRead(gvm)
	if err != nil {
		return
	}

	networksPlan, diag := networks.ToPlan()
	if diag.HasError() {
		err = fmt.Errorf("error converting networks to plan: %s", diag[0].Detail())
		return
	}

	resource, err := commonvm.ResourceRead(gvm)
	if err != nil {
		return
	}

	disks, err := DisksRead(gvm)
	if err != nil {
		return
	}

	guestproperties, err := commonvm.GuestPropertiesRead(gvm)
	if err != nil {
		return
	}

	customization, err := commonvm.CustomizationRead(gvm)
	if err != nil {
		return
	}

	customizationState, d := commonvm.CustomizationFromPlan(ctx, x.Customization)
	if d.HasError() {
		err = fmt.Errorf("error convert customization: %s", d[0].Detail())
		return
	}

	customization.Force = customizationState.Force

	statusText, err := gvm.GetStatus()
	if err != nil {
		statusText = vmUnknownStatus
	}

	vapp, err := gvm.GetParentVApp()
	if err != nil {
		return
	}

	err = gvm.Refresh()
	if err != nil {
		return
	}

	newPlan = &vmResourceModel{
		ID:  types.StringValue(gvm.VM.ID),
		VDC: types.StringValue(vdc.Vdc.Name),

		VappName:       types.StringValue(vapp.VApp.Name),
		VappTemplateID: x.VappTemplateID,

		Resource: resource.ToPlan(),

		VMName:           types.StringValue(gvm.VM.Name),
		VMNameInTemplate: x.VMNameInTemplate,

		Description:    types.StringValue(gvm.VM.Description),
		Href:           types.StringValue(gvm.VM.HREF),
		AcceptAllEulas: x.AcceptAllEulas,

		PowerON:               x.PowerON,
		PreventUpdatePowerOff: x.PreventUpdatePowerOff,

		BootImageID: x.BootImageID,
		// OverrideTemplateDisks: x.OverrideTemplateDisks,

		Networks:               networksPlan,
		NetworkDhcpWaitSeconds: x.NetworkDhcpWaitSeconds,

		ExposeHardwareVirtualization: types.BoolValue(gvm.VM.NestedHypervisorEnabled),
		GuestProperties:              guestproperties.ToPlan(),
		Customization:                customization.ToPlan(),
		StatusCode:                   types.Int64Value(int64(gvm.VM.Status)),
		StatusText:                   types.StringValue(statusText),
	}

	if gvm.VM.VmSpecSection != nil {
		newPlan.OsType = types.StringValue(gvm.VM.VmSpecSection.OsType)
	}

	if gvm.VM.ComputePolicy != nil {
		if gvm.VM.ComputePolicy.VmSizingPolicy != nil {
			newPlan.SizingPolicyID = types.StringValue(gvm.VM.ComputePolicy.VmSizingPolicy.ID)
		}

		if gvm.VM.ComputePolicy.VmPlacementPolicy != nil {
			newPlan.PlacementPolicyID = types.StringValue(gvm.VM.ComputePolicy.VmPlacementPolicy.ID)
		}
	}

	if gvm.VM.GuestCustomizationSection != nil {
		newPlan.ComputerName = types.StringValue(gvm.VM.GuestCustomizationSection.ComputerName)
	} else if !x.ComputerName.IsNull() {
		newPlan.ComputerName = x.ComputerName
	}

	if gvm.VM.StorageProfile != nil {
		newPlan.StorageProfile = types.StringValue(gvm.VM.StorageProfile.Name)
	}

	xDisks, d := disks.ToPlan(ctx)
	if d.HasError() {
		err = fmt.Errorf("error convert disks: %s", d[0].Detail())
		return
	}
	newPlan.Disks = xDisks

	return newPlan, nil
}
