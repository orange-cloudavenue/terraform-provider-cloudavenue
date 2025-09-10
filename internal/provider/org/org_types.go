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
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	OrgModel struct { //nolint:revive
		ID                   supertypes.StringValue                                  `tfsdk:"id"`
		Name                 supertypes.StringValue                                  `tfsdk:"name"`
		Description          supertypes.StringValue                                  `tfsdk:"description"`
		FullName             supertypes.StringValue                                  `tfsdk:"full_name"`
		Enabled              supertypes.BoolValue                                    `tfsdk:"enabled"`
		Resources            supertypes.SingleNestedObjectValueOf[OrgModelResources] `tfsdk:"resources"`
		Email                supertypes.StringValue                                  `tfsdk:"email"`
		InternetBillingModel supertypes.StringValue                                  `tfsdk:"internet_billing_mode"`
	}

	OrgModelResources struct { //nolint:revive
		CountVDC       supertypes.Int64Value `tfsdk:"count_vdc"`
		CountCatalog   supertypes.Int64Value `tfsdk:"count_catalog"`
		CountVApp      supertypes.Int64Value `tfsdk:"count_vapp"`
		CountRunningVM supertypes.Int64Value `tfsdk:"count_running_vm"`
		CountUser      supertypes.Int64Value `tfsdk:"count_user"`
		CountDisk      supertypes.Int64Value `tfsdk:"count_disk"`
	}
)

func (rm *OrgModel) Copy() *OrgModel {
	x := &OrgModel{}
	utils.ModelCopy(rm, x)
	return x
}
