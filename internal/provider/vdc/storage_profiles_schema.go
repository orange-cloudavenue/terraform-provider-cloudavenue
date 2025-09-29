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
	"context"

	"github.com/orange-cloudavenue/common-go/urn"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func storageProfilesSuperSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The storage profiles data source retrieves a list of all storage profiles available within a VDC.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The auto-generated identifier for this data source instance.",
					Computed:            true,
				},
			},
			"vdc_name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "The name of the VDC containing the storage profiles.",
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vdc_id"), path.MatchRoot("vdc_name")),
					},
				},
			},
			"vdc_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The unique identifier of the VDC containing the storage profiles.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vdc_id"), path.MatchRoot("vdc_name")),
						fstringvalidator.IsURN(),
						fstringvalidator.PrefixContains(urn.VDC.String()),
					},
				},
			},
			"storage_profiles": superschema.SuperListNestedAttributeOf[storageProfileDataSourceModelStorageProfile]{
				DataSource: &schemaD.ListNestedAttribute{
					Computed:            true,
					MarkdownDescription: "A collection of storage profiles available within the VDC.",
				},
				Attributes: superschema.Attributes{
					"id": &superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The unique identifier of the storage profile.",
							Computed:            true,
						},
					},
					"class": &superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The storage class type of the storage profile.",
							Computed:            true,
						},
					},
					"limit": &superschema.SuperInt64Attribute{
						DataSource: &schemaD.Int64Attribute{
							MarkdownDescription: "The maximum storage limit (in GB) allocated to this storage profile.",
							Computed:            true,
						},
					},
					"used": &superschema.SuperInt64Attribute{
						DataSource: &schemaD.Int64Attribute{
							MarkdownDescription: "The amount of storage (in GB) currently used within this storage profile.",
							Computed:            true,
						},
					},
					"default": &superschema.SuperBoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "Indicates whether this storage profile is set as the default profile for the VDC.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
