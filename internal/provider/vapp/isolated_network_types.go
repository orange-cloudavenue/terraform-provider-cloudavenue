package vapp

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	isolatedNetworkModel struct {
		ID                 supertypes.StringValue                                              `tfsdk:"id"`
		VDC                supertypes.StringValue                                              `tfsdk:"vdc"`
		Name               supertypes.StringValue                                              `tfsdk:"name"`
		Description        supertypes.StringValue                                              `tfsdk:"description"`
		VAppName           supertypes.StringValue                                              `tfsdk:"vapp_name"`
		VAppID             supertypes.StringValue                                              `tfsdk:"vapp_id"`
		Netmask            supertypes.StringValue                                              `tfsdk:"netmask"`
		Gateway            supertypes.StringValue                                              `tfsdk:"gateway"`
		DNS1               supertypes.StringValue                                              `tfsdk:"dns1"`
		DNS2               supertypes.StringValue                                              `tfsdk:"dns2"`
		DNSSuffix          supertypes.StringValue                                              `tfsdk:"dns_suffix"`
		GuestVLANAllowed   supertypes.BoolValue                                                `tfsdk:"guest_vlan_allowed"`
		RetainIPMacEnabled supertypes.BoolValue                                                `tfsdk:"retain_ip_mac_enabled"`
		StaticIPPool       supertypes.SetNestedObjectValueOf[isolatedNetworkModelStaticIPPool] `tfsdk:"static_ip_pool"`
	}

	isolatedNetworkModelStaticIPPool struct {
		StartAddress supertypes.StringValue `tfsdk:"start_address"`
		EndAddress   supertypes.StringValue `tfsdk:"end_address"`
	}
)

// Copy returns a copy of the backupModel.
func (rm *isolatedNetworkModel) Copy() *isolatedNetworkModel {
	x := &isolatedNetworkModel{}
	utils.ModelCopy(rm, x)
	return x
}
