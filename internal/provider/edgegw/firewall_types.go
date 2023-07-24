package edgegw

import "github.com/hashicorp/terraform-plugin-framework/types"

type firewallModel struct {
	ID              types.String `tfsdk:"id"`
	EdgeGatewayID   types.String `tfsdk:"edge_gateway_id"`
	EdgeGatewayName types.String `tfsdk:"edge_gateway_name"`
	Rules           types.List   `tfsdk:"rules"`
}

type firewallModelRules []firewallModelRule

type firewallModelRule struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Direction         types.String `tfsdk:"direction"`
	IPProtocol        types.String `tfsdk:"ip_protocol"`
	Action            types.String `tfsdk:"action"`
	Logging           types.Bool   `tfsdk:"logging"`
	SourceIDs         types.Set    `tfsdk:"source_ids"`
	DestinationIDs    types.Set    `tfsdk:"destination_ids"`
	AppPortProfileIDs types.Set    `tfsdk:"app_port_profile_ids"`
}
