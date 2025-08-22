/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package publicip

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	"golang.org/x/net/context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func publicIPsSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The public IP data source displays the list of public IP addresses.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The ID of the Public IP.",
					Computed:            true,
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway.",
					Computed:            true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway.",
					Computed:            true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"public_ips": superschema.SuperListNestedAttributeOf[publicIPNetworkConfigModel]{
				DataSource: &schemaD.ListNestedAttribute{
					MarkdownDescription: "A list of public IPs.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The ID of the Public IP.",
							Computed:            true,
						},
					},
					"public_ip": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The Public IP Address.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
