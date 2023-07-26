package vapp

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type isolatedNetworkDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	VDC                types.String `tfsdk:"vdc"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	VAppName           types.String `tfsdk:"vapp_name"`
	VAppID             types.String `tfsdk:"vapp_id"`
	Netmask            types.String `tfsdk:"netmask"`
	Gateway            types.String `tfsdk:"gateway"`
	DNS1               types.String `tfsdk:"dns1"`
	DNS2               types.String `tfsdk:"dns2"`
	DNSSuffix          types.String `tfsdk:"dns_suffix"`
	GuestVLANAllowed   types.Bool   `tfsdk:"guest_vlan_allowed"`
	RetainIPMacEnabled types.Bool   `tfsdk:"retain_ip_mac_enabled"`
	StaticIPPool       types.Set    `tfsdk:"static_ip_pool"`
}

type isolatedNetworkResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	VDC                types.String `tfsdk:"vdc"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	VAppName           types.String `tfsdk:"vapp_name"`
	VAppID             types.String `tfsdk:"vapp_id"`
	Netmask            types.String `tfsdk:"netmask"`
	Gateway            types.String `tfsdk:"gateway"`
	DNS1               types.String `tfsdk:"dns1"`
	DNS2               types.String `tfsdk:"dns2"`
	DNSSuffix          types.String `tfsdk:"dns_suffix"`
	GuestVLANAllowed   types.Bool   `tfsdk:"guest_vlan_allowed"`
	RetainIPMacEnabled types.Bool   `tfsdk:"retain_ip_mac_enabled"`
	StaticIPPool       types.Set    `tfsdk:"static_ip_pool"`
}

type staticIPPoolModel struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

var staticIPPoolModelAttrTypes = map[string]attr.Type{
	"start_address": types.StringType,
	"end_address":   types.StringType,
}
