// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vapp

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
)

// resolvedVDC is the VDC a vApp resource or data source operates in.
type resolvedVDC struct {
	// ID is the VDC URN (urn:vcloud:vdc:<uuid>), as expected by the vApp SDK params (e.g. ParamsComposeVApp.VDCID).
	ID string
	// Name is the VDC name, used for the `vdc` attribute in state.
	Name string
}

// resolveVDC resolves a VDC name to its ID using SDK v2.
// If vdcName is empty, the provider-level default VDC is used; with no default
// either, an error diagnostic wrapping client.ErrEmptyVDCNameProvided is returned.
// Replaces the govcd-based VDC init for the vapp package.
func (c vAppClientV2) resolveVDC(ctx context.Context, vdcName string) (resolvedVDC, diag.Diagnostics) {
	var diags diag.Diagnostics
	if vdcName == "" {
		if !c.client.DefaultVDCExist() {
			cerrs.AddError(&diags, cerrs.ActionRead, "VDC", client.ErrEmptyVDCNameProvided)
			return resolvedVDC{}, diags
		}
		vdcName = c.client.GetDefaultVDC()
	}

	// ListVDC (exact name filter) is required: GET /vdc/{name} does not return the
	// VDC ID (see vapp-sdk-v2-issues.md, B5), only the listing exposes it.
	vdcs, err := c.vdcClient.ListVDC(ctx, types.ParamsListVDC{Name: vdcName})
	if err != nil {
		cerrs.AddError(&diags, cerrs.ActionRead, fmt.Sprintf("VDC %s", vdcName), err)
		return resolvedVDC{}, diags
	}
	if len(vdcs.VDCS) == 0 {
		cerrs.AddError(&diags, cerrs.ActionRead, fmt.Sprintf("VDC %s", vdcName), fmt.Errorf("vdc %q not found", vdcName))
		return resolvedVDC{}, diags
	}

	// The VDC listing may expose a bare UUID: the vApp SDK params expect the URN.
	return resolvedVDC{ID: urn.Normalize(urn.VDC, vdcs.VDCS[0].ID).String(), Name: vdcs.VDCS[0].Name}, diags
}

// vappInVDC reports whether the vApp belongs to the VDC.
// TODO(sdk): ModelVApp.VDC is empty with the real API (see vapp-sdk-v2-issues.md, B2):
// the check is skipped in that case until the SDK fills it.
func vappInVDC(app *types.ModelVApp, vdc resolvedVDC) bool {
	return app.VDC == "" || app.VDC == vdc.ID
}
