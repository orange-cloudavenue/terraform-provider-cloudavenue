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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// GetVDCGroupByNameOrID returns the VDC group using the name or ID provided in the argument.
// Deprecated: Use GetVdcGroupByName or GetVdcGroupById instead.
func (ao *AdminOrg) GetVDCGroupByNameOrID(nameOrID string) (*govcd.VdcGroup, error) {
	if urn.IsVDCGroup(nameOrID) {
		return ao.GetVdcGroupById(nameOrID)
	}
	return ao.GetVdcGroupByName(nameOrID)
}
