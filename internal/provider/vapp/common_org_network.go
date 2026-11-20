// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vapp

import (
	"fmt"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

type orgNetworkModel struct {
	ID          types.String `tfsdk:"id"`
	VAppName    types.String `tfsdk:"vapp_name"`
	VAppID      types.String `tfsdk:"vapp_id"`
	VDC         types.String `tfsdk:"vdc"`
	NetworkName types.String `tfsdk:"network_name"`
}

// vAppNetworkID returns the identifier of a vApp network: its own ID, then the UUID of
// its HREF, then the org network it is attached to (ParentNetwork). The runtime
// networkConfigSection does not always expose an id or an href on the NetworkConfig.
func vAppNetworkID(network sdktypes.ModelVAppRuntimeNetworkConfig) string {
	if network.ID != "" {
		return network.ID
	}
	if uuid := urn.ExtractUUID(network.HREF); uuid != "" {
		return uuid
	}
	if parent := network.Configuration.ParentNetwork; parent != nil {
		if parent.ID != "" {
			return parent.ID
		}
		if uuid := urn.ExtractUUID(parent.HREF); uuid != "" {
			return uuid
		}
	}
	return ""
}

func (s *orgNetworkModel) findOrgNetwork(vAppNetworkConfig []sdktypes.ModelVAppRuntimeNetworkConfig) (*sdktypes.ModelVAppRuntimeNetworkConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	for _, networkConfig := range vAppNetworkConfig {
		networkID := vAppNetworkID(networkConfig)

		idMatch := s.ID.ValueString() != "" && networkID != "" && urn.ExtractUUID(networkID) == urn.ExtractUUID(s.ID.ValueString())
		nameMatch := !s.NetworkName.IsNull() && networkConfig.NetworkName == s.NetworkName.ValueString()
		if !idMatch && !nameMatch {
			continue
		}

		if networkID == "" {
			diags.AddError("Network ID not found", fmt.Sprintf("vApp network %q was found but exposes no ID, href or parent network", networkConfig.NetworkName))
			return nil, diags
		}

		networkConfig.ID = networkID
		return &networkConfig, nil
	}

	diags.AddError("Network not found in vApp", "No vApp network matched ID or name")
	return nil, diags
}
