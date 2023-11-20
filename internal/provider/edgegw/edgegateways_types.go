// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

type (
	edgeGatewaysDataSourceModel struct {
		ID           supertypes.StringValue                                                    `tfsdk:"id"`
		EdgeGateways supertypes.ListNestedObjectValueOf[edgeGatewayDataSourceModelEdgeGateway] `tfsdk:"edge_gateways"`
	}
	edgeGatewayDataSourceModelEdgeGateway struct {
		Tier0VrfName supertypes.StringValue `tfsdk:"tier0_vrf_name"`
		Name         supertypes.StringValue `tfsdk:"name"`
		ID           supertypes.StringValue `tfsdk:"id"`
		OwnerType    supertypes.StringValue `tfsdk:"owner_type"`
		OwnerName    supertypes.StringValue `tfsdk:"owner_name"`
		Description  supertypes.StringValue `tfsdk:"description"`
		LbEnabled    supertypes.BoolValue   `tfsdk:"lb_enabled"`
	}
)
