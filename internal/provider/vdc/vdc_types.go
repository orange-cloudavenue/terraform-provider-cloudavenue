/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

type (
	vdcModel struct {
		ID          supertypes.StringValue `tfsdk:"id"`
		Name        supertypes.StringValue `tfsdk:"name"`
		Description supertypes.StringValue `tfsdk:"description"`

		// * Availability properties
		ServiceClass       supertypes.StringValue `tfsdk:"service_class"`
		DisponibilityClass supertypes.StringValue `tfsdk:"disponibility_class"`

		// * Billing properties
		BillingModel        supertypes.StringValue `tfsdk:"billing_model"`
		StorageBillingModel supertypes.StringValue `tfsdk:"storage_billing_model"`

		// * Resource properties
		VCPU            supertypes.Int64Value                                     `tfsdk:"vcpu"`
		Memory          supertypes.Int64Value                                     `tfsdk:"memory"`
		StorageProfiles supertypes.SetNestedObjectValueOf[vdcModelStorageProfile] `tfsdk:"storage_profiles"`

		// * Deprecated fields - Maintain for backward compatibility
		VCPUInMhz       supertypes.Int64Value `tfsdk:"cpu_speed_in_mhz"`
		CPUAllocated    supertypes.Int64Value `tfsdk:"cpu_allocated"`
		MemoryAllocated supertypes.Int64Value `tfsdk:"memory_allocated"`
	}

	vdcModelStorageProfile struct {
		ID      supertypes.StringValue `tfsdk:"id"`
		Class   supertypes.StringValue `tfsdk:"class"`
		Limit   supertypes.Int64Value  `tfsdk:"limit"`
		Default supertypes.BoolValue   `tfsdk:"default"`
		Used    supertypes.Int64Value  `tfsdk:"used"`
	}
)

// fromSDK converts SDK object to model
func (rm *vdcModel) fromSDK(ctx context.Context, vdc *types.ModelGetVDC, sp *types.ModelListStorageProfiles, diags *diag.Diagnostics) {
	if rm == nil || vdc == nil || sp == nil {
		diags.AddError("Error in vdcModel.fromSDK", "Cannot convert to model from SDK: vdcModel, ModelGetVDC or ModelListStorageProfiles is nil")
		return
	}

	rm.ID.Set(vdc.ID)
	rm.Name.Set(vdc.Name)
	rm.Description.Set(vdc.Description)

	// * Resource properties
	rm.ServiceClass.Set(vdc.Properties.ServiceClass)
	rm.DisponibilityClass.Set(vdc.Properties.DisponibilityClass)

	// * Billing properties
	rm.BillingModel.Set(vdc.Properties.BillingModel)
	rm.StorageBillingModel.Set(vdc.Properties.StorageBillingModel)

	// * Resource properties
	rm.VCPU.SetInt(vdc.ComputeCapacity.CPU.Limit)
	rm.Memory.SetInt(vdc.ComputeCapacity.Memory.Limit)

	sps := []*vdcModelStorageProfile{}
	for _, sp := range sp.VDCS[0].StorageProfiles {
		sptf := &vdcModelStorageProfile{}
		sptf.ID.Set(sp.ID)
		sptf.Class.Set(sp.Class)
		sptf.Limit.SetInt(sp.Limit)
		sptf.Default.Set(sp.Default)
		sptf.Used.SetInt(sp.Used)
		sps = append(sps, sptf)
	}

	diags.Append(rm.StorageProfiles.Set(ctx, sps)...)
	if diags.HasError() {
		return
	}

	// ! Deprecated fields - Maintain for backward compatibility
	rm.VCPUInMhz.SetInt(vdc.ComputeCapacity.CPU.VCPUFrequency)
	rm.CPUAllocated.SetInt(vdc.ComputeCapacity.CPU.FrequencyLimit)
	rm.MemoryAllocated.SetInt(vdc.ComputeCapacity.Memory.Limit)
}
