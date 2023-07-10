package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type catalogMediasDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Medias      types.Map    `tfsdk:"medias"`
	MediasName  types.List   `tfsdk:"medias_name"`
	CatalogName types.String `tfsdk:"catalog_name"`
	CatalogID   types.String `tfsdk:"catalog_id"`
}
