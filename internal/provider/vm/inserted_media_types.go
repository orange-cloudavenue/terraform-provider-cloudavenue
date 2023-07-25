package vm

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// vmInsertedMediaResource is the resource implementation.
type vmInsertedMediaResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	vapp   vapp.VAPP
	org    org.Org
}

type vmInsertedMediaResourceModel struct {
	ID       types.String `tfsdk:"id"`
	VDC      types.String `tfsdk:"vdc"`
	Catalog  types.String `tfsdk:"catalog"`
	Name     types.String `tfsdk:"name"`
	VAppName types.String `tfsdk:"vapp_name"`
	VAppID   types.String `tfsdk:"vapp_id"`
	VMName   types.String `tfsdk:"vm_name"`
	// EjectForce types.Bool   `tfsdk:"eject_force"` - Disable attributes - Issue referrer: vmware/go-vcloud-director#552
}
