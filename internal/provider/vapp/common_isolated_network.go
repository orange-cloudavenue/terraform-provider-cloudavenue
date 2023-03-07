package vapp

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
)

func (s *isolatedNetworkResourceModel) getVDCName(client *client.CloudAvenue) *diagnosticError {
	// If VDC is not defined at data source level, use the one defined at provider level
	if s.VDC.IsNull() || s.VDC.IsUnknown() {
		if client.DefaultVDCExist() {
			s.VDC = types.StringValue(client.GetDefaultVDC())
		} else {
			return &diagnosticError{
				Summary: "Missing VDC",
				Detail:  "VDC is required when not defined at provider level",
			}
		}
	}
	return nil
}

func (s *isolatedNetworkResourceModel) initNetworkQuery(ctx context.Context, client *client.CloudAvenue, lockVApp bool) (*networkRef, *diagnosticError) {
	if err := s.getVDCName(client); err != nil {
		return nil, err
	}

	_, vdc, err := client.GetOrgAndVDC(client.GetOrg(), s.VDC.ValueString())
	if err != nil {
		return nil, &diagnosticError{
			Summary: "Error retrieving VDC",
			Detail:  err.Error(),
		}
	}

	networkRef := &networkRef{
		VDC: vdc,
	}

	if lockVApp {
		networkRef.VAppRef = vapp.Ref{
			Name:  s.VAppName.ValueString(),
			Org:   client.GetOrg(),
			VDC:   s.VDC.ValueString(),
			TFCtx: ctx,
		}

		if errLock := networkRef.VAppRef.LockParentVApp(); errors.Is(errLock, vapp.ErrVAppRefEmpty) {
			return nil, &diagnosticError{
				Summary: "Error locking vApp",
				Detail:  "Empty name, org or vdc in vapp.VAppRef",
			}
		}
		networkRef.VAppLocked = true
		networkRef.VAppUnlockF = func() {
			if errUnlock := networkRef.VAppRef.UnLockParentVApp(); errUnlock != nil {
				// tflog print error is enough ?
				tflog.Warn(ctx, "Error unlock vapp")
			}
		}
	}

	vapp, err := vdc.GetVAppByName(s.VAppName.ValueString(), false)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, &diagnosticError{
				Summary: ErrVAppNotFound,
				Detail:  err.Error(),
			}
		}
		return nil, &diagnosticError{
			Summary: "Error retrieving vApp",
			Detail:  err.Error(),
		}
	}

	networkRef.VApp = vapp

	return networkRef, nil
}
