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

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func tier0VrfSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Tier-0 VRF",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source retrieve informations about a Tier-0 VRF.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Tier-0 VRF.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Tier-0 VRF.",
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"tier0_provider": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Tier-0 provider info.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"class_service": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "List of Tags for the Tier-0 VRF.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"services": superschema.SuperListNestedAttributeOf[segmentModel]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "Services list of the Tier-0 VRF.",
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"service": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Service of the segment.",
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"vlan_id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "VLAN ID of the segment.",
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
