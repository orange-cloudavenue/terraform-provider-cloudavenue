package vdc

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	categoryName = "vdc"
)

// Struct for VDC on Cloud Avenue API

type vdcStorageProfileModel struct {
	Class   types.String `tfsdk:"class"`
	Limit   types.Int64  `tfsdk:"limit"`
	Default types.Bool   `tfsdk:"default"`
}
