package bms

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	bmsModelDatasource struct {
		ID       supertypes.StringValue                                       `tfsdk:"id"`
		Timeouts timeoutsD.Value                                              `tfsdk:"timeouts"`
		Network  supertypes.SetNestedObjectValueOf[bmsModelDatasourceNetwork] `tfsdk:"network"`
		BMS      supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMS]     `tfsdk:"bms"`
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

// NewBMSModelDatasource returns a new BMSModelDatasource.
func NewBMSModelDatasource() *bmsModelDatasource {
	return &bmsModelDatasource{}
}

// SetID sets the ID.
func (m *bmsModelDatasource) SetID(id *string) {
	m.ID.Set(*id)
}

// SetNetwork sets the Network of BMS listed.
func (m *bmsModelDatasource) SetNetwork(ctx context.Context, bms *v1.BMS) error {
	bmsService, err := bms.List()
	if err != nil {
		return err
	}
	networks := bmsService.GetNetworks()
	if len(networks) == 0 {
		return nil
	}
	var net []*bmsModelDatasourceNetwork
	for _, network := range networks {
		var x bmsModelDatasourceNetwork
		x.VLANID.Set(network.VLANID)
		x.Subnet.Set(network.Subnet)
		x.Prefix.Set(network.Prefix)
		net = append(net, &x)
	}
	m.Network.Set(ctx, net)
	return nil
}

// SetBMS sets the BMS information.
func (m *bmsModelDatasource) SetBMS(ctx context.Context, bms *v1.BMS) error {
	bmsService, err := bms.List()
	if err != nil {
		return err
	}
	bmsDetails := bmsService.GetBMSDetails()
	var bmsList []*bmsModelDatasourceBMS
	for _, bms := range bmsDetails {
		x := &bmsModelDatasourceBMS{}
		x.Hostname.Set(bms.Hostname)
		x.BMSType.Set(bms.BMSType)
		x.OS.Set(bms.OS)
		x.BiosConfiguration.Set(bms.BiosConfiguration)
		x.Storage.Set(ctx, setStorage(ctx, bms.GetBMSStorage()))
		bmsList = append(bmsList, x)
	}
	m.BMS.Set(ctx, bmsList)
	return nil
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

// Copy returns a new BMSModelDatasource.
func (m *bmsModelDatasource) Copy() any {
	x := &bmsModelDatasource{}
	utils.ModelCopy(m, x)
	return x
}
