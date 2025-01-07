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
