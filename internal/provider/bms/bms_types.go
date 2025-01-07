/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package bms

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	bmsModelDatasource struct {
		ID       supertypes.StringValue                                   `tfsdk:"id"`
		Timeouts timeoutsD.Value                                          `tfsdk:"timeouts"`
		Env      supertypes.SetNestedObjectValueOf[bmsModelDatasourceEnv] `tfsdk:"env"`
	}

	bmsModelDatasourceEnv struct {
		Network supertypes.SetNestedObjectValueOf[bmsModelDatasourceNetwork] `tfsdk:"network"`
		BMS     supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMS]     `tfsdk:"bms"`
	}

	bmsModelDatasourceNetwork struct {
		VLANID supertypes.StringValue `tfsdk:"vlan_id"`
		Subnet supertypes.StringValue `tfsdk:"subnet"`
		Prefix supertypes.StringValue `tfsdk:"prefix"`
	}

	bmsModelDatasourceBMS struct {
		Hostname          supertypes.StringValue                                             `tfsdk:"hostname"`
		BMSType           supertypes.StringValue                                             `tfsdk:"type"`
		OS                supertypes.StringValue                                             `tfsdk:"os"`
		BiosConfiguration supertypes.StringValue                                             `tfsdk:"bios_configuration"`
		Storage           supertypes.SingleNestedObjectValueOf[bmsModelDatasourceBMSStorage] `tfsdk:"storage"`
	}

	bmsModelDatasourceBMSStorage struct {
		Local  supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMSStorageDetail] `tfsdk:"local"`
		System supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMSStorageDetail] `tfsdk:"system"`
		Data   supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMSStorageDetail] `tfsdk:"data"`
		Shared supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMSStorageDetail] `tfsdk:"shared"`
	}

	bmsModelDatasourceBMSStorageDetail struct {
		Size         supertypes.StringValue `tfsdk:"size"`
		StorageClass supertypes.StringValue `tfsdk:"storage_class"`
	}
)

// NewbmsModelDatasource returns a new bmsModelDatasource.
func NewBMSModelDatasource(ctx context.Context) *bmsModelDatasource {
	return &bmsModelDatasource{
		ID:       supertypes.NewStringNull(),
		Timeouts: timeoutsD.Value{},
		Env:      supertypes.NewSetNestedObjectValueOfNull[bmsModelDatasourceEnv](ctx),
	}
}

// Put API Network information to Terraform Object.
func NetworkToTerraform(bms *v1.BMS) (net []*bmsModelDatasourceNetwork) {
	for _, network := range bms.GetNetworks() {
		var x bmsModelDatasourceNetwork
		x.VLANID.Set(network.VLANID)
		x.Subnet.Set(network.Subnet)
		x.Prefix.Set(network.Prefix)
		net = append(net, &x)
	}

	return
}

// Put API BMS information to Terraform Object.
func BMSToTerraform(ctx context.Context, bms *v1.BMS) (bmsList []*bmsModelDatasourceBMS) { //nolint: revive
	for _, bms := range bms.GetBMS() {
		x := &bmsModelDatasourceBMS{}
		x.Hostname.Set(bms.Hostname)
		x.BMSType.Set(bms.BMSType)
		x.OS.Set(bms.OS)
		x.BiosConfiguration.Set(bms.BiosConfiguration)
		x.Storage.Set(ctx, setStorage(ctx, bms.GetStorages()))
		bmsList = append(bmsList, x)
	}
	return bmsList
}

// SetStorage sets the Storage of BMS.
func setStorage(ctx context.Context, storages v1.BMSStorage) (storage *bmsModelDatasourceBMSStorage) {
	storage = &bmsModelDatasourceBMSStorage{}
	storage.Local.Set(ctx, setStorageDetail(storages.GetLocal()))
	storage.System.Set(ctx, setStorageDetail(storages.GetSystem()))
	storage.Data.Set(ctx, setStorageDetail(storages.GetData()))
	storage.Shared.Set(ctx, setStorageDetail(storages.GetShared()))
	return
}

func setStorageDetail(storage []v1.BMSStorageDetail) (storageDetail []*bmsModelDatasourceBMSStorageDetail) {
	// for each storage detail, set the size and storage class
	for _, stor := range storage {
		x := &bmsModelDatasourceBMSStorageDetail{}
		x.Size.Set(stor.GetSize())
		x.StorageClass.Set(stor.GetStorageClass())
		storageDetail = append(storageDetail, x)
	}
	return
}

// Copy returns a new bmsModelDatasource.
func (m *bmsModelDatasource) Copy() any {
	x := &bmsModelDatasource{}
	utils.ModelCopy(m, x)
	return x
}

// New bmsModelDatasourceEnv returns a new bmsModelDatasourceEnv.
func newBMSModelDatasourceEnv(ctx context.Context) *bmsModelDatasourceEnv {
	return &bmsModelDatasourceEnv{
		Network: supertypes.NewSetNestedObjectValueOfNull[bmsModelDatasourceNetwork](ctx),
		BMS:     supertypes.NewSetNestedObjectValueOfNull[bmsModelDatasourceBMS](ctx),
	}
}
