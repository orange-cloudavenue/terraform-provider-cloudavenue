package iam

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

type user interface {
	GetUser(refresh bool) (*govcd.OrgUser, error)
}

type commonUser struct {
	ID   types.String
	Name types.String
}

// GetUser.
func (c *commonUser) GetUser(a adminorg.AdminOrg, refresh bool) (*govcd.OrgUser, error) {
	return a.GetUserByNameOrId(c.GetIDOrName(), refresh)
}

// GetIDOrName.
func (c *commonUser) GetIDOrName() string {
	if c.ID.ValueString() != "" {
		return c.ID.ValueString()
	}
	return c.Name.ValueString()
}
