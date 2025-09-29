/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vrf

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func tier0VrfSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Tier-0",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source retrieve informations about a Tier-0. It can be used to retrieve the Tier-0 by its name, edge gateway ID or edge gateway name.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Tier-0.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Tier-0.",
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("edge_gateway_id"), path.MatchRoot("edge_gateway_name")),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway where the Tier-0 is located.",
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("edge_gateway_id"), path.MatchRoot("edge_gateway_name")),
					},
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway where the Tier-0 is located.",
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("edge_gateway_id"), path.MatchRoot("edge_gateway_name")),
					},
				},
			},
			"class_service": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The class service associated with the Tier-0.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"bandwidth": superschema.SuperSingleNestedAttributeOf[tier0VrfDataSourceModelBandwidth]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The bandwidth information for the Tier-0.",
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"capacity": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The total bandwidth capacity for the Tier-0 in Mbps.",
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"provisioned": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The total bandwidth provisioned for the Tier-0 across all edge gateways in Mbps.",
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"remaining": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The remaining bandwidth that can be allocated to the new edge gateway in Mbps.",
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"allowed_bandwidth_values": superschema.SuperListAttributeOf[int64]{
						Common: &schemaR.ListAttribute{
							MarkdownDescription: "The allowed bandwidth values for the new edge gateway in Mbps.",
						},
						DataSource: &schemaD.ListAttribute{
							Computed: true,
						},
					},
					"allow_unlimited": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Indicates if unlimited bandwidth is allowed for the Tier-0.",
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
			"edgegateways": superschema.SuperListNestedAttributeOf[tier0VrfDataSourceModelEdgeGateway]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "List of Edge Gateways where the Tier-0 is located.",
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the Edge Gateway.",
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the Edge Gateway.",
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"bandwidth": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The bandwidth allocated to the Edge Gateway in Mbps.",
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"allowed_bandwidth_values": superschema.SuperListAttributeOf[int64]{
						Common: &schemaR.ListAttribute{
							MarkdownDescription: "The allowed bandwidth values for the Edge Gateway in Mbps.",
						},
						DataSource: &schemaD.ListAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
