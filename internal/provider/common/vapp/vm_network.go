package vapp

import (
	"fmt"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

// GetNetworkConnection give you the network connection of the vApp.
func (v VAPP) constructNetworkConnection(networks vm.Networks) (networkConnection govcdtypes.NetworkConnectionSection, err error) {
	for index, network := range networks {
		netCon := &govcdtypes.NetworkConnection{
			Network:                 network.Name.ValueString(),
			IsConnected:             network.Connected.ValueBool(),
			IPAddressAllocationMode: network.IPAllocationMode.ValueString(),
			IPAddress:               network.IP.ValueString(),
			NetworkConnectionIndex:  index,
		}

		switch network.Type.ValueString() {
		case "vapp":
			if ok, err := v.IsVAPPNetwork(network.Name.ValueString()); err != nil {
				return govcdtypes.NetworkConnectionSection{}, err
			} else if !ok {
				return govcdtypes.NetworkConnectionSection{}, fmt.Errorf("vApp network : %s is not found", network.Name.ValueString())
			}
		case "org":
			if ok, err := v.IsVAPPOrgNetwork(network.Name.ValueString()); err != nil {
				return govcdtypes.NetworkConnectionSection{}, err
			} else if !ok {
				return govcdtypes.NetworkConnectionSection{}, fmt.Errorf("vApp Org network : %s is not found", network.Name.ValueString())
			}
		}

		if network.Mac.ValueString() != "" {
			netCon.MACAddress = network.Mac.ValueString()
		}

		if network.AdapterType.ValueString() != "" {
			netCon.NetworkAdapterType = network.AdapterType.ValueString()
		}

		networkConnection.NetworkConnection = append(networkConnection.NetworkConnection, netCon)

		if network.IsPrimary.ValueBool() {
			networkConnection.PrimaryNetworkConnectionIndex = index
		}

	}

	return networkConnection, nil
}
