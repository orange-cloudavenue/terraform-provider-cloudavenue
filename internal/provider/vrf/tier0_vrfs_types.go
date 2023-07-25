package vrf

import "github.com/hashicorp/terraform-plugin-framework/types"

type tier0VrfsDataSourceModel struct {
	ID    types.String   `tfsdk:"id"`
	Names []types.String `tfsdk:"names"`
}
