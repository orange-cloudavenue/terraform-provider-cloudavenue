package vrf

import "github.com/hashicorp/terraform-plugin-framework/types"

type tier0VrfDataSourceModel struct {
	ID           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Provider     types.String   `tfsdk:"tier0_provider"`
	ClassService types.String   `tfsdk:"class_service"`
	Services     []segmentModel `tfsdk:"services"`
}

type segmentModel struct {
	Service types.String `tfsdk:"service"`
	VLANID  types.String `tfsdk:"vlan_id"`
}
