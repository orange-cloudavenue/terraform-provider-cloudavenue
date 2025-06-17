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
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
					MarkdownDescription: "The ID of the Firewall Edge Gateway.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
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
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"rules": superschema.SuperSetNestedAttributeOf[firewallModelRule]{
				// All default values are set is commented in the attributes below because
				// a known bug in the Terraform Framework does not allow to set default values
				// for set nested attributes. See https://github.com/hashicorp/terraform-plugin-framework/issues/783
				// All default are applied in the ModifyPlan method of the resource.
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "The collection of rules for configuring the firewall.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the rule.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
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
							MarkdownDescription: "Direction in firewall rules specifies whether the rule applies to incoming (inbound), outgoing (outbound), or both types of traffic. This attribute is crucial for defining how the firewall handles network traffic based on its direction.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								fstringvalidator.OneOfWithDescription(
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "IN",
										Description: "Inbound (ingress) traffic: Data packets coming into a network or device from external sources. An inbound rule controls which external connections are allowed to reach internal systems.",
									},
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "OUT",
										Description: "Outbound (egress) traffic: Data packets leaving a network or device to external destinations. An outbound rule controls which internal connections are allowed to reach external systems.",
									},
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "IN_OUT",
										Description: "Both inbound and outbound traffic.",
									},
								),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"ip_protocol": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The IP protocol of the rule. Default value is `IPV4`. ",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							// Default:  stringdefault.StaticString("IPV4"),
							Validators: []validator.String{
								stringvalidator.OneOf("IPV4", "IPV6", "IPV4_IPV6"),
							},
						},
					},
					"action": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Defines the behavior of the rule.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								fstringvalidator.OneOfWithDescription(
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "ALLOW",
										Description: "The firewall permits the matching traffic to pass through. For example, if an inbound rule with action ALLOW matches a packet, that packet is allowed into the network.",
									},
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "DROP",
										Description: "The firewall silently discards the matching traffic without notifying the sender. For example, if an outbound rule with action DROP matches a packet, that packet is dropped without any response sent back to the sender.",
									},
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "REJECT",
										Description: "The firewall blocks the matching traffic and sends a response to the sender, indicating that the connection was refused. This informs the sender that their traffic was intentionally blocked.",
									},
								),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"priority": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The priority of the rule. Lower values have higher priority. If lots of rules have the same priority, the alphabetical order of the rule name is used to determine the priority.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Validators: []validator.Int64{
								int64validator.Between(1, 1000),
							},
							// Default: int64default.StaticInt64(1),
						},
					},
					"enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Defines if the rule is enabled or not. Default value is `true`.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							// Default:  booldefault.StaticBool(true),
						},
					},
					"logging": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Defines if the rule should log matching traffic. Default value is `false`.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							// Default:  booldefault.StaticBool(false),
						},
					},
					"source_ip_addresses": superschema.SuperSetAttributeOf[string]{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "A set of source IP addresses, IP Range or CIDR. If `source_ids` attribute and this attribute are both empty, it means `Any` (all).",
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									fstringvalidator.IsNetwork([]fstringvalidator.NetworkValidatorType{
										fstringvalidator.IPV4,
										fstringvalidator.IPV4Range,
										fstringvalidator.IPV4WithCIDR,
									}, false),
								),
							},
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},

					"source_ids": superschema.SuperSetAttributeOf[string]{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "A set of Source Firewall Group IDs (`IP Sets` or `Security Groups`). If `source_ip_addresses` attribute and this attribute are both empty, it means `Any` (all).",
							ElementType:         supertypes.StringType{},
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
					"destination_ip_addresses": superschema.SuperSetAttributeOf[string]{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "A set of destination IP addresses, IP Range or CIDR. If `destination_ids` attribute and this attribute are both empty, it means `Any` (all).",
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									fstringvalidator.IsNetwork([]fstringvalidator.NetworkValidatorType{
										fstringvalidator.IPV4,
										fstringvalidator.IPV4Range,
										fstringvalidator.IPV4WithCIDR,
									}, false),
								),
							},
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
