package vdc

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

type vdcsDataSourceModel struct {
	ID   supertypes.StringValue                     `tfsdk:"id"`
	VDCs supertypes.ListNestedObjectValueOf[vdcRef] `tfsdk:"vdcs"`
}

type vdcRef struct {
	ID   supertypes.StringValue `tfsdk:"id"`
	Name supertypes.StringValue `tfsdk:"name"`
}
