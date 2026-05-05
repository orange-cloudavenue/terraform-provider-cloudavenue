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
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func catalogsSuperSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The catalogs datasource show the details of all the catalogs.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Generated ID of the catalogs.",
					Computed:            true,
				},
			},
			"catalogs_name": superschema.ListAttribute{
				DataSource: &schemaD.ListAttribute{
					MarkdownDescription: "List of catalogs name.",
					Computed:            true,
					ElementType:         types.StringType,
				},
			},
			catalogsAttr: superschema.MapNestedAttribute{
				DataSource: &schemaD.MapNestedAttribute{
					MarkdownDescription: "Map of catalogs.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: catalogIDDescription,
							Computed:            true,
						},
					},
					name: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: catalogNameDescription,
							Computed:            true,
						},
					},
					createdAt: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descCatalogCreatedAt,
							Computed:            true,
						},
					},
					description: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descCatalogDescription,
							Computed:            true,
						},
					},
					ownerName: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descCatalogOwnerName,
							Computed:            true,
						},
					},
					preserveIdentityInformation: superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: descCatalogPreserveIdentityInfo,
							Computed:            true,
						},
					},
					numberOfMedia: superschema.Int64Attribute{
						DataSource: &schemaD.Int64Attribute{
							MarkdownDescription: descCatalogNumberOfMedia,
							Computed:            true,
						},
					},
					mediaItemList: superschema.ListAttribute{
						DataSource: &schemaD.ListAttribute{
							MarkdownDescription: descCatalogMediaItemList,
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
					isShared: superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: descCatalogIsShared,
							Computed:            true,
						},
					},
					isLocal: superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: descCatalogIsLocal,
							Computed:            true,
						},
					},
					isPublished: superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: descCatalogIsPublished,
							Computed:            true,
						},
					},
					isCached: superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: descCatalogIsCached,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
