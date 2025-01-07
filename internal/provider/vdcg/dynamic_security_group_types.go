/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	DynamicSecurityGroupModel struct {
		ID           supertypes.StringValue                                                `tfsdk:"id"`
		Name         supertypes.StringValue                                                `tfsdk:"name"`
		Description  supertypes.StringValue                                                `tfsdk:"description"`
		VDCGroupName supertypes.StringValue                                                `tfsdk:"vdc_group_name"`
		VDCGroupID   supertypes.StringValue                                                `tfsdk:"vdc_group_id"`
		Criteria     supertypes.ListNestedObjectValueOf[DynamicSecurityGroupModelCriteria] `tfsdk:"criteria"`
	}

	DynamicSecurityGroupModelCriteria struct {
		Rules supertypes.ListNestedObjectValueOf[DynamicSecurityGroupModelRule] `tfsdk:"rules"`
	}

	DynamicSecurityGroupModelRule struct {
		Type     supertypes.StringValue `tfsdk:"type"`
		Value    supertypes.StringValue `tfsdk:"value"`
		Operator supertypes.StringValue `tfsdk:"operator"`
	}
)

func (rm *DynamicSecurityGroupModel) Copy() *DynamicSecurityGroupModel {
	x := &DynamicSecurityGroupModel{}
	utils.ModelCopy(rm, x)
	return x
}

// ToSDKDynamicSecurityGroupModel converts the model to the SDK model.
func (rm *DynamicSecurityGroupModel) ToSDKDynamicSecurityGroupModel(ctx context.Context) (*v1.FirewallGroupDynamicSecurityGroupModel, diag.Diagnostics) {
	sdkCriterias := make(v1.FirewallGroupDynamicSecurityGroupModelCriterias, 0)

	vmCriteria, d := rm.Criteria.Get(ctx)
	if d.HasError() {
		return nil, d
	}

	for _, criteria := range vmCriteria {
		sdkRules := make(v1.FirewallGroupDynamicSecurityGroupModelCriteriaRules, 0)

		rules, d := criteria.Rules.Get(ctx)
		if d.HasError() {
			return nil, d
		}

		for _, rule := range rules {
			sdkRules = append(sdkRules, v1.FirewallGroupDynamicSecurityGroupModelCriteriaRule{
				RuleType: v1.FirewallGroupDynamicSecurityGroupModelCriteriaRuleType(rule.Type.Get()),
				Value:    rule.Value.Get(),
				Operator: v1.FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator(rule.Operator.Get()),
			})
		}

		sdkCriterias = append(sdkCriterias, v1.FirewallGroupDynamicSecurityGroupModelCriteria{
			Rules: sdkRules,
		})
	}

	return &v1.FirewallGroupDynamicSecurityGroupModel{
		FirewallGroupModel: v1.FirewallGroupModel{
			ID:          rm.ID.Get(),
			Name:        rm.Name.Get(),
			Description: rm.Description.Get(),
		},
		Criteria: sdkCriterias,
	}, nil
}
