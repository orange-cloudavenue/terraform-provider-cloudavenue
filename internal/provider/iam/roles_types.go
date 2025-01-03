package iam

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type RolesModel struct {
	ID    supertypes.StringValue                                 `tfsdk:"id"`
	Roles supertypes.MapNestedObjectValueOf[RoleDataSourceModel] `tfsdk:"roles"`
}
