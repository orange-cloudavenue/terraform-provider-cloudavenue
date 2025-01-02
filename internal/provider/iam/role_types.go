package iam

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type RoleResourceModel struct {
	ID          supertypes.StringValue        `tfsdk:"id"`
	Name        supertypes.StringValue        `tfsdk:"name"`
	Description supertypes.StringValue        `tfsdk:"description"`
	Rights      supertypes.SetValueOf[string] `tfsdk:"rights"`
}

type RoleDataSourceModel struct {
	ID          supertypes.StringValue        `tfsdk:"id"`
	Name        supertypes.StringValue        `tfsdk:"name"`
	Description supertypes.StringValue        `tfsdk:"description"`
	ReadOnly    supertypes.BoolValue          `tfsdk:"read_only"`
	Rights      supertypes.SetValueOf[string] `tfsdk:"rights"`
}
