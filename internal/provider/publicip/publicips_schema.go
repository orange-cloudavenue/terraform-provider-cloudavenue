package publicip

import (
	"golang.org/x/net/context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func publicIPsSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "The public IP data source displays the list of public IP addresses.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"public_ips": schema.ListNestedAttribute{
				MarkdownDescription: "A list of public IPs.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                publicIPSchema().GetDataSource(ctx).Attributes["id"],
						"public_ip":         publicIPSchema().GetDataSource(ctx).Attributes["public_ip"],
						"edge_gateway_name": publicIPSchema().GetDataSource(ctx).Attributes["edge_gateway_name"],
						"edge_gateway_id":   publicIPSchema().GetDataSource(ctx).Attributes["edge_gateway_id"],
					},
				},
			},
		},
	}
}
