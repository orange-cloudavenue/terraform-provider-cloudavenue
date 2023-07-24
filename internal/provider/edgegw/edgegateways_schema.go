// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func edgeGatewaysSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "The edge gateways data source show the list of edge gateways of an organization.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"edge_gateways": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of Edge Gateways.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tier0_vrf_name": edgegwSchema().GetDataSource(ctx).Attributes["tier0_vrf_name"],
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the Edge Gateway.",
							Computed:            true,
						},
						"id":          edgegwSchema().GetDataSource(ctx).Attributes["id"],
						"owner_type":  edgegwSchema().GetDataSource(ctx).Attributes["owner_type"],
						"owner_name":  edgegwSchema().GetDataSource(ctx).Attributes["owner_name"],
						"description": edgegwSchema().GetDataSource(ctx).Attributes["description"],
						"lb_enabled":  edgegwSchema().GetDataSource(ctx).Attributes["lb_enabled"],
					},
				},
			},
		},
	}
}
