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
	"github.com/orange-cloudavenue/common-go/utils"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/validators"
)

/*
vdcSchema

This function is used to create the superschema for the vdc resource and datasource.
*/
func vdcSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "Provides a Cloud Avenue vDC (Virtual Data Center) ",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource. This can be used to create, update and delete vDC.\n\n",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source. This can be used to reference a vDC and use its data within other resources or data sources.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the vDC.",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the vDC.",
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					Validators: []validator.String{
						validators.ResourceName("vdc"),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A description of the vDC.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},

			// * Availability properties
			"service_class": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Defines the service class tier for the Virtual Data Center (vDC), indicating its level of service and performance.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					// TODO: Add link to devex documentation when available
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Default: stringdefault.StaticString("STD"),
				},
			},
			"disponibility_class": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies the service class tier for the Virtual Data Center (vDC), defining its performance and service level.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					// TODO: Add link to devex documentation when available
					MarkdownDescription: "",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Default: stringdefault.StaticString("ONE-ROOM"),
				},
			},

			// * Billing properties
			"billing_model": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Choose Billing model of compute resources.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					// TODO: Add link to devex documentation when available
					MarkdownDescription: "",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Default: stringdefault.StaticString("PAYG"),
				},
			},

			"storage_billing_model": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Choose Billing model of storage resources.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					// TODO: Add link to devex documentation when available
					MarkdownDescription: "",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Default: stringdefault.StaticString("PAYG"),
				},
			},

			// * Resource properties
			"vcpu": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Number of virtual CPU allocated to the vDC.",
					// TODO: Add link to devex documentation when available

					Computed: true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true, // ? Make optional the time that the attribute cpu_allocated is deprecated.
					Validators: []validator.Int64{
						int64validator.ExactlyOneOf(path.MatchRoot("vcpu"), path.MatchRoot("cpu_allocated")),
					},
				},
			},
			"memory": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Amount of memory in *GiB* allocated to the vDC.",
					// TODO: Add link to devex documentation when available

					Computed: true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true, // ? Make optional the time that the attribute memory_allocated is deprecated.
					Validators: []validator.Int64{
						int64validator.ExactlyOneOf(path.MatchRoot("memory"), path.MatchRoot("memory_allocated")),
					},
				},
			},

			"storage_profiles": superschema.SuperSetNestedAttributeOf[vdcModelStorageProfile]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "List of storage profiles for this vDC.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Required: true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the storage profile.",
							Computed:            true,
						},
					},
					"class": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Defines the classification tier of the storage profile, indicating its performance and intended use case.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							// TODO: Add link to devex documentation when available
							MarkdownDescription: "",
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"limit": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Max number in *GiB* of units allocated for this storage profile.",
						},
						Resource: &schemaR.Int64Attribute{
							Required: true,
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"default": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Set this storage profile as default for this vDC. Only one storage profile can be default per vDC.",
						},
						Resource: &schemaR.BoolAttribute{
							Required: true,
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
					"used": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Number in *GiB* of units used for this storage profile.",
							Computed:            true,
						},
					},
				},
			},

			// ! Deprecated fields - Maintain for backward compatibility
			"cpu_speed_in_mhz": superschema.SuperInt64Attribute{
				Deprecated: &superschema.Deprecated{
					DeprecationMessage:         "The attribute `cpu_speed_in_mhz` is no longer mandatory and will be returned as read-only starting from version v1.0.0.",
					MarkdownDeprecationMessage: "The attribute `cpu_speed_in_mhz` is no longer mandatory and will be returned as read-only starting from version v1.0.0.",
					OnlyResource:               utils.ToPTR(true),
				},
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Frequency of the VCPUs in MHz.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
				},
				DataSource: &schemaD.Int64Attribute{},
			},

			"cpu_allocated": superschema.SuperInt64Attribute{
				Deprecated: &superschema.Deprecated{
					DeprecationMessage:                "The attribute `cpu_allocated` has been deprecated and will be removed in a future release. Please use the `vcpu` attribute instead.",
					ComputeMarkdownDeprecationMessage: true,
					Removed:                           true,
					FromAttributeName:                 "cpu_allocated",
					TargetAttributeName:               "vcpu",
					TargetRelease:                     "v1.0.0",
					LinkToIssue:                       "",
				},
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "CPU capacity in *MHz* that is committed to be available or used as a limit in PAYG mode.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Validators: []validator.Int64{
						int64validator.ExactlyOneOf(path.MatchRoot("vcpu"), path.MatchRoot("cpu_allocated")),
					},
				},
			},
			"memory_allocated": superschema.SuperInt64Attribute{
				Deprecated: &superschema.Deprecated{
					DeprecationMessage:                "The attribute `memory_allocated` has been deprecated and will be removed in a future release. Please use the `memory` attribute instead.",
					ComputeMarkdownDeprecationMessage: true,
					Renamed:                           true,
					FromAttributeName:                 "memory_allocated",
					TargetAttributeName:               "memory",
					TargetRelease:                     "v1.0.0",
					LinkToIssue:                       "",
				},
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Memory capacity in GiB that is committed to be available or used as a limit in PAYG mode.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Validators: []validator.Int64{
						int64validator.ExactlyOneOf(path.MatchRoot("memory"), path.MatchRoot("memory_allocated")),
					},
				},
			},
		},
	}
}
