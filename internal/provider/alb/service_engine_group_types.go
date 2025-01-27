package alb

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type serviceEngineGroupModel struct {
	ID                      supertypes.StringValue `tfsdk:"id"`
	Name                    supertypes.StringValue `tfsdk:"name"`
	EdgeGatewayID           supertypes.StringValue `tfsdk:"edge_gateway_id"`
	EdgeGatewayName         supertypes.StringValue `tfsdk:"edge_gateway_name"`
	MaxVirtualServices      supertypes.Int64Value  `tfsdk:"max_virtual_services"`
	ReservedVirtualServices supertypes.Int64Value  `tfsdk:"reserved_virtual_services"`
	DeployedVirtualServices supertypes.Int64Value  `tfsdk:"deployed_virtual_services"`
}
