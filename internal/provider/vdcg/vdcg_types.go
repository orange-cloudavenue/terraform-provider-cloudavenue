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

package vdcg

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type vdcgModel struct {
	ID          supertypes.StringValue        `tfsdk:"id"`
	Name        supertypes.StringValue        `tfsdk:"name"`
	Description supertypes.StringValue        `tfsdk:"description"`
	VDCIDs      supertypes.SetValueOf[string] `tfsdk:"vdc_ids"`
	Type        supertypes.StringValue        `tfsdk:"type"`
	Status      supertypes.StringValue        `tfsdk:"status"`
}

func (rm *vdcgModel) Copy() *vdcgModel {
	x := &vdcgModel{}
	utils.ModelCopy(rm, x)
	return x
}
