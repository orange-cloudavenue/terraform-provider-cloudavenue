package vm

import "github.com/hashicorp/terraform-plugin-framework/types"

type VMDataSourceModel struct { //nolint:revive
	ID          types.String `tfsdk:"id"`
	VDC         types.String `tfsdk:"vdc"`
	Name        types.String `tfsdk:"name"`
	VappName    types.String `tfsdk:"vapp_name"`
	VappID      types.String `tfsdk:"vapp_id"`
	Description types.String `tfsdk:"description"`
	State       types.Object `tfsdk:"state"`
	Resource    types.Object `tfsdk:"resource"`
	Settings    types.Object `tfsdk:"settings"`
}
