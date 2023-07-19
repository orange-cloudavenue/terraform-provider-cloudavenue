package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func catalogsSchema() schema.Schema {
	return schema.Schema{
		Description: "The catalogs datasource show the details of all the catalogs.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"catalogs_name": schema.ListAttribute{
				MarkdownDescription: "List of catalogs name.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"catalogs": schema.MapNestedAttribute{
				MarkdownDescription: "Map of catalogs.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: catalogDatasourceAttributes(),
				},
			},
		},
	}
}
