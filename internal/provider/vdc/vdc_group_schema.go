package vdc

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

func groupSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdc_group` resource allows you to manage VDC Group.",
			Deprecated: superschema.DeprecatedResource{
				DeprecationMessage:                "The `cloudavenue_vdc_group` resource is deprecated. Please use the `cloudavenue_vdcg` resource instead.",
				ComputeMarkdownDeprecationMessage: true,
				Renamed:                           true,
				TargetResourceName:                "cloudavenue_vdcg",
				TargetRelease:                     "v0.30.0",
				LinkToMigrationGuide:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc_group#how-to-migrate-existing-resources",
				LinkToNewResourceDoc:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg",
				LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/18",
				LinkToIssue:                       "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/869",
			},
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vdc_group` data source allows you to retrieve informations about an existing VDC Group.",
			Deprecated: superschema.DeprecatedResource{
				DeprecationMessage:                "The `cloudavenue_vdc_group` data source is deprecated. Please use the `cloudavenue_vdcg` data source instead.",
				ComputeMarkdownDeprecationMessage: true,
				Renamed:                           true,
				TargetResourceName:                "cloudavenue_vdcg",
				TargetRelease:                     "v0.30.0",
				LinkToNewResourceDoc:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/data-sources/vdcg",
				LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/18",
				LinkToIssue:                       "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/869",
			},
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the VDC Group.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the VDC Group.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the VDC Group.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"vdc_ids": superschema.SuperSetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "The list of VDC IDs of the VDC Group.",
					ElementType:         supertypes.StringType{},
				},
				Resource: &schemaR.SetAttribute{
					Required: true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
			"status": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The status of the VDC Group.",
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("SAVING", "SAVED", "CONFIGURING", "REALIZED", "REALIZATION_FAILED", "DELETING", "DELETE_FAILED", "OBJECT_NOT_FOUND", "UNCONFIGURED"),
					},
				},
			},
			"type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of the VDC Group.",
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("LOCAL", "UNIVERSAL"),
					},
				},
			},
		},
	}
}
