/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package storage

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

func (d *profilesDataSource) superSchema(ctx context.Context) superschema.Schema {
	pDS := profileDataSource{}
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_storage_profile` data source can be used to access information about a storage profiles in a VDC.",
			Deprecated: superschema.DeprecatedResource{
				DeprecationMessage:                "This data source is deprecated and will be removed in a future release. Please use the `cloudavenue_vdc_storage_profiles` data source instead.",
				ComputeMarkdownDeprecationMessage: true,
				Renamed:                           true,
				TargetResourceName:                "cloudavenue_vdc_storage_profiles",
				TargetRelease:                     "v1.0.0",
				LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/28",
				LinkToIssue:                       "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1175",
				LinkToNewResourceDoc:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/data-sources/vdc_storage_profiles",
			},
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "ID of storage profile.",
					Computed:            true,
				},
			},
			"vdc": vdc.SuperSchema(),
			"storage_profiles": superschema.ListNestedAttribute{
				DataSource: &schemaD.ListNestedAttribute{
					MarkdownDescription: "List of storage profiles.",
					Computed:            true,
				},
				Attributes: pDS.superSchema(ctx).Attributes,
			},
		},
	}
}
