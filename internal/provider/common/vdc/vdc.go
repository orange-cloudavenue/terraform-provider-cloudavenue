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
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type VDC struct {
	*v1.VDC
}

/*
SuperSchema

	For the resource :
	Optional: true
	Computed: true
	RequiresReplace
	UseStateForUnknown

	For the data source :
	Optional: true
	Computed: true
*/
func SuperSchema() superschema.StringAttribute {
	return superschema.StringAttribute{
		Common: &schemaR.StringAttribute{
			MarkdownDescription: "The name of vDC to use, optional if defined at provider level.",
			Optional:            true,
			Computed:            true,
		},
		Resource: &schemaR.StringAttribute{
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
	}
}

func SuperSchemaSuperType() superschema.SuperStringAttribute {
	return superschema.SuperStringAttribute{
		Common: &schemaR.StringAttribute{
			MarkdownDescription: "The name of vDC to use, optional if defined at provider level.",
			Optional:            true,
			Computed:            true,
		},
		Resource: &schemaR.StringAttribute{
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
	}
}

/*
Init

If vDC is not defined at data source level, use the one defined at provider level.
*/
func Init(c *client.CloudAvenue, vdc types.String) (VDC, diag.Diagnostics) {
	var (
		d = diag.Diagnostics{}
		v = VDC{}

		err error
	)

	v.VDC, err = c.GetVDC(vdc.ValueString())
	if err != nil {
		d.AddError("Unable to get VDC", err.Error())
		return VDC{}, d
	}

	return v, nil
}
