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
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func vdcsSchema() superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Retrieve all Virtual Data Centers (vDCs) within an organization.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The ID of the resource. This value is system-generated.",
					Computed:            true,
				},
			},
			"vdcs": superschema.SuperListNestedAttributeOf[vdcsDataSourceModelVDC]{
				DataSource: &schemaD.ListNestedAttribute{
					MarkdownDescription: "A list of Virtual Data Centers available within the organization.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The unique identifier of the Virtual Data Center.",
							Computed:            true,
						},
					},
					"name": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the vDC.",
							Computed:            true,
						},
					},
					"description": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The description of the vDC.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
