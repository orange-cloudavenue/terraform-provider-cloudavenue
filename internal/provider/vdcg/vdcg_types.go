package vdcg

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type vdcgModel struct {
	ID          supertypes.StringValue        `tfsdk:"id"`
	Name        supertypes.StringValue        `tfsdk:"name"`
	Description supertypes.StringValue        `tfsdk:"description"`
	VDCIDs      supertypes.SetValueOf[string] `tfsdk:"vdc_ids"`
	Type        supertypes.StringValue        `tfsdk:"type"`
	Status      supertypes.StringValue        `tfsdk:"status"`
}

func (rm *vdcgModel) Copy() *vdcgModel {
	x := &vdcgModel{}
	utils.ModelCopy(rm, x)
	return x
}
