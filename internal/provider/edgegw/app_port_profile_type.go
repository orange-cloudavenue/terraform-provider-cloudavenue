package edgegw

import (
	"context"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type AppPortProfileModel struct {
	ID              supertypes.StringValue                                         `tfsdk:"id"`
	Name            supertypes.StringValue                                         `tfsdk:"name"`
	EdgeGatewayID   supertypes.StringValue                                         `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue                                         `tfsdk:"edge_gateway_name"`
	Description     supertypes.StringValue                                         `tfsdk:"description"`
	AppPorts        supertypes.ListNestedObjectValueOf[AppPortProfileModelAppPort] `tfsdk:"app_ports"`
}

type AppPortProfileModelAppPort struct {
	Protocol supertypes.StringValue        `tfsdk:"protocol"`
	Ports    supertypes.SetValueOf[string] `tfsdk:"ports"`
}

func (rm *AppPortProfileModel) Copy() *AppPortProfileModel {
	x := &AppPortProfileModel{}
	utils.ModelCopy(rm, x)
	return x
}

// toNsxtAppPortProfile converts the AppPortProfileModel to the NSX-T API representation.
func (rm *AppPortProfileModel) toNsxtAppPortProfilePorts(ctx context.Context) (nsxtAppPortProfilePorts []govcdtypes.NsxtAppPortProfilePort, diags diag.Diagnostics) {
	appPorts, d := rm.AppPorts.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	nsxtAppPortProfilePorts = make([]govcdtypes.NsxtAppPortProfilePort, 0)
	for _, appPort := range appPorts {
		destPorts, d := appPort.Ports.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		nsxtAppPortProfilePorts = append(nsxtAppPortProfilePorts, govcdtypes.NsxtAppPortProfilePort{
			Protocol:         appPort.Protocol.Get(),
			DestinationPorts: destPorts,
		})
	}

	return nsxtAppPortProfilePorts, diags
}
