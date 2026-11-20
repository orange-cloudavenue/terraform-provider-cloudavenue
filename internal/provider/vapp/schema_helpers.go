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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Schema helpers copied from the shared vdc and vapp wrapper packages so that this package no longer
// imports them (the govcd-based wrappers are still used by the vm package).

/*
vdcSuperSchema is the `vdc` attribute (superschema.StringAttribute).

	For the resource :
	Optional: true
	Computed: true
	RequiresReplace
	UseStateForUnknown

	For the data source :
	Optional: true
	Computed: true
*/
func vdcSuperSchema() superschema.StringAttribute {
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

const (
	schemaVAppID   = "vapp_id"
	schemaVAppName = "vapp_name"
	descVAppID     = "ID of the vApp."
	descVAppName   = "Name of the vApp."
)

// vappRefSuperSchema returns the superschema of `vapp_id` and `vapp_name`.
func vappRefSuperSchema() map[string]superschema.Attribute {
	return map[string]superschema.Attribute{
		schemaVAppID: superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: descVAppID,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
				},
				Optional: true,
			},
			Resource: &schemaR.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		schemaVAppName: superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: descVAppName,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_id"), path.MatchRoot("vapp_name")),
				},
				Optional: true,
			},
			Resource: &schemaR.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
	}
}
