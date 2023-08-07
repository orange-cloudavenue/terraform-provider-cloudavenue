package catalog

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

const (
	categoryName = "catalog"

	catalogID   = "catalog_id"
	catalogName = "catalog_name"
)

type catalog interface {
	GetID() string
	GetName() string
	GetIDOrName() string
	GetCatalog() (*govcd.AdminCatalog, error)
}

type base struct {
	id   string
	name string
}

func (b base) GetID() string {
	return b.id
}

func (b base) GetName() string {
	return b.name
}

func (b base) GetIDOrName() string {
	if b.id != "" {
		return b.id
	}
	return b.name
}
