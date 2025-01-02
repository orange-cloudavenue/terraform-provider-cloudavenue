package vdc

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	networkIsolatedModel struct {
		ID               supertypes.StringValue                                              `tfsdk:"id"`
		Name             supertypes.StringValue                                              `tfsdk:"name"`
		Description      supertypes.StringValue                                              `tfsdk:"description"`
		VDC              supertypes.StringValue                                              `tfsdk:"vdc"`
		Gateway          supertypes.StringValue                                              `tfsdk:"gateway"`
		PrefixLength     supertypes.Int64Value                                               `tfsdk:"prefix_length"`
		DNS1             supertypes.StringValue                                              `tfsdk:"dns1"`
		DNS2             supertypes.StringValue                                              `tfsdk:"dns2"`
		DNSSuffix        supertypes.StringValue                                              `tfsdk:"dns_suffix"`
		StaticIPPool     supertypes.SetNestedObjectValueOf[networkIsolatedModelStaticIPPool] `tfsdk:"static_ip_pool"`
		GuestVLANAllowed supertypes.BoolValue                                                `tfsdk:"guest_vlan_allowed"`
	}

	networkIsolatedModelStaticIPPool struct {
		StartAddress supertypes.StringValue `tfsdk:"start_address"`
		EndAddress   supertypes.StringValue `tfsdk:"end_address"`
	}
)

func (rm *networkIsolatedModel) Copy() *networkIsolatedModel {
	x := &networkIsolatedModel{}
	utils.ModelCopy(rm, x)
	return x
}

// ToSDK converts the model to the SDK model.
func (rm *networkIsolatedModel) ToSDK(ctx context.Context) (values *v1.VDCNetworkIsolatedModel, diags diag.Diagnostics) {
	values = &v1.VDCNetworkIsolatedModel{
		ID:          rm.ID.Get(),
		Name:        rm.Name.Get(),
		Description: rm.Description.Get(),
		Subnet: func() v1.VDCNetworkModelSubnet {
			return v1.VDCNetworkModelSubnet{
				Gateway:      rm.Gateway.Get(),
				PrefixLength: rm.PrefixLength.GetInt(),
				DNSServer1:   rm.DNS1.Get(),
				DNSServer2:   rm.DNS2.Get(),
				DNSSuffix:    rm.DNSSuffix.Get(),
				IPRanges: func() v1.VDCNetworkModelSubnetIPRanges {
					var ipRanges v1.VDCNetworkModelSubnetIPRanges

					ipPools, d := rm.StaticIPPool.Get(ctx)
					if d.HasError() {
						diags.Append(d...)
						return ipRanges
					}

					for _, ipRange := range ipPools {
						ipRanges = append(ipRanges, v1.VDCNetworkModelSubnetIPRange{
							StartAddress: ipRange.StartAddress.Get(),
							EndAddress:   ipRange.EndAddress.Get(),
						})
					}
					return ipRanges
				}(),
			}
		}(),
	}

	return values, diags
}
