/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package client

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type VAPP struct {
	*govcd.VApp
}

// GetName give you the name of the vApp.
func (v VAPP) GetName() string {
	return v.VApp.VApp.Name
}

// GetID give you the ID of the vApp.
func (v VAPP) GetID() string {
	return v.VApp.VApp.ID
}

// GetStatusCode give you the status code of the vApp.
func (v VAPP) GetStatusCode() int {
	return v.VApp.VApp.Status
}

// GetHREF give you the HREF of the vApp.
func (v VAPP) GetHREF() string {
	return v.VApp.VApp.HREF
}

// GetDescription give you the status code of the vApp.
func (v VAPP) GetDescription() string {
	return v.VApp.VApp.Description
}

// GetDeploymentLeaseInSeconds retrieves the lease duration in seconds for a deployment.
func (v VAPP) GetDeploymentLeaseInSeconds() int {
	return v.VApp.VApp.LeaseSettingsSection.DeploymentLeaseInSeconds
}

// GetStorageLeaseInSeconds retrieves the lease duration in seconds for a storage resource.
func (v VAPP) GetStorageLeaseInSeconds() int {
	return v.VApp.VApp.LeaseSettingsSection.StorageLeaseInSeconds
}

// IsVAPPOrgNetwork check if it is a vApp Org Network (not vApp network).
func (v VAPP) IsVAPPOrgNetwork(networkName string) (bool, error) {
	vAppNetworkConfig, err := v.GetNetworkConfig()
	if err != nil {
		return false, fmt.Errorf("error getting vApp networks: %w", err)
	}

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.NetworkName == networkName &&
			!govcd.IsVappNetwork(networkConfig.Configuration) {
			return true, nil
		}
	}

	return false, nil
}

// IsVAPPNetwork check if it is a vApp network (not vApp Org Network).
func (v VAPP) IsVAPPNetwork(networkName string) (bool, error) {
	x, err := v.IsVAPPOrgNetwork(networkName)
	return !x, err
}
