package org

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

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
// 1. Define the schema in the file internal/provider/org/certificate_library_datasources_go_schema.go
// 2. Add the resource or data source to the file internal/provider/provider_resources.go or internal/provider/provider_data_sources.go respectively
// 3. Launch the following command to generate golang structs for the schema:
// go run ./cmd/types-generator/*.go -file internal/provider/org/certificate_library_datasources_go_schema.go -resource cloudavenue_org_certificate_library -is-resource.
func certificateLibrarySchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_org_certificate_library` resource allows you to manage certificate in your organization library.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_org_certificate_library` data source allows you to retrieve information about an certificate in your organization library.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the certificate library.",
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
					MarkdownDescription: "The name of the certificate library.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the certificate library.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"certificate": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The certificate content. It can be a PEM encoded certificate or a certificate chain. Contains all strings including the BEGIN CERTIFICATE and END CERTIFICATE lines. No empty lines are allowed.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"private_key": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The private key of the certificate in PEM format.",
					Optional:            true,
					Sensitive:           true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"passphrase": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The passphrase of the private key.",
					Optional:            true,
					Sensitive:           true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
		},
	}
}
