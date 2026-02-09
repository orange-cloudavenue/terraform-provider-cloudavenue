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
			"catalogs": superschema.MapNestedAttribute{
				DataSource: &schemaD.MapNestedAttribute{
					MarkdownDescription: "Map of catalogs.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The ID of the catalog.",
							Computed:            true,
						},
					},
					"name": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the catalog.",
							Computed:            true,
						},
					},
					"created_at": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The creation date of the catalog.",
							Computed:            true,
						},
					},
					"description": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The description of the catalog.",
							Computed:            true,
						},
					},
					"owner_name": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The owner name of the catalog.",
							Computed:            true,
						},
					},
					"preserve_identity_information": superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "Include BIOS UUIDs and MAC addresses in the downloaded OVF package. Keep in mind that preserving this identity information reduces the package's portability, so only include it when necessary.",
							Computed:            true,
						},
					},
					"number_of_media": superschema.Int64Attribute{
						DataSource: &schemaD.Int64Attribute{
							MarkdownDescription: "The number of media in the catalog.",
							Computed:            true,
						},
					},
					"media_item_list": superschema.ListAttribute{
						DataSource: &schemaD.ListAttribute{
							MarkdownDescription: "The list of media items in the catalog.",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
					"is_shared": superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "Indicates whether the catalog is shared.",
							Computed:            true,
						},
					},
					"is_local": superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "Indicates whether the catalog is local.",
							Computed:            true,
						},
					},
					"is_published": superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "Indicates whether the catalog is published.",
							Computed:            true,
						},
					},
					"is_cached": superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "Indicates whether the catalog is cached.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
