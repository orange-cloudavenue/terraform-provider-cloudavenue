/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func vdcgsSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Retrieves information about one or more existing Virtual Data Center Group (VDC Group) in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Unique identifier generated for the VDC Group list.",
				},
			},
			"filter_by_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Specifies the identifier of the Virtual Data Center Group (VDCG) used to filter VDC Groups. If no filter is apply, all vdcgroup will be listed.",
					Computed:            true,
					Optional:            true,
					Validators: []validator.String{
						fstringvalidator.IsURN(),
						fstringvalidator.PrefixContains(urn.VDCGroup.String()),
					},
				},
			},
			"filter_by_name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Name of the Virtual Data Center Group (VDCG) used to filter VDC Groups. Supports partial matches using the `*` wildcard. If no filter is apply, all vdcgroup will be listed.",
					Optional:            true,
				},
			},
			"vdc_groups": superschema.SuperSetNestedAttributeOf[vdcgsModelVDCG]{
				DataSource: &schemaD.SetNestedAttribute{
					MarkdownDescription: "A set of Virtual Data Center (VDC) IDs that are members of this VDC Group.",
					Computed:            true,
				},
				Attributes: map[string]superschema.Attribute{
					"id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "Unique identifier for the Virtual Data Center Group (VDCG).",
							Computed:            true,
						},
					},
					"name": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "Name assigned to the Virtual Data Center Group (VDCG).",
							Computed:            true,
						},
					},
					"description": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "Detailed description of the Virtual Data Center Group (VDCG) and its purpose.",
							Computed:            true,
						},
					},
					"number_of_vdcs": superschema.SuperInt64Attribute{
						DataSource: &schemaD.Int64Attribute{
							MarkdownDescription: "The total number of Virtual Data Centers (VDCs) that are part of this VDC Group.",
							Computed:            true,
						},
					},
					"vdcs": superschema.SuperSetNestedAttributeOf[vdcgsModelVDC]{
						DataSource: &schemaD.SetNestedAttribute{
							MarkdownDescription: "A set of Virtual Data Center (VDC) IDs that are members of this VDC Group.",
							Computed:            true,
						},
						Attributes: superschema.Attributes{
							"id": superschema.SuperStringAttribute{
								DataSource: &schemaD.StringAttribute{
									MarkdownDescription: "Unique identifier for the Virtual Data Center (VDC).",
									Computed:            true,
								},
							},
							"name": superschema.SuperStringAttribute{
								DataSource: &schemaD.StringAttribute{
									MarkdownDescription: "Name assigned to the Virtual Data Center (VDC).",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}
