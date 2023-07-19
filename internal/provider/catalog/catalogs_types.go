package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type catalogsDataSourceModel struct {
	ID           types.String                      `tfsdk:"id"`
	Catalogs     map[string]catalogDataSourceModel `tfsdk:"catalogs"`
	CatalogsName types.List                        `tfsdk:"catalogs_name"`
}
