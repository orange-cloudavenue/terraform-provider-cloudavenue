/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package draas

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type draasIPResourceModel struct {
	ID        supertypes.StringValue `tfsdk:"id"`
	IPAddress supertypes.StringValue `tfsdk:"ip_address"`
}

func (rm *draasIPResourceModel) Copy() *draasIPResourceModel {
	x := &draasIPResourceModel{}
	utils.ModelCopy(rm, x)
	return x
}
