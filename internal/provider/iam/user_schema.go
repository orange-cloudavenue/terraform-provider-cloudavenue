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
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

/*
userSchema

This function is used to create the schema for the user resource and datasource.
*/
func userSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The user",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource allows you to manage local users in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source allows you to read users in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the user.",
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
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the user.",
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
			"role_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The role assigned to the user.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"full_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The user's full name.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"email": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The user's email address.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"telephone": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The user's telephone number.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"enabled": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "`true` if the user is enabled and can log in.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(true),
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"deployed_vm_quota": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Quota of vApps that this user can deploy. A value of `0` specifies an unlimited quota.",
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Computed: true,
					Default:  int64default.StaticInt64(0),
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"stored_vm_quota": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Quota of vApps that this user can store. A value of `0` specifies an unlimited quota.",
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Computed: true,
					Default:  int64default.StaticInt64(0),
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"password": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The user's password. This value is never returned on read.",
					Required:            true,
					Sensitive:           true,
				},
			},
			"take_ownership": superschema.SuperBoolAttribute{
				Resource: &schemaR.BoolAttribute{
					MarkdownDescription: "`true` if the user should take ownership of all vApps and media that are currently owned by the user that is being deleted.",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(true),
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"provider_type": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Identity provider type for this this user.",
					Computed:            true,
				},
			},
		},
	}
}
