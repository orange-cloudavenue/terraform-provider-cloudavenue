package vdc

import "github.com/hashicorp/terraform-plugin-framework/types"

type vdcsDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	VDCs []vdcRef     `tfsdk:"vdcs"`
}

type vdcRef struct {
	VDCName types.String `tfsdk:"vdc_name"`
	VDCUuid types.String `tfsdk:"vdc_uuid"`
}
