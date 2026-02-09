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

package vrf

import supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

type tier0VrfDataSourceModel struct {
	ID           supertypes.StringValue                           `tfsdk:"id"`
	Name         supertypes.StringValue                           `tfsdk:"name"`
	Provider     supertypes.StringValue                           `tfsdk:"tier0_provider"`
	ClassService supertypes.StringValue                           `tfsdk:"class_service"`
	Services     supertypes.ListNestedObjectValueOf[segmentModel] `tfsdk:"services"`
}

type segmentModel struct {
	Service supertypes.StringValue `tfsdk:"service"`
	VLANID  supertypes.StringValue `tfsdk:"vlan_id"`
}
