/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	OrgModel struct { //nolint:revive
		ID                  supertypes.StringValue `tfsdk:"id"`
		Name                supertypes.StringValue `tfsdk:"name"`
		Description         supertypes.StringValue `tfsdk:"description"`
		FullName            supertypes.StringValue `tfsdk:"full_name"`
		Enabled             supertypes.BoolValue   `tfsdk:"enabled"`
		Email               supertypes.StringValue `tfsdk:"email"`
		InternetBillingMode supertypes.StringValue `tfsdk:"internet_billing_mode"`
	}
)

func (rm *OrgModel) Copy() *OrgModel {
	x := &OrgModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (data *OrgModel) fromModel(ctx context.Context, o *types.ModelGetOrganization) (diags diag.Diagnostics) {
	// ctx kept for future use; avoid unused param linter error after resources.* removal
	_ = ctx
	if o == nil {
		diags.AddError("Error reading organization", "Received nil organization from API")
		return diags
	}
	data.ID.Set(o.ID)
	data.Name.Set(o.Name)
	data.Description.Set(o.Description)
	data.FullName.Set(o.FullName)
	data.Enabled.Set(o.Enabled)
	data.Email.Set(o.Email)
	data.InternetBillingMode.Set(o.InternetBillingMode)

	return diags
}
