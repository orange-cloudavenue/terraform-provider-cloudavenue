package vapp

import "github.com/hashicorp/terraform-plugin-framework/types"

type aclResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	VDC                 types.String `tfsdk:"vdc"`
	VAppID              types.String `tfsdk:"vapp_id"`
	VAppName            types.String `tfsdk:"vapp_name"`
	EveryoneAccessLevel types.String `tfsdk:"everyone_access_level"`
	SharedWith          types.Set    `tfsdk:"shared_with"`
}
