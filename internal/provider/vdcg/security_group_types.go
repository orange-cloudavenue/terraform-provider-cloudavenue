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

type SecurityGroupModel struct {
	ID           supertypes.StringValue        `tfsdk:"id"`
	Name         supertypes.StringValue        `tfsdk:"name"`
	Description  supertypes.StringValue        `tfsdk:"description"`
	VDCGroupName supertypes.StringValue        `tfsdk:"vdc_group_name"`
	VDCGroupID   supertypes.StringValue        `tfsdk:"vdc_group_id"`
	Members      supertypes.SetValueOf[string] `tfsdk:"member_org_network_ids"`
}

func (rm *SecurityGroupModel) ToSDKSecurityGroupModel(ctx context.Context) (*v1.FirewallGroupSecurityGroupModel, diag.Diagnostics) {
	members, d := rm.Members.Get(ctx)
	if d.HasError() {
		return nil, d
	}

	return &v1.FirewallGroupSecurityGroupModel{
		FirewallGroupModel: v1.FirewallGroupModel{
			ID:          rm.ID.Get(),
			Name:        rm.Name.Get(),
			Description: rm.Description.Get(),
		},
		Members: utils.SliceIDToOpenAPIReference(members),
	}, nil
}

func (rm *SecurityGroupModel) Copy() *SecurityGroupModel {
	x := &SecurityGroupModel{}
	utils.ModelCopy(rm, x)
	return x
}
