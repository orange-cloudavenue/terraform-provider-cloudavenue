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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func mediaSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Catalog media allows you to",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "manage a media in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "retrieve information about a media in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: descMediaID,
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot(name), path.MatchRoot("id")),
					},
				},
			},
			catalogID: superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: catalogIDDescription,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot(catalogName), path.MatchRoot(catalogID)),
					},
				},
			},
			catalogName: superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: catalogNameDescription,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot(catalogName), path.MatchRoot(catalogID)),
					},
				},
			},
			name: superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: descMediaName,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot(name), path.MatchRoot("id")),
					},
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
	}
}
