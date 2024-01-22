package vdc

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

type vdcsDataSourceModel struct {
	ID   supertypes.StringValue                     `tfsdk:"id"`
	VDCs supertypes.ListNestedObjectValueOf[vdcRef] `tfsdk:"vdcs"`
}

type vdcRef struct {
	VDCName supertypes.StringValue `tfsdk:"vdc_name"` // Deprecated
	VDCUUID supertypes.StringValue `tfsdk:"vdc_uuid"` // Deprecated
	Name    supertypes.StringValue `tfsdk:"name"`
	ID      supertypes.StringValue `tfsdk:"id"`
}
