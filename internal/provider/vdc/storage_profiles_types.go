/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type (
	storageProfilesDataSourceModel struct {
		ID              supertypes.StringValue                                                          `tfsdk:"id"`
		VDCName         supertypes.StringValue                                                          `tfsdk:"vdc_name"`
		VDCID           supertypes.StringValue                                                          `tfsdk:"vdc_id"`
		StorageProfiles supertypes.ListNestedObjectValueOf[storageProfileDataSourceModelStorageProfile] `tfsdk:"storage_profiles"`
	}
	storageProfileDataSourceModelStorageProfile struct {
		ID      supertypes.StringValue `tfsdk:"id"`
		Class   supertypes.StringValue `tfsdk:"class"`
		Limit   supertypes.Int64Value  `tfsdk:"limit"`
		Used    supertypes.Int64Value  `tfsdk:"used"`
		Default supertypes.BoolValue   `tfsdk:"default"`
	}
)
