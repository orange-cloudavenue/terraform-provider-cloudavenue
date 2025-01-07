/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

/*
vappSchema

This function is used to create the superschema for the vapp resource and datasource.
*/
func vappSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "Provides a Cloud Avenue vApp",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource. This can be used to create, modify, and delete vApps.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source. This can be used to reference vApps.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "ID of the vApp.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
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
					MarkdownDescription: "Name of the vApp.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"vdc": vdc.SuperSchemaSuperType(),
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Description of the vApp.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"guest_properties": superschema.SuperMapAttribute{
				Common: &schemaR.MapAttribute{
					MarkdownDescription: "Key/value settings for guest properties.",
					ElementType:         types.StringType,
				},
				Resource: &schemaR.MapAttribute{
					Optional: true,
				},
				DataSource: &schemaD.MapAttribute{
					Computed: true,
				},
			},
			"lease": superschema.SuperSingleNestedAttributeOf[vappResourceModelLease]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Informations about vApp lease",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Computed: true,
					Default: objectdefault.StaticValue(types.ObjectValueMust(
						map[string]attr.Type{
							"runtime_lease_in_sec": types.Int64Type,
							"storage_lease_in_sec": types.Int64Type,
						},
						map[string]attr.Value{
							"runtime_lease_in_sec": types.Int64Value(0),
							"storage_lease_in_sec": types.Int64Value(0),
						})),
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"runtime_lease_in_sec": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "How long any of the VMs in the vApp can run before the vApp is automatically powered off or suspended. Allowed values are 3600 to 31536000 seconds (1 hour to 365 days) or 0 means never expires.",
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Computed: true,
							Validators: []validator.Int64{
								int64validator.Any(
									int64validator.OneOf(0),                // 0 is never expires
									int64validator.Between(3600, 31536000), // 1 hour to 365 days
								),
							},
							Default: int64default.StaticInt64(0),
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"storage_lease_in_sec": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "How long the vApp is available before being automatically deleted or marked as expired. Allowed values are 3600 to 31536000 seconds (1 hour to 365 days) or 0 means never expires.",
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Computed: true,
							Validators: []validator.Int64{
								int64validator.Any(
									int64validator.OneOf(0),                // 0 is never expires
									int64validator.Between(3600, 31536000), // 1 hour to 365 days
								),
							},
							Default: int64default.StaticInt64(0),
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
