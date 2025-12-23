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

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi/rules"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	vdcResourceModel struct {
		Timeouts            timeouts.Value                                                       `tfsdk:"timeouts"`
		ID                  supertypes.StringValue                                               `tfsdk:"id"`
		Name                supertypes.StringValue                                               `tfsdk:"name"`
		Description         supertypes.StringValue                                               `tfsdk:"description"`
		ServiceClass        supertypes.StringValue                                               `tfsdk:"service_class"`
		DisponibilityClass  supertypes.StringValue                                               `tfsdk:"disponibility_class"`
		BillingModel        supertypes.StringValue                                               `tfsdk:"billing_model"`
		VCPUInMhz           supertypes.Int64Value                                                `tfsdk:"cpu_speed_in_mhz"`
		CPUAllocated        supertypes.Int64Value                                                `tfsdk:"cpu_allocated"`
		MemoryAllocated     supertypes.Int64Value                                                `tfsdk:"memory_allocated"`
		StorageBillingModel supertypes.StringValue                                               `tfsdk:"storage_billing_model"`
		StorageProfiles     supertypes.SetNestedObjectValueOf[vdcResourceModelVDCStorageProfile] `tfsdk:"storage_profiles"`
	}

	vdcResourceModelVDCStorageProfile struct {
		Class   supertypes.StringValue `tfsdk:"class"`
		Limit   supertypes.Int64Value  `tfsdk:"limit"`
		Default supertypes.BoolValue   `tfsdk:"default"`
	}

	vdcDataSourceModel struct {
		ID                  supertypes.StringValue                                               `tfsdk:"id"`
		Name                supertypes.StringValue                                               `tfsdk:"name"`
		Description         supertypes.StringValue                                               `tfsdk:"description"`
		ServiceClass        supertypes.StringValue                                               `tfsdk:"service_class"`
		DisponibilityClass  supertypes.StringValue                                               `tfsdk:"disponibility_class"`
		BillingModel        supertypes.StringValue                                               `tfsdk:"billing_model"`
		VCPUInMhz           supertypes.Int64Value                                                `tfsdk:"cpu_speed_in_mhz"`
		CPUAllocated        supertypes.Int64Value                                                `tfsdk:"cpu_allocated"`
		MemoryAllocated     supertypes.Int64Value                                                `tfsdk:"memory_allocated"`
		StorageBillingModel supertypes.StringValue                                               `tfsdk:"storage_billing_model"`
		StorageProfiles     supertypes.SetNestedObjectValueOf[vdcResourceModelVDCStorageProfile] `tfsdk:"storage_profiles"`
	}
)

func (rm *vdcResourceModel) Copy() *vdcResourceModel {
	x := &vdcResourceModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *vdcResourceModel) ToCAVVirtualDataCenter(ctx context.Context) (obj *infrapi.CAVVirtualDataCenter, diags diag.Diagnostics) {
	// Prepare the body to create a VDC.
	obj = &infrapi.CAVVirtualDataCenter{
		VDC: infrapi.CAVVirtualDataCenterVDC{
			Name:                rm.Name.Get(),
			Description:         rm.Description.Get(),
			ServiceClass:        rules.ServiceClass(rm.ServiceClass.Get()),
			DisponibilityClass:  rules.DisponibilityClass(rm.DisponibilityClass.Get()),
			BillingModel:        rules.BillingModel(rm.BillingModel.Get()),
			VCPUInMhz:           rm.VCPUInMhz.GetInt(),
			CPUAllocated:        rm.CPUAllocated.GetInt(),
			MemoryAllocated:     rm.MemoryAllocated.GetInt(),
			StorageBillingModel: rules.BillingModel(rm.StorageBillingModel.Get()),
		},
	}

	storageProfiles, d := rm.StorageProfiles.Get(ctx)
	diags.Append(d...)
	if d.HasError() {
		return obj, diags
	}

	// Iterate over the storage profiles and add them to the body.
	for _, storageProfile := range storageProfiles {
		obj.VDC.StorageProfiles = append(obj.VDC.StorageProfiles, infrapi.StorageProfile{
			Class:   infrapi.StorageProfileClass(storageProfile.Class.Get()),
			Limit:   storageProfile.Limit.GetInt(),
			Default: storageProfile.Default.Get(),
		})
	}

	return obj, diags
}
