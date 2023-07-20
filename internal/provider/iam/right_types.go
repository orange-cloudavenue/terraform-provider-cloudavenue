package iam

import "github.com/hashicorp/terraform-plugin-framework/types"

type iamRightDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CategoryID    types.String `tfsdk:"category_id"`
	BundleKey     types.String `tfsdk:"bundle_key"`
	RightType     types.String `tfsdk:"right_type"`
	ImpliedRights types.Set    `tfsdk:"implied_rights"`
}
