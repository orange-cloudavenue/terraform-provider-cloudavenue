/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type orgNetworkModel struct {
	ID          types.String `tfsdk:"id"`
	VAppName    types.String `tfsdk:"vapp_name"`
	VAppID      types.String `tfsdk:"vapp_id"`
	VDC         types.String `tfsdk:"vdc"`
	NetworkName types.String `tfsdk:"network_name"`
}

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
			continue
		}

		if (networkConfig.ID == s.ID.ValueString() && !s.ID.IsNull()) || (networkConfig.NetworkName == s.NetworkName.ValueString() && !s.NetworkName.IsNull()) {
			return &networkConfig, &networkConfig.ID, nil
		}
	}

	diags.AddError("Unable to find network in the VApp", "The network was not found in the VApp")
	return nil, nil, diags
}
