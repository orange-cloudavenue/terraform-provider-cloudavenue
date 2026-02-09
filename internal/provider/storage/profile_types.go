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

package storage

import "github.com/hashicorp/terraform-plugin-framework/types"

type profileDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	VDC                 types.String `tfsdk:"vdc"`
	Name                types.String `tfsdk:"name"`
	Limit               types.Int64  `tfsdk:"limit"`
	UsedStorage         types.Int64  `tfsdk:"used_storage"`
	Default             types.Bool   `tfsdk:"default"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	IopsAllocated       types.Int64  `tfsdk:"iops_allocated"`
	Units               types.String `tfsdk:"units"`
	IopsLimitingEnabled types.Bool   `tfsdk:"iops_limiting_enabled"`
	MaximumDiskIops     types.Int64  `tfsdk:"maximum_disk_iops"`
	DefaultDiskIops     types.Int64  `tfsdk:"default_disk_iops"`
	DiskIopsPerGbMax    types.Int64  `tfsdk:"disk_iops_per_gb_max"`
	IopsLimit           types.Int64  `tfsdk:"iops_limit"`
}
