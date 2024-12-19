package alb

import supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

type serviceEngineGroupsModel struct {
	ID                  supertypes.StringValue                                      `tfsdk:"id"`
	ServiceEngineGroups supertypes.ListNestedObjectValueOf[serviceEngineGroupModel] `tfsdk:"service_engine_groups"`
	EdgeGatewayID       supertypes.StringValue                                      `tfsdk:"edge_gateway_id"`
	EdgeGatewayName     supertypes.StringValue                                      `tfsdk:"edge_gateway_name"`
}
