package backup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

// TODO : Remove unused imports.
// ! This is outside import block because golangci-lint remove commented import.
// * Hashicorp Validators
// "github.com/Hashicorp/terraform-plugin-framework-validators/stringvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/boolvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/int64validator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/float64validator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/listvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/mapvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/setvalidator"

// * Hashicorp Plan Modifiers Resource
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"

// * Hashicorp Plan Modifiers DataSource
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/stringplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/boolplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/int64planmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/listplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/mapplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/setplanmodifier"

// * Hashicorp Default Values
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

// * FrangipaneTeam Custom Validators
// fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
// fboolvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/boolvalidator"
// fint64validator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/int64validator"
// flistvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/listvalidator"
// fmapvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/mapvalidator"
// fsetvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/setvalidator"

// * FrangipaneTeam Custom Plan Modifiers
// fstringplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/stringplanmodifier"
// fboolplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/boolplanmodifier"
// fint64planmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/int64planmodifier"
// flistplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/listplanmodifier"
// fmapplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/mapplanmodifier"
// fsetplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/setplanmodifier"

// How to use types generator:
// 1. Define the schema in the file internal/provider/backup/backup_schema.go
// 2. Add the resource or data source to the file internal/provider/provider_resources.go or internal/provider/provider_data_sources.go respectively
// 3. Launch the following command to generate golang structs for the schema:
// go run ./cmd/types-generator/*.go -file internal/provider/backup/backup_schema.go -resource cloudavenue_backup -is-resource.
func backupSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_backup` resource allows you to manage ...",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_backup` data source allows you to retrieve information about an ...",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the backup.",
				},
			},
			"type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Scope of the backup.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("vdc", "VDC", "VAPP", "vapp", "VM", "vm"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"target_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the target.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("target_id"), path.MatchRoot("target_name")),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"target_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the target.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("target_id"), path.MatchRoot("target_name")),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"policies": superschema.SuperSetNestedAttribute{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "The backup policies of the target.",
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
					"policy_id": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The ID of the backup policy.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							// Validators: []validator.Int64{
							// 	int64validator.NullIfAttributeIsSet(path.MatchRoot("policy_name")),
							// 	// int64validator.ExactlyOneOf(path.MatchRoot("policy_id"), path.MatchRoot("policy_name")),
							// },
						},
					},
					"policy_name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the backup policy.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							// Validators: []validator.String{
							// 	stringvalidator.ExactlyOneOf(path.MatchRoot("policy_id"), path.MatchRoot("policy_name")),
							// },
						},
					},
					"enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether the backup policy is enabled or not.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(true),
						},
					},
				},
			},
		},
	}
}
