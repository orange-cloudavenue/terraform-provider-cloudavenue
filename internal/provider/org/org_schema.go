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
			MarkdownDescription: "The `cloudavenue_org` resource allows you to manage the properties of an organization.",
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
				Resource: &schemaR.StringAttribute{
					Optional: true,
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
			"email": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The email of the organization.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
			},
			"internet_billing_mode": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The internet billing mode of the organization.",
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
		},
	}
}
