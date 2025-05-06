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
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

/*
edgegwSchema

This function is used to create the schema for the edgegateway resource and datasource.
*/
func edgegwSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Edge Gateway ",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource allows you to create and delete Edge Gateways in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source allows you to show the details of an Edge Gateways in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"timeouts": &superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Delete: true,
					Update: true,
				},
			},
			"id": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"tier0_vrf_name": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Tier-0 VRF to which the Edge Gateway is attached.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional:            true,
					MarkdownDescription: "If not specified, the Edge Gateway will be created if only one Tier-0 VRF is available.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"owner_name": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway owner. It can be a VDC or a VDC Group name.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"description": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"bandwidth": &superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The bandwidth in `Mbps` of the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "If no value is specified, the bandwidth is automatically calculated based on the remaining bandwidth of the Tier-0 VRF. More information can be found [here](#bandwidth-attribute).\n\n!> **Warning** This attribute is not supported if your Tier-0 VRF has a class of service `DEDICATED`. This is due to a bug in the API (See [#1068](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1069))",
				},
			},
		},
	}
}
