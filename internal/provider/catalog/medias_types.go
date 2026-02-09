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

type (
	catalogMediasDataSourceModel struct {
		ID          types.String `tfsdk:"id"`
		Medias      types.Map    `tfsdk:"medias"`
		MediasName  types.List   `tfsdk:"medias_name"`
		CatalogName types.String `tfsdk:"catalog_name"`
		CatalogID   types.String `tfsdk:"catalog_id"`
	}

	catalogMediasDataSourceModelMedias     map[string]catalogMediaDataSourceModel
	catalogMediasDataSourceModelMediasName []string
)
