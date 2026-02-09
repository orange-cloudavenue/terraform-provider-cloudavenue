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

package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type catalogDataSourceModel struct {
	// BASE
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	OwnerName   types.String `tfsdk:"owner_name"`

	// SPECIFIC DATA SOURCE
	PreserveIdentityInformation types.Bool  `tfsdk:"preserve_identity_information"`
	NumberOfMedia               types.Int64 `tfsdk:"number_of_media"`
	MediaItemList               types.List  `tfsdk:"media_item_list"`
	IsShared                    types.Bool  `tfsdk:"is_shared"`
	IsPublished                 types.Bool  `tfsdk:"is_published"`
	IsLocal                     types.Bool  `tfsdk:"is_local"`
	IsCached                    types.Bool  `tfsdk:"is_cached"`
}

type catalogResourceModel struct {
	// BASE
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	OwnerName   types.String `tfsdk:"owner_name"`

	// SPECIFIC RESOURCE
	StorageProfile  types.String `tfsdk:"storage_profile"`
	DeleteForce     types.Bool   `tfsdk:"delete_force"`
	DeleteRecursive types.Bool   `tfsdk:"delete_recursive"`
}
