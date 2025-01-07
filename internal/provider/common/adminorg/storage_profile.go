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
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/storageprofile"
)

var _ storageprofile.Handler = (*AdminOrg)(nil)

// GetStorageProfile returns the storage profile.
func (ao *AdminOrg) GetStorageProfile(storageProfileID string, refresh bool) (*govcdtypes.CatalogStorageProfiles, error) {
	if storageProfileID == "" {
		return nil, storageprofile.ErrStorageProfileNameIsEmpty
	}

	// Get the storage profile
	storageProfileReference, err := ao.GetStorageProfileReferenceById(storageProfileID, refresh)
	if err != nil {
		return nil, err
	}

	return &govcdtypes.CatalogStorageProfiles{VdcStorageProfile: []*govcdtypes.Reference{storageProfileReference}}, nil
}

// GetStorageProfileReference returns the storage profile reference.
func (ao *AdminOrg) GetStorageProfileReference(storageProfileID string, refresh bool) (*govcdtypes.Reference, error) {
	if storageProfileID == "" {
		return nil, storageprofile.ErrStorageProfileNameIsEmpty
	}

	return ao.GetStorageProfileReferenceById(storageProfileID, refresh)
}

// GetStorageProfileID returns the storage profile ID.
func (ao *AdminOrg) FindStorageProfileName(storageProfileName string) (string, error) {
	refs, err := ao.GetAllStorageProfileReferences(true)
	if err != nil {
		return "", err
	}

	for _, ref := range refs {
		if ref.Name == storageProfileName {
			return ref.ID, nil
		}
	}

	return "", storageprofile.ErrStorageProfileNotFound
}
