package vcda

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type vcdaIPResourceModel struct {
	ID        supertypes.StringValue `tfsdk:"id"`
	IPAddress supertypes.StringValue `tfsdk:"ip_address"`
}

func (rm *vcdaIPResourceModel) Copy() *vcdaIPResourceModel {
	x := &vcdaIPResourceModel{}
	utils.ModelCopy(rm, x)
	return x
}
