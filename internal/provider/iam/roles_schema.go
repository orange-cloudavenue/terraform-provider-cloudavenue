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

package iam

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func rolesSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_iam_roles` data source allows you to retrieve information about the roles available in the organization.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Generated ID of the roles.",
				},
			},
			"roles": superschema.SuperMapNestedAttributeOf[RoleDataSourceModel]{
				DataSource: &schemaD.MapNestedAttribute{
					MarkdownDescription: "Map of the roles available in the organization.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The ID of the role.",
							Computed:            true,
						},
					},
					"name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{},
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the role.",
							Computed:            true,
						},
					},
					"description": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "A description of the role.",
							Computed:            true,
						},
					},
					"rights": superschema.SuperSetAttributeOf[string]{
						DataSource: &schemaD.SetAttribute{
							MarkdownDescription: "A list of rights for the role.",
							ElementType:         supertypes.StringType{},
							Computed:            true,
						},
					},
					"read_only": superschema.SuperBoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "Indicates if the role is read only",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
