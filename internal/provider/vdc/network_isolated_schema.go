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

package vdc

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
)

func networkIsolatedSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdc_network_isolated` resource allows you to manage an isolated network in a `VDC`.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdc_network_isolated` data source allows you to retrieve information about an isolated network in a `VDC`.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the isolated network.",
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the network. This value must be unique within the `VDC` that owns the network.",
					Required:            true,
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A description of the network.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"vdc": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of vDC to use",
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"gateway": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The gateway IP address for the network. This value define also the network IP range with the prefix length.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						fstringvalidator.IsIP(),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"guest_vlan_allowed": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Indicates if the network allows guest VLANs.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Default:  booldefault.StaticBool(false),
				},
			},
			"prefix_length": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The prefix length for the network. This value must be a valid prefix length for the network IP range. (e.g. /24 for netmask 255.255.255.0). For more information, see [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing).",
				},
				Resource: &schemaR.Int64Attribute{
					Required: true,
					Validators: []validator.Int64{
						int64validator.Between(1, 32),
					},
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"dns1": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The primary DNS server IP address for the network.",
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
			"dns2": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The secondary DNS server IP address for the network.",
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
			"dns_suffix": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The DNS suffix for the network.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"static_ip_pool": superschema.SuperSetNestedAttributeOf[networkIsolatedModelStaticIPPool]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "A set of static IP pools to be used for this network.",
					Computed:            true,
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"start_address": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The start address of the IP pool. This value must be a valid IP address in the network IP range.",
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
					"end_address": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The end address of the IP pool. This value must be a valid IP address in the network IP range.",
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
				},
			},
		},
	}
}
