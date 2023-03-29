package vapp

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
)

func (s *orgNetworkResourceModel) findOrgNetwork(vAppNetworkConfig *govcdtypes.NetworkConfigSection) (*govcdtypes.VAppNetworkConfiguration, *string, diag.Diagnostics) {
	var (
		vAppNetwork govcdtypes.VAppNetworkConfiguration
		networkID   string
		diags       diag.Diagnostics
	)

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.Link != nil {
			id, err := govcd.GetUuidFromHref(networkConfig.Link.HREF, false)
			if err != nil {
				diags.AddError("Unable to get network ID from HREF", err.Error())
				return nil, nil, diags
			}
			// name check needed for datasource to find network as don't have ID
			if common.ExtractUUID(s.ID.ValueString()) == common.ExtractUUID(id) || networkConfig.NetworkName == s.NetworkName.ValueString() {
				networkID = id
				vAppNetwork = networkConfig
				break
			}
		}
	}
	return &vAppNetwork, &networkID, nil
}
