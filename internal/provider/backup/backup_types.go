/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package backup

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type backupModel struct {
	ID         supertypes.Int64Value     `tfsdk:"id"`
	Policies   supertypes.SetNestedValue `tfsdk:"policies"`
	TargetID   supertypes.StringValue    `tfsdk:"target_id"`
	TargetName supertypes.StringValue    `tfsdk:"target_name"`
	Type       supertypes.StringValue    `tfsdk:"type"`
}

// * Policies.
type backupModelPolicies []backupModelPolicy

// * Policy.
type backupModelPolicy struct {
	PolicyID   supertypes.Int64Value  `tfsdk:"policy_id"`
	PolicyName supertypes.StringValue `tfsdk:"policy_name"`
}

// NewBackup returns a new backupModel.
func newBackup() *backupModel {
	return &backupModel{
		ID:         supertypes.NewInt64Unknown(),
		Policies:   supertypes.NewSetNestedNull(types.ObjectType{AttrTypes: map[string]attr.Type{"policy_id": types.Int64Type, "policy_name": types.StringType}}),
		TargetID:   supertypes.NewStringNull(),
		TargetName: supertypes.NewStringNull(),
		Type:       supertypes.NewStringNull(),
	}
}

// Copy returns a copy of the backupModel.
func (rm *backupModel) Copy() *backupModel {
	x := &backupModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetPolicies returns the value of the Policies field.
func (rm *backupModel) getPolicies(ctx context.Context) (values *backupModelPolicies, diags diag.Diagnostics) {
	values = &backupModelPolicies{}
	d := rm.Policies.Get(ctx, &values, false)
	return values, d
}

// Get target ID or Name.
func (rm *backupModel) getTargetIDOrName() string {
	if rm.TargetID.IsKnown() {
		return rm.TargetID.Get()
	}
	return rm.TargetName.Get()
}
