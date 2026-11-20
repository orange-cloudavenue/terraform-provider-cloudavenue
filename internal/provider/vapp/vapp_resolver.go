// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vapp

import (
	"cmp"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// resolveVApp returns the parent vApp identified by its ID or, if the ID is empty, by its name.
// The vApp must belong to the given VDC.
func (c vAppClientV2) resolveVApp(ctx context.Context, vdc resolvedVDC, vappID, vappName string) (*types.ModelVApp, diag.Diagnostics) {
	var diags diag.Diagnostics

	p := types.ParamsGetVApp{}
	if vappID != "" {
		p.ID = vappID
	} else {
		p.Name = vappName
	}

	app, err := c.vappClient.GetVApp(ctx, p)
	if err != nil {
		cerrs.AddError(&diags, cerrs.ActionRead, fmt.Sprintf("parent vApp %s", cmp.Or(p.ID, p.Name)), err)
		return nil, diags
	}

	if !vappInVDC(app, vdc) {
		diags.AddError(
			"Parent vApp not found in VDC",
			fmt.Sprintf("vApp %s(%s) is in VDC %s, expected VDC %s(%s)", app.Name, app.ID, app.VDC, vdc.Name, vdc.ID),
		)
		return nil, diags
	}

	return app, diags
}

// lockVApp locks the vApp and returns the function that unlocks it.
// The returned function is always safe to call.
func lockVApp(ctx context.Context, vdcName, vappName string) (unlock func(), diags diag.Diagnostics) {
	if err := mutex.LockVApp(ctx, vdcName, vappName); err != nil {
		diags.AddError("Invalid lock arguments", err.Error())
		return func() {}, diags
	}

	return func() {
		_ = mutex.UnlockVApp(ctx, vdcName, vappName)
	}, diags
}
