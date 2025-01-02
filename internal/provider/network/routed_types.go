package network

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	networkRoutedModel struct {
		ID              supertypes.StringValue                                            `tfsdk:"id"`
		Name            supertypes.StringValue                                            `tfsdk:"name"`
		Description     supertypes.StringValue                                            `tfsdk:"description"`
		EdgeGatewayID   supertypes.StringValue                                            `tfsdk:"edge_gateway_id"`
		EdgeGatewayName supertypes.StringValue                                            `tfsdk:"edge_gateway_name"`
		InterfaceType   supertypes.StringValue                                            `tfsdk:"interface_type"`
		Gateway         supertypes.StringValue                                            `tfsdk:"gateway"`
		PrefixLength    supertypes.Int64Value                                             `tfsdk:"prefix_length"`
		DNS1            supertypes.StringValue                                            `tfsdk:"dns1"`
		DNS2            supertypes.StringValue                                            `tfsdk:"dns2"`
		DNSSuffix       supertypes.StringValue                                            `tfsdk:"dns_suffix"`
		StaticIPPool    supertypes.SetNestedObjectValueOf[networkRoutedModelStaticIPPool] `tfsdk:"static_ip_pool"`
	}
	networkRoutedModelStaticIPPool struct {
		StartAddress supertypes.StringValue `tfsdk:"start_address"`
		EndAddress   supertypes.StringValue `tfsdk:"end_address"`
	}
)

func (rm *networkRoutedModel) Copy() *networkRoutedModel {
	x := &networkRoutedModel{}
	utils.ModelCopy(rm, x)
	return x
}
