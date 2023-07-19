package catalog

import "github.com/hashicorp/terraform-plugin-framework/types"

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
