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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func vdcgSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Manages a Virtual Data Center Group (VDC Group) in Cloud Avenue. A VDC Group allows you to aggregate multiple Virtual Data Centers for unified management and networking.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Retrieves information about an existing Virtual Data Center Group (VDC Group) in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Unique identifier of the VDC Group.",
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Name assigned to the VDC Group.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Detailed description of the VDC Group and its purpose.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"vdc_ids": superschema.SuperSetAttributeOf[string]{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "A set of Virtual Data Center (VDC) IDs that are members of this VDC Group.",
				},
				Resource: &schemaR.SetAttribute{
					Required: true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
		},
	}
}
