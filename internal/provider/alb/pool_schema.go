/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package alb

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

/*
albPoolSchema

This function is used to create the schema for the ALB Pool resource and datasource.
*/
func albPoolSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a resource to manage Advanced Load Balancer Pools. Pools maintain the list of assigned servers and perform health monitoring, load balancing, and persistence. A pool may be used or referenced by only one virtual service at a time.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a data source to manage Advanced Load Balancer Pools. Pools maintain the list of assigned servers and perform health monitoring, load balancing, and persistence.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "ID of ALB Pool.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Name of ALB Pool.",
					Required:            true,
				},
			},
			"edge_gateway_id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Edge gateway ID in which ALB Pool",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_id"), path.MatchRoot("edge_gateway_name")),
					},
				},
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "should be created.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "was created.",
				},
			},
			"edge_gateway_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Edge gateway Name in which ALB Pool",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_id"), path.MatchRoot("edge_gateway_name")),
					},
				},
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "should be created.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "was created.",
				},
			},
			"enabled": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Define if ALB Pool is enabled or not.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Default:  booldefault.StaticBool(true),
				},
			},
			"description": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Description of ALB Pool.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"algorithm": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Algorithm for selecting members within a pool.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Default:  stringdefault.StaticString("LEAST_CONNECTIONS"),
					Validators: []validator.String{
						stringvalidator.OneOf("ROUND_ROBIN", "CONSISTENT_HASH", "LEAST_CONNECTIONS"),
					},
				},
			},
			"default_port": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Destination server port used for traffic sent to a member.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Default:  int64default.StaticInt64(80),
					Validators: []validator.Int64{
						int64validator.Between(1, 65535),
					},
				},
			},
			"graceful_timeout_period": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Maximum time in minutes allowed for gracefully disabling a pool member.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Default:  int64default.StaticInt64(1),
				},
			},
			"members": superschema.SetNestedAttribute{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "ALB Pool Member(s).",
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"enabled": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Indicates whether a pool member accepts traffic.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(true),
						},
					},
					"ip_address": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "IP address of a pool member.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								fstringvalidator.IsIP(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"port": superschema.Int64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Member port.",
						},
						Resource: &schemaR.Int64Attribute{
							Required: true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"ratio": superschema.Int64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Ratio of selecting eligible servers in the pool.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Default:  int64default.StaticInt64(1),
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
					},
				},
			},
			"health_monitors": superschema.SetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "List of health monitors type to activate.",
					ElementType:         types.StringType,
				},
				Resource: &schemaR.SetAttribute{
					Optional: true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(stringvalidator.OneOf("HTTP", "HTTPS", "TCP", "UDP", "PING")),
					},
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
			"persistence_profile": superschema.SingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Persistence profile ensures that a user remains connected to the same server for a specified duration. If the persistence profile is unmanaged by Cloud Avenue, updates with unchanged values will continue using the same unmanaged profile. However, any changes to the persistence profile will prompt Cloud Avenue to switch the pool to a profile it manages.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Type of persistence strategy.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("CLIENT_IP", "HTTP_COOKIE"),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"value": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Value of attribute based on persistence type.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("persistence_profile").AtName("type"), []attr.Value{types.StringValue("HTTP_COOKIE")}),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"passive_monitoring_enabled": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Monitors if the traffic is accepted by node.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Default:  booldefault.StaticBool(true),
				},
			},
		},
	}
}
