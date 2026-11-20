// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vapp

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	vappResourceModel struct {
		VAppName        supertypes.StringValue                                       `tfsdk:"name"`
		VAppID          supertypes.StringValue                                       `tfsdk:"id"`
		VDC             supertypes.StringValue                                       `tfsdk:"vdc"`
		Description     supertypes.StringValue                                       `tfsdk:"description"`
		GuestProperties supertypes.MapValue                                          `tfsdk:"guest_properties"`
		Lease           supertypes.SingleNestedObjectValueOf[vappResourceModelLease] `tfsdk:"lease"`
	}

	vappResourceModelLease struct {
		RuntimeLeaseInSec supertypes.Int64Value `tfsdk:"runtime_lease_in_sec"`
		StorageLeaseInSec supertypes.Int64Value `tfsdk:"storage_lease_in_sec"`
	}
)

func (rm *vappResourceModel) Copy() *vappResourceModel {
	x := &vappResourceModel{}
	utils.ModelCopy(rm, x)
	return x
}

// fromSDK fills the model from an SDK vApp.
func (rm *vappResourceModel) fromSDK(ctx context.Context, v *sdktypes.ModelVApp, vdcName string) (diags diag.Diagnostics) {
	rm.VAppID.Set(v.ID)
	rm.VAppName.Set(v.Name)
	rm.VDC.Set(vdcName)
	rm.Description.Set(v.Description)

	if v.DeploymentLeaseInSeconds != nil && v.StorageLeaseInSeconds != nil {
		lease := new(vappResourceModelLease)
		lease.RuntimeLeaseInSec.SetIntPtr(v.DeploymentLeaseInSeconds)
		lease.StorageLeaseInSec.SetIntPtr(v.StorageLeaseInSeconds)
		diags.Append(rm.Lease.Set(ctx, lease)...)
	}

	return diags
}

// guestPropertiesFromSDK fills GuestProperties with key/value pairs of every section.
// When the API returns no property the attribute is set to null (an explicit empty map
// is kept), so a property that was not stored is a visible difference and not hidden by
// the value of the plan.
func (rm *vappResourceModel) guestPropertiesFromSDK(ctx context.Context, gp *sdktypes.ModelVAppGuestProperties) (diags diag.Diagnostics) {
	if gp == nil {
		return diags
	}

	properties := make(map[string]string)
	for _, section := range gp.Sections {
		for key, prop := range section.Properties {
			properties[key] = prop.Value
		}
	}

	if len(properties) > 0 {
		diags.Append(rm.GuestProperties.Set(ctx, properties)...)
		return diags
	}

	// An explicit empty map in the configuration is kept as is.
	if !rm.GuestProperties.IsNull() && !rm.GuestProperties.IsUnknown() && len(rm.GuestProperties.Elements()) == 0 {
		return diags
	}

	rm.GuestProperties.SetNull(ctx)

	return diags
}

// guestPropertiesToSDK builds the whole-list replacement sent to UpdateVAppGuestProperties.
func (rm *vappResourceModel) guestPropertiesToSDK(ctx context.Context) (gp sdktypes.ModelVAppGuestProperties, diags diag.Diagnostics) {
	// One ProductSection without class nor instance: the SDK sends the "Custom properties" info.
	section := sdktypes.ModelVAppGuestPropertySection{
		Properties: map[string]sdktypes.ModelVAppGuestProperty{},
	}

	if !rm.GuestProperties.IsNull() {
		values := make(map[string]string)
		diags.Append(rm.GuestProperties.Get(ctx, &values, true)...)
		if diags.HasError() {
			return gp, diags
		}

		for k, v := range values {
			section.Properties[k] = sdktypes.ModelVAppGuestProperty{
				Key:              k,
				Label:            k,
				Type:             "string",
				Value:            v,
				UserConfigurable: true,
			}
		}
	}

	gp.Sections = []sdktypes.ModelVAppGuestPropertySection{section}
	return gp, diags
}
