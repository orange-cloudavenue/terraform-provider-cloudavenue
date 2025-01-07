/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/storageprofile"
)

var _ storageprofile.Handler = (*VDC)(nil)

// GetStorageProfile returns the storage profile.
func (v *VDC) GetStorageProfile(storageProfileName string, refresh bool) (*govcdtypes.CatalogStorageProfiles, error) {
	if storageProfileName == "" {
		return nil, storageprofile.ErrStorageProfileNameIsEmpty
	}

	// Get the storage profile
	storageProfileReference, err := v.GetStorageProfileReference(storageProfileName, refresh)
	if err != nil {
		return nil, err
	}

	return &govcdtypes.CatalogStorageProfiles{VdcStorageProfile: []*govcdtypes.Reference{storageProfileReference}}, nil
}

// GetStorageProfileReference returns the storage profile reference.
func (v *VDC) GetStorageProfileReference(storageProfileName string, refresh bool) (*govcdtypes.Reference, error) {
	if storageProfileName == "" {
		return nil, storageprofile.ErrStorageProfileNameIsEmpty
	}

	for _, sp := range v.Vdc.Vdc.VdcStorageProfiles.VdcStorageProfile {
		if sp.Name == storageProfileName {
			return &govcdtypes.Reference{HREF: sp.HREF, Name: sp.Name, ID: sp.ID}, nil
		}
	}

	return nil, storageprofile.ErrStorageProfileNotFound
}

// GetStorageProfileID returns the storage profile ID.
func (v *VDC) FindStorageProfileName(storageProfileName string) (string, error) {
	refs, err := v.FindStorageProfileReference(storageProfileName)
	if err != nil {
		return "", err
	}

	return refs.ID, nil
}
