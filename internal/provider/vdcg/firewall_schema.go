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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

func firewallSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdcg_firewall` resource allows you to manage a firewall in a VDC Group. The firewall is a set of rules that are applied to the VDC Group networks to control the traffic.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdcg_firewall` data source allows you to retrieve information about an existing firewall in a VDC Group.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the firewall. The firewall does not have an ID. It is identified by the VDC Group ID.",
				},
			},
			"vdc_group_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the VDC Group to which the firewall belongs.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot("vdc_group_name"), path.MatchRoot("vdc_group_id")),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
				},
			},
			"vdc_group_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the VDC Group to which the firewall belongs.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot("vdc_group_name"), path.MatchRoot("vdc_group_id")),
						fstringvalidator.PrefixContains(urn.VDCGroup.String()),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
				},
			},
			"enabled": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					Computed:            true,
					MarkdownDescription: "Defines if the firewall is enabled or not.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Default:  booldefault.StaticBool(true),
				},
			},
			"rules": superschema.SuperListNestedAttributeOf[FirewallModelRule]{
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
								stringvalidator.OneOf(
									func() (values []string) {
										for _, value := range v1.VDCGroupFirewallTypeRuleDirections {
											values = append(values, string(value))
										}
										return
									}()...,
								),
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
							Default:  stringdefault.StaticString(string(v1.VDCGroupFirewallTypeRuleIPProtocolIPv4)),
							Validators: []validator.String{
								stringvalidator.OneOf(
									func() (values []string) {
										for _, value := range v1.VDCGroupFirewallTypeRuleIPProtocols {
											values = append(values, string(value))
										}
										return
									}()...,
								),
							},
						},
					},
					"action": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Defines if the rule should matching traffic.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									func() (values []string) {
										for _, value := range v1.VDCGroupFirewallTypeRuleActions {
											values = append(values, string(value))
										}
										return
									}()...,
								),
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
							MarkdownDescription: "A set of Source Firewall Group IDs ([`IP Sets`](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg_ip_set), [`Security Groups`](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg_security_group) or [`Dynamic Security Group`](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg_dynamic_security_group)). Leaving it empty means `Any` (all).",
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
							MarkdownDescription: "A set of Destination Firewall Group IDs ([`IP Sets`](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg_ip_set), [`Security Groups`](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg_security_group) or [`Dynamic Security Group`](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg_dynamic_security_group)). Leaving it empty means `Any` (all).",
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
							MarkdownDescription: "A set of [Application Port Profile](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg_app_port_profile) IDs. Leaving it empty means `Any` (all).",
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
					"source_groups_excluded": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Reverses value of `source_ids` for the rule to match everything except specified IDs.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(false),
						},
					},
					"destination_groups_excluded": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Reverses value of `destination_ids` for the rule to match everything except specified IDs.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}
