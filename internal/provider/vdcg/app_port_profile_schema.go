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
	"github.com/orange-cloudavenue/terraform-plugin-framework-validators/setvalidator"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

func appPortProfileSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdcg_app_port_profile` resource allows you to manage an application port profile.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdcg_app_port_profile` data source allows you to retrieve information about an application port profile.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the App Port profile.",
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
					MarkdownDescription: "Application Port Profile name.",
				},
				Resource: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "The name is unique within your organization.",
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Application Port Profile description.",
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
					MarkdownDescription: "The name of vDC group that owns the application port profile.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot("vdc_group_name"), path.MatchRoot("vdc_group_id")),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"vdc_group_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of vDC group that owns the application port profile.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot("vdc_group_name"), path.MatchRoot("vdc_group_id")),
						fstringvalidator.PrefixContains(urn.VDCGroup.String()),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"scope": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The scope of the application port profile. This attribute is required only if the terraform apply return an error with the message `Multiple App Port Profiles found with the same name`. In this case, you must specify the scope of the application port profile.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(func() (resp []string) {
							for _, scope := range v1.FirewallGroupAppPortProfileModelScopes {
								resp = append(resp, string(scope))
							}
							return
						}()...),
					},
				},
			},
			"app_ports": superschema.SuperListNestedAttributeOf[AppPortProfileModelAppPort]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "List of application ports.",
				},
				Resource: &schemaR.ListNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"ports": superschema.SuperSetAttributeOf[string]{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "Set of destination ports or destination ports ranges.",
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
							Validators: []validator.Set{
								setvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("protocol"), []attr.Value{types.StringValue("ICMPv4"), types.StringValue("ICMPv6")}),
								setvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("protocol"), []attr.Value{types.StringValue("TCP"), types.StringValue("UDP")}),
							},
						},
					},
					"protocol": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Protocol of the application port.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(func() (resp []string) {
									for _, protocol := range v1.FirewallGroupAppPortProfileModelPortProtocols {
										resp = append(resp, string(protocol))
									}
									return
								}()...),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
