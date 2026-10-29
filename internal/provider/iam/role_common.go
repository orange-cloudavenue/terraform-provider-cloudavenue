// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	iam "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/iam/v1"
	v2errors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

// roleClientV2 holds the SDK v2 IAM client for org role operations.
type roleClientV2 struct {
	client    *client.CloudAvenue
	iamClient *iam.Client
}

// InitRoleClientV2 initializes the v2 IAM client for role operations.
func InitRoleClientV2(c *client.CloudAvenue) (roleClientV2, diag.Diagnostics) {
	var diags diag.Diagnostics
	iamClient, err := iam.New(c.CAVSDKV2)
	if err != nil {
		diags.AddError("Error initializing IAM client", err.Error())
		return roleClientV2{}, diags
	}
	return roleClientV2{
		client:    c,
		iamClient: iamClient,
	}, diags
}

// GetOrgRole retrieves an org-scoped role by ID or name.
// If the value looks like a URN (contains "urn:vcloud:role:"), it extracts
// the ID and looks up by ID; otherwise lists all roles and matches by name.
func (r *roleClientV2) GetOrgRole(ctx context.Context, idOrName string) (*sdktypes.ModelRole, error) {
	// Check if it looks like a role URN
	if strings.Contains(idOrName, "urn:vcloud:role:") {
		// Extract ID from URN if needed, then lookup by ID
		roleID := idOrName
		return r.iamClient.GetOrgRole(ctx, sdktypes.ParamsGetOrgRole{ID: roleID})
	}

	// Not a URN, look up by name
	roles, err := r.iamClient.ListOrgRoles(ctx)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.Name == idOrName {
			return role, nil
		}
	}

	return nil, v2errors.ErrNotFound
}

// ListOrgRoles lists all org-scoped roles.
func (r *roleClientV2) ListOrgRoles(ctx context.Context) ([]*sdktypes.ModelRole, error) {
	return r.iamClient.ListOrgRoles(ctx)
}

// ListOrgRoleRights lists the rights of an org-scoped role.
func (r *roleClientV2) ListOrgRoleRights(ctx context.Context, roleID string) ([]*sdktypes.ModelRight, error) {
	return r.iamClient.ListOrgRoleRights(ctx, sdktypes.ParamsListOrgRoleRights{ID: roleID})
}

// CreateOrgRole creates a new org-scoped role.
func (r *roleClientV2) CreateOrgRole(ctx context.Context, name, description string) (*sdktypes.ModelRole, error) {
	return r.iamClient.CreateOrgRole(ctx, sdktypes.ParamsCreateOrgRole{
		Name:        name,
		Description: description,
	})
}

// UpdateOrgRole updates an existing org-scoped role.
func (r *roleClientV2) UpdateOrgRole(ctx context.Context, id, name, description string) (*sdktypes.ModelRole, error) {
	return r.iamClient.UpdateOrgRole(ctx, sdktypes.ParamsUpdateOrgRole{
		ID:          id,
		Name:        name,
		Description: description,
	})
}

// DeleteOrgRole deletes an org-scoped role.
func (r *roleClientV2) DeleteOrgRole(ctx context.Context, id string) error {
	return r.iamClient.DeleteOrgRole(ctx, sdktypes.ParamsDeleteOrgRole{ID: id})
}

// ReplaceOrgRoleRights replaces the rights of an org-scoped role.
func (r *roleClientV2) ReplaceOrgRoleRights(ctx context.Context, roleID string, rights []sdktypes.ParamsEntityReference) ([]*sdktypes.ModelRight, error) {
	return r.iamClient.ReplaceOrgRoleRights(ctx, sdktypes.ParamsReplaceOrgRoleRights{
		ID:     roleID,
		Rights: rights,
	})
}

// RemoveAllOrgRoleRights removes all rights from an org-scoped role.
func (r *roleClientV2) RemoveAllOrgRoleRights(ctx context.Context, roleID string) error {
	return r.iamClient.RemoveAllOrgRoleRights(ctx, sdktypes.ParamsDeleteOrgRole{ID: roleID})
}

// ListRights lists all rights, optionally filtered.
func (r *roleClientV2) ListRights(ctx context.Context, filter string) ([]*sdktypes.ModelRight, error) {
	return r.iamClient.ListRights(ctx, sdktypes.ParamsListRights{Filter: filter})
}
