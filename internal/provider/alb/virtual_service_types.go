package alb

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

type (
	VirtualServiceModel struct {
		ID                     supertypes.StringValue           `tfsdk:"id"`
		Name                   supertypes.StringValue           `tfsdk:"name"`
		EdgeGatewayName        supertypes.StringValue           `tfsdk:"edge_gateway_name"`
		EdgeGatewayID          supertypes.StringValue           `tfsdk:"edge_gateway_id"`
		Description            supertypes.StringValue           `tfsdk:"description"`
		Enabled                supertypes.BoolValue             `tfsdk:"enabled"`
		PoolName               supertypes.StringValue           `tfsdk:"pool_name"`
		PoolID                 supertypes.StringValue           `tfsdk:"pool_id"`
		ServiceEngineGroupName supertypes.StringValue           `tfsdk:"service_engine_group_name"`
		VirtualIP              supertypes.StringValue           `tfsdk:"virtual_ip"`
		CertificateID          supertypes.StringValue           `tfsdk:"certificate_id"`
		ServicePorts           []VirtualServiceModelServicePort `tfsdk:"service_ports"`
		PreserveClientIP       supertypes.BoolValue             `tfsdk:"preserve_client_ip"`
	}

	VirtualServiceModelServicePort struct {
		PortStart supertypes.Int64Value  `tfsdk:"port_start"`
		PortEnd   supertypes.Int64Value  `tfsdk:"port_end"`
		PortType  supertypes.StringValue `tfsdk:"port_type"`
		PortSSL   supertypes.BoolValue   `tfsdk:"port_ssl"`
	}
)
