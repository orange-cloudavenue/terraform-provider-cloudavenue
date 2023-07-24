// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type edgeGatewaysDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	EdgeGateways types.List   `tfsdk:"edge_gateways"`
}

var edgeGatewayDataSourceModelAttrTypes = map[string]attr.Type{
	"tier0_vrf_name": types.StringType,
	"name":           types.StringType,
	"id":             types.StringType,
	"owner_type":     types.StringType,
	"owner_name":     types.StringType,
	"description":    types.StringType,
	"lb_enabled":     types.BoolType,
}
