package vdc

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func vdcsSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "List all vDC inside an Organization.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"vdcs": schema.ListNestedAttribute{
				MarkdownDescription: "VDC list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vdc_name": schema.StringAttribute{
							MarkdownDescription: "VDC name.",
							Computed:            true,
						},
						"vdc_uuid": schema.StringAttribute{
							MarkdownDescription: "VDC UUID.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
