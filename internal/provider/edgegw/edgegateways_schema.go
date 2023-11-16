// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

func edgeGatewaysSuperSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The edge gateways data source show the list of edge gateways of an organization.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Generated ID of the resource.",
				},
			},
			"edge_gateways": superschema.SuperListNestedAttributeOf[edgeGatewayDataSourceModelEdgeGateway]{
				DataSource: &schemaD.ListNestedAttribute{
					Computed:            true,
					MarkdownDescription: "A list of Edge Gateways.",
				},
				Attributes: superschema.Attributes{
					"tier0_vrf_name": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the Tier-0 VRF to which the Edge Gateway is attached.",
						},
					},
					"name": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the Edge Gateway.",
						},
					},
					"id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the Edge Gateway.",
						},
					},
					"owner_type": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the Edge Gateway owner. Must be vdc or vdc-group.",
						},
					},
					"owner_name": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the Edge Gateway owner.",
						},
					},
					"description": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the Edge Gateway.",
						},
					},
					"lb_enabled": superschema.SuperBoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Load Balancing state on the Edge Gateway.",
						},
					},
				},
			},
		},
	}
}
