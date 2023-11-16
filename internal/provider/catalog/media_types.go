package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type catalogMediaDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	CatalogID      types.String `tfsdk:"catalog_id"`
	CatalogName    types.String `tfsdk:"catalog_name"`
	Description    types.String `tfsdk:"description"`
	IsISO          types.Bool   `tfsdk:"is_iso"`
	OwnerName      types.String `tfsdk:"owner_name"`
	IsPublished    types.Bool   `tfsdk:"is_published"`
	CreatedAt      types.String `tfsdk:"created_at"`
	Size           types.Int64  `tfsdk:"size"`
	Status         types.String `tfsdk:"status"`
	StorageProfile types.String `tfsdk:"storage_profile"`
}

type MediaModel struct {
	ID types.String `tfsdk:"id"`
}
