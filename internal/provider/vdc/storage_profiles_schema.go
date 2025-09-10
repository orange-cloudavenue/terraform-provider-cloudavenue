/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package edgegw provides a Terraform resource to manage edge gateways.
package vdc

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func storageProfilesSuperSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The storage profiles data source show the list of storage prfile of a VDC.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Generated ID of the resource storage profile.",
					Computed:            true,
				},
			},
			"vdc_name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "The VDC name of storage profiles contains.",
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vdc_id"), path.MatchRoot("vdc_name")),
					},
				},
			},
			"vdc_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The VDC ID of storage profiles contains.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vdc_id"), path.MatchRoot("vdc_name")),
					},
				},
			},
			"storage_profiles": superschema.SuperListNestedAttributeOf[storageProfilesDataSourceModelStorageProfile]{
				DataSource: &schemaD.ListNestedAttribute{
					Computed:            true,
					MarkdownDescription: "A list of Storage Profiles.",
				},
				Attributes: superschema.Attributes{
					"id": &superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The ID of the Storage Profile.",
							Computed:            true,
						},
					},
					"class": &superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The class name of the Storage Profile.",
							Computed:            true,
						},
					},
					"limit": &superschema.SuperInt64Attribute{
						DataSource: &schemaD.Int64Attribute{
							MarkdownDescription: "The limit of the Storage Profile in GiB.",
							Computed:            true,
						},
					},
					"used": &superschema.SuperInt64Attribute{
						DataSource: &schemaD.Int64Attribute{
							MarkdownDescription: "The used space of the Storage Profile in GiB.",
							Computed:            true,
						},
					},
					"default": &superschema.SuperBoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "Whether the Storage Profile is the default one.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
