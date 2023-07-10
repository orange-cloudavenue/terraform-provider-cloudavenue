package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func catalogMediaDataSourceModelType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"name":            types.StringType,
		"catalog_id":      types.StringType,
		"catalog_name":    types.StringType,
		"description":     types.StringType,
		"is_iso":          types.BoolType,
		"owner_name":      types.StringType,
		"is_published":    types.BoolType,
		"created_at":      types.StringType,
		"size":            types.Int64Type,
		"status":          types.StringType,
		"storage_profile": types.StringType,
	}
}
