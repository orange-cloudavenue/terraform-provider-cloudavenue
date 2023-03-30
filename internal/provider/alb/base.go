package alb

import "github.com/vmware/go-vcloud-director/v2/govcd"

const (
	categoryName = "alb"
)

type albPool interface {
	GetID() string
	GetName() string
	GetAlbPool(string) (*govcd.NsxtAlbPool, error)
}

type base struct {
	id   string
	name string
}
