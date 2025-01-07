/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func firewallSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The firewall resource allows you to manage rules on an Firewall.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The firewall data source allows you to retrieve information about an Firewall.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the Firewall Edge Gateway Service.",
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"rules": superschema.SuperListNestedAttributeOf[firewallModelRule]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "The list of rules to apply to the firewall.",
				},
				Resource: &schemaR.ListNestedAttribute{
					Required: true,
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the rule.",
						},
					},
					"name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the rule.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"direction": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The direction of the rule.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("IN", "OUT", "IN_OUT"),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"ip_protocol": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The IP protocol of the rule.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Default:  stringdefault.StaticString("IPV4"),
							Validators: []validator.String{
								stringvalidator.OneOf("IPV4", "IPV6", "IPV4_IPV6"),
							},
						},
					},
					"action": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Defines if the rule should `ALLOW` or `DROP` matching traffic.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("ALLOW", "DROP"),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Defines if the rule is enabled or not.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(true),
						},
					},
					"logging": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Defines if the rule should log matching traffic.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(false),
						},
					},
					"source_ids": superschema.SuperSetAttributeOf[string]{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "A set of Source Firewall Group IDs (`IP Sets` or `Security Groups`). Leaving it empty means `Any` (all).",
							ElementType:         supertypes.StringType{},
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
					"destination_ids": superschema.SuperSetAttributeOf[string]{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "A set of Destination Firewall Group IDs (`IP Sets` or `Security Groups`). Leaving it empty means `Any` (all).",
							ElementType:         supertypes.StringType{},
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
					"app_port_profile_ids": superschema.SuperSetAttributeOf[string]{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "A set of Application Port Profile IDs. Leaving it empty means `Any` (all).",
							ElementType:         types.StringType,
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
