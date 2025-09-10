/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func orgSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_org` resource allows you to manage the properties of an organization. You can update certain attributes of the organization, such as its `description`, `full_name`, `email`, and `internet_billing_mode`. However, please note that creating or deleting an organization is not supported through this resource.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_org` data source allows you to retrieve information about an organization.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the organization.",
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the organization.",
					Computed:            true,
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the organization.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
			},
			"full_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The full name of the organization, visible in the VCloud Director IHM.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
			},
			"email": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Your organization's contact email.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
			},
			"internet_billing_mode": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The organization's Internet bandwidth billing method.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "For more information, see the [documentation](https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/services/internet-access/).",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("PAYG", "TRAFFIC_VOLUME"),
					},
				},
			},
			"enabled": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Indicates whether the organization is enabled.",
					Computed:            true,
				},
			},
			"resources": superschema.SuperSingleNestedAttributeOf[OrgModelResources]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The resource usage of the organization.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"vdc": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The number of VDCs in the organization.",
							Computed:            true,
						},
					},
					"catalog": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The number of catalogs in the organization.",
							Computed:            true,
						},
					},
					"vapp": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The number of vApps in the organization.",
							Computed:            true,
						},
					},
					"vm_running": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The number of running VMs in the organization.",
							Computed:            true,
						},
					},
					"user": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The number of users in the organization.",
							Computed:            true,
						},
					},
					"disk": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The number of standalone disks in the organization.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
