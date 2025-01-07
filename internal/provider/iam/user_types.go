/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import "github.com/hashicorp/terraform-plugin-framework/types"

type userResourceModel struct {
	// Base
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RoleName        types.String `tfsdk:"role_name"`
	FullName        types.String `tfsdk:"full_name"`
	Email           types.String `tfsdk:"email"`
	Telephone       types.String `tfsdk:"telephone"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	DeployedVMQuota types.Int64  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   types.Int64  `tfsdk:"stored_vm_quota"`

	// Specific
	Password      types.String `tfsdk:"password"`
	TakeOwnership types.Bool   `tfsdk:"take_ownership"`
}

type userDataSourceModel struct {
	// Base
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RoleName        types.String `tfsdk:"role_name"`
	FullName        types.String `tfsdk:"full_name"`
	Email           types.String `tfsdk:"email"`
	Telephone       types.String `tfsdk:"telephone"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	DeployedVMQuota types.Int64  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   types.Int64  `tfsdk:"stored_vm_quota"`

	// Specific
	ProviderType types.String `tfsdk:"provider_type"`
}
