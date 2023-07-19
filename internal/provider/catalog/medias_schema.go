package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

func mediasSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Catalog medias allows you to",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "manage a medias in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "retrieve information about a medias in Cloud Avenue.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The ID of the medias.",
					Computed:            true,
				},
			},
			"medias": superschema.MapNestedAttribute{
				DataSource: &schemaD.MapNestedAttribute{
					MarkdownDescription: "The map of medias.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The ID of the media.",
							Computed:            true,
						},
					},
					"catalog_id": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The ID of the catalog.",
							Computed:            true,
						},
					},
					"catalog_name": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the catalog.",
							Computed:            true,
						},
					},
					"name": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the media.",
							Computed:            true,
						},
					},
					"description": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The description of the media.",
							Computed:            true,
						},
					},
					"is_iso": superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "`True` if the media is an ISO.",
							Computed:            true,
						},
					},
					"owner_name": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the owner of the media.",
							Computed:            true,
						},
					},
					"is_published": superschema.BoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							MarkdownDescription: "`True` if the media is published.",
							Computed:            true,
						},
					},
					"created_at": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The date and time when the media was created.",
							Computed:            true,
						},
					},
					"size": superschema.Int64Attribute{
						DataSource: &schemaD.Int64Attribute{
							MarkdownDescription: "The size of the media in bytes.",
							Computed:            true,
						},
					},
					"status": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The status of the media.",
							Computed:            true,
						},
					},
					"storage_profile": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The storage profile of the media.",
							Computed:            true,
						},
					},
				},
			},
			"medias_name": superschema.ListAttribute{
				DataSource: &schemaD.ListAttribute{
					MarkdownDescription: "The list of medias name.",
					Computed:            true,
					ElementType:         types.StringType,
				},
			},
			catalogName: mediaSchema().Attributes[catalogName],
			catalogID:   mediaSchema().Attributes[catalogID],
		},
	}
}
