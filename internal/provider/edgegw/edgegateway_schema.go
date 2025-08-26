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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
			MarkdownDescription: "resource allows you to create and delete Edge Gateway in Cloud Avenue. EdgeGateway is a virtualized network appliance designed to provide secure connectivity, routing, and network services at the edge of a virtualized environment. It acts as a critical component for managing network traffic between internal virtual networks and external networks, such as the internet or remote sites.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source allows you to show the details of an Edge Gateway in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
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
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					},
				},
			},
			"name": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway.",
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
						stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					},
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
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Tier-0 VRF to which the Edge Gateway is attached.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional:            true,
					MarkdownDescription: "If not specified, the Edge Gateway will be created if only one Tier-0 VRF is available.",
				},
			},
			"t0_name": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the T0 Name to which the Edge Gateway is attached.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional:            true,
					MarkdownDescription: "If not specified, the Edge Gateway will be created on the T0 available in your organization. Works if only one T0 Name is available.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"t0_id": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the T0 to which the Edge Gateway is attached.",
					Computed:            true,
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
			"bandwidth": &superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The bandwidth in `Mbps` of the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "The bandwidth limit in Mbps for the edge gateway. If t0 is `SHARED`, it must be one of the available values for the T0 router, if no value is specified, the bandwidth is automatically calculated based on the remaining bandwidth of the T0. If t0 is `DEDICATED`, unlimited bandwidth is allowed (0 = unlimited). More information can be found [here](#bandwidth-attribute).",
					PlanModifiers: []planmodifier.Int64{
						// This is used because the create edge gateway operation may not have all the information available
						int64planmodifier.UseStateForUnknown(),
					},
				},
			},
		},
	}
}
