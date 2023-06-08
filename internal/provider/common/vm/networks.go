package vm

import (
	"context"
	"fmt"
	"sort"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type VMResourceModelResourceNetworks []VMResourceModelResourceNetwork //nolint:revive

type VMResourceModelResourceNetwork struct { //nolint:revive
	Type             types.String `tfsdk:"type"`
	IPAllocationMode types.String `tfsdk:"ip_allocation_mode"`
	Name             types.String `tfsdk:"name"`
	IP               types.String `tfsdk:"ip"`
	IsPrimary        types.Bool   `tfsdk:"is_primary"`
	Mac              types.String `tfsdk:"mac"`
	AdapterType      types.String `tfsdk:"adapter_type"`
	Connected        types.Bool   `tfsdk:"connected"`
}

// attrTypes() returns the types of the attributes of the Networks attribute.
func (n *VMResourceModelResourceNetworks) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":               types.StringType,
		"ip_allocation_mode": types.StringType,
		"ip":                 types.StringType,
		"name":               types.StringType,
		"is_primary":         types.BoolType,
		"mac":                types.StringType,
		"adapter_type":       types.StringType,
		"connected":          types.BoolType,
	}
}

// ObjectType() returns the type of the Networks attribute.
func (n *VMResourceModelResourceNetworks) ObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: n.AttrTypes(),
	}
}

// toAttrValues() returns the values of the attributes of the Networks attribute.
func (n *VMResourceModelResourceNetwork) toAttrValues() map[string]attr.Value { //nolint:unused
	return map[string]attr.Value{
		"type":               n.Type,
		"ip_allocation_mode": n.IPAllocationMode,
		"ip":                 n.IP,
		"name":               n.Name,
		"is_primary":         n.IsPrimary,
		"mac":                n.Mac,
		"adapter_type":       n.AdapterType,
		"connected":          n.Connected,
	}
}

// ToPlan returns the value of the Networks attribute, if set, as a types.Object.
func (n *VMResourceModelResourceNetworks) ToPlan(ctx context.Context) (basetypes.ListValue, diag.Diagnostics) {
	if n == nil {
		return types.ListNull(n.ObjectType()), diag.Diagnostics{}
	}
	return types.ListValueFrom(context.Background(), n.ObjectType(), n)
}

// Equal returns true if the two VMResourceModelResourceNetwork are equal.
func (n *VMResourceModelResourceNetwork) Equal(other VMResourceModelResourceNetwork) bool {
	return n.Type.Equal(other.Type) &&
		n.IPAllocationMode.Equal(other.IPAllocationMode) &&
		n.IP.Equal(other.IP) &&
		n.Name.Equal(other.Name) &&
		n.IsPrimary.Equal(other.IsPrimary) &&
		n.Mac.Equal(other.Mac) &&
		n.AdapterType.Equal(other.AdapterType) &&
		n.Connected.Equal(other.Connected)
}

// ConvertToNetworkConnection converts a VMResourceModelResourceNetworks to a NetworkConnection.
func (n *VMResourceModelResourceNetwork) ConvertToNetworkConnection() NetworkConnection {
	return NetworkConnection{
		Name:             n.Name,
		Connected:        n.Connected,
		IPAllocationMode: n.IPAllocationMode,
		IP:               n.IP,
		IsPrimary:        n.IsPrimary,
		Mac:              n.Mac,
		AdapterType:      n.AdapterType,
		Type:             n.Type,
	}
}

// NetworksRead returns network configuration for saving into statefile.
func (v VM) NetworksRead() (*VMResourceModelResourceNetworks, error) {
	vapp, err := v.GetParentVApp()
	if err != nil {
		return nil, fmt.Errorf("error getting vApp: %w", err)
	}

	// Determine type for all networks in vApp
	vAppNetworkConfig, err := vapp.GetNetworkConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting vApp networks: %w", err)
	}
	// If vApp network is "isolated" and has no ParentNetwork - it is a vApp network.
	// https://code.vmware.com/apis/72/vcloud/doc/doc/types/NetworkConfigurationType.html
	vAppNetworkTypes := make(map[string]string, 0)
	for _, netConfig := range vAppNetworkConfig.NetworkConfig {
		switch {
		case netConfig.NetworkName == govcdtypes.NoneNetwork:
			vAppNetworkTypes[netConfig.NetworkName] = govcdtypes.NoneNetwork
		case govcd.IsVappNetwork(netConfig.Configuration):
			vAppNetworkTypes[netConfig.NetworkName] = "vapp"
		default:
			vAppNetworkTypes[netConfig.NetworkName] = "org"
		}
	}

	nets := make(VMResourceModelResourceNetworks, 0)

	if v.NetworkConnectionIsDefined() {
		// Sort NIC cards by their virtual slot numbers as the API returns them in random order
		sort.SliceStable(v.GetNetworkConnection(), func(i, j int) bool {
			return v.GetNetworkConnection()[i].NetworkConnectionIndex <
				v.GetNetworkConnection()[j].NetworkConnectionIndex
		})

		for _, vmNet := range v.GetNetworkConnection() {
			singleNIC := VMResourceModelResourceNetwork{
				IPAllocationMode: types.StringValue(vmNet.IPAddressAllocationMode),
				IP:               utils.StringValueOrNull(vmNet.IPAddress),
				Mac:              types.StringValue(vmNet.MACAddress),
				AdapterType:      types.StringValue(vmNet.NetworkAdapterType),
				Connected:        types.BoolValue(vmNet.IsConnected),
				IsPrimary:        types.BoolValue(false),
				Type:             types.StringValue(vAppNetworkTypes[vmNet.Network]),
			}

			if vmNet.Network != govcdtypes.NoneNetwork {
				singleNIC.Name = types.StringValue(vmNet.Network)
			}

			if vmNet.NetworkConnectionIndex == v.VM.VM.VM.NetworkConnectionSection.PrimaryNetworkConnectionIndex {
				singleNIC.IsPrimary = types.BoolValue(true)
			}

			nets = append(nets, singleNIC)
		}
	}

	return &nets, nil
}
