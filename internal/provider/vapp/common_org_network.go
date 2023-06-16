package vapp

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (s *orgNetworkModel) findOrgNetwork(vAppNetworkConfig *govcdtypes.NetworkConfigSection) (*govcdtypes.VAppNetworkConfiguration, *string, diag.Diagnostics) {
	// vAppNetwork govcdtypes.VAppNetworkConfiguration
	// networkID   string
	var diags diag.Diagnostics

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.ID == "" && networkConfig.Link != nil {
			// Get the network id from the HREF
			id, err := govcd.GetUuidFromHref(networkConfig.Link.HREF, false)
			if err != nil {
				break
			}
			networkConfig.ID = id
		} else if networkConfig.ID == "" && networkConfig.Link == nil {
			break
		}

		if (networkConfig.ID == s.ID.ValueString() && !s.ID.IsNull()) || (networkConfig.NetworkName == s.NetworkName.ValueString() && !s.NetworkName.IsNull()) {
			return &networkConfig, &networkConfig.ID, nil
		}
	}

	diags.AddError("Unable to find network ID or Name", "networkConfig.NetworkName or networkConfig.NetworkID is unknow")
	return nil, nil, diags
}
