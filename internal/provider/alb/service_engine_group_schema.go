package alb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

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
// 1. Define the schema in the file internal/provider/cloudavenue/alb_service_engine_group_datasource_name_schema.go
// 2. Add the resource or data source to the file internal/provider/provider_resources.go or internal/provider/provider_data_sources.go respectively
// 3. Launch the following command to generate golang structs for the schema:
// go run ./cmd/types-generator/*.go -file internal/provider/cloudavenue/alb_service_engine_group_datasource_name_schema.go -resource cloudavenue_cloudavenue_alb_service_engine_group_datasource_name -is-resource.
func albServiceEngineGroupSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_alb_service_engine_group` data source allows you to retrieve information about an ALB Service Engine Group.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Computed:            true,
					Optional:            true,
					MarkdownDescription: "The ID of the ALB Service Engine Group.",
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The name of the ALB Service Engine Group.",
					Computed:            true,
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Edge gateway ID in which ALB Service Engine Group should be located.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_id"), path.MatchRoot("edge_gateway_name")),
					},
					Computed: true,
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Edge gateway Name in which ALB Service Engine Group should be located.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_id"), path.MatchRoot("edge_gateway_name")),
					},
					Computed: true,
				},
			},
			"max_virtual_services": superschema.SuperInt64Attribute{
				DataSource: &schemaD.Int64Attribute{
					Computed:            true,
					MarkdownDescription: "The maximum number of virtual services that can be deployed on the ALB Service Engine Group.",
				},
			},
			"reserved_virtual_services": superschema.SuperInt64Attribute{
				DataSource: &schemaD.Int64Attribute{
					Computed:            true,
					MarkdownDescription: "The number of reserved virtual services for the ALB Service Engine Group.",
				},
			},
			"deployed_virtual_services": superschema.SuperInt64Attribute{
				DataSource: &schemaD.Int64Attribute{
					Computed:            true,
					MarkdownDescription: "The number of deployed virtual services on the ALB Service Engine Group.",
				},
			},
		},
	}
}
