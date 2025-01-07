/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type cloudavenueProviderModel struct {
	URL               types.String `tfsdk:"url"`
	User              types.String `tfsdk:"user"`
	Password          types.String `tfsdk:"password"`
	Org               types.String `tfsdk:"org"`
	VDC               types.String `tfsdk:"vdc"`
	NetBackupURL      types.String `tfsdk:"netbackup_url"`
	NetBackupUser     types.String `tfsdk:"netbackup_user"`
	NetBackupPassword types.String `tfsdk:"netbackup_password"`
}
