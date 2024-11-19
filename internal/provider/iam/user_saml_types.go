package iam

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type UserSAMLModel struct {
	ID              supertypes.StringValue `tfsdk:"id"`
	UserName        supertypes.StringValue `tfsdk:"user_name"`
	RoleName        supertypes.StringValue `tfsdk:"role_name"`
	Enabled         supertypes.BoolValue   `tfsdk:"enabled"`
	DeployedVMQuota supertypes.Int64Value  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   supertypes.Int64Value  `tfsdk:"stored_vm_quota"`
	TakeOwnership   supertypes.BoolValue   `tfsdk:"take_ownership"`
}

func (rm *UserSAMLModel) Copy() *UserSAMLModel {
	x := &UserSAMLModel{}
	utils.ModelCopy(rm, x)
	return x
}
