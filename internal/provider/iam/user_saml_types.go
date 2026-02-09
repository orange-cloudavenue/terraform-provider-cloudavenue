/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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

type UserSAMLModel struct {
	ID              supertypes.StringValue `tfsdk:"id"`
	UserName        supertypes.StringValue `tfsdk:"user_name"`
	RoleName        supertypes.StringValue `tfsdk:"role_name"`
	Enabled         supertypes.BoolValue   `tfsdk:"enabled"`
	DeployedVMQuota supertypes.Int64Value  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   supertypes.Int64Value  `tfsdk:"stored_vm_quota"`
	TakeOwnership   supertypes.BoolValue   `tfsdk:"take_ownership"`
}

func (rm *UserSAMLModel) Copy() *UserSAMLModel {
	x := &UserSAMLModel{}
	utils.ModelCopy(rm, x)
	return x
}
