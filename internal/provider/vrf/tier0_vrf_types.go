package vrf

import supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

type tier0VrfDataSourceModel struct {
	ID           supertypes.StringValue                           `tfsdk:"id"`
	Name         supertypes.StringValue                           `tfsdk:"name"`
	Provider     supertypes.StringValue                           `tfsdk:"tier0_provider"`
	ClassService supertypes.StringValue                           `tfsdk:"class_service"`
	Services     supertypes.ListNestedObjectValueOf[segmentModel] `tfsdk:"services"`
}

type segmentModel struct {
	Service supertypes.StringValue `tfsdk:"service"`
	VLANID  supertypes.StringValue `tfsdk:"vlan_id"`
}
