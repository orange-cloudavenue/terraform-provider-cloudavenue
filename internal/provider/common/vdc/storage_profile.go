package vdc

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/storageprofile"
)

var _ storageprofile.Handler = (*VDC)(nil)

// GetStorageProfile returns the storage profile.
func (v *VDC) GetStorageProfile(storageProfileID string, refresh bool) (*govcdtypes.CatalogStorageProfiles, error) {
	if storageProfileID == "" {
		return nil, storageprofile.ErrStorageProfileIDIsEmpty
	}

	// Get the storage profile
	storageProfileReference, err := v.GetStorageProfileReference(storageProfileID, refresh)
	if err != nil {
		return nil, err
	}

	return &govcdtypes.CatalogStorageProfiles{VdcStorageProfile: []*govcdtypes.Reference{storageProfileReference}}, nil
}

// GetStorageProfileReference returns the storage profile reference.
func (v *VDC) GetStorageProfileReference(storageProfileID string, refresh bool) (*govcdtypes.Reference, error) {
	if storageProfileID == "" {
		return nil, storageprofile.ErrStorageProfileIDIsEmpty
	}

	for _, sp := range v.Vdc.Vdc.VdcStorageProfiles.VdcStorageProfile {
		if sp.ID == storageProfileID {
			return &govcdtypes.Reference{HREF: sp.HREF, Name: sp.Name, ID: sp.ID}, nil
		}
	}

	return nil, storageprofile.ErrStorageProfileNotFound
}

// GetStorageProfileID returns the storage profile ID.
func (v *VDC) FindStorageProfileID(storageProfileName string) (string, error) {
	refs, err := v.FindStorageProfileReference(storageProfileName)
	if err != nil {
		return "", err
	}

	return refs.ID, nil
}
