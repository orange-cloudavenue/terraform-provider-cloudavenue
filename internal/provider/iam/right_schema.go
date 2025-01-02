package iam

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func iamRightSuperSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a data source for available rights in Cloud Avenue.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The id of the right.",
					Computed:            true,
				},
			},
			"name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The name of the right.",
					Required:            true,
				},
			},
			"description": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "A description for the right.",
					Computed:            true,
				},
			},
			"category_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The category id for the right.",
					Computed:            true,
				},
			},
			// * Remove
			"bundle_key": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The bundle key for the right.",
					Computed:            true,
				},
			},
			"right_type": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The right type for the right.",
					Computed:            true,
				},
			},
			"implied_rights": superschema.SuperSetNestedAttribute{
				DataSource: &schemaD.SetNestedAttribute{
					MarkdownDescription: "The list of rights that are implied with this one.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"name": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "Name of the implied right.",
							Computed:            true,
						},
					},
					"id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "ID of the implied right.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
