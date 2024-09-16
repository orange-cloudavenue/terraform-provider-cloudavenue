package bms

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	BMSModelDatasource struct {
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
		BMSType           supertypes.StringValue                                             `tfsdk:"bms_type"`
		OS                supertypes.StringValue                                             `tfsdk:"os"`
		BiosConfiguration supertypes.StringValue                                             `tfsdk:"bios_configuration"`
		Storage           supertypes.SingleNestedObjectValueOf[bmsModelDatasourceBMSStorage] `tfsdk:"storage"`
	}

	bmsModelDatasourceBMSStorage struct {
		Local  supertypes.SingleNestedObjectValueOf[bmsModelDatasourceBMSStorageGen] `tfsdk:"local"`
		System supertypes.SingleNestedObjectValueOf[bmsModelDatasourceBMSStorageGen] `tfsdk:"system"`
		Data   supertypes.SingleNestedObjectValueOf[bmsModelDatasourceBMSStorageGen] `tfsdk:"data"`
		Shared supertypes.SingleNestedObjectValueOf[bmsModelDatasourceBMSStorageGen] `tfsdk:"shared"`
	}

	bmsModelDatasourceBMSStorageGen struct {
		Size         supertypes.StringValue `tfsdk:"size"`
		StorageClass supertypes.StringValue `tfsdk:"storage_class"`
	}
)

// NewbmsModelDatasource returns a new bmsModelDatasource.
func NewbmsModelDatasource() *BMSModelDatasource {
	return &BMSModelDatasource{}
}

// SetID sets the ID.
func (m *BMSModelDatasource) SetID(id *string) {
	m.ID.Set(*id)
}

// SetNetwork sets the Network of BMS listed.
func (m *BMSModelDatasource) SetNetwork(ctx context.Context, bms *v1.BMS) (net []*bmsModelDatasourceNetwork, err error) {
	networks := bms.GetNetworks()
	if len(networks) == 0 {
		return make([]*bmsModelDatasourceNetwork, 0), nil
	}
	for _, network := range networks {
		var x bmsModelDatasourceNetwork
		x.VLANID.Set(network.VLANID)
		x.Subnet.Set(network.Subnet)
		x.Prefix.Set(network.Prefix)
		net = append(net, &x)
	}

	return net, nil
}

// SetBMS sets the BMS information.
func (m *BMSModelDatasource) SetBMS(ctx context.Context, bms *v1.BMS) (bmsList []*bmsModelDatasourceBMS, err error) {
	bmsDetails := bms.GetBMSDetails()
	if len(bmsDetails) == 0 {
		return make([]*bmsModelDatasourceBMS, 0), nil
	}
	for _, bms := range bmsDetails {
		x := &bmsModelDatasourceBMS{}
		x.Hostname.Set(bms.Hostname)
		x.BMSType.Set(bms.BMSType)
		x.OS.Set(bms.OS)
		x.BiosConfiguration.Set(bms.BiosConfiguration)
		x.Storage.Set(ctx, setStorage(ctx, bms.GetBMSStorage()))
		bmsList = append(bmsList, x)
	}
	return bmsList, nil
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

func setStorageDetail(storage v1.BMSStorageDetail) (storageDetail *bmsModelDatasourceBMSStorageGen) {
	storageDetail = &bmsModelDatasourceBMSStorageGen{}
	storageDetail.Size.Set(storage.Size)
	storageDetail.StorageClass.Set(storage.StorageClass)
	return
}

// Copy returns a new bmsModelDatasource.
func (m *BMSModelDatasource) Copy() any {
	x := &BMSModelDatasource{}
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

// New bmsModelDatasourceNetwork returns a new bmsModelDatasourceNetwork.
func newBMSModelDatasourceNetwork() *bmsModelDatasourceNetwork {
	return &bmsModelDatasourceNetwork{
		VLANID: supertypes.NewStringNull(),
		Subnet: supertypes.NewStringNull(),
		Prefix: supertypes.NewStringNull(),
	}
}

// New bmsModelDatasourceBMS returns a new bmsModelDatasourceBMS.
func newBMSModelDatasourceBMS(ctx context.Context) *bmsModelDatasourceBMS {
	return &bmsModelDatasourceBMS{
		Hostname:          supertypes.NewStringNull(),
		BMSType:           supertypes.NewStringNull(),
		OS:                supertypes.NewStringNull(),
		BiosConfiguration: supertypes.NewStringNull(),
		Storage:           supertypes.NewSingleNestedObjectValueOfNull[bmsModelDatasourceBMSStorage](ctx),
	}
}

// New bmsModelDatasourceBMSStorage returns a new bmsModelDatasourceBMSStorage.
func newBMSModelDatasourceBMSStorage(ctx context.Context) *bmsModelDatasourceBMSStorage {
	return &bmsModelDatasourceBMSStorage{
		Local:  supertypes.NewSingleNestedObjectValueOfNull[bmsModelDatasourceBMSStorageGen](ctx),
		System: supertypes.NewSingleNestedObjectValueOfNull[bmsModelDatasourceBMSStorageGen](ctx),
		Data:   supertypes.NewSingleNestedObjectValueOfNull[bmsModelDatasourceBMSStorageGen](ctx),
		Shared: supertypes.NewSingleNestedObjectValueOfNull[bmsModelDatasourceBMSStorageGen](ctx),
	}
}
