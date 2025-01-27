package alb

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func serviceEngineGroupSchema(_ context.Context) superschema.Schema {
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
