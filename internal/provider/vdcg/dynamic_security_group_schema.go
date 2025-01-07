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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

func dynamicSecurityGroupSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdcg_dynamic_security_group` resource allows you to manage dynamic security groups in the VDC Group. A dynamic security group is a group of VMs that share the same security rules. The VMs are dynamically added or removed from the group based on the criteria defined in the security group. The dynamic security group will be attached to the VDC Group firewall",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdcg_dynamic_security_group` data source allows you to retrieve information about an existing dynamic security group in the VDC Group.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the dynamic security group.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the dynamic security group.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the dynamic security group.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"vdc_group_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the VDC Group to which the dynamic security group belongs.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot("vdc_group_name"), path.MatchRoot("vdc_group_id")),
					},
				},
			},
			"vdc_group_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the VDC Group to which the dynamic security group belongs.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot("vdc_group_name"), path.MatchRoot("vdc_group_id")),
						fstringvalidator.PrefixContains(urn.VDCGroup.String()),
					},
				},
			},
			"criteria": superschema.SuperListNestedAttributeOf[DynamicSecurityGroupModelCriteria]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "The list of dynamic criteria that determines whether a VM belongs to a dynamic firewall group. A VM needs to meet at least one criteria to belong to the firewall group. In other words, the logical AND is used for rules within a single criteria and the logical OR is used in between each criteria.",
				},
				Resource: &schemaR.ListNestedAttribute{
					Optional: true,
					Validators: []validator.List{
						listvalidator.SizeAtMost(3),
					},
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"rules": superschema.SuperListNestedAttributeOf[DynamicSecurityGroupModelRule]{
						Common: &schemaR.ListNestedAttribute{
							MarkdownDescription: "The list of rules that determine whether a VM belongs to a dynamic firewall group. A VM needs to meet all rules within a single criteria to belong to the firewall group. In other words, the logical AND is used for rules within a single criteria and the logical OR is used in between each criteria.",
						},
						Resource: &schemaR.ListNestedAttribute{
							Optional: true,
							Validators: []validator.List{
								listvalidator.SizeAtMost(4),
							},
						},
						DataSource: &schemaD.ListNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"type": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The type of the rule.",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										fstringvalidator.OneOfWithDescription(func() (values []fstringvalidator.OneOfWithDescriptionValues) {
											values = make([]fstringvalidator.OneOfWithDescriptionValues, 0, len(v1.FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypes))

											for _, v := range v1.FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypes {
												values = append(values, fstringvalidator.OneOfWithDescriptionValues{
													Value:       string(v.Type),
													Description: v.Description,
												},
												)
											}

											return values
										}()...),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"value": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "String to evaluate by given `type` and `operator`.",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"operator": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The operator to use to evaluate the `value`.",
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										fstringvalidator.OneOfWithDescriptionIfAttributeIsOneOf(
											path.MatchRelative().AtParent().AtName("type"),
											[]attr.Value{types.StringValue("VM_NAME")},
											func() (values []fstringvalidator.OneOfWithDescriptionIfAttributeIsOneOfValues) {
												values = make([]fstringvalidator.OneOfWithDescriptionIfAttributeIsOneOfValues, 0, len(v1.FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMNameOperator))

												for _, v := range v1.FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMNameOperator {
													values = append(values, fstringvalidator.OneOfWithDescriptionIfAttributeIsOneOfValues{
														Value:       string(v.Operator),
														Description: v.Description,
													},
													)
												}

												return values
											}()...),
										fstringvalidator.OneOfWithDescriptionIfAttributeIsOneOf(
											path.MatchRelative().AtParent().AtName("type"),
											[]attr.Value{types.StringValue("VM_TAG")},
											func() (values []fstringvalidator.OneOfWithDescriptionIfAttributeIsOneOfValues) {
												values = make([]fstringvalidator.OneOfWithDescriptionIfAttributeIsOneOfValues, 0, len(v1.FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTagOperator))

												for _, v := range v1.FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTagOperator {
													values = append(values, fstringvalidator.OneOfWithDescriptionIfAttributeIsOneOfValues{
														Value:       string(v.Operator),
														Description: v.Description,
													},
													)
												}

												return values
											}()...),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
