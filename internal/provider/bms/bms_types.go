package bms

import (
	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

type bmsModelDatasource struct {
	ID       supertypes.StringValue                                       `tfsdk:"id"`
	Timeouts timeoutsD.Value                                              `tfsdk:"timeouts"`
	Network  supertypes.SetNestedObjectValueOf[bmsModelDatasourceNetwork] `tfsdk:"network"`
	BMS      supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMS]     `tfsdk:"bms"`
}

type bmsModelDatasourceNetwork struct {
	VLANID supertypes.StringValue `tfsdk:"vlan_id"`
	Subnet supertypes.StringValue `tfsdk:"subnet"`
	Prefix supertypes.StringValue `tfsdk:"prefix"`
}

type bmsModelDatasourceBMS struct {
	Hostname          supertypes.StringValue                                          `tfsdk:"hostname"`
	BMSType           supertypes.StringValue                                          `tfsdk:"bms_type"`
	OS                supertypes.StringValue                                          `tfsdk:"os"`
	BiosConfiguration supertypes.StringValue                                          `tfsdk:"bios_configuration"`
	Storage           supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMSStorage] `tfsdk:"storage"`
}

type bmsModelDatasourceBMSStorage struct {
	Local  supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMSStorageGen] `tfsdk:"local"`
	System supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMSStorageGen] `tfsdk:"system"`
	Data   supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMSStorageGen] `tfsdk:"data"`
	Shared supertypes.SetNestedObjectValueOf[bmsModelDatasourceBMSStorageGen] `tfsdk:"shared"`
}

type bmsModelDatasourceBMSStorageGen struct {
	Size         supertypes.StringValue `tfsdk:"size"`
	StorageClass supertypes.StringValue `tfsdk:"storage_class"`
}
