package publicip

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type publicIPResourceModel struct {
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	ID              types.String   `tfsdk:"id"`
	PublicIP        types.String   `tfsdk:"public_ip"`
	EdgeGatewayName types.String   `tfsdk:"edge_gateway_name"`
	EdgeGatewayID   types.String   `tfsdk:"edge_gateway_id"`
}
