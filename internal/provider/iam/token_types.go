package iam

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type TokenModel struct {
	FileName      supertypes.StringValue `tfsdk:"file_name"`
	ID            supertypes.StringValue `tfsdk:"id"`
	Name          supertypes.StringValue `tfsdk:"name"`
	PrintToken    supertypes.BoolValue   `tfsdk:"print_token"`
	SaveInFile    supertypes.BoolValue   `tfsdk:"save_in_file"`
	SaveInTfstate supertypes.BoolValue   `tfsdk:"save_in_tfstate"`
	Token         supertypes.StringValue `tfsdk:"token"`
}

func (rm *TokenModel) Copy() *TokenModel {
	x := &TokenModel{}
	utils.ModelCopy(rm, x)
	return x
}
