/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
