package edgegw

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type edgeGatewaysResourceModel struct {
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
	ID                  types.String   `tfsdk:"id"`
	Tier0VrfID          types.String   `tfsdk:"tier0_vrf_name"`
	Name                types.String   `tfsdk:"name"`
	OwnerType           types.String   `tfsdk:"owner_type"`
	OwnerName           types.String   `tfsdk:"owner_name"`
	Description         types.String   `tfsdk:"description"`
	EnableLoadBalancing types.Bool     `tfsdk:"lb_enabled"`
}

type edgeGatewayDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Tier0VrfID          types.String `tfsdk:"tier0_vrf_name"`
	Name                types.String `tfsdk:"name"`
	OwnerType           types.String `tfsdk:"owner_type"`
	OwnerName           types.String `tfsdk:"owner_name"`
	Description         types.String `tfsdk:"description"`
	EnableLoadBalancing types.Bool   `tfsdk:"lb_enabled"`
}
