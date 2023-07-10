package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type catalogDataSourceModel struct {
	// BASE
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	OwnerName   types.String `tfsdk:"owner_name"`

	// SPECIFIC DATA SOURCE
	PreserveIdentityInformation types.Bool  `tfsdk:"preserve_identity_information"`
	NumberOfMedia               types.Int64 `tfsdk:"number_of_media"`
	MediaItemList               types.List  `tfsdk:"media_item_list"`
	IsShared                    types.Bool  `tfsdk:"is_shared"`
	IsPublished                 types.Bool  `tfsdk:"is_published"`
	IsLocal                     types.Bool  `tfsdk:"is_local"`
	IsCached                    types.Bool  `tfsdk:"is_cached"`
}

type catalogResourceModel struct {
	// BASE
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	OwnerName   types.String `tfsdk:"owner_name"`

	// SPECIFIC RESOURCE
	StorageProfile  types.String `tfsdk:"storage_profile"`
	DeleteForce     types.Bool   `tfsdk:"delete_force"`
	DeleteRecursive types.Bool   `tfsdk:"delete_recursive"`
}

type vAppTemplateDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	TemplateName types.String `tfsdk:"template_name"`
	TemplateID   types.String `tfsdk:"template_id"`
	CatalogID    types.String `tfsdk:"catalog_id"`
	CatalogName  types.String `tfsdk:"catalog_name"`
	Description  types.String `tfsdk:"description"`
	CreatedAt    types.String `tfsdk:"created_at"`
	VMNames      types.List   `tfsdk:"vm_names"`
}
