package vm

import "github.com/hashicorp/terraform-plugin-framework/types"

type vmAffinityRuleResourceModel struct {
	ID       types.String `tfsdk:"id"`
	VDC      types.String `tfsdk:"vdc"`
	Name     types.String `tfsdk:"name"`
	Polarity types.String `tfsdk:"polarity"`
	Required types.Bool   `tfsdk:"required"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	VMIDs    types.Set    `tfsdk:"vm_ids"`
}

type vmAffinityRuleDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	VDC      types.String `tfsdk:"vdc"`
	Name     types.String `tfsdk:"name"`
	Polarity types.String `tfsdk:"polarity"`
	Required types.Bool   `tfsdk:"required"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	VMIDs    types.Set    `tfsdk:"vm_ids"`
}
