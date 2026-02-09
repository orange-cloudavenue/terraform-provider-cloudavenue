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
)

type RoleResourceModel struct {
	ID          supertypes.StringValue        `tfsdk:"id"`
	Name        supertypes.StringValue        `tfsdk:"name"`
	Description supertypes.StringValue        `tfsdk:"description"`
	Rights      supertypes.SetValueOf[string] `tfsdk:"rights"`
}

type RoleDataSourceModel struct {
	ID          supertypes.StringValue        `tfsdk:"id"`
	Name        supertypes.StringValue        `tfsdk:"name"`
	Description supertypes.StringValue        `tfsdk:"description"`
	ReadOnly    supertypes.BoolValue          `tfsdk:"read_only"`
	Rights      supertypes.SetValueOf[string] `tfsdk:"rights"`
}
