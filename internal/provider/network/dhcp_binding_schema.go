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

package network

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
)

func dhcpBindingSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `network_dhcp_binding` resource allows you to manage DHCP bindings.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `network_dhcp_binding` data source allows you to retrieve information about an existing DHCP binding.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the DHCP Binding.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"org_network_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "The ID of the Org Network.<br/>**Note** (`.id` field) of `cloudavenue_vdc_network_isolated`, `cloudavenue_edgegateway_network_routed` or `cloudavenue_network_dhcp` can be referenced here. It is more convenient to use reference to `cloudavenue_network_dhcp` ID because it makes sure that DHCP is enabled before configuring pools.",
				},
				Resource: &schemaR.StringAttribute{
					Validators: []validator.String{
						fstringvalidator.IsURN(),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "The name of the DHCP Binding.",
				},
			},
			"ip_address": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The IP address of the DHCP Binding.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						fstringvalidator.IsIP(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"mac_address": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The MAC address of the DHCP Binding.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						fstringvalidator.IsMacAddress(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the DHCP Binding.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"dns_servers": superschema.SuperListAttribute{ // TODO use dns_servers attribute from network_dhcp resource
				Common: &schemaR.ListAttribute{
					MarkdownDescription: "The DNS server IPs to be assigned by this DHCP service.",
					ElementType:         supertypes.StringType{},
				},
				Resource: &schemaR.ListAttribute{
					Optional: true,
					Validators: []validator.List{
						listvalidator.SizeAtMost(2),
					},
				},
				DataSource: &schemaD.ListAttribute{
					Computed: true,
				},
			},
			"lease_time": superschema.SuperInt64Attribute{ // TODO use lease_time attribute from network_dhcp resource
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The lease time in seconds for the DHCP service.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Default:  int64default.StaticInt64(86400),
					Validators: []validator.Int64{
						int64validator.AtLeast(60),
					},
				},
			},
			"dhcp_v4_config": superschema.SuperSingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The DHCPv4 configuration for the DHCP Binding.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"gateway_address": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The gateway address to be assigned by this DHCP service.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								fstringvalidator.IsIP(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"hostname": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The hostname to be assigned by this DHCP service.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
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
