package publicip

import "github.com/hashicorp/terraform-plugin-framework/types"

type publicIPDataSourceModel struct {
	ID        types.String                 `tfsdk:"id"`
	PublicIPs []publicIPNetworkConfigModel `tfsdk:"public_ips"`
}

type publicIPNetworkConfigModel struct {
	ID              types.String `tfsdk:"id"`
	PublicIP        types.String `tfsdk:"public_ip"`
	EdgeGatewayName types.String `tfsdk:"edge_gateway_name"`
	EdgeGatewayID   types.String `tfsdk:"edge_gateway_id"`
}
