package iam

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

type role interface {
	GetRole() (*govcd.Role, error)
}

type commonRole struct {
	ID   types.String
	Name types.String
}

// GetRole.
func (c *commonRole) GetRole(a adminorg.AdminOrg) (*govcd.Role, error) {
	var (
		role *govcd.Role
		err  error
	)
	// Get the role
	if c.ID.IsNull() {
		role, err = a.GetRoleByName(c.Name.ValueString())
	} else {
		role, err = a.GetRoleById(c.ID.ValueString())
	}
	if err != nil {
		return nil, err
	}

	return role, nil
}
