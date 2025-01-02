package s3

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type CredentialModel struct {
	ID            supertypes.StringValue `tfsdk:"id"`
	Username      supertypes.StringValue `tfsdk:"username"`
	FileName      supertypes.StringValue `tfsdk:"file_name"`
	SaveInFile    supertypes.BoolValue   `tfsdk:"save_in_file"`
	PrintToken    supertypes.BoolValue   `tfsdk:"print_token"`
	SaveInTFState supertypes.BoolValue   `tfsdk:"save_in_tfstate"`
	AccessKey     supertypes.StringValue `tfsdk:"access_key"`
	SecretKey     supertypes.StringValue `tfsdk:"secret_key"`
}

func (rm *CredentialModel) Copy() *CredentialModel {
	x := &CredentialModel{}
	utils.ModelCopy(rm, x)
	return x
}
