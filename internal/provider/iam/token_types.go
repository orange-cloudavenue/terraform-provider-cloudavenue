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

type TokenModel struct {
	FileName      supertypes.StringValue `tfsdk:"file_name"`
	ID            supertypes.StringValue `tfsdk:"id"`
	Name          supertypes.StringValue `tfsdk:"name"`
	PrintToken    supertypes.BoolValue   `tfsdk:"print_token"`
	SaveInFile    supertypes.BoolValue   `tfsdk:"save_in_file"`
	SaveInTfstate supertypes.BoolValue   `tfsdk:"save_in_tfstate"`
	Token         supertypes.StringValue `tfsdk:"token"`
}

func (rm *TokenModel) Copy() *TokenModel {
	x := &TokenModel{}
	utils.ModelCopy(rm, x)
	return x
}
