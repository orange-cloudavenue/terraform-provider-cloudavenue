package vdc

import (
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
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
					"vdc_name": superschema.SuperStringAttribute{
						Deprecated: &superschema.Deprecated{
							DeprecationMessage:                "Use `name` instead.",
							ComputeMarkdownDeprecationMessage: true,
							Renamed:                           true,
							FromAttributeName:                 "vdc_name",
							TargetAttributeName:               "name",
							TargetRelease:                     "v0.19.0",
							OnlyDataSource:                    utils.TakeBoolPointer(true),
						},
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "VDC name.",
							Computed:            true,
						},
					},
					"vdc_uuid": superschema.SuperStringAttribute{
						Deprecated: &superschema.Deprecated{
							DeprecationMessage:                "Use `id` instead.",
							ComputeMarkdownDeprecationMessage: true,
							Renamed:                           true,
							FromAttributeName:                 "id",
							TargetAttributeName:               "name",
							TargetRelease:                     "v0.19.0",
							OnlyDataSource:                    utils.TakeBoolPointer(true),
						},
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "VDC UUID.",
							Computed:            true,
						},
					},
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
