package network

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
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

type networkRoutedModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	EdgeGatewayID   types.String `tfsdk:"edge_gateway_id"`
	EdgeGatewayName types.String `tfsdk:"edge_gateway_name"`
	InterfaceType   types.String `tfsdk:"interface_type"`
	Gateway         types.String `tfsdk:"gateway"`
	PrefixLength    types.Int64  `tfsdk:"prefix_length"`
	DNS1            types.String `tfsdk:"dns1"`
	DNS2            types.String `tfsdk:"dns2"`
	DNSSuffix       types.String `tfsdk:"dns_suffix"`
	StaticIPPool    types.Set    `tfsdk:"static_ip_pool"`
}

var networkMutexKV = mutex.NewKV()

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

// Set data to network routed model.
func SetDataToNetworkRoutedModel(network *govcd.OpenApiOrgVdcNetwork) networkRoutedModel {
	return networkRoutedModel{
		ID:              types.StringValue(network.OpenApiOrgVdcNetwork.ID),
		Name:            types.StringValue(network.OpenApiOrgVdcNetwork.Name),
		Description:     utils.StringValueOrNull(network.OpenApiOrgVdcNetwork.Description),
		EdgeGatewayID:   types.StringValue(network.OpenApiOrgVdcNetwork.Connection.RouterRef.ID),
		EdgeGatewayName: types.StringValue(network.OpenApiOrgVdcNetwork.Connection.RouterRef.Name),
		InterfaceType:   types.StringValue(network.OpenApiOrgVdcNetwork.Connection.ConnectionType),
		Gateway:         types.StringValue(network.OpenApiOrgVdcNetwork.Subnets.Values[0].Gateway),
		PrefixLength:    types.Int64Value(int64(network.OpenApiOrgVdcNetwork.Subnets.Values[0].PrefixLength)),
		DNS1:            utils.StringValueOrNull(network.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1),
		DNS2:            utils.StringValueOrNull(network.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2),
		DNSSuffix:       utils.StringValueOrNull(network.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix),
	}
}
