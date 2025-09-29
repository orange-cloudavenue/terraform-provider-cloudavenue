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

func (d *profileDataSource) superSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_storage_profile` data source can be used to access information about a storage profile in a VDC.",
			Deprecated: superschema.DeprecatedResource{
				DeprecationMessage:                "This data source is deprecated and will be removed in a future release. Please use the `cloudavenue_vdc_storage_profile` data source instead.",
				ComputeMarkdownDeprecationMessage: true,
				Renamed:                           true,
				TargetResourceName:                "cloudavenue_vdc_storage_profile",
				TargetRelease:                     "v1.0.0",
				LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/28",
				LinkToIssue:                       "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1175",
				LinkToNewResourceDoc:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/data-sources/vdc_storage_profile",
			},
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "ID of storage profile.",
					Computed:            true,
				},
			},
			"name": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Name of storage profile.",
					Required:            true,
				},
			},
			"vdc": vdc.SuperSchema(),
			"limit": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Maximum number of storage bytes (scaled by 'units' field) allocated for this profile. `0` means `maximum possible`",
					Computed:            true,
				},
			},
			// used_storage
			"used_storage": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Storage used, in Megabytes, by the storage profile.",
					Computed:            true,
				},
			},
			// default
			"default": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Indicates whether this is the default storage profile for the VDC.",
					Computed:            true,
				},
			},
			// enabled
			"enabled": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Indicates whether this storage profile is enabled for the VDC.",
					Computed:            true,
				},
			},
			// iops_allocated
			"iops_allocated": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Total IOPS currently allocated to this storage profile.",
					Computed:            true,
				},
			},
			// units
			"units": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Scale used to define Limit.",
					Computed:            true,
				},
			},
			// iops_limiting_enabled
			"iops_limiting_enabled": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "True if this storage profile is IOPS-based placement enabled.",
					Computed:            true,
				},
			},
			// maximum_disk_iops
			"maximum_disk_iops": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "The maximum IOPS value that this storage profile is permitted to deliver. Value of 0 means this max setting is disabled and there is no max disk IOPS restriction.",
					Computed:            true,
				},
			},
			// default_disk_iops
			"default_disk_iops": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Value of 0 for disk IOPS means that no IOPS would be reserved or provisioned for that virtual disk.",
					Computed:            true,
				},
			},
			// disk_iops_per_gb_max
			"disk_iops_per_gb_max": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "The maximum IOPS per GB value that this storage profile is permitted to deliver. Value of 0 means this max setting is disabled and there is no max disk IOPS per GB restriction.",
					Computed:            true,
				},
			},
			// iops_limit
			"iops_limit": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Maximum number of IOPs that can be allocated for this profile. `0` means `maximum possible`.",
					Computed:            true,
				},
			},
		},
	}
}
