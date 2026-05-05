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
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func mediasSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Catalog medias allows you to",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "manage a medias in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "retrieve information about a medias in Cloud Avenue.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The ID of the medias.",
					Computed:            true,
				},
			},
			"medias": superschema.MapNestedAttribute{
				DataSource: &schemaD.MapNestedAttribute{
					MarkdownDescription: "The map of medias.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descMediaID,
							Computed:            true,
						},
					},
					catalogID: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: catalogIDDescription,
							Computed:            true,
						},
					},
					catalogName: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: catalogNameDescription,
							Computed:            true,
						},
					},
					name: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descMediaName,
							Computed:            true,
						},
					},
					description: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descMediaDescription,
							Computed:            true,
						},
					},
					isISO: superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: descMediaIsISO,
							Computed:            true,
						},
					},
					ownerName: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descMediaOwnerName,
							Computed:            true,
						},
					},
					isPublished: superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: descMediaIsPublished,
							Computed:            true,
						},
					},
					createdAt: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descMediaCreatedAt,
							Computed:            true,
						},
					},
					size: superschema.Int64Attribute{
						DataSource: &schemaD.Int64Attribute{
							MarkdownDescription: descMediaSize,
							Computed:            true,
						},
					},
					status: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descMediaStatus,
							Computed:            true,
						},
					},
					storageProfile: superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: descMediaStorageProfile,
							Computed:            true,
						},
					},
				},
			},
			"medias_name": superschema.ListAttribute{
				DataSource: &schemaD.ListAttribute{
					MarkdownDescription: "The list of medias name.",
					Computed:            true,
					ElementType:         types.StringType,
				},
			},
			catalogName: mediaSchema().Attributes[catalogName],
			catalogID:   mediaSchema().Attributes[catalogID],
		},
	}
}
