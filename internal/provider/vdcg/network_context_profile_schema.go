/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
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
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	sdkv1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

func networkContextProfileSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdcg_network_context_profile` resource allows you to manage a custom (TENANT-scoped) Network Context Profile on a VDC Group. Context profiles define Layer 7 traffic criteria (application identifiers and/or domain names) that can be referenced in firewall rules via `network_context_profile_ids`.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdcg_network_context_profile` data source allows you to retrieve information about a Network Context Profile (Layer 7) available on a VDC Group. Use this to reference SYSTEM or PROVIDER profiles by name in firewall rules.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Network Context Profile.",
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
						stringvalidator.ExactlyOneOf(path.MatchRoot(name), path.MatchRoot("id")),
					},
				},
			},
			name: superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Network Context Profile.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot(name), path.MatchRoot("id")),
					},
				},
			},
			description: superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A human-readable description of the Network Context Profile.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			vdcGroupID: superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the VDC Group.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot(vdcGroupID), path.MatchRoot(vdcGroupName)),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
				},
			},
			vdcGroupName: superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the VDC Group.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.AtLeastOneOf(path.MatchRoot(vdcGroupID), path.MatchRoot(vdcGroupName)),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
				},
			},
			"scope": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The scope of the Network Context Profile (`SYSTEM`, `PROVIDER` or `TENANT`). Resources are always created as `TENANT`.",
					Computed:            true,
				},
			},
			"app_id": superschema.SuperSingleNestedAttributeOf[networkContextProfileModelAppID]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Layer 7 App ID attribute. Defines a set of application identifiers to match (e.g. `SSL`, `CIFS`, `HTTP`).\n\n" +
						"  ~> **Note:** Sub-attributes (`sub_attribute`) are only supported when `app_id.values` contains exactly **one** entry.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					attrValues: superschema.SuperSetAttributeOf[string]{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "The set of App ID values to match (e.g. `[\"SSL\", \"CIFS\"]`).",
							ElementType:         supertypes.StringType{},
						},
						Resource: &schemaR.SetAttribute{
							Required: true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf(sdkv1.NetworkContextProfileKnownAppIDs...),
								),
							},
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
					"sub_attribute": superschema.SuperListNestedAttributeOf[networkContextProfileModelSubAttribute]{
						Common: &schemaR.ListNestedAttribute{
							MarkdownDescription: "Optional sub-attributes to refine the App ID match (e.g. TLS version, cipher suites, SMB version).\n\n" +
								"  ~> **Note:** Only supported when `app_id.values` contains exactly one entry.",
						},
						Resource: &schemaR.ListNestedAttribute{
							Optional: true,
						},
						DataSource: &schemaD.ListNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							attrType: superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The sub-attribute type. Allowed values: `TLS_VERSION`, `TLS_CIPHER_SUITE`, `CIFS_SMB_VERSION`.",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf(sdkv1.NetworkContextProfileKnownSubAttributeTypes...),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							attrValues: superschema.SuperSetAttributeOf[string]{
								Common: &schemaR.SetAttribute{
									MarkdownDescription: "The set of allowed values for this sub-attribute type.",
									ElementType:         supertypes.StringType{},
								},
								Resource: &schemaR.SetAttribute{
									Required: true,
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.OneOf(sdkv1.NetworkContextProfileKnownSubAttributeValues...),
										),
									},
								},
								DataSource: &schemaD.SetAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
			"domain_name": superschema.SuperSingleNestedAttributeOf[networkContextProfileModelDomainName]{
				DataSource: &schemaD.SingleNestedAttribute{
					MarkdownDescription: "Domain Name (FQDN) attribute. Present on SYSTEM profiles that match traffic by fully-qualified domain name.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					attrValues: superschema.SuperSetAttributeOf[string]{
						DataSource: &schemaD.SetAttribute{
							MarkdownDescription: "The set of domain name values for this profile.",
							Computed:            true,
							ElementType:         supertypes.StringType{},
						},
					},
				},
			},
		},
	}
}
