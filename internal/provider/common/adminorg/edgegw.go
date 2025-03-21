/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package adminorg

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
)

var _ edgegw.Handler = (*AdminOrg)(nil)

// GetEdgeGateway returns the edge gateway.
func (ao *AdminOrg) GetEdgeGateway(egw edgegw.BaseEdgeGW) (edgegw.EdgeGateway, error) {
	if egw.GetIDOrName() == "" {
		return edgegw.EdgeGateway{}, edgegw.ErrEdgeGatewayIDOrNameIsEmpty
	}

	var (
		govcdValues *govcd.NsxtEdgeGateway
		err         error
	)

	if egw.GetID() != "" {
		govcdValues, err = ao.GetNsxtEdgeGatewayById(egw.GetID())
	} else {
		govcdValues, err = ao.GetNsxtEdgeGatewayByName(egw.GetName())
	}
	if err != nil {
		return edgegw.EdgeGateway{}, err
	}

	return edgegw.EdgeGateway{
		Client:          ao.c,
		NsxtEdgeGateway: govcdValues,
	}, err
}
