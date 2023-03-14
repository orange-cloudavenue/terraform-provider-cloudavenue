package client

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type AdminOrg struct {
	*govcd.AdminOrg
}

// GetAdminOrg return the admin org using the name provided in the provider.
func (c *CloudAvenue) GetAdminOrg() (*AdminOrg, error) {
	x, err := c.Vmware.GetAdminOrgByNameOrId(c.Org)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRetrievingOrg, err)
	}

	return &AdminOrg{x}, nil
}

/*
ListCatalogs

Get the catalogs list from the admin org.
*/
func (ao *AdminOrg) ListCatalogs() *govcdtypes.CatalogsList {
	return ao.AdminOrg.AdminOrg.Catalogs
}
