// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package mutex

import (
	"context"
	"errors"
	"fmt"
)

// ErrInvalidVAppLockArgs is returned when the vDC name or the vApp name is empty.
var ErrInvalidVAppLockArgs = errors.New("invalid vApp lock arguments")

// vappMutex is shared by every vApp lock user (vapp and vm packages).
var vappMutex = NewKV()

// vappLockKey builds the lock key. The scheme is shared with `common/vapp`
// (vm package) and must stay name based to remain compatible.
func vappLockKey(vdcName, vappName string) (string, error) {
	if vdcName == "" || vappName == "" {
		return "", fmt.Errorf("%w: vDC %q vApp %q", ErrInvalidVAppLockArgs, vdcName, vappName)
	}
	return fmt.Sprintf("vdc:%s|vapp:%s", vdcName, vappName), nil
}

// LockVApp locks the vApp identified by its vDC name and vApp name.
// The caller must call UnlockVApp with the same arguments.
func LockVApp(ctx context.Context, vdcName, vappName string) error {
	if ctx == nil {
		return fmt.Errorf("%w: nil context", ErrInvalidVAppLockArgs)
	}
	key, err := vappLockKey(vdcName, vappName)
	if err != nil {
		return err
	}
	vappMutex.KvLock(ctx, key)
	return nil
}

// UnlockVApp unlocks the vApp identified by its vDC name and vApp name.
func UnlockVApp(ctx context.Context, vdcName, vappName string) error {
	if ctx == nil {
		return fmt.Errorf("%w: nil context", ErrInvalidVAppLockArgs)
	}
	key, err := vappLockKey(vdcName, vappName)
	if err != nil {
		return err
	}
	vappMutex.KvUnlock(ctx, key)
	return nil
}
