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

package s3

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type CredentialModel struct {
	ID            supertypes.StringValue `tfsdk:"id"`
	Username      supertypes.StringValue `tfsdk:"username"`
	FileName      supertypes.StringValue `tfsdk:"file_name"`
	SaveInFile    supertypes.BoolValue   `tfsdk:"save_in_file"`
	PrintToken    supertypes.BoolValue   `tfsdk:"print_token"`
	SaveInTFState supertypes.BoolValue   `tfsdk:"save_in_tfstate"`
	AccessKey     supertypes.StringValue `tfsdk:"access_key"`
	SecretKey     supertypes.StringValue `tfsdk:"secret_key"`
}

func (rm *CredentialModel) Copy() *CredentialModel {
	x := &CredentialModel{}
	utils.ModelCopy(rm, x)
	return x
}
