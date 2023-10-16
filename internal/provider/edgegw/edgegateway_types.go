package edgegw

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type edgeGatewayResourceModel struct {
	Timeouts            timeouts.Value         `tfsdk:"timeouts"`
	ID                  supertypes.StringValue `tfsdk:"id"`
	Tier0VrfID          supertypes.StringValue `tfsdk:"tier0_vrf_name"`
	Name                supertypes.StringValue `tfsdk:"name"`
	OwnerType           supertypes.StringValue `tfsdk:"owner_type"`
	OwnerName           supertypes.StringValue `tfsdk:"owner_name"`
	Description         supertypes.StringValue `tfsdk:"description"`
	EnableLoadBalancing supertypes.BoolValue   `tfsdk:"lb_enabled"`
	Bandwidth           supertypes.Int64Value  `tfsdk:"bandwidth"`
}

type edgeGatewayDatasourceModel struct {
	ID                  supertypes.StringValue `tfsdk:"id"`
	Tier0VrfID          supertypes.StringValue `tfsdk:"tier0_vrf_name"`
	Name                supertypes.StringValue `tfsdk:"name"`
	OwnerType           supertypes.StringValue `tfsdk:"owner_type"`
	OwnerName           supertypes.StringValue `tfsdk:"owner_name"`
	Description         supertypes.StringValue `tfsdk:"description"`
	EnableLoadBalancing supertypes.BoolValue   `tfsdk:"lb_enabled"`
	Bandwidth           supertypes.Int64Value  `tfsdk:"bandwidth"`
}

// Copy returns a copy of the edgeGatewayResourceModel.
func (rm *edgeGatewayResourceModel) Copy() *edgeGatewayResourceModel {
	x := &edgeGatewayResourceModel{}
	utils.ModelCopy(rm, x)
	return x
}

// Copy returns a copy of the edgeGatewayDatasourceModel.
func (dm *edgeGatewayDatasourceModel) Copy() *edgeGatewayDatasourceModel {
	x := &edgeGatewayDatasourceModel{}
	utils.ModelCopy(dm, x)
	return x
}
