/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vm

import "github.com/hashicorp/terraform-plugin-framework/types"

type vmAffinityRuleResourceModel struct {
	ID       types.String `tfsdk:"id"`
	VDC      types.String `tfsdk:"vdc"`
	Name     types.String `tfsdk:"name"`
	Polarity types.String `tfsdk:"polarity"`
	Required types.Bool   `tfsdk:"required"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	VMIDs    types.Set    `tfsdk:"vm_ids"`
}

type vmAffinityRuleDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	VDC      types.String `tfsdk:"vdc"`
	Name     types.String `tfsdk:"name"`
	Polarity types.String `tfsdk:"polarity"`
	Required types.Bool   `tfsdk:"required"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	VMIDs    types.Set    `tfsdk:"vm_ids"`
}
