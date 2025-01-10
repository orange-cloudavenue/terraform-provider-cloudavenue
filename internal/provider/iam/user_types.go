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
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type userResourceModel struct {
	// Base
	ID              supertypes.StringValue `tfsdk:"id"`
	Name            supertypes.StringValue `tfsdk:"name"`
	RoleName        supertypes.StringValue `tfsdk:"role_name"`
	FullName        supertypes.StringValue `tfsdk:"full_name"`
	Email           supertypes.StringValue `tfsdk:"email"`
	Telephone       supertypes.StringValue `tfsdk:"telephone"`
	Enabled         supertypes.BoolValue   `tfsdk:"enabled"`
	DeployedVMQuota supertypes.Int64Value  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   supertypes.Int64Value  `tfsdk:"stored_vm_quota"`

	// Specific
	Password      supertypes.StringValue `tfsdk:"password"`
	TakeOwnership supertypes.BoolValue   `tfsdk:"take_ownership"`
}

type userDataSourceModel struct {
	// Base
	ID              supertypes.StringValue `tfsdk:"id"`
	Name            supertypes.StringValue `tfsdk:"name"`
	RoleName        supertypes.StringValue `tfsdk:"role_name"`
	FullName        supertypes.StringValue `tfsdk:"full_name"`
	Email           supertypes.StringValue `tfsdk:"email"`
	Telephone       supertypes.StringValue `tfsdk:"telephone"`
	Enabled         supertypes.BoolValue   `tfsdk:"enabled"`
	DeployedVMQuota supertypes.Int64Value  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   supertypes.Int64Value  `tfsdk:"stored_vm_quota"`

	// Specific
	ProviderType supertypes.StringValue `tfsdk:"provider_type"`
}

func (rm *userResourceModel) Copy() *userResourceModel {
	x := &userResourceModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (dm *userDataSourceModel) Copy() *userDataSourceModel {
	x := &userDataSourceModel{}
	utils.ModelCopy(dm, x)
	return x
}
