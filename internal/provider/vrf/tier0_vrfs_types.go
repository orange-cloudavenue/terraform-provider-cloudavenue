package vrf

import supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

type tier0VrfsDataSourceModel struct {
	ID    supertypes.StringValue         `tfsdk:"id"`
	Names supertypes.ListValueOf[string] `tfsdk:"names"`
}
