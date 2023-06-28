package edgegw

import "github.com/hashicorp/terraform-plugin-framework/types"

//go:generate terraform-framework-types-generator -file $GOFILE -struct portProfilesResourceModel
type portProfilesResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	VDC         types.String `tfsdk:"vdc"`
	Description types.String `tfsdk:"description"`
	AppPorts    types.List   `tfsdk:"app_ports" tftype:"portProfilesResourceModelAppPorts"`
}

type portProfilesResourceModelAppPorts []portProfilesResourceModelAppPort

type portProfilesResourceModelAppPort struct {
	Protocol types.String `tfsdk:"protocol"`
	Ports    types.Set    `tfsdk:"ports"`
}
