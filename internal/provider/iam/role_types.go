package iam

import "github.com/hashicorp/terraform-plugin-framework/types"

type roleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Rights      types.Set    `tfsdk:"rights"`
}

type roleDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	Rights      types.Set    `tfsdk:"rights"`
}
