package network

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

type staticIPPool struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

var staticIPPoolAttrTypes = map[string]attr.Type{
	"start_address": types.StringType,
	"end_address":   types.StringType,
}

type networkIsolatedModel struct {
	ID           types.String `tfsdk:"id"`
	VDC          types.String `tfsdk:"vdc"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Gateway      types.String `tfsdk:"gateway"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	DNS1         types.String `tfsdk:"dns1"`
	DNS2         types.String `tfsdk:"dns2"`
	DNSSuffix    types.String `tfsdk:"dns_suffix"`
	StaticIPPool types.Set    `tfsdk:"static_ip_pool"`
}

// Get parent edge gateway ID.
func GetParentEdgeGatewayID(org org.Org, edgeGatewayID string) (*string, diag.Diagnostic) {
	anyEdgeGateway, err := org.GetAnyTypeEdgeGatewayById(edgeGatewayID)
	if err != nil {
		return nil, diag.NewErrorDiagnostic("error retrieving edge gateway", err.Error())
	}
	if anyEdgeGateway == nil {
		return nil, diag.NewErrorDiagnostic("error retrieving edge gateway", "edge gateway is a nil object")
	}
	id := anyEdgeGateway.EdgeGateway.OwnerRef.ID

	return &id, nil
}

// Get IP Pool information data from network.
func GetIPRanges(network *govcd.OpenApiOrgVdcNetwork) []staticIPPool {
	ipPools := []staticIPPool{}

	for _, ipRange := range network.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
		ipPool := staticIPPool{
			StartAddress: types.StringValue(ipRange.StartAddress),
			EndAddress:   types.StringValue(ipRange.EndAddress),
		}
		ipPools = append(ipPools, ipPool)
	}
	return ipPools
}
