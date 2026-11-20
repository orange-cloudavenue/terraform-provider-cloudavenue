// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vapp

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Schema helpers copied from the shared vdc wrapper package so that this package no longer
// imports it (the govcd-based wrapper is still used by the vm package).

// vdcSuperSchemaSuperType is the `vdc` attribute (superschema.SuperStringAttribute).
func vdcSuperSchemaSuperType() superschema.SuperStringAttribute {
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
