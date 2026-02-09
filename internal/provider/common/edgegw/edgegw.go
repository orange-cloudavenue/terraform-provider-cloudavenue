/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	sdkv1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var ErrEdgeGatewayIDOrNameIsEmpty = fmt.Errorf("edge gateway ID or name is empty")

type Handler interface {
	// GetEdgeGateway allows retrieving NSX-T edge gateway by ID Or Name.
	GetEdgeGateway(egw BaseEdgeGW) (EdgeGateway, error)
}

type BaseEdgeGW struct {
	ID   types.String
	Name types.String
}

type EdgeGateway struct {
	// Client is the CloudAvenue client.

	Client *client.CloudAvenue

	// EdgeClient is the EdgeGateway client.
	*sdkv1.EdgeClient

	// NsxtEdgeGateway is the NSX-T edge gateway.
	//
	// Deprecated: Use EdgeClient instead.
	*govcd.NsxtEdgeGateway
}

// GetID returns the ID of the Edge Gateway.
func (e BaseEdgeGW) GetID() string {
	return e.ID.ValueString()
}

// GetName returns the Edge Gateway Name.
func (e BaseEdgeGW) GetName() string {
	return e.Name.ValueString()
}

// GetIDOrName returns the ID or the name of the Edge Gateway.
func (e BaseEdgeGW) GetIDOrName() string {
	if e.GetID() != "" {
		return e.GetID()
	}
	return e.GetName()
}

// GetParent returns the parent of the Edge Gateway.
// The parent is VDC Or VDC Group object.
func (e EdgeGateway) GetParent() (sdkv1.VDCOrVDCGroupInterface, error) {
	return e.Client.CAVSDK.V1.VDC().GetVDCOrVDCGroup(e.EdgeGateway.OwnerRef.Name)
}

// GetName returns the name of the Edge Gateway.
func (e EdgeGateway) GetName() string {
	return e.EdgeGateway.Name
}

// GetID returns the ID of the Edge Gateway.
func (e EdgeGateway) GetID() string {
	return e.EdgeGateway.ID
}

// GetSecurityGroupByNameOrID.
func (e EdgeGateway) GetSecurityGroupByNameOrID(nsxtFirewallGroupNameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if err := e.Refresh(); err != nil {
		return nil, err
	}

	if urn.IsValid(nsxtFirewallGroupNameOrID) {
		return e.GetNsxtFirewallGroupById(nsxtFirewallGroupNameOrID)
	}

	return e.GetNsxtFirewallGroupByName(nsxtFirewallGroupNameOrID, govcdtypes.FirewallGroupTypeSecurityGroup)
}

// GetIPSetByID.
func (e EdgeGateway) GetIPSetByID(nsxtFirewallGroupID string) (*govcd.NsxtFirewallGroup, error) {
	if err := e.Refresh(); err != nil {
		return nil, err
	}

	return e.GetNsxtFirewallGroupById(nsxtFirewallGroupID)
}

// GetIPSetByName.
func (e EdgeGateway) GetIPSetByName(nsxtFirewallGroupName string) (*govcd.NsxtFirewallGroup, error) {
	if err := e.Refresh(); err != nil {
		return nil, err
	}

	return e.GetNsxtFirewallGroupByName(nsxtFirewallGroupName, govcdtypes.FirewallGroupTypeIpSet)
}

// GetIPSetByNameOrID.
func (e EdgeGateway) GetIPSetByNameOrID(nsxtFirewallGroupNameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if err := e.Refresh(); err != nil {
		return nil, err
	}

	if urn.IsValid(nsxtFirewallGroupNameOrID) {
		return e.GetIPSetByID(nsxtFirewallGroupNameOrID)
	}

	return e.GetIPSetByName(nsxtFirewallGroupNameOrID)
}

// SetIPSet.
func (e EdgeGateway) SetIPSet(ipSetConfig *govcdtypes.NsxtFirewallGroup) (*govcd.NsxtFirewallGroup, error) {
	return e.CreateNsxtFirewallGroup(ipSetConfig)
}
