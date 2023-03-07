package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type catalogMediaDataStruct struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	IsISO          types.Bool   `tfsdk:"is_iso"`
	OwnerName      types.String `tfsdk:"owner_name"`
	IsPublished    types.Bool   `tfsdk:"is_published"`
	CreatedAt      types.String `tfsdk:"created_at"`
	Size           types.Int64  `tfsdk:"size"`
	Status         types.String `tfsdk:"status"`
	StorageProfile types.String `tfsdk:"storage_profile"`
}

func schemaCatalogDataSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The ID of the catalog media.",
		},
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The name of the media.",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The description of the media.",
		},
		"is_iso": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this media file is an Iso.",
		},
		"owner_name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The name of the owner.",
		},
		"is_published": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this media file is in a published catalog.",
		},
		"created_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The creation date of the media.",
		},
		"size": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "The size of the media in bytes.",
		},
		"status": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The media status.",
		},
		"storage_profile": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The name of the storage profile.",
		},
	}
}
