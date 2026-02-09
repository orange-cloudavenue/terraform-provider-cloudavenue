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

package vm

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

func vmInsertedMediaSuperSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The inserted_media resource resource for inserting or ejecting media (ISO) file for the VM. Create this resource for inserting the media, and destroy it for ejecting.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the inserted media. This is the vm Id where the media is inserted.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"vdc":       vdc.SuperSchema(),
			"vapp_id":   vapp.SuperSchema()["vapp_id"],
			"vapp_name": vapp.SuperSchema()["vapp_name"],
			"catalog": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "The name of the catalog where to find media file",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"name": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "Media file name in catalog which will be inserted to VM",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"vm_name": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "VM name where media will be inserted or ejected",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			// "eject_force": schema.BoolAttribute{ - Disable attributes - Issue referrer: vmware/go-vcloud-director#552
			//	Optional:            true,
			//	MarkdownDescription: "Allows to pass answer to question in vCD when ejecting from a VM which is powered on. True means 'Yes' as answer to question. Default is true",
			//	PlanModifiers: []planmodifier.Bool{
			//		boolpm.SetDefault(true),
			//	},
			// },
		},
	}
}
