package vdc

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type vdcGroupDataSourceModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	ErrorMessage               types.String `tfsdk:"error_message"`
	DFWEnabled                 types.Bool   `tfsdk:"dfw_enabled"`
	LocalEgress                types.Bool   `tfsdk:"local_egress"`
	NetworkPoolID              types.String `tfsdk:"network_pool_id"`
	NetworkPoolUniversalID     types.String `tfsdk:"network_pool_universal_id"`
	NetworkProviderType        types.String `tfsdk:"network_provider_type"`
	Status                     types.String `tfsdk:"status"`
	Type                       types.String `tfsdk:"type"`
	UniversalNetworkingEnabled types.Bool   `tfsdk:"universal_networking_enabled"`
	Vdcs                       types.List   `tfsdk:"vdcs"`
}

type vdcModel struct {
	FaultDomainTag       types.String `tfsdk:"fault_domain_tag"`
	NetworkProviderScope types.String `tfsdk:"network_provider_scope"`
	IsRemoteOrg          types.Bool   `tfsdk:"is_remote_org"`
	Status               types.String `tfsdk:"status"`
	SiteID               types.String `tfsdk:"site_name"`
	SiteName             types.String `tfsdk:"site_id"`
	Name                 types.String `tfsdk:"name"`
	ID                   types.String `tfsdk:"id"`
}

var vdcModelAttrTypes = map[string]attr.Type{
	"fault_domain_tag":       types.StringType,
	"network_provider_scope": types.StringType,
	"is_remote_org":          types.BoolType,
	"status":                 types.StringType,
	"site_id":                types.StringType,
	"site_name":              types.StringType,
	"name":                   types.StringType,
	"id":                     types.StringType,
}
