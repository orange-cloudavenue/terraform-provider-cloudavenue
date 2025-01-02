package vapp

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

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
