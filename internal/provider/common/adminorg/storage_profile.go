package adminorg

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/storageprofile"
)

var _ storageprofile.Handler = (*AdminOrg)(nil)

// GetStorageProfile returns the storage profile.
func (ao *AdminOrg) GetStorageProfile(storageProfileID string, refresh bool) (*govcdtypes.CatalogStorageProfiles, error) {
	if storageProfileID == "" {
		return nil, storageprofile.ErrStorageProfileIDIsEmpty
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
		return nil, storageprofile.ErrStorageProfileIDIsEmpty
	}

	return ao.GetStorageProfileReferenceById(storageProfileID, refresh)
}

// GetStorageProfileID returns the storage profile ID.
func (ao *AdminOrg) FindStorageProfileID(storageProfileName string) (string, error) {
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
