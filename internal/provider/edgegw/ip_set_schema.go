/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func ipSetSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_edgegateway_ip_set` resource allows you to manage an IP Set rule on an Edge Gateway. IP Sets are groups of objects to which the firewall rules apply. Combining multiple objects into IP Sets helps reduce the total number of firewall rules to be created.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_edgegateway_ip_set` data source allows you to retrieve information about an IP Set rule on an Edge Gateway.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the IP Set.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			name: superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the IP Set.",
					Required:            true,
				},
			},
			edgeGatewayName: superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: edgeGatewayNameDescription,
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot(edgeGatewayName), path.MatchRoot(edgeGatewayID)),
					},
				},
			},
			edgeGatewayID: superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: edgeGatewayIDDescription,
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot(edgeGatewayName), path.MatchRoot(edgeGatewayID)),
					},
				},
			},
			description: superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the IP Set.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"ip_addresses": superschema.SuperSetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "A set of IP address, CIDR or IP range.",
					ElementType:         supertypes.StringType{},
				},
				Resource: &schemaR.SetAttribute{
					Optional: true,
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
		},
	}
}
