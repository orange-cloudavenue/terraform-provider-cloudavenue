/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vcda

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func VcdaIPSchema() superschema.Schema { //nolint:revive
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The VCDa",
			Deprecated: superschema.DeprecatedResource{
				DeprecationMessage:                "The resource 'cloudavenue_vcda_ip' has been deprecated in favor of 'cloudavenue_draas_onpremise'.\n\n" + "Please update your configuration accordingly.",
				ComputeMarkdownDeprecationMessage: true,
				Renamed:                           true,
				TargetResourceName:                "cloudavenue_draas_onpremise",
				TargetRelease:                     "v1.0.0",
				LinkToNewResourceDoc:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/draas_onpremise",
				LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/28",
				LinkToMigrationGuide:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vcda_ip",
			},
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource allows you to declare or remove your on-premises IP address for the DRaaS service.\n\n" +
				" -> Note: For more information, please refer to the [Cloud Avenue DRaaS documentation](https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/services/draas-with-vcda-2/).",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the VCDa resource.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"ip_address": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The on-premises IP address refers to the IP address of your local infrastructure running vCloud Extender.",
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
			},
		},
	}
}
