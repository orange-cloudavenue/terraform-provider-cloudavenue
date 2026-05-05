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
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

/*
catalogSchema

This function is used to create the schema for the catalog resource and datasource.
*/
func catalogSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Catalog allows you to",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "manage a catalog in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "retrieve information about a catalog in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: catalogIDDescription,
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
			name: superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: catalogNameDescription,
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot(name), path.MatchRoot("id")),
					},
				},
			},
			createdAt: superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: descCatalogCreatedAt,
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			description: superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: descCatalogDescription,
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			ownerName: superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: descCatalogOwnerName,
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
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
			storageProfile: superschema.StringAttribute{
				// TODO - this is a reference to a storage profile, not a string
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "Storage profile to override the VM default one.",
					Optional:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"delete_force": superschema.BoolAttribute{
				Resource: &schemaR.BoolAttribute{
					Required:            true,
					MarkdownDescription: "When destroying a catalog, use `delete_force=True` along with `delete_recursive=True` to remove the catalog and any contained objects, regardless of their state.",
				},
			},
			"delete_recursive": superschema.BoolAttribute{
				Resource: &schemaR.BoolAttribute{
					Required:            true,
					MarkdownDescription: "When destroying a catalog, use `delete_recursive=True to remove the catalog and any contained objects that are in a state permitting removal.",
				},
			},
		},
	}
}
