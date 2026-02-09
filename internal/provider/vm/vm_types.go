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

package vm

import "github.com/hashicorp/terraform-plugin-framework/types"

type VMDataSourceModel struct { //nolint:revive
	ID          types.String `tfsdk:"id"`
	VDC         types.String `tfsdk:"vdc"`
	Name        types.String `tfsdk:"name"`
	VappName    types.String `tfsdk:"vapp_name"`
	VappID      types.String `tfsdk:"vapp_id"`
	Description types.String `tfsdk:"description"`
	State       types.Object `tfsdk:"state"`
	Resource    types.Object `tfsdk:"resource"`
	Settings    types.Object `tfsdk:"settings"`
}
