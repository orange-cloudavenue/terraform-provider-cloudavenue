package vapp

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
)

func (s *orgNetworkModel) findOrgNetwork(vAppNetworkConfig *govcdtypes.NetworkConfigSection) (*govcdtypes.VAppNetworkConfiguration, *string, diag.Diagnostics) {
	var (
		vAppNetwork govcdtypes.VAppNetworkConfiguration
		networkID   string
		diags       diag.Diagnostics
	)

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		// If Link is not nill then we can get the ID from the HREF
		if networkConfig.Link != nil {
			// Get the network id from the HREF
			id, err := govcd.GetUuidFromHref(networkConfig.Link.HREF, false)
			if err != nil {
				diags.AddError("Unable to get network ID from HREF", err.Error())
				return nil, nil, diags
			}
			networkID = id
			// name check needed for datasource to find network as don't have ID
			if (common.ExtractUUID(s.ID.ValueString()) == common.ExtractUUID(id) && !s.ID.IsNull()) || (networkConfig.NetworkName == s.NetworkName.ValueString() && !s.NetworkName.IsNull()) {
				vAppNetwork = networkConfig
				break
			} else { // return error when networkConfig.NetworkName or networkConfig.NetworkID is unknow
				diags.AddError("Unable to find network ID or Name", "networkConfig.NetworkName or networkConfig.NetworkID is unknow")
				return nil, nil, diags
			}
		} else { // return error when networkConfig.Link is nil
			diags.AddError("Unable to get network ID from HREF", "networkConfig.Link is nil")
			return nil, nil, diags
		}
	}
	return &vAppNetwork, &networkID, nil
}
