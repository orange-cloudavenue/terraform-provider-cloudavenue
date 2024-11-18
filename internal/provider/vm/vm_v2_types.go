package vm

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

type (
	VMV2Model struct {
		ID          supertypes.StringValue                                  `tfsdk:"id"`
		VDC         supertypes.StringValue                                  `tfsdk:"vdc"`
		Name        supertypes.StringValue                                  `tfsdk:"name"`
		VappName    supertypes.StringValue                                  `tfsdk:"vapp_name"`
		VappID      supertypes.StringValue                                  `tfsdk:"vapp_id"`
		Description supertypes.StringValue                                  `tfsdk:"description"`
		DeployOS    supertypes.SingleNestedObjectValueOf[VMV2ModelDeployOS] `tfsdk:"deploy_os"`
		State       supertypes.SingleNestedObjectValueOf[VMV2ModelState]    `tfsdk:"state"`
		Resource    supertypes.SingleNestedObjectValueOf[VMV2ModelResource] `tfsdk:"resource"`
		Settings    supertypes.SingleNestedObjectValueOf[VMV2ModelSettings] `tfsdk:"settings"`
	}

	VMV2ModelDeployOS struct {
		VAppTemplateID   supertypes.StringValue `tfsdk:"vapp_template_id"`
		VMNameInTemplate supertypes.StringValue `tfsdk:"vm_name_in_template"`
		BootImageID      supertypes.StringValue `tfsdk:"boot_image_id"`
		AcceptAllEulas   supertypes.BoolValue   `tfsdk:"accept_all_eulas"`
	}

	VMV2ModelState struct {
		PowerON supertypes.BoolValue   `tfsdk:"power_on"`
		Status  supertypes.StringValue `tfsdk:"status"`
	}

	VMV2ModelResource struct {
		CPUs                supertypes.Int64Value                                         `tfsdk:"cpus"`
		CPUsCores           supertypes.Int64Value                                         `tfsdk:"cpus_cores"`
		CPUHotAddEnabled    supertypes.BoolValue                                          `tfsdk:"cpu_hot_add_enabled"`
		Memory              supertypes.Int64Value                                         `tfsdk:"memory"`
		MemoryHotAddEnabled supertypes.BoolValue                                          `tfsdk:"memory_hot_add_enabled"`
		Networks            supertypes.ListNestedObjectValueOf[VMV2ModelResourceNetworks] `tfsdk:"networks"`
	}

	VMV2ModelResourceNetworks struct {
		Type             supertypes.StringValue `tfsdk:"type"`
		IPAllocationMode supertypes.StringValue `tfsdk:"ip_allocation_mode"`
		Name             supertypes.StringValue `tfsdk:"name"`
		IP               supertypes.StringValue `tfsdk:"ip"`
		IsPrimary        supertypes.BoolValue   `tfsdk:"is_primary"`
		Mac              supertypes.StringValue `tfsdk:"mac"`
		AdapterType      supertypes.StringValue `tfsdk:"adapter_type"`
		Connected        supertypes.BoolValue   `tfsdk:"connected"`
	}

	VMV2ModelSettings struct {
		ExposeHardwareVirtualization supertypes.BoolValue                                                 `tfsdk:"expose_hardware_virtualization"`
		OsType                       supertypes.StringValue                                               `tfsdk:"os_type"`
		StorageProfile               supertypes.StringValue                                               `tfsdk:"storage_profile"`
		GuestProperties              supertypes.MapValueOf[VMVMModelSettingsGuestProperties]              `tfsdk:"guest_properties"`
		AffinityRuleID               supertypes.StringValue                                               `tfsdk:"affinity_rule_id"`
		Customization                supertypes.SingleNestedObjectValueOf[VMV2ModelSettingsCustomization] `tfsdk:"customization"`
	}

	VMVMModelSettingsGuestProperties map[string]string

	VMV2ModelSettingsCustomization struct {
		Force                          supertypes.BoolValue   `tfsdk:"force"`
		Enabled                        supertypes.BoolValue   `tfsdk:"enabled"`
		ChangeSID                      supertypes.BoolValue   `tfsdk:"change_sid"`
		AllowLocalAdminPassword        supertypes.BoolValue   `tfsdk:"allow_local_admin_password"`
		MustChangePasswordOnFirstLogin supertypes.BoolValue   `tfsdk:"must_change_password_on_first_login"`
		AdminPassword                  supertypes.StringValue `tfsdk:"admin_password"`
		AutoGeneratePassword           supertypes.BoolValue   `tfsdk:"auto_generate_password"`
		NumberOfAutoLogons             supertypes.Int64Value  `tfsdk:"number_of_auto_logons"`
		JoinDomain                     supertypes.BoolValue   `tfsdk:"join_domain"`
		JoinOrgDomain                  supertypes.BoolValue   `tfsdk:"join_org_domain"`
		JoinDomainName                 supertypes.StringValue `tfsdk:"join_domain_name"`
		JoinDomainUser                 supertypes.StringValue `tfsdk:"join_domain_user"`
		JoinDomainPassword             supertypes.StringValue `tfsdk:"join_domain_password"`
		JoinDomainAccountOU            supertypes.StringValue `tfsdk:"join_domain_account_ou"`
		InitScript                     supertypes.StringValue `tfsdk:"init_script"`
		Hostname                       supertypes.StringValue `tfsdk:"hostname"`
	}
)

// * Resource Network
// ConvertToNetworkConnection converts a VMResourceModelResourceNetworks to a NetworkConnection.
func (n *VMV2ResourceModelResourceNetwork) ConvertToNetworkConnection() NetworkConnection {
	return NetworkConnection{
		Name:             n.Name,
		Connected:        n.Connected,
		IPAllocationMode: n.IPAllocationMode,
		IP:               n.IP,
		IsPrimary:        n.IsPrimary,
		Mac:              n.Mac,
		AdapterType:      n.AdapterType,
		Type:             n.Type,
	}
}
