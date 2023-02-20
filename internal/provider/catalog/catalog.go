package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type catalogDataSourceStruct struct {
	ID                          types.String `tfsdk:"id"`
	CatalogName                 types.String `tfsdk:"catalog_name"`
	CreatedAt                   types.String `tfsdk:"created_at"`
	Description                 types.String `tfsdk:"description"`
	PreserveIdentityInformation types.Bool   `tfsdk:"preserve_identity_information"`
	Href                        types.String `tfsdk:"href"`
	OwnerName                   types.String `tfsdk:"owner_name"`
	NumberOfMedia               types.Int64  `tfsdk:"number_of_media"`
	MediaItemList               types.List   `tfsdk:"media_item_list"`
	IsShared                    types.Bool   `tfsdk:"is_shared"`
	IsPublished                 types.Bool   `tfsdk:"is_published"`
	IsLocal                     types.Bool   `tfsdk:"is_local"`
	IsCached                    types.Bool   `tfsdk:"is_cached"`
}

// schemaDataSource returns the catalog schema for the data source.
func schemaDataSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"catalog_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Name of the catalog.",
		},
		"created_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Time stamp of when the catalog was created",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Description of the catalog.",
		},
		"preserve_identity_information": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Include BIOS UUIDs and MAC addresses in the downloaded OVF package. Preserving the identity information limits the portability of the package and you should use it only when necessary.",
		},
		"href": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Catalog HREF",
		},
		"owner_name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Owner name from the catalog.",
		},
		"number_of_media": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "Number of Medias this catalog contains.",
		},
		"media_item_list": schema.ListAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			MarkdownDescription: "List of Media items in this catalog",
		},
		"is_shared": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this catalog is shared.",
		},
		"is_local": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this catalog belongs to the current organization.",
		},
		"is_published": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this catalog is shared to all organizations.",
		},
		"is_cached": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this catalog is cached.",
		},
	}
}
