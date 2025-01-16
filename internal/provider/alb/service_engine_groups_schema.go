package alb

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func serviceEngineGroupsSchema(ctx context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_alb_service_engine_groups` data source allows you to retrieve information about all the Service Engine Group of an Edge Gateway.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the service engine groups.",
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
			"service_engine_groups": superschema.SuperListNestedAttributeOf[serviceEngineGroupModel]{
				DataSource: &schemaD.ListNestedAttribute{
					Computed:            true,
					MarkdownDescription: "The list of service engine groups.",
				},
				Attributes: map[string]superschema.Attribute{
					"id":                        serviceEngineGroupSchema(ctx).Attributes["id"],
					"name":                      serviceEngineGroupSchema(ctx).Attributes["name"],
					"edge_gateway_id":           serviceEngineGroupSchema(ctx).Attributes["edge_gateway_id"],
					"edge_gateway_name":         serviceEngineGroupSchema(ctx).Attributes["edge_gateway_name"],
					"max_virtual_services":      serviceEngineGroupSchema(ctx).Attributes["max_virtual_services"],
					"reserved_virtual_services": serviceEngineGroupSchema(ctx).Attributes["reserved_virtual_services"],
					"deployed_virtual_services": serviceEngineGroupSchema(ctx).Attributes["deployed_virtual_services"],
				},
			},
		},
	}
}
