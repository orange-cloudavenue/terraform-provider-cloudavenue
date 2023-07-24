package vapp

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type vappLeaseModel struct {
	RuntimeLeaseInSec types.Int64 `tfsdk:"runtime_lease_in_sec"`
	StorageLeaseInSec types.Int64 `tfsdk:"storage_lease_in_sec"`
}

var vappLeaseAttrTypes = map[string]attr.Type{
	"runtime_lease_in_sec": types.Int64Type,
	"storage_lease_in_sec": types.Int64Type,
}

type orgNetworkModel struct {
	ID                 types.String `tfsdk:"id"`
	VAppName           types.String `tfsdk:"vapp_name"`
	VAppID             types.String `tfsdk:"vapp_id"`
	VDC                types.String `tfsdk:"vdc"`
	NetworkName        types.String `tfsdk:"network_name"`
	IsFenced           types.Bool   `tfsdk:"is_fenced"`
	RetainIPMacEnabled types.Bool   `tfsdk:"retain_ip_mac_enabled"`
}
