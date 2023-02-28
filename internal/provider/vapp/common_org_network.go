package vapp

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
)

const (
	ErrVAppNotFound = "VApp not found"
)

type diagnosticError struct {
	Summary string
	Detail  string
}

type orgNetworkRef struct {
	VDC         *govcd.Vdc
	VApp        *govcd.VApp
	VAppRef     vapp.Ref
	VAppLocked  bool
	VAppUnlockF func()
}

func (s *orgNetworkResourceModel) getVDCName(client *client.CloudAvenue) *diagnosticError {
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

func (s *orgNetworkResourceModel) initOrgNetworkQuery(ctx context.Context, client *client.CloudAvenue, lockVApp bool) (*orgNetworkRef, *diagnosticError) {
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

	orgNetworkRef := &orgNetworkRef{
		VDC: vdc,
	}

	if lockVApp {
		orgNetworkRef.VAppRef = vapp.Ref{
			Name:  s.VAppName.ValueString(),
			Org:   client.GetOrg(),
			VDC:   s.VDC.ValueString(),
			TFCtx: ctx,
		}

		if errLock := orgNetworkRef.VAppRef.LockParentVApp(); errors.Is(errLock, vapp.ErrVAppRefEmpty) {
			return nil, &diagnosticError{
				Summary: "Error locking vApp",
				Detail:  "Empty name, org or vdc in vapp.VAppRef",
			}
		}
		orgNetworkRef.VAppLocked = true
		orgNetworkRef.VAppUnlockF = func() {
			if errUnlock := orgNetworkRef.VAppRef.UnLockParentVApp(); errUnlock != nil {
				// tflog print error is enough ?
				tflog.Warn(ctx, "Error unlock vapp")
			}
		}
	}

	vappName := s.VAppName.ValueString()
	vapp, err := vdc.GetVAppByName(vappName, false)
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

	orgNetworkRef.VApp = vapp

	return orgNetworkRef, nil
}

func (s *orgNetworkResourceModel) findOrgNetwork(vAppNetworkConfig *govcdtypes.NetworkConfigSection) (*govcdtypes.VAppNetworkConfiguration, *string, *diagnosticError) {
	var vAppNetwork govcdtypes.VAppNetworkConfiguration
	var networkID string

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.Link != nil {
			id, err := govcd.GetUuidFromHref(networkConfig.Link.HREF, false)
			if err != nil {
				return nil, nil, &diagnosticError{
					Summary: "Unable to get network ID from HREF",
					Detail:  err.Error(),
				}
			}
			// name check needed for datasource to find network as don't have ID
			if common.ExtractUUID(s.ID.ValueString()) == common.ExtractUUID(id) || networkConfig.NetworkName == s.NetworkName.ValueString() {
				networkID = id
				vAppNetwork = networkConfig
				break
			}
		}
	}
	return &vAppNetwork, &networkID, nil
}
