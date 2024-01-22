package publicip

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type publicIPResourceModel struct {
	Timeouts        timeouts.Value         `tfsdk:"timeouts"`
	ID              supertypes.StringValue `tfsdk:"id"`
	PublicIP        supertypes.StringValue `tfsdk:"public_ip"`
	EdgeGatewayName supertypes.StringValue `tfsdk:"edge_gateway_name"`
	EdgeGatewayID   supertypes.StringValue `tfsdk:"edge_gateway_id"`
}

func (rm *publicIPResourceModel) Copy() *publicIPResourceModel {
	x := &publicIPResourceModel{}
	utils.ModelCopy(rm, x)
	return x
}
