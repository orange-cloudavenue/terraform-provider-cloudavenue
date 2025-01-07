/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	powerON = "POWERED_ON"
	// powerOFF = "POWERED_OFF".
)

type VMResourceModelState struct { //nolint:revive
	PowerON types.Bool   `tfsdk:"power_on"`
	Status  types.String `tfsdk:"status"`
}

// attrTypes() returns the types of the attributes of the State attribute.
func (s *VMResourceModelState) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"power_on": types.BoolType,
		"status":   types.StringType,
	}
}

// toAttrValues() returns the values of the attributes of the State attribute.
func (s *VMResourceModelState) toAttrValues() map[string]attr.Value {
	return map[string]attr.Value{
		"power_on": s.PowerON,
		"status":   s.Status,
	}
}

// ToPlan returns the value of the State attribute, if set, as a types.Object.
func (s *VMResourceModelState) ToPlan(ctx context.Context) types.Object {
	if s == nil {
		return types.Object{}
	}

	// If the VM is powered on, set the power_on attribute to true, otherwise false.
	if s.Status.ValueString() == powerON {
		s.PowerON = types.BoolValue(true)
	} else {
		s.PowerON = types.BoolValue(false)
	}

	return types.ObjectValueMust(s.attrTypes(), s.toAttrValues())
}

// StateRead returns the value of the State attribute, if set, as a *VMResourceModelState.
func (v VM) StateRead(ctx context.Context) (*VMResourceModelState, error) {
	status, err := v.GetStatus()
	if err != nil {
		return nil, fmt.Errorf("error getting status: %w", err)
	}

	return &VMResourceModelState{
		PowerON: types.BoolValue(v.IsPoweredON()),
		Status:  types.StringValue(status),
	}, nil
}
