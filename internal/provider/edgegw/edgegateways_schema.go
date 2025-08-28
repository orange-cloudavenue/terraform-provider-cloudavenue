/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func edgeGatewaysSuperSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The edge gateways data source show the list of edge gateways of an organization.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Generated ID of the resource.",
				},
			},
			"edge_gateways": superschema.SuperListNestedAttributeOf[edgeGatewaysDataSourceModelEdgeGateway]{
				DataSource: &schemaD.ListNestedAttribute{
					Computed:            true,
					MarkdownDescription: "A list of Edge Gateways.",
				},
				Attributes: superschema.Attributes{
					"id": &superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The ID of the Edge Gateway.",
							Computed:            true,
						},
					},
					"name": &superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the Edge Gateway.",
							Computed:            true,
						},
					},
					"tier0_vrf_name": &superschema.SuperStringAttribute{
						Deprecated: &superschema.Deprecated{
							DeprecationMessage:                "This field is deprecated and will be removed in future versions. Please use 't0_name' instead.",
							ComputeMarkdownDeprecationMessage: true,
							Renamed:                           true,
							FromAttributeName:                 "tier0_vrf_name",
							TargetAttributeName:               "t0_name",
							TargetRelease:                     "1.0.0",
							LinkToIssue:                       "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1165",
							LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/28",
							LinkToResourceDoc:                 "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/edgegateway",
						},
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the Tier-0 VRF to which the Edge Gateway is attached.",
							Computed:            true,
						},
					},
					"t0_name": &superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the T0 Name to which the Edge Gateway is attached.",
							Computed:            true,
						},
					},
					"t0_id": &superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the T0 to which the Edge Gateway is attached.",
							Computed:            true,
						},
					},
					"owner_name": &superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the Edge Gateway owner. It can be a VDC or a VDC Group name.",
							Computed:            true,
						},
					},
					"owner_id": &superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the Edge Gateway owner. It can be a VDC or a VDC Group ID.",
							Computed:            true,
						},
					},
					"description": &superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The description of the Edge Gateway.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
