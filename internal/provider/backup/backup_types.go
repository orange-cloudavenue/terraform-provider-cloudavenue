package backup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type backupModel struct {
	ID         supertypes.StringValue    `tfsdk:"id"`
	Policies   supertypes.SetNestedValue `tfsdk:"policies"`
	TargetID   supertypes.StringValue    `tfsdk:"target_id"`
	TargetName supertypes.StringValue    `tfsdk:"target_name"`
	Type       supertypes.StringValue    `tfsdk:"type"`
}

// * Policies.
type backupModelPolicies []backupModelPolicy

// * Policies.
type backupModelPolicy struct {
	Enabled    supertypes.BoolValue   `tfsdk:"enabled"`
	PolicyID   supertypes.Int64Value  `tfsdk:"policy_id"`
	PolicyName supertypes.StringValue `tfsdk:"policy_name"`
}

func (rm *backupModel) Copy() *backupModel {
	x := &backupModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetPolicies returns the value of the Policies field.
func (rm *backupModel) GetPolicies(ctx context.Context) (values backupModelPolicies, diags diag.Diagnostics) {
	values = make(backupModelPolicies, 0)
	d := rm.Policies.Get(ctx, &values, false)
	return values, d
}
