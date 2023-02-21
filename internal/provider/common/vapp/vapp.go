// Package vapp provides common functionality for vApp resources.
package vapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ref is a reference to a vApp.
type Ref struct {
	Org   string
	VDC   string
	Name  string
	TFCtx context.Context
}

var (
	// ErrVAppRefEmpty is returned when a vApp reference is missing information.
	ErrVAppRefEmpty = errors.New("missing information in vapp ref")
	vcdMutexKV      = mutex.NewKV()
)

// LockParentVApp locks the parent vApp.
func (v *Ref) LockParentVApp() error {
	if v.Org == "" || v.VDC == "" || v.Name == "" || v.TFCtx == nil {
		return ErrVAppRefEmpty
	}
	key := fmt.Sprintf("org:%s|vdc:%s|vapp:%s", v.Org, v.VDC, v.Name)
	vcdMutexKV.KvLock(v.TFCtx, key)
	return nil
}

// UnLockParentVApp unlocks the parent vApp.
func (v *Ref) UnLockParentVApp() error {
	if v.Org == "" || v.VDC == "" || v.Name == "" || v.TFCtx == nil {
		return ErrVAppRefEmpty
	}
	key := fmt.Sprintf("org:%s|vdc:%s|vapp:%s", v.Org, v.VDC, v.Name)
	vcdMutexKV.KvUnlock(v.TFCtx, key)
	return nil
}
