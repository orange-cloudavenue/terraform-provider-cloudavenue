package vapp

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
)

const (
	ErrVAppNotFound = "VApp not found"
)

type diagnosticError struct {
	Summary string
	Detail  string
}

type networkRef struct {
	VDC         *govcd.Vdc
	VApp        *govcd.VApp
	VAppRef     vapp.Ref
	VAppLocked  bool
	VAppUnlockF func()
}
