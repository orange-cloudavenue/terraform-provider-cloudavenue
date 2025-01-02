package vdc

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func vdcsSchema() superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "List all vDC inside an Organization.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The ID of the resource. This value is system-generated.",
					Computed:            true,
				},
			},
			"vdcs": superschema.SuperListNestedAttributeOf[vdcRef]{
				DataSource: &schemaD.ListNestedAttribute{
					MarkdownDescription: "VDC list.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"name": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the vDC.",
							Computed:            true,
						},
					},
					"id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The ID of the vDC.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
