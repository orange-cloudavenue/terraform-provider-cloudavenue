package publicip

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

type publicIPDataSourceModel struct {
	ID        supertypes.StringValue                                         `tfsdk:"id"`
	PublicIPs supertypes.ListNestedObjectValueOf[publicIPNetworkConfigModel] `tfsdk:"public_ips"`
}

type publicIPNetworkConfigModel struct {
	ID              supertypes.StringValue `tfsdk:"id"`
	PublicIP        supertypes.StringValue `tfsdk:"public_ip"`
	EdgeGatewayName supertypes.StringValue `tfsdk:"edge_gateway_name"`
	EdgeGatewayID   supertypes.StringValue `tfsdk:"edge_gateway_id"`
}
