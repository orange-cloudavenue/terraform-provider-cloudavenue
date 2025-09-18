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

	"github.com/orange-cloudavenue/common-go/regex"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func storageProfileSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The storage profile data source displays the state of a storage profile contained within a VDC.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The unique identifier of the storage profile.",
					Computed:            true,
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("class")),
						stringvalidator.RegexMatches(regex.URNWithUUID4Regex(), "URN with UUID4 (urn:...-....-4...-...)"),
					},
				},
			},
			"class": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The storage class of the storage profile.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("class")),
					},
				},
			},
			"vdc_name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Computed:            true,
					Optional:            true,
					MarkdownDescription: "The name of the VDC containing the storage profile.",
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vdc_id"), path.MatchRoot("vdc_name")),
						stringvalidator.RegexMatches(regex.VDCNameRegex(), "VDC name (<alphanumeric> with hyphen and minus, with max length 27 and min length 2) - https://regex101.com/r/NgL6X0/1"),
					},
				},
			},
			"vdc_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The unique identifier of the VDC containing the storage profile.",
					Computed:            true,
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vdc_id"), path.MatchRoot("vdc_name")),
						stringvalidator.RegexMatches(regex.URNWithUUID4Regex(), "URN with UUID4 (urn:...-....-4...-...)"),
					},
				},
			},
			"limit": &superschema.SuperInt64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "The storage limit allocated to the storage profile in GiB.",
					Computed:            true,
				},
			},
			"used": &superschema.SuperInt64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "The amount of storage space currently used by the storage profile in GiB.",
					Computed:            true,
				},
			},
			"default": &superschema.SuperBoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Indicates whether this storage profile is set as the default for the VDC.",
					Computed:            true,
				},
			},
		},
	}
}
