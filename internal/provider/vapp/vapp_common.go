package vapp

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
)

type vappResourceModel struct {
	VAppName        types.String `tfsdk:"name"`
	VAppID          types.String `tfsdk:"id"`
	VDC             types.String `tfsdk:"vdc"`
	Description     types.String `tfsdk:"description"`
	PowerON         types.Bool   `tfsdk:"power_on"`
	GuestProperties types.Map    `tfsdk:"guest_properties"`
	Lease           types.Object `tfsdk:"lease"`
}

func processGuestProperties(vapp vapp.VAPP) (properties map[string]attr.Value, d diag.Diagnostics) {
	guestProperties, err := vapp.GetProductSectionList()
	if err != nil {
		d.AddError("Error retrieving guest properties", err.Error())
		return
	}

	properties = make(map[string]attr.Value)
	if guestProperties != nil && guestProperties.ProductSection != nil && guestProperties.ProductSection.Property != nil {
		for _, guestProperty := range guestProperties.ProductSection.Property {
			if guestProperty.Value != nil {
				properties[guestProperty.Key] = types.StringValue(guestProperty.Value.Value)
			}
		}
	}
	return
}
